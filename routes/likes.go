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
		return c.JSON(400, models.Response{Success: false, Message: "Vous devez renseigner un token."})
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
		c.JSON(400, models.Response{Success: false, Message: "Token invalide"})
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

	viper.SetConfigName("posts")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Variable pour dire si on a retiré ou ajouté le like
	var wordResponse string

	// Rajouter le like dans le fichier posts.json

	// Rajouter le like dans le fichier posts.json
	posts := viper.Get("posts").([]interface{})
	// Boucle for pour parcourir les posts
	for _, post := range posts {
		p := post.(map[string]interface{}) // On cast le post en map[string]interface{} afin de pouvoir accéder à ses valeurs
		if p["id"].(string) == id {        // Si l'id du post est égal à l'id donné en paramètre
			likedby, ok := p["likedby"].([]interface{}) // On récupère les personnes qui ont liké le post
			if !ok {
				likedby = []interface{}{} // Si likedby est vide, on met une nouvelle liste vide
			}
			index := -1
			for i, user := range likedby { // On parcourt les personnes qui ont liké le post
				if user.(string) == username { // Si le username est déjà dans likedby
					index = i // On récupère l'index de ce username
					break
				}
			}
			if index != -1 {
				// Si le username est déjà dans likedby, on le retire
				likedby = append(likedby[:index], likedby[index+1:]...)
				wordResponse = "retiré"
			} else {
				// Sinon, on ajoute le username à likedby
				likedby = append(likedby, username)
				wordResponse = "ajouté"
			}
			p["likedby"] = likedby
			break
		}
	}

	viper.Set("posts", posts)

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{Success: true, Message: fmt.Sprintf("Like %s avec succès.", wordResponse)})

}

func GetLikes(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")

	if authorization == "" {
		return c.JSON(400, models.Response{Success: false, Message: "Vous devez renseigner un token."})
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

	// Configuration de Viper pour lire le fichier posts.json
	viper.SetConfigName("posts")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Trouver le post avec l'id spécifié
	posts := viper.Get("posts").([]interface{})

	for _, post := range posts {
		p := post.(map[string]interface{}) // On cast le post en map[string]interface{} afin de pouvoir accéder à ses valeurs

		if p["id"].(string) == id {
			likedby, ok := p["likedby"].([]interface{}) // On récupère les personnes qui ont liké le post
			if !ok {
				likedby = []interface{}{} // Si likedby est vide, on met une nouvelle liste vide
			}
			return c.JSON(http.StatusOK, models.Response{Success: true, Message: "Liste des personnes qui ont liké le post.", Data: likedby})
		}
	}

	return nil
}
