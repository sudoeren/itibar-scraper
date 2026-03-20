package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	srv *http.Server
	svc *Service
}

func New(svc *Service, addr string) (*Server, error) {
	ans := Server{
		svc: svc,
		srv: &http.Server{
			Addr:              addr,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/v1/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			ans.apiScrape(w, r)
		case http.MethodGet:
			ans.apiGetJobs(w, r)
		default:
			renderJSON(w, http.StatusMethodNotAllowed, apiError{
				Code:    http.StatusMethodNotAllowed,
				Message: "Method not allowed",
			})
		}
	})

	mux.HandleFunc("/api/v1/jobs/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasSuffix(path, "/download") {
			ans.handleDownload(w, r)
			return
		}

		id := extractID(path)
		if id == "" {
			id = r.URL.Query().Get("id")
		}

		parsed, err := uuid.Parse(id)
		if err != nil {
			renderJSON(w, http.StatusBadRequest, apiError{
				Code:    http.StatusBadRequest,
				Message: "Invalid ID",
			})
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), idCtxKey, parsed))

		switch r.Method {
		case http.MethodGet:
			ans.apiGetJob(w, r)
		case http.MethodDelete:
			ans.apiDeleteJob(w, r)
		default:
			renderJSON(w, http.StatusMethodNotAllowed, apiError{
				Code:    http.StatusMethodNotAllowed,
				Message: "Method not allowed",
			})
		}
	})

	handler := corsHeaders(mux)
	ans.srv.Handler = handler

	return &ans, nil
}

func extractID(path string) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) >= 5 {
		id := parts[len(parts)-1]
		if !strings.HasSuffix(id, "download") {
			return id
		}
	}
	return ""
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		err := s.srv.Shutdown(context.Background())
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("server stopped")
	}()

	log.Printf("API server listening on %s\n", s.srv.Addr)

	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

type ctxKey string

const idCtxKey ctxKey = "id"

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type apiScrapeRequest struct {
	Name string
	JobData
}

type apiScrapeResponse struct {
	ID string `json:"id"`
}

func (s *Server) apiScrape(w http.ResponseWriter, r *http.Request) {
	var req apiScrapeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, apiError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	newJob := Job{
		ID:     uuid.New().String(),
		Name:   req.Name,
		Date:   time.Now().UTC(),
		Status: StatusPending,
		Data:   req.JobData,
	}

	if newJob.Data.MaxTime == 0 {
		newJob.Data.MaxTime = 600 * time.Second
	} else {
		newJob.Data.MaxTime *= time.Second
	}

	err = newJob.Validate()
	if err != nil {
		renderJSON(w, http.StatusUnprocessableEntity, apiError{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	err = s.svc.Create(r.Context(), &newJob)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	renderJSON(w, http.StatusCreated, apiScrapeResponse{ID: newJob.ID})
}

func (s *Server) apiGetJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := s.svc.All(r.Context())
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	renderJSON(w, http.StatusOK, jobs)
}

func (s *Server) apiGetJob(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(idCtxKey).(uuid.UUID)
	if !ok {
		renderJSON(w, http.StatusBadRequest, apiError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	job, err := s.svc.Get(r.Context(), id.String())
	if err != nil {
		renderJSON(w, http.StatusNotFound, apiError{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		})
		return
	}

	renderJSON(w, http.StatusOK, job)
}

func (s *Server) apiDeleteJob(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(idCtxKey).(uuid.UUID)
	if !ok {
		renderJSON(w, http.StatusBadRequest, apiError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	err := s.svc.Delete(r.Context(), id.String())
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		renderJSON(w, http.StatusBadRequest, apiError{
			Code:    http.StatusBadRequest,
			Message: "Invalid path",
		})
		return
	}

	id := parts[len(parts)-2]
	if id == "" {
		id = r.URL.Query().Get("id")
	}

	parsed, err := uuid.Parse(id)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, apiError{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		})
		return
	}

	_ = parsed

	filePath, err := s.svc.GetCSV(r.Context(), id)
	if err != nil {
		renderJSON(w, http.StatusNotFound, apiError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to open file",
		})
		return
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "text/csv")

	_, err = io.Copy(w, file)
	if err != nil {
		renderJSON(w, http.StatusInternalServerError, apiError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to send file",
		})
		return
	}
}

func renderJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func corsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
