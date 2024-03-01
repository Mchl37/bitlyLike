package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Charger les variables d'environnement à partir du fichier .env
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Servir les fichiers statiques
	workDir, _ := os.Getwd()
	filesDir := http.Dir(workDir + "/public")
	FileServer(r, "/assets", filesDir)

	// Routes
	r.Get("/", GetRoot)
	r.Post("/shorten", CreateShortURL)
	r.Get("/{shortURL}", RedirectToLongURL)

	// Démarrer le serveur
	err = http.ListenAndServe(":1234", r)
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

// FileServer crée un gestionnaire de fichiers statiques
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer ne prend pas de caractères de substitution ou de pièces globales.")
	}
	fs := http.StripPrefix(path, http.FileServer(root))
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
