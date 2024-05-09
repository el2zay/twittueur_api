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

/*
On va copier le code du fichier likes.go et le modifier pour qu'il corresponde à la gestion des bookmarks.

Le code est donc très similaire à celui de likes.go, mais on remplace les mots "like" par "bookmark" et on change les noms des fonctions et des variables.
*/

func PostBookmarks(c echo.Context) error {
	// Récupérer le token du header
	authorization := c.Request().Header.Get("Authorization")

	// S'il est vide on retourne une erreur
	if authorization == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un token."})
	}

	err := utils.IsTokenExists(c, authorization)
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
		username = claims["username"].(string)
	}

	id := c.FormValue("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un id."})
	}

	err = utils.IsPostExists(c, id)
	if err != nil {
		return err
	}

	// Variable pour dire si on a retiré ou ajouté le bookmark
	var wordResponse string

	var data models.PostRequest
	if err := viper.Unmarshal(&data); err != nil {
		return err
	}

	// Parcourir les posts pour trouver celui avec l'id spécifié
	for i, post := range data.Posts {
		if post.ID == id {
			// Trouvé le post, récupérer la liste des bookmarks directement
			bookmarkedby := post.Bookmarkedby
			// Vérifier si l'utilisateur a déjà bookmark le post
			for j, user := range bookmarkedby {
				if user == username {
					// Si oui, retirer le bookmark
					bookmarkedby = append(bookmarkedby[:j], bookmarkedby[j+1:]...)
					wordResponse = "retiré"
					break
				}
			}
			// Si l'utilisateur n'a pas bookmark le post, ajouter le bookmark
			if wordResponse == "" {
				bookmarkedby = append(bookmarkedby, username)
				wordResponse = "ajouté"
			}
			data.Posts[i].Bookmarkedby = bookmarkedby
			break
		}
	}

	// Réécrire le fichier avec le nouveau contenu
	viper.Set("posts", data.Posts)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Message: fmt.Sprintf("Bookmark %s avec succès.", wordResponse)})

}

func GetBookmarks(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")

	if authorization == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un token."})
	}
	err := utils.IsTokenExists(c, authorization)
	if err != nil {
		return err
	}

	id := c.FormValue("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Success: false, Message: "Vous devez renseigner un id."})
	}

	err = utils.IsPostExists(c, id)
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
			// Trouvé le post, retourner la liste des bookmarks
			return c.JSON(http.StatusOK, models.Response{Success: true, Data: post.Bookmarkedby})
		}
	}

	// Si le post n'est pas trouvé
	return c.JSON(http.StatusNotFound, models.Response{Success: false, Message: "Post non trouvé."})
}
