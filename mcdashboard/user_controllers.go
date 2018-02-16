package mcdashboard

import (
	"mevericcore/mcecho"
	"github.com/labstack/echo"
	"net/http"
	"mevericcore/mccommon"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserController struct {
	mcecho.ModelControllerBase
}

func (this *UserController) Auth(c echo.Context) error {
	userData := mccommon.UserModel{}
	if err := this.GetReqModelsData(&userData, &c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad parameters")
	}
	if userData.Login == "" || userData.Password == "" {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Email and password are required")
	}

	user := new(mccommon.UserModel)

	if err := UsersCollectionManager.FindModelByLogin(userData.Login, user); err != nil {
		if err == UsersCollectionManager.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotAcceptable, "Invalid email or password")
		} else {
			return echo.NewHTTPError(http.StatusNotAcceptable, "Try again")
		}
	}

	if !user.CheckPasswordHash(userData.Password) {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Invalid email or password")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))

	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Try again")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}