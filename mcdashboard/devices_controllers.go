package mcdashboard

import (
	"mevericcore/mcecho"
	"github.com/labstack/echo"
	"net/http"
)

type UserDevicesControllerSt struct {
	mcecho.ModelControllerBase
}

func (this *UserDevicesControllerSt) List(c echo.Context) error {
	userId, err := mcecho.GetContextClientId(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Haven't got user id")
	}

	devices := &DevicesListModelSt{}

	if err := devicesCollectionManager.FindByOwnerId(userId, devices); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return mcecho.SendJSON(devices, &c)
}

func (this *UserDevicesControllerSt) getDevice(c echo.Context) (*DeviceModelSt, error) {
	userId, err := mcecho.GetContextClientId(&c)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Haven't got user id")
	}

	device := &DeviceModelSt{}
	deviceShadowId := c.Param("id")

	if err := devicesCollectionManager.FindByShadowId(deviceShadowId, device); err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Device not found")
	}

	if isOwner, err := device.IsOwnerStringId(userId); err != nil {
		return nil, echo.NewHTTPError(http.StatusServiceUnavailable, "Try again")
	} else if !isOwner {
		return nil, echo.NewHTTPError(http.StatusForbidden, "You can use only your own devices")
	}

	return device, nil
}

func (this *UserDevicesControllerSt) Retrieve(c echo.Context) error {
	device, err := this.getDevice(c)

	if err != nil {
		return err
	}

	return mcecho.SendJSON(device, &c)
}

func (this *UserDevicesControllerSt) Destroy(c echo.Context) error {
	device, err := this.getDevice(c)

	if err != nil {
		return err
	}

	if err := devicesCollectionManager.DeleteModel(device); err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Can't delete. Try again")
	}

	return c.NoContent(http.StatusOK)
}