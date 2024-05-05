package routes

// Les librairies à importer
import (
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"twittueur_api/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Générer une passphrase qui permettra à l'utilisateur de se connecter
// Et lui éviter de faire un mot de passe.

func generatePassphrase() string {
	// Lire le fichier words.json
	file, err := os.Open("words.json")
	if err != nil {
		print(err.Error())
		return ""
	}
	defer file.Close()

	// Décodeur JSON pour lire le fichier
	decoder := json.NewDecoder(file)

	// Les rendre lisible pour Go
	var words []string
	err = decoder.Decode(&words)
	if err != nil {
		print(err.Error())
		return ""
	}

	var passphraseWords []string
	// Sélectionner 20 mots aléatoires
	for i := 0; i < 20; i++ {
		index := rand.Intn(len(words))
		passphraseWords = append(passphraseWords, words[index])
	}

	// Retirer le dernier espace de passphraseWords
	passphrase := strings.Join(passphraseWords, " ")

	// Renvoyer la passphrase
	return passphrase
}

/*
On créé une fonction Register qui sera appelée dans le server.go
Le paramètre c nous permet de gérer les erreurs, les informations
que l'on doit renvoyer à l'utilisateur.
Le type à renvoyer (que l'on spécifie après les paramètres) est error.
La requête doit avoir un username, un nom et un avatar.
A la fin la fonction renverra donc un message/nil s'il n'y a pas d'erreur, ou l'erreur
à l'utilisateur.
*/

func Register(c echo.Context) error {
	passphrase := generatePassphrase()

	username := c.FormValue("username")
	name := c.FormValue("name")

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username":   username,
		"passphrase": passphrase,
	})

	tokenString, _err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if _err != nil {
		return c.JSON(500, models.Response{Message: "Une erreur s'est passé de notre coté, réessayez plus tard.", Success: false})
	}

	// Vérifier si le nom d'utilisateur et le nom sont fournis
	if username == "" {
		return c.JSON(400, models.Response{Message: "Le nom d'utilisateur est requis", Success: false})
	}

	if name == "" {
		return c.JSON(400, models.Response{Message: "Le nom est requis", Success: false})
	}

	// Créer un nouvel utilisateur
	user := &models.User{
		Username:   username,
		Name:       name,
		Passphrase: passphrase,
	}

	// On vérifie si le fichier users.json existe, sinon on le créé
	if _, err := os.Stat("db/users.json"); os.IsNotExist(err) {
		file, err := os.Create("db/users.json")
		if err != nil {
			return err
		}
		defer file.Close()
		file.WriteString(`{"users": []}`)
	}

	// Configuration de Viper
	viper.SetConfigName("users")
	viper.SetConfigType("json")
	viper.AddConfigPath("db")

	// On lit le fichier json
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Convertir la donnée en une structur go
	var data models.Users
	if err := viper.Unmarshal(&data); err != nil {
		return err
	}

	user.Passphrase = passphrase

	// On vérifie si l'utilisateur existe déjà
	for _, u := range data.Users {
		if u.Username == user.Username {
			// Si oui on renvoie une erreur
			return echo.NewHTTPError(400, "Le nom d'utilisateur existe déjà")
		}
	}

	// Lire l'avatar
	avatar, _ := c.FormFile("avatar")

	if avatar == nil {
		// Ouvrir l'image par défaut
		src, err := os.Open("assets/empty.png")
		if err != nil {
			return err
		}
		defer src.Close()

		// Créer un nouveau fichier avec le nom de l'utilisateur
		dst, err := os.Create("db/avatars/" + username + ".png")
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copier l'image par défaut vers le nouveau fichier
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

		user.Avatar = "db/avatars/" + username + ".png"
	} else {

		// Enregistrer l'avatar dans la db
		// Lire le fichier
		src, err := avatar.Open()
		if err != nil {
			return err
		}

		defer src.Close()

		// Récupérer l'extension du fichier
		ext := filepath.Ext(avatar.Filename)
		// Si l'extension n'est pas png, convertir l'image en png
		if ext != ".png" {
			img, _, err := image.Decode(src)
			if err != nil {
				// Afficher une erreur compréhensible
				print(err.Error())
				return err
			}

			// Créer un nouveau fichier png
			dst, err := os.Create("db/avatars/" + username + ".png")
			if err != nil {
				return err
			}
			defer dst.Close()

			// Écrire l'image en png
			err = png.Encode(dst, img)
			if err != nil {
				print(err)
				return err
			}

			user.Avatar = "db/avatars/" + username + ".png"
		} else {
			// Créé un fichier dans le dossier db/avatars avec comme nom l'username
			dst, err := os.Create("db/avatars/" + username + ext)
			if err != nil {
				return err
			}
			defer dst.Close()

			// Copier l'image de la raquête vers celle du serveur.
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

			user.Avatar = "db/avatars/" + username + ext
		}

	}
	// Ajouter un nouvel utilisateur au JSON.
	data.Users = append(data.Users, *user)
	viper.Set("users", data.Users)

	if err := viper.WriteConfig(); err != nil {
		return err
	}
	// Si tout est bon on envoie un message à l'aide de la structure response.
	return c.JSON(http.StatusOK, models.Response{Message: tokenString, Success: true})
}
