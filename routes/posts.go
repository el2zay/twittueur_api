package routes

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"twittueur_api/models"
	"twittueur_api/src/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Fonction pour créé un id pour chaque post.
func generatePostId() string {
	// On récupère l'heure actuelle au format Unix
	timestamp := time.Now().Unix()
	// Transformer la variable en Hex
	hexId := fmt.Sprintf("%x", timestamp)

	return hexId
}

// Fonction pour poster des posts
func PostData(c echo.Context) error {
	id := generatePostId()
	body := c.FormValue("body")
	date := c.FormValue("date")
	device := c.FormValue("device")
	comment := c.FormValue("comment")

	// Récupérer le token du header
	authorization := c.Request().Header.Get("Authorization")

	if authorization == "" {
		return c.JSON(400, models.Response{Success: false, Message: "Vous devez renseigner un token."})
	}

	tokenString := authorization[7:] // On ignore les 7 premières lettres (Bearer )
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return c.JSON(400, models.Response{Success: false, Message: "Token invalide"})
	}

	// Récupérer l'passphrase à partir du token
	var passphrase string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		passphrase = claims["passphrase"].(string)
	}

	// Vérifier si l'utilisateur existe
	err = utils.IsPassphraseExists(c, passphrase)
	if err != nil {
		return err
	}

	// Les vérifications
	if body == "" {
		return c.JSON(400, models.Response{Success: false, Message: "Le body est requis."})
	}

	if date == "" {
		return c.JSON(400, models.Response{Success: false, Message: "La date est requise."})
	}

	if device == "" {
		return c.JSON(400, models.Response{Success: false, Message: "L'appareil est requis.'"})
	}

	if len(body) > 1000 {
		return c.JSON(400, models.Response{Success: false, Message: "Le texte du body ne doit pas dépasser les 1000 caractères"})
	}

	if comment != "" {
		// Détecter si le post existe
		err := utils.IsPostExists(c, comment)
		if err != nil {
			return errors.New("le post n'existe pas")
		}

		var data models.PostRequest
		if err := viper.Unmarshal(&data); err != nil {
			return err
		}

		// Trouver le post avec l'ID correspondant à la valeur de comment
		for i, post := range data.Posts {
			if post.ID == comment {
				// Ajouter l'ID du nouveau commentaire au tableau des commentaires du post
				data.Posts[i].Comments = append(data.Posts[i].Comments, id)
				break
			}
		}

		// Réécrire le fichier avec le nouveau contenu
		viper.Set("posts", data.Posts)
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}

	post := &models.Post{
		ID:         id,
		Body:       body,
		Date:       date,
		Device:     device,
		Passphrase: passphrase,
		Likedby:    []string{},
		IsComment:  comment != "",
	}

	// Lire l'image
	image, err := c.FormFile("image")
	if err != nil {
		if err != http.ErrMissingFile {
			// Si une autre erreur se produit, retournez une réponse d'erreur
			return c.JSON(400, models.Response{Success: false, Message: "Une erreur s'est produite lors de la lecture de l'image"})
		}
		// S'il n'y a pas d'image on continue.
	} else {
		// S'il y a une image on la lit.
		src, err := image.Open()
		if err != nil {
			// L'erreur 500 signifie qu'il y a une erreur côté serveur
			return c.JSON(500, models.Response{Success: false, Message: "Une erreur s'est produite lors de l'ouverture de l'imagee"})
		}

		defer src.Close()

		// Récupérer l'extension du fichier
		ext := filepath.Ext(image.Filename)
		imagePath := "db/images/" + id + ext

		// Créé un fichier dans le dossier db/images avec l'id du post
		dst, err := os.Create(imagePath)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copier l'image de la raquête vers celle du serveur.
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		post.Image = imagePath
	}

	// On vérifie si le fichier users.json existe, sinon on le créé
	if _, err := os.Stat("db/posts.json"); os.IsNotExist(err) {
		file, err := os.Create("db/posts.json")
		if err != nil {
			return err
		}
		defer file.Close()
		file.WriteString(`{"posts": []}`)
	}

	// Configuration de Viper
	viper.SetConfigName("posts")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Convertir la donnée en une structure go
	var data models.PostRequest
	if err := viper.Unmarshal(&data); err != nil {
		return err
	}
	// Ajouter le nouveau post à la liste des posts
	data.Posts = append(data.Posts, *post)

	// Réécrire le fichier avec le nouveau contenu
	viper.Set("posts", data.Posts)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Message: "Le post a bien été créé"})
}

// Fonction pour récupérer les posts

func GetPosts(c echo.Context) error {
	/* Paramètre pour ignorer les posts déjà chargé.
	 Le client enverra une liste d'id qu'il a déjà chargé et afficher
	Grâce à cette liste on pourra filtrer les posts à afficher pour éviter
	les doublons.*/

	idsParam := c.QueryParam("ids")
	showComments := c.QueryParam("showComments")
	ids := strings.Split(idsParam, ",") // Convertir la chaîne d'ids en parties

	// Configuration de Viper pour lire le fichier posts.json
	viper.SetConfigName("posts")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Convertir la donnée en une structure go
	var data models.PostRequest
	if err := viper.Unmarshal(&data); err != nil {
		return err
	}
	// Faire un random
	// Cette fonction échange chaque posts entre eux pour les "mélanger".
	rand.Shuffle(len(data.Posts), func(i, j int) { data.Posts[i], data.Posts[j] = data.Posts[j], data.Posts[i] })

	var posts []models.Post // Créer une slice vide pour les posts

	// Boucle pour ajouter les 10 premiers post à la List
	for _, post := range data.Posts {
		if contains(ids, post.ID) || (showComments == "false" && post.IsComment) { //Si l'id est dans la slice ou que c'est un commentaire
			continue //  alors que l'utilisateur n'en veut pas, alors on ignore ce post
		}

		if len(posts) == 10 {
			break
		}
		posts = append(posts, post) // Ajouter le post à la liste
	}

	// Retourner les 10 premiers posts
	return c.JSON(http.StatusOK, posts)
}

// la fonction contains vérifie si une partie contient une certaine valeur
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
