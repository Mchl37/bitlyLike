package main

import (
	"context"
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDBConnection(t *testing.T) {
	// Récupérer l'URI MongoDB depuis les variables d'environnement
	err := godotenv.Load()
	if err != nil {
		println("le fichier n'a pas chargé")
		panic("Error loading .env file")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		t.Fatal("MONGO_URI not set in environment variables")
	}

	// Connexion à MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Vérifier l'état de la connexion
	err = client.Ping(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to ping MongoDB server: %v", err)
	}

	t.Log("Successfully connected to MongoDB")
}
func TestRedirectToLongURL(t *testing.T) {
	// Crée une requête HTTP GET avec un shortURL comme paramètre d'URL
	req, err := http.NewRequest("GET", "/shortURL", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crée un enregistreur de réponse HTTP factice
	rr := httptest.NewRecorder()

	// Simule les données de la base de données
	mockURLMapping := URLMapping{
		OriginalURL:    "https://www.dealabs.com/bons-plans",
		ShortURL:       "shortURL",
		ExpirationDate: time.Now().Add(1 * time.Hour),
	}

	// Exécute la fonction à tester
	RedirectToLongURL(rr, req)

	// Vérifie le code de statut de la réponse
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

	// Vérifie l'en-tête de localisation de la réponse
	expectedURL := mockURLMapping.OriginalURL
	if rr.Header().Get("Location") != expectedURL {
		t.Errorf("handler returned wrong location header: got %v want %v",
			rr.Header().Get("Location"), expectedURL)
	}
}

func TestGenerateShortURL(t *testing.T) {
	// URL d'origine
	originalURL := "https://www.facebook.com"

	// Générer l'URL courte
	shortURL := GenerateShortURL(originalURL)

	// Vérifier la longueur de l'URL courte
	expectedLength := 6
	if len(shortURL) != expectedLength {
		t.Errorf("expected length of shortURL to be %d, got %d", expectedLength, len(shortURL))
	}

	// Vérifier si l'URL courte est non vide
	if shortURL == "" {
		t.Error("expected non-empty shortURL, got empty string")
	}

	// Réessaye avec une autre URL d'origine
	originalURL = "https://www.google.com"
	shortURL = GenerateShortURL(originalURL)

	// Vérifie à nouveau la longueur de l'URL courte
	if len(shortURL) != expectedLength {
		t.Errorf("expected length of shortURL to be %d, got %d", expectedLength, len(shortURL))
	}

	// Vérifie à nouveau si l'URL courte est non vide
	if shortURL == "" {
		t.Error("expected non-empty shortURL, got empty string")
	}
}
