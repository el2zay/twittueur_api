package routes

import (
	"twittueur_api/src/utils"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	passphrase := c.FormValue("passphrase")
	userData := utils.IsPassphraseExists(c, passphrase, true)
	if userData != nil {
		return userData
	}

	return userData
}
