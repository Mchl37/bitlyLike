# BitlyLike

BitlyLike est une application de raccourcissement d'URL écrite en Go. Elle vous permet de créer des versions raccourcies de vos URLs longues pour faciliter le partage sur les réseaux sociaux, les messages, ou tout autre moyen de communication. Cette application utilise une base de données MongoDB pour stocker les mappings entre les URLs longues et courtes.

## Installation

Pour installer et exécuter BitlyLike sur votre propre machine, suivez ces étapes :
1 - Cloner le dépôt GitHub :<br/>
`git clone https://github.com/Mchl37/bitlyLike.git`

2 - Accéder au répertoire du projet :<br/>
`cd bitlyLike`

3 - Assurez-vous d'avoir Go installé sur votre système. Si ce n'est pas le cas, suivez les instructions d'installation depuis le site officiel de Go : https://golang.org/doc/install

4 - Installez les dépendances :<br/>
`go mod tidy`

5 - Créez un fichier .env à la racine du projet et définissez les variables d'environnement nécessaires. Vous pouvez vous référer au fichier .env.example pour connaître les variables requises.

6 - Assurez-vous d'avoir MongoDB installé sur votre système. Si ce n'est pas le cas, suivez les instructions d'installation depuis le site officiel de MongoDB : https://docs.mongodb.com/manual/installation/

7 - Lancez l'application :<br/>
`go run .`

8 - L'application devrait maintenant être accessible à l'adresse http://localhost:1234.

## Modèle MongoDB
ID             string    `json:"_id" bson:"_id,omitempty"`<br/>
OriginalURL    string    `json:"original_url" bson:"original_url"`<br/>
ShortURL       string    `json:"short_url" bson:"short_url"`<br/>
CreatedAt      time.Time `json:"created_at" bson:"created_at"`<br/>
ExpirationDate time.Time `json:"expiration_date" bson:"expiration_date"`

## Utilisation

Une fois l'application en cours d'exécution, vous pouvez accéder à l'interface utilisateur à l'adresse http://localhost:1234 dans votre navigateur. Vous pouvez entrer une URL longue dans le champ prévu à cet effet, puis cliquer sur le bouton "Shorten" pour obtenir une version raccourcie de l'URL. La version raccourcie sera affichée sous le champ de saisie une fois générée.

## Auteurs

- Johan MARIN
- Michel GUELIN
- Anass AOULAD
