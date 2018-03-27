package mcecho

import "github.com/labstack/echo"

func CreateModelControllerRoutes(
	e *echo.Group,
	groupName string,
	controller ModelControllerBaseInterface,
) (*echo.Group, *echo.Group) {

	controllerListG := e.Group(groupName)

	if controller.GetListMiddlewares() != nil {
		controllerListG.GET("", controller.List, controller.GetListMiddlewares()...)
		controllerListG.POST("", controller.Create, controller.GetListMiddlewares()...)
	} else {
		controllerListG.GET("", controller.List)
		controllerListG.POST("", controller.Create)
	}

	controllerDetailG := controllerListG.Group("/:id")

	if controller.GetDetailMiddlewares() != nil {
		controllerDetailG.GET("", controller.Retrieve, controller.GetDetailMiddlewares()...)
		controllerDetailG.DELETE("", controller.Destroy, controller.GetDetailMiddlewares()...)
		controllerDetailG.PUT("", controller.Update, controller.GetDetailMiddlewares()...)
	} else {
		controllerDetailG.GET("", controller.Retrieve)
		controllerDetailG.DELETE("", controller.Destroy)
		controllerDetailG.PUT("", controller.Update)
	}

	return controllerListG, controllerDetailG
}
