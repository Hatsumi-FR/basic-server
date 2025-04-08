package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no file, continue...")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL is missing in .env")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		email, password, ok := parseBasicAuth(r)
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var storedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&storedPassword)
		if err != nil || storedPassword != password {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("authenticated as %s", email)))
	})

	fmt.Println("server running on :8080")
	_ = http.ListenAndServe(":8080", r)
}

func parseBasicAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		return
	}
	payload, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
	if err != nil {
		return
	}
	parts := strings.SplitN(string(payload), ":", 2)
	if len(parts) != 2 {
		return
	}
	return parts[0], parts[1], true
}
