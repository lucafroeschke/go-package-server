package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lucafroeschke/go-package-server/config"
	"github.com/lucafroeschke/go-package-server/logger"
	"github.com/lucafroeschke/go-package-server/templates"
	"html/template"
	"net/http"
)

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()

	tmpl, err := template.ParseFS(templates.Templates, "index.html")
	if err != nil {
		logger.WriteLog(logger.ERROR, fmt.Sprintf("Failed to parse template: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, cfg)
}

func handlePackagePage(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()

	vars := mux.Vars(r)
	packageName := vars["package"]

	pkg, exists := config.GetPackage(packageName)

	if !exists {
		http.NotFound(w, r)
		return
	}

	if r.URL.Query().Get("go-get") != "1" {
		http.Redirect(w, r, pkg.Repository, http.StatusFound)
		return
	}

	tmpl, err := template.ParseFS(templates.Templates, "package.html")
	if err != nil {
		logger.WriteLog(logger.ERROR, fmt.Sprintf("Failed to parse template: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, PackageResponse{Config: *cfg, Package: *pkg})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.WriteLog(logger.INFO, fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func Start() error {
	cfg := config.GetConfig()
	addr := fmt.Sprintf("%s:%d", cfg.ListeningAddress, cfg.ListeningPort)

	r := mux.NewRouter()
	r.HandleFunc("/", handleIndexPage)
	r.HandleFunc("/{package}", handlePackagePage)

	if cfg.LogRequests {
		r.Use(logRequest)
	}

	http.Handle("/", r)

	logger.WriteLog(logger.INFO, fmt.Sprintf("Listening on %s", addr))
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		return err
	}

	return nil
}
