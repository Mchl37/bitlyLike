package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectToLongURL(t *testing.T) {
	// Créer une requête GET avec un shortURL fictif
	req, err := http.NewRequest("GET", "/MGJhND", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Créer un enregistreur de réponse factice
	rr := httptest.NewRecorder()

	// Créer un gestionnaire HTTP à partir de la fonction redirectToLongURL
	handler := http.HandlerFunc(redirectToLongURL)

	// Appeler la fonction de gestionnaire HTTP, en passant la requête fictive et l'enregistreur de réponse fictif
	handler.ServeHTTP(rr, req)

	// Vérifier le code de statut de la réponse
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMovedPermanently)
	}

	// Vérifier que la redirection se fait vers l'URL longue attendue
	expectedURL := "https://www.dealabs.com/bons-plans"
	if location := rr.Header().Get("Location"); location != expectedURL {
		t.Errorf("handler returned unexpected redirect location: got %v want %v",
			location, expectedURL)
	}
}

func TestCreateShortURL(t *testing.T) {
	// Corps de la requête JSON avec une URL longue
	requestBody := `{"longURL": "https://www.dealabs.com/bons-plans"}`

	// Crée une requête POST avec le corps JSON
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Crée un enregistreur de réponse factice
	rr := httptest.NewRecorder()

	// Crée un gestionnaire HTTP à partir de la fonction createShortURL
	handler := http.HandlerFunc(createShortURL)

	// Appel la fonction de gestionnaire HTTP, en passant la requête fictive et l'enregistreur de réponse fictif
	handler.ServeHTTP(rr, req)

	// Afficher les détails de la réponse enregistrée
	fmt.Println("Response Body:", rr.Body.String())
	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Headers:", rr.Header())
}
