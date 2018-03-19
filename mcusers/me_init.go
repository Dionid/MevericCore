package mcusers

import "github.com/labstack/echo"

func initMeModule(e *echo.Group) {
	meG := e.Group("/me")
	meG.Use(jwtMdlw)
	initMeRoutes(meG)
}
