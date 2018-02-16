package mcdashboard

import (
	"mevericcore/mcecho"
	"github.com/labstack/echo"
	"tztatom/tztcore"
	"net/http"
	"mevericcore/mccommon"
	"fmt"
)

type DevicesBaseControllerSt struct {
	mcecho.ModelControllerBase
}

func (this *DevicesBaseControllerSt) getDevice(c echo.Context) (*mccommon.DeviceBaseModel, error) {
	userM := c.Get("userM").(*mccommon.UserModel)
	fmt.Println(userM)

	device := &mccommon.DeviceBaseModel{}
	deviceShadowId := c.Param("id")
	if err := DevicesCollectionManager.FindByShadowId(deviceShadowId, device); err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Device not found")
	}

	if isOwner, err := device.IsOwner(userM.ID); err != nil {
		return nil, echo.NewHTTPError(http.StatusServiceUnavailable, "Try again")
	} else if !isOwner {
		return nil, echo.NewHTTPError(http.StatusForbidden, "You can use only your own devices")
	}

	return device, nil
}

func (this *DevicesBaseControllerSt) List(c echo.Context) error {
	userId, err := mcecho.GetContextClientId(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Haven't got user id")
	}

	devices := new(mccommon.DevicesListBaseModel)

	if err := DevicesCollectionManager.FindByOwnerId(userId, devices); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return tztcore.SendJSON(devices, &c)
}

func (this *DevicesBaseControllerSt) Retrieve(c echo.Context) error {
	device, err := this.getDevice(c)

	if err != nil {
		return err
	}

	return tztcore.SendJSON(device, &c)
}

func (this *DevicesBaseControllerSt) Destroy(c echo.Context) error {
	device, err := this.getDevice(c)

	if err != nil {
		return err
	}

	if err := DevicesCollectionManager.DeleteModel(device); err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Can't delete. Try again")
	}

	return c.NoContent(http.StatusOK)
}

