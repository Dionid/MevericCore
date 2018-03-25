package mcdashboard_old

import (
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"mevericcore/mccommon"
)

func GetUserMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt := c.Get("client").(*jwt.Token)
		claims := userJwt.Claims.(jwt.MapClaims)
		userId := claims["id"].(string)

		userM := new(mccommon.UserModel)
		if err := UsersCollectionManager.FindModelByStringId(userId, userM); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		c.Set("userM", userM)
		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}