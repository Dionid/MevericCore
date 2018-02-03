package tztusers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func GetUserMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt := c.Get("client").(*jwt.Token)
		claims := userJwt.Claims.(jwt.MapClaims)
		userId := claims["id"].(string)

		userM := new(UserModel)
		if err := UsersCollectionManager.FindByStringId(userId, userM); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		c.Set("userM", userM)
		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}
