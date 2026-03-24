package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/KamerFa/quran-api/internal/db"
	"github.com/KamerFa/quran-api/internal/prefs"
	"github.com/KamerFa/quran-api/internal/quran"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("data", "quran.db")
	}

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()

	// Register API handlers
	prefsStore := prefs.NewStore(database)
	prefsHandler := prefs.NewHandler(prefsStore)
	prefsHandler.RegisterRoutes(mux)

	quranHandler := quran.NewHandler()
	quranHandler.RegisterRoutes(mux)

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Serve static files (the frontend)
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = ".."
	}
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))

	// Wrap with middleware
	handler := sessionMiddleware(corsMiddleware(mux))

	fmt.Printf("Quran API server starting on :%s\n", port)
	fmt.Printf("Database: %s\n", dbPath)
	fmt.Printf("Static files: %s\n", staticDir)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// sessionMiddleware ensures every request has a session cookie.
func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		if err != nil {
			id := generateSessionID()
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    id,
				Path:     "/",
				MaxAge:   365 * 24 * 60 * 60, // 1 year
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			r.AddCookie(&http.Cookie{Name: "session_id", Value: id})
		}
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
