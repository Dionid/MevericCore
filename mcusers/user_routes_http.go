package mcusers

import "github.com/labstack/echo"

var (
	UserCtrl = &UserController{}
)

func initUsersRoutes(group *echo.Group) {
	group.POST("", UserCtrl.Create)
}

func initAuthRoutes(group *echo.Group) {
	group.POST("/login", UserCtrl.Auth)
}

//func initWsRoute(group *echo.Group) {
//	group.GET("/ws", UserCtrl.WSHandler)
//}
