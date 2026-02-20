package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

//go:embed static
var staticFiles embed.FS

type Server struct {
	port    int
	host    string
	dataDir string
	mux     *http.ServeMux
}

func New(port int, host string, dataDir string) *Server {
	return &Server{
		port:    port,
		host:    host,
		dataDir: dataDir,
		mux:     http.NewServeMux(),
	}
}

func (s *Server) Setup() error {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}

	// Ensure attachments directory exists
	attachDir := filepath.Join(s.dataDir, "attachments")
	if err := os.MkdirAll(attachDir, 0755); err != nil {
		return fmt.Errorf("failed to create attachments dir: %w", err)
	}

	// Ensure data.json exists
	dataPath := filepath.Join(s.dataDir, "data.json")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		emptyCanvas := []byte(`{"nodes":[],"edges":[]}`)
		if err := os.WriteFile(dataPath, emptyCanvas, 0644); err != nil {
			return fmt.Errorf("failed to create data.json: %w", err)
		}
		log.Printf("INFO [server] Created empty canvas at %s", dataPath)
	}

	// Static assets
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// API routes
	s.mux.HandleFunc("GET /api/canvas", s.handleGetCanvas)
	s.mux.HandleFunc("PUT /api/canvas", s.handlePutCanvas)
	s.mux.HandleFunc("GET /api/md-files", s.handleListMDFiles)
	s.mux.HandleFunc("GET /api/files/{path...}", s.handleGetFile)
	s.mux.HandleFunc("POST /api/upload", s.handleUpload)

	// Serve index.html at root
	s.mux.HandleFunc("/", s.handleIndex)

	return nil
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	log.Printf("INFO [server] Starting on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func (s *Server) handleGetCanvas(w http.ResponseWriter, r *http.Request) {
	dataPath := filepath.Join(s.dataDir, "data.json")
	data, err := os.ReadFile(dataPath)
	if err != nil {
		http.Error(w, "Failed to read canvas data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) handlePutCanvas(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 10<<20)) // 10MB max
	if err != nil {
		http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Validate JSON structure
	var canvas map[string]interface{}
	if err := json.Unmarshal(body, &canvas); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dataPath := filepath.Join(s.dataDir, "data.json")
	if err := os.WriteFile(dataPath, body, 0644); err != nil {
		http.Error(w, "Failed to save canvas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleListMDFiles(w http.ResponseWriter, r *http.Request) {
	var files []string
	filepath.WalkDir(s.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		// Skip hidden directories
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		// Skip attachments directory
		if d.IsDir() && d.Name() == "attachments" {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			relPath, _ := filepath.Rel(s.dataDir, path)
			files = append(files, relPath)
		}
		return nil
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (s *Server) handleGetFile(w http.ResponseWriter, r *http.Request) {
	reqPath := r.PathValue("path")
	if reqPath == "" {
		http.Error(w, "Path required", http.StatusBadRequest)
		return
	}

	// Prevent path traversal
	cleanPath := filepath.Clean(reqPath)
	if strings.Contains(cleanPath, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(s.dataDir, cleanPath)

	// Verify the resolved path is within dataDir
	absData, _ := filepath.Abs(s.dataDir)
	absFile, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFile, absData) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Detect content type
	ext := strings.ToLower(filepath.Ext(fullPath))
	switch ext {
	case ".md":
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".pdf":
		w.Header().Set("Content-Type", "application/pdf")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Write(data)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20) // 50MB max
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	newName := uuid.New().String() + ext
	attachDir := filepath.Join(s.dataDir, "attachments")
	destPath := filepath.Join(attachDir, newName)

	dst, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	relativePath := filepath.Join("attachments", newName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path":     relativePath,
		"filename": header.Filename,
	})
}
