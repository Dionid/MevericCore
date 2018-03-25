package mcdashboard_old

import "github.com/labstack/echo"

func initMeModule(e *echo.Echo) {
	meG := e.Group("/me")
	meG.Use(jwtMdlw)
	initMeRoutes(meG)
}
