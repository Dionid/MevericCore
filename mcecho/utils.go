package mcecho

import (
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"github.com/dgrijalva/jwt-go"
)

func UnmarshalRequestData(model interface {
	UnmarshalJSON([]byte) error
}, c echo.Context) error {
	raw, _ := ioutil.ReadAll(c.Request().Body)
	if err := model.UnmarshalJSON(raw); err != nil {
		return err
	}
	return nil
}

func SendJSON(data interface {
	MarshalJSON() ([]byte, error)
}, c *echo.Context) error {
	res, err := data.MarshalJSON()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return (*c).JSONBlob(http.StatusOK, res)
}

func GetContextClientId(c *echo.Context) (string, error) {
	userJwt := (*c).Get("client").(*jwt.Token)
	claims := userJwt.Claims.(jwt.MapClaims)
	userGuid := claims["id"].(string)
	return userGuid, nil
}