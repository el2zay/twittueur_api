package routes

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"twittueur_api/models"
	"twittueur_api/src/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func PostLikes(c echo.Context) error {
	// Récupérer le token du header
	authorization := c.Request().Header.Get("Authorization")

	// S'il est vide on retourne une erreur
	if authorization == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un token."})
	}

	err := utils.IsTokenExists(c, authorization) // Vérifier si le token existe
	if err != nil {
		return err
	}

	tokenString := authorization[7:] // On ignore les 7 premières lettres (Bearer )
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Token invalide"})
		return errors.New("token invalide")
	}

	var username string

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username = claims["username"].(string) // On récupère le username depuis le token
	}

	id := c.FormValue("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un id."})
	}

	err = utils.IsPostExists(c, id)
	if err != nil {
		return err
	}

	// Variable pour dire si on a retiré ou ajouté le like
	var wordResponse string

	var data models.PostRequest
	if err := viper.Unmarshal(&data); err != nil {
		return err
	}

	// Parcourir les posts pour trouver celui avec l'id spécifié
	for i, post := range data.Posts {
		if post.ID == id {
			// Trouvé le post, récupérer la liste des likes directement
			likedby := post.Likedby
			// Vérifier si l'utilisateur a déjà liké le post
			for j, user := range likedby {
				if user == username {
					// Si oui, retirer le like
					likedby = append(likedby[:j], likedby[j+1:]...)
					wordResponse = "retiré"
					break
				}
			}
			// Si l'utilisateur n'a pas liké le post, ajouter le like
			if wordResponse == "" {
				likedby = append(likedby, username)
				wordResponse = "ajouté"
			}
			data.Posts[i].Likedby = likedby
			break
		}
	}

	// Réécrire le fichier avec le nouveau contenu
	viper.Set("posts", data.Posts)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Message: fmt.Sprintf("Like %s avec succès.", wordResponse)})

}

func GetLikes(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization") // Récupérer le token du header

	if authorization == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un token."})
	}
	err := utils.IsTokenExists(c, authorization) // Vérifier si le token existe
	if err != nil {
		return err
	}

	id := c.FormValue("id") // Récupérer l'id du post
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un id."})
	}

	err = utils.IsPostExists(c, id) // Vérifier si le post existe
	if err != nil {
		return err
	}

	var data models.PostRequest

	if err := viper.Unmarshal(&data); err != nil {
		return err
	}

	// Parcourir les posts pour trouver celui avec l'id spécifié
	for _, post := range data.Posts {
		if post.ID == id {
			// Trouvé le post, retourner la liste des likes
			return c.JSON(http.StatusOK, models.Response{Success: true, Data: post.Likedby})
		}
	}

	// Si le post n'est pas trouvé
	return c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "Post non trouvé."})
}
