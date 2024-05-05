package routes

import (
	"net/http"
	"strconv"
	"twittueur_api/models"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

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
