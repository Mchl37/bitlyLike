package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"time"
)

type URLMapping struct {
	ID             string    `json:"_id" bson:"_id,omitempty"`
	OriginalURL    string    `json:"original_url" bson:"original_url"`
	ShortURL       string    `json:"short_url" bson:"short_url"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	ExpirationDate time.Time `json:"expiration_date" bson:"expiration_date"`
}

func GenerateShortURL(originalURL string) string {
	// Génération du hachage unique pour l'URL d'origine
	hash := md5.Sum([]byte(originalURL))
	// Convertion du hachage en une chaîne hexadécimale
	hashStr := hex.EncodeToString(hash[:])
	// Convertion de la chaîne hexadécimale en base64 pour obtenir une chaîne plus courte
	shortBytes := base64.StdEncoding.EncodeToString([]byte(hashStr))
	// Retourne les 6 premiers caractères du hash
	return shortBytes[:6]
}
