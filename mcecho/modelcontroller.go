package mcecho

import (
	"net/http"

	"github.com/labstack/echo"
)

type ModelControllerBase struct{}

type ModelControllerBaseInterface interface {
	GetListMiddlewares() []echo.MiddlewareFunc   // return middlewares for List Routes
	GetDetailMiddlewares() []echo.MiddlewareFunc // return middlewares for Detail Routes
	List(c echo.Context) error                   // method = GET / processing
	Create(c echo.Context) error                 // method = POST / processing
	Retrieve(c echo.Context) error               // method = GET /:id handling
	Update(c echo.Context) error                 // method = PUT /:id processing
	Destroy(c echo.Context) error                // method = DELETE /:id processing
	GetReqModelsData(model interface {
		UnmarshalJSON([]byte) error
	}, c *echo.Context) error
}

func (this *ModelControllerBase) GetListMiddlewares() []echo.MiddlewareFunc {
	return nil
}

func (this *ModelControllerBase) GetDetailMiddlewares() []echo.MiddlewareFunc {
	return nil
}

func (this *ModelControllerBase) List(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotAcceptable, "List Method Not Allowed")
}

func (this *ModelControllerBase) Create(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotAcceptable, "Create Method Not Allowed")
}

func (this *ModelControllerBase) Retrieve(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotAcceptable, "Retrieve Method Not Allowed")
}

func (this *ModelControllerBase) Update(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotAcceptable, "Update Method Not Allowed")
}

func (this *ModelControllerBase) Destroy(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotAcceptable, "Destroy Method Not Allowed")
}

func (this *ModelControllerBase) GetReqModelsData(model interface {
	UnmarshalJSON([]byte) error
}, c *echo.Context) error {
	return UnmarshalRequestData(model, *c)
}