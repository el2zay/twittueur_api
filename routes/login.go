package routes

import (
	"net/http"
	"os"
	"twittueur_api/models"
	"twittueur_api/src/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func Login(c echo.Context) error {
	username := c.FormValue("username")
	passphrase := c.FormValue("passphrase")

	if passphrase == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Vous devez renseigner la passphrase", Success: false})
	}

	if username == "" {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Vous devez renseigner le nom d'utilisateur", Success: false})
	}

	// Vérifier si l'utilisateur existe
	err := utils.IsPassphraseExists(c, passphrase)
	if err != nil {
		return err
	}

	// Vérifier que l'username est associé à la passphrase dans la base de données
	viper.SetConfigName("users")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return c.JSON(500, models.Response{Message: "Une erreur s'est produite.", Success: false})
	}

	var users []models.User

	if err := viper.UnmarshalKey("users", &users); err != nil {
		return c.JSON(500, models.Response{Message: "Une erreur s'est produite.", Success: false})
	}

	userExists := false
	for _, user := range users { // Parcourir les utilisateurs
		if user.Username == username && user.Passphrase == passphrase {
			userExists = true // Si l'utilisateur existe, on sort de la boucle
			break
		}
	}

	if !userExists {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "L'utilisateur n'existe pas. Vérifiez l'username.", Success: false})
	}

	// Regénérer le token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username":   username,
		"passphrase": passphrase,
	})

	tokenString, _err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // "Signer" le token
	if _err != nil {
		return c.JSON(500, models.Response{Message: "Une erreur s'est passé de notre coté, réessayez plus tard.", Success: false})
	}

	return c.JSON(http.StatusAccepted, models.Response{Success: true, Message: tokenString})
}
