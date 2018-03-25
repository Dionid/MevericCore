package mcplantainer

import (
	"github.com/labstack/echo"
	"mevericcore/mcecho"
	"net/http"
	"mevericcore/mccommon"
	"gopkg.in/mgo.v2/bson"
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

	//deviceShadowId := mccommon.RandString(13)

	device := &PlantainerModelSt{
		DeviceBaseModel: mccommon.DeviceBaseModel{
			OwnersIds: []bson.ObjectId{bson.ObjectIdHex(userId)},
		},
		Shadow: PlantainerShadowSt{},
	}

	if err := plantainerCollectionManager.SaveModel(device); err != nil {
		// TODO: Must check if problem with shadowId
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := plantainerCollectionManager.FindModelById(device.ID, device); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return mcecho.SendJSON(device, &c)
}
