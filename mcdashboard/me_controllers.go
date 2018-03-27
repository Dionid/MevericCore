package mcdashboard

import (
	"mevericcore/mclibs/mcecho"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type MeController struct {
	mcecho.ModelControllerBase
}

func (this *MeController) Me(c echo.Context) error {
	userJwt := c.Get("client").(*jwt.Token)
	claims := userJwt.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	userM := new(MeModel)
	if err := usersCollectionManager.FindModelByStringId(userId, userM); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(200, map[string]interface{}{
		"login": userM.Login,
		"email": userM.Email,
		"id": userM.ID,
	})
}
