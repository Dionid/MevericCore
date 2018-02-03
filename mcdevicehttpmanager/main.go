package mcdevicehttpmanager

import (
	"github.com/labstack/echo"
	"mevericcore/mccommon"
	"net/http"
	"tztatom/tztcore"
	"io/ioutil"
	"mevericcore/mcecho"
)

type DeviceHTTPManagerSt struct {
	reqHandler mccommon.DeviceToServerReqHandler
}

func (this *DeviceHTTPManagerSt) SetReqHandler(handler mccommon.DeviceToServerReqHandler) {
	this.reqHandler = handler
}

func (this *DeviceHTTPManagerSt) HandleReq(msg *mccommon.DeviceToServerReqSt) (resMsg mccommon.JSONData, sendBack bool, errMsg mccommon.JSONData) {
	if this.reqHandler != nil {
		return this.reqHandler(msg)
	}

	// ToDo: Send req to QueueManager and return

	return nil, false, nil
}



func (this *DeviceHTTPManagerSt) ReqPostHandler(c echo.Context) error {
	publisherDeviceId, cErr := mcecho.GetContextClientId(&c)

	if cErr != nil {
		// ToDo: LOG
		return echo.NewHTTPError(http.StatusUnauthorized, "Haven't got user id")
	}

	channelId := c.Param("channelId")
	deviceId := c.Param("deviceId")

	if publisherDeviceId == deviceId {
		// ToDo: LOG
		return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to change other device data")
	}

	raw, _ := ioutil.ReadAll(c.Request().Body)

	res, sendBack, err := this.HandleReq(&mccommon.DeviceToServerReqSt{
		DeviceId:  deviceId,
		ChannelId: channelId,
		Protocol:  "HTTP",
		Msg:       &raw,
	})

	if err != nil {
		// ToDo: LOG
		c.Logger().Error(err)
		b, jerr := res.MarshalJSON()
		if jerr != nil {
			// ToDo: LOG
			c.Logger().Error(jerr)
			return echo.NewHTTPError(http.StatusInternalServerError, jerr.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, b)
	}

	if sendBack && res != nil {
		return tztcore.SendJSON(res, &c)
	}

	// ToDo: LOG
	return echo.NewHTTPError(http.StatusInternalServerError, "")
}


//func main() {
//	// Init echo server
//	// Make .post
//}
