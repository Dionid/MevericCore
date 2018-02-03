package tztusers

import (
	"net/http"
	"tztatom/tztcore"

	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserController struct {
	tztcore.ModelControllerBase
}

func (this *UserController) Create(c echo.Context) error {
	userData := UserModel{}
	if err := this.GetReqModelsData(&userData, &c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad parameters")
	}

	if userData.Email == "" || userData.Password == "" {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Email and password are required")
	}

	findUser := new(UserModel)

	var err error

	if err = UsersCollectionManager.FindByEmail(userData.Email, findUser); err == nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "This email is already in use")
	} else if err != UsersCollectionManager.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Try later")
	}

	user := &UserModel{
		Email: userData.Email,
		Password: userData.Password,
		IsAdmin: false,
	}

	if err := UsersCollectionManager.Save(user); err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error())
	}

	return c.NoContent(200)
}

func (this *UserController) Auth(c echo.Context) error {
	userData := UserModel{}
	if err := this.GetReqModelsData(&userData, &c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad parameters")
	}
	if userData.Email == "" || userData.Password == "" {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Email and password are required")
	}

	user := new(UserModel)

	if err := UsersCollectionManager.FindByEmail(userData.Email, user); err != nil {
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
