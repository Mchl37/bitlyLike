package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDBConnection(t *testing.T) {
	println("TEST DB START")
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

	println("TEST DB END")
}

func TestGenerateShortURL(t *testing.T) {
	println("TEST GENERATE SHORT URL Start")
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
	println("TEST GENERATE SHORT URL END")
}
func TestRedirectToLongURL(t *testing.T) {
	println("TEST RESDIRECT TO URL STart")
	// Créer une requête HTTP GET avec un shortURL comme paramètre d'URL
	req, err := http.NewRequest("GET", "/shortURL", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Créer un enregistreur de réponse HTTP factice
	rr := httptest.NewRecorder()

	// Exécuter la fonction à tester
	RedirectToLongURL(rr, req)

	// Vérifier le code de statut de la réponse
	assert.Equal(t, http.StatusMovedPermanently, rr.Code)

	// Vérifier l'en-tête de localisation de la réponse
	expectedURL := "https://www.dealabs.com/bons-plans" // L'URL longue attendue
	assert.Equal(t, expectedURL, rr.Header().Get("Location"))
	println("TEST RESDIRECT TO URL END")
}
func TestGetRoot(t *testing.T) {
	println("TEST GET ROOT START")
	// Crée une requête HTTP GET pour la racine "/"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crée un enregistreur de réponse HTTP factice
	rr := httptest.NewRecorder()

	// Exécute la fonction à tester
	GetRoot(rr, req)

	// Vérifie le code de statut de la réponse
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Vérifie le contenu de la réponse
	expectedContentType := "text/html; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned unexpected content type: got %v want %v",
			contentType, expectedContentType)
	}
	println("TEST GET ROOT END")
}
func TestFileServer(t *testing.T) {
	println("TEST FILE SERVER START")
	// Crée un routeur Chi
	r := chi.NewRouter()

	// Appele la fonction FileServer avec un chemin et un système de fichiers factices
	// Pour ce test, nous utilisons un système de fichiers simple avec un seul fichier
	testFS := http.Dir("./testdata")
	FileServer(r, "/assets", testFS)

	// Créeune requête HTTP GET pour un fichier dans le répertoire de test
	req, err := http.NewRequest("GET", "/assets/testfile.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Créer un enregistreur de réponse HTTP factice
	rr := httptest.NewRecorder()

	// Exécuter la requête via le routeur
	r.ServeHTTP(rr, req)

	// Vérifie le code de statut de la réponse
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Vérifie le contenu de la réponse
	expectedContent := "Test file content."
	if rr.Body.String() != expectedContent {
		t.Errorf("handler returned unexpected content: got %v want %v",
			rr.Body.String(), expectedContent)
	}
	println("TEST FILE SERVER END")
}
func TestRedirectToNonExistentURL(t *testing.T) {
	println("TEST REDIRECT T NON EXISTENT URL START")
	// Crée une requête HTTP GET avec un shortURL inexistant comme paramètre d'URL
	req, err := http.NewRequest("GET", "/nonexistentURL", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crée un enregistreur de réponse HTTP factice
	rr := httptest.NewRecorder()

	// Exécuter la fonction à tester
	RedirectToLongURL(rr, req)

	// Vérifiele code de statut de la réponse
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Vérifie le corps de la réponse pour s'assurer qu'il contient un message d'erreur approprié
	expectedError := "URL not found"
	if rr.Body.String() != expectedError {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedError)
	}

	println("TEST REDIRECT T NON EXISTENT URL END")
}
