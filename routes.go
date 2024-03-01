package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	http.ServeFile(w, r, "./public/index.html")
}

func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		LongURL string `json:"longURL"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier si un lien existant pour cette URL est déjà expiré
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database("url_shortener").Collection("urls")
	var existingURLMapping URLMapping
	err = collection.FindOne(ctx, bson.M{"original_url": requestBody.LongURL}).Decode(&existingURLMapping)
	if err == nil && existingURLMapping.ExpirationDate.Before(time.Now()) {
		// Le lien existe déjà et il est expiré, donc mettez à jour le lien existant
		updateResult, err := collection.UpdateOne(ctx, bson.M{"original_url": requestBody.LongURL}, bson.M{
			"$set": bson.M{
				"short_url":       GenerateShortURL(requestBody.LongURL),
				"created_at":      time.Now(),
				"expiration_date": time.Now().Add(1 * time.Hour),
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Updated %d document(s)\n", updateResult.ModifiedCount)

		// Répondre avec l'URL raccourcie existante
		response := map[string]string{"shortURL": existingURLMapping.ShortURL}

		// Définir le type de contenu de la réponse comme JSON
		w.Header().Set("Content-Type", "application/json")

		// Écrire le code de réponse 200 OK
		w.WriteHeader(http.StatusOK)

		// Écrire la réponse JSON dans le corps de la réponse
		json.NewEncoder(w).Encode(response)

		// Afficher l'URL raccourcie dans les logs
		fmt.Println("Shortened URL:", existingURLMapping.ShortURL)
	} else if err == mongo.ErrNoDocuments || (err == nil && existingURLMapping.ExpirationDate.After(time.Now())) {
		// Si le lien n'existe pas ou s'il n'est pas expiré, créez un nouveau lien
		// Générer une URL raccourcie unique à partir de l'URL d'origine
		shortURL := GenerateShortURL(requestBody.LongURL)

		// Enregistrer le mapping dans la base de données
		_, err := collection.InsertOne(ctx, bson.M{
			"original_url":    requestBody.LongURL,
			"short_url":       shortURL,
			"created_at":      time.Now(),
			"expiration_date": time.Now().Add(1 * time.Hour),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Répondre avec l'URL raccourcie générée
		response := map[string]string{"shortURL": "http://localhost:1234/" + shortURL}

		// Définir le type de contenu de la réponse comme JSON
		w.Header().Set("Content-Type", "application/json")

		// Écrire le code de réponse 200 OK
		w.WriteHeader(http.StatusOK)

		// Écrire la réponse JSON dans le corps de la réponse
		json.NewEncoder(w).Encode(response)

		// Afficher l'URL raccourcie dans les logs
		fmt.Println("Shortened URL:", "http://localhost:1234/"+shortURL)
	} else if err != nil && err != mongo.ErrNoDocuments {
		// Une erreur s'est produite lors de la recherche du lien existant
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RedirectToLongURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	fmt.Println("Short URL from request:", shortURL)

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
		// Afficher l'erreur dans les logs
		fmt.Println("Error fetching URL mapping:", err)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Vérifier si le lien a expiré
	if result.ExpirationDate.Before(time.Now()) {
		fmt.Println("URL has expired:", result.ExpirationDate)
		http.Error(w, "URL has expired", http.StatusNotFound)
		return
	}

	// Afficher l'URL longue dans les logs
	fmt.Println("Original URL:", result.OriginalURL)

	// Rediriger vers l'URL longue correspondante
	http.Redirect(w, r, result.OriginalURL, http.StatusMovedPermanently)
}
