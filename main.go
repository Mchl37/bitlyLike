package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLMapping struct {
	ID          string    `json:"_id" bson:"_id,omitempty"`
	OriginalURL string    `json:"original_url" bson:"original_url"`
	ShortURL    string    `json:"short_url" bson:"short_url"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	http.ServeFile(w, r, "./public/index.html")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		LongURL string `json:"longURL"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Générer une URL raccourcie unique à partir de l'URL d'origine
	shortURL := generateShortURL(requestBody.LongURL)

	// Enregistrer le mapping dans la base de données
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database("url_shortener").Collection("urls")
	_, err = collection.InsertOne(ctx, bson.M{"original_url": requestBody.LongURL, "short_url": shortURL, "created_at": time.Now()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Répondre avec l'URL raccourcie générée
	response := map[string]string{"shortURL": shortURL}

	// Définir le type de contenu de la réponse comme JSON
	w.Header().Set("Content-Type", "application/json")

	// Écrire le code de réponse 200 OK
	w.WriteHeader(http.StatusOK)

	// Écrire la réponse JSON dans le corps de la réponse
	json.NewEncoder(w).Encode(response)
}

func redirectToLongURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")

	// Rechercher le mapping dans la base de données
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database("url_shortener").Collection("urls")
	var result URLMapping
	err = collection.FindOne(ctx, bson.M{"short_url": shortURL}).Decode(&result)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Rediriger vers l'URL longue correspondante
	http.Redirect(w, r, result.OriginalURL, http.StatusMovedPermanently)
}

func generateShortURL(originalURL string) string {
	// Générer un hachage unique pour l'URL d'origine
	hash := md5.Sum([]byte(originalURL))
	// Convertir le hachage en une chaîne hexadécimale
	hashStr := hex.EncodeToString(hash[:])
	// Convertir la chaîne hexadécimale en base64 pour obtenir une chaîne plus courte
	shortBytes := base64.StdEncoding.EncodeToString([]byte(hashStr))
	// Formater l'URL raccourcie
	shortURL := "http://localhost:1234/" + shortBytes[:6] // Utilisez les 6 premiers caractères du hash
	return shortURL
}

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
	r.Get("/", getRoot)
	r.Get("/hello", getHello)
	r.Post("/shorten", createShortURL)
	r.Get("/{shortURL}", redirectToLongURL)

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
