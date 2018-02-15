package tztusers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"time"
	"mevericcore/mcecho"
	"mevericcore/mccommon"
)

type UserController struct {
	mcecho.ModelControllerBase
}

func (this *UserController) Create(c echo.Context) error {
	userData := mccommon.UserModel{}
	if err := this.GetReqModelsData(&userData, &c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad parameters")
	}

	if userData.Login == "" || userData.Password == "" {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Email and password are required")
	}

	findUser := new(mccommon.UserModel)

	var err error

	if err = UsersCollectionManager.FindModelByLogin(userData.Login, findUser); err == nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "This email is already in use")
	} else if err != UsersCollectionManager.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Try later")
	}

	user := &mccommon.UserModel{
		Email: userData.Email,
		Password: userData.Password,
		IsAdmin: false,
	}

	if err := UsersCollectionManager.SaveModel(user); err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, err.Error())
	}

	return c.NoContent(200)
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

type CompanyControllerSt struct {
	mcecho.ModelControllerBase
}

//func (this *CompanyControllerSt) List(c echo.Context) error {
//	userId, err := mcecho.GetContextClientId(&c)
//	if err != nil {
//		return echo.NewHTTPError(http.StatusUnauthorized, "Haven't got user id")
//	}
//
//	companies := new(CompanyListModel)
//
//	if err := CompaniesCollectionManager.FindModelByStringId(userId, companies); err != nil {
//		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
//	}
//
//	return tztcore.SendJSON(companies, &c)
//}
