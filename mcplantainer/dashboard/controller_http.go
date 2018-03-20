package dashboard

import (
	"mevericcore/mcecho"
	"github.com/labstack/echo"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"mevericcore/mccommon"
	"tztatom/tztcore"
	"mevericcore/mcplantainer/common"
)

type UserPlantainerControllerSt struct {
	mcecho.ModelControllerBase
}

func (this *UserPlantainerControllerSt) Create(c echo.Context) error {
	userId, err := mcecho.GetContextClientId(&c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Haven't got user id")
	}

	// TODO: ONLY ADMIN CHECK

	deviceShadowId := mccommon.RandString(13)

	device := &common.PlantainerModelSt{
		DeviceBaseModel: mccommon.DeviceBaseModel{
			Shadow: mccommon.ShadowModelSt{
				Id: deviceShadowId,
			},
			OwnersIds: []bson.ObjectId{bson.ObjectIdHex(userId)},
		},
	}

	if err := common.PlantainerCollectionManager.SaveModel(device); err != nil {
		// TODO: Must check if problem with shadowId
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := common.PlantainerCollectionManager.FindModelById(device.ID, device); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return tztcore.SendJSON(device, &c)
}
