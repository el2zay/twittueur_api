package routes

import (
	"net/http"
	"strconv"
	"twittueur_api/models"
	"twittueur_api/src/utils"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Posts length

func GlobalPostsLength(c echo.Context) error {
	// Configuration de Viper pour lire le fichier posts.json
	viper.SetConfigName("posts")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Convertir data en une structure go
	var data models.PostRequest
	if err := viper.Unmarshal(&data); err != nil {
		return c.JSON(500, models.Response{Message: "Une erreur s'est produite.", Success: false})
	}

	// Retourner la longueur du slice de posts
	return c.JSON(http.StatusOK, models.Response{Message: strconv.Itoa(len(data.Posts)), Success: true})
}

func UserPostsLength(c echo.Context) error {
	passphrase := c.FormValue("passphrase")
	if passphrase == "" {
		return c.JSON(400, models.Response{Message: "Vous devez renseigner la passphrase", Success: false})
	}

	// Configuration de Viper pour lire le fichier posts.json
	viper.SetConfigName("posts")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Convertir data en une structure go
	var data models.PostRequest
	if err := viper.Unmarshal(&data); err != nil {
		return c.JSON(500, models.Response{Message: "Une erreur s'est produite.", Success: false})
	}

	// Filtrer les posts de l'utilisateur
	var userPosts []models.Post
	for _, post := range data.Posts {
		if post.Passphrase == passphrase {
			userPosts = append(userPosts, post)
		}
	}

	// Retourner la longueur du slice de posts de l'utilisateur
	return c.JSON(http.StatusOK, models.Response{Message: strconv.Itoa(len(userPosts)), Success: true})
}


func GetLikesByPost(c echo.Context) error {
	// Récupérer le token du header
	authorization := c.Request().Header.Get("Authorization")

	if authorization == "" {
		return c.JSON(400, models.Response{Success: false, Message: "Vous devez renseigner un token."})
	}

	err := utils.IsTokenExists(c, authorization)
	if err != nil {
		return err
	}

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un id."})
	}

	viper.SetConfigName("posts")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	posts := viper.Get("posts").([]interface{})
	// Boucle for pour parcourir les posts
	for _, post := range posts {
		p := post.(map[string]interface{}) // On cast le post en map[string]interface{} afin de pouvoir accéder à ses valeurs
		if p["id"].(string) == id {        // Si l'id du post est égal à l'id donné en paramètre
			likedby := p["likedby"].([]interface{}) // On récupère les personnes qui ont liké le post
			return c.JSON(http.StatusOK, models.Response{Success: true, Message: strconv.Itoa(len(likedby))})
		}
	}

	return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Le post n'existe pas."})
}
