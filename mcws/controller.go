package mcws

//import (
//	"github.com/labstack/echo"
//	"net/http"
//	"github.com/gorilla/websocket"
//)
//
//var (
//	wsUpgrader = websocket.Upgrader{
//		CheckOrigin: func(r *http.Request) bool {
//			return true
//		},
//	}
//)
//
//func WSHandler(c echo.Context) error {
//	clientId := c.QueryParam("clientId")
//
//	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
//	if _, ok := err.(websocket.HandshakeError); ok {
//		return echo.NewHTTPError(http.StatusBadRequest, "WS problem")
//	} else if err != nil {
//		return echo.NewHTTPError(http.StatusBadRequest, "WS problem")
//	}
//
//	defer ws.Close()
//
//
//
//	return c.NoContent(200)
//}