package mcdashboard_old

import (
	"mevericcore/mcecho"
	"github.com/labstack/echo"
	"net/http"
	"mevericcore/mccommon"
	"fmt"
)

type DevicesBaseControllerSt struct {
	mcecho.ModelControllerBase
}

func (this *DevicesBaseControllerSt) getDevice(c echo.Context) (*mccommon.DeviceWithCustomDataBaseModel, error) {
	userM := c.Get("userM").(*mccommon.UserModel)
	fmt.Println(userM)

	device := &mccommon.DeviceWithCustomDataBaseModel{}
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

	devices := new(mccommon.DevicesWithCustomDataListBaseModel)

	if err := DevicesCollectionManager.FindByOwnerId(userId, devices); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return mcecho.SendJSON(devices, &c)
}

func (this *DevicesBaseControllerSt) Retrieve(c echo.Context) error {
	device, err := this.getDevice(c)

	if err != nil {
		return err
	}

	return mcecho.SendJSON(device, &c)
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

