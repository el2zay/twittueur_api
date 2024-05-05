package utils

import (
	"errors"
	"os"
	"twittueur_api/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func IsTokenExists(c echo.Context, authorization string) error {
	tokenString := authorization[7:] // On ignore les 7 premières lettres (Bearer )
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(400, models.Response{Success: false, Message: "Token invalide"})
		return errors.New("token invalide")
	}

	var passphrase string

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		passphrase = claims["passphrase"].(string)
	}

	viper.SetConfigName("users")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// On vérifie si l'utilisateur existe
	var users []models.User

	if err := viper.UnmarshalKey("users", &users); err != nil {
		c.JSON(500, models.Response{Message: "Une erreur s'est produite.", Success: false})
		return errors.New("une erreur s'est produite")
	}

	userExists := false
	for _, user := range users {
		if user.Passphrase == passphrase {
			userExists = true
			break
		}
	}

	if !userExists {
		c.JSON(400, models.Response{Success: false, Message: "L'utilisateur n'existe pas."})
		return errors.New("l'utilisateur n'existe pas")
	}

	return nil
}

func IsPassphraseExists(c echo.Context, passphrase string, info ...bool) error {
	// Configuration de Viper pour lire le fichier users.json
	viper.SetConfigName("users")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		c.JSON(400, models.Response{Success: false, Message: "Erreur lors de la lecture du fichier users.json"})
		return errors.New("erreur lors de la lecture du fichier users.json")
	}

	// Récupérer les données du fichier users.json
	var users models.Users
	if err := viper.Unmarshal(&users); err != nil {
		c.JSON(400, models.Response{Success: false, Message: "Erreur lors de la lecture du fichier users.json"})
		return errors.New("erreur lors de la lecture du fichier users.json")
	}

	// Vérifier si le passphrase n'existe pas
	userExists := false
	for _, user := range users.Users {
		if user.Passphrase == passphrase {
			userExists = true
			// Si info n'est pas null, retourner les données de l'utilisateur
			if len(info) > 0 && info[0] {
				return c.JSON(200, models.Response{Success: true, Data: user})
			}
		}
	}

	if !userExists {
		c.JSON(404, models.Response{Success: false, Message: "L'utilisateur n'existe pas."})
		return errors.New("l'utilisateur n'existe pas")
	}

	return nil
}

func IsPostExists(c echo.Context, id string) error {
	viper.SetConfigName("posts")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// On vérifie si le post existe
	var posts []models.Post

	if err := viper.UnmarshalKey("posts", &posts); err != nil {
		c.JSON(500, models.Response{Message: "Une erreur s'est produite.", Success: false})
		return errors.New("une erreur s'est produite")
	}

	postExists := false
	for _, post := range posts {
		if post.ID == id {
			postExists = true
			break
		}
	}

	if !postExists {
		c.JSON(400, models.Response{Success: false, Message: "Le post n'existe pas."})
		return errors.New("le post n'existe pas")
	}

	return nil
}
