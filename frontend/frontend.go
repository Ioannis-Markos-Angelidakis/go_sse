package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	buildPath, err := filepath.Abs("./build")
	if err != nil {
		log.Fatal("Error resolving build path:", err)
	}

	fs := http.FileServer(http.Dir(buildPath))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "public, max-age=31536000")
		fs.ServeHTTP(w, r)
	})

	port := os.Getenv("FRONTEND_PORT")
	if port == "" {
		port = "80"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Frontend server running on port %s", port)
	log.Fatal(server.ListenAndServe())
}
