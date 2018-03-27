package mcdashboard_old

import (
	"mevericcore/mclibs/mcws"
	"mevericcore/mclibs/mccommon"
	"github.com/labstack/echo"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func InitWsManager() {
	authG := WSocketsResources.Group("auth")
	authG.AddHandler("login", func(req *mcws.ReqSt) error {
		tokenMsg := &WsTokenActionReqSt{}
		if err := tokenMsg.UnmarshalJSON(req.Msg); err != nil {
			return err
		}
		if tokenMsg.Login == "" || tokenMsg.Password == "" {
			req.Ws.SendErrorMsg("Login and password are required", tokenMsg.Method, 503, tokenMsg.Id)
			return nil
		}
		// Create token
		user := new(mccommon.UserModel)

		if err := UsersCollectionManager.FindModelByLogin(tokenMsg.Login, user); err != nil {
			if err == UsersCollectionManager.ErrNotFound {
				return echo.NewHTTPError(http.StatusNotAcceptable, "Invalid email or password")
			} else {
				return echo.NewHTTPError(http.StatusNotAcceptable, "Try again")
			}
		}

		if !user.CheckPasswordHash(tokenMsg.Password) {
			return echo.NewHTTPError(http.StatusNotAcceptable, "Invalid email or password")
		}

		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))

		if err != nil {
			req.Ws.SendErrorMsg("Token creation problem", tokenMsg.Method, 503, tokenMsg.Id)
			return nil
		}

		// Send success
		CreateAndSendWsTokenActionRes(req.Ws, t, tokenMsg.Method, tokenMsg.Id)
		return nil
	})
	authG.AddHandler("auth", func(req *mcws.ReqSt) error {
		authMsg := &WsAuthenticateActionReqSt{}
		if err := authMsg.UnmarshalJSON(req.Msg); err != nil {
			return err
		}
		if authMsg.Token == "" {
			req.Ws.SendErrorMsg("Token is required", authMsg.Method, 503, authMsg.Id)
			return nil
		}
		// Auth user
		t, err := jwt.Parse(authMsg.Token, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			req.Ws.SendErrorMsg("Problem with token", authMsg.Method, 503, authMsg.Id)
		}

		userTokenId := t.Claims.(jwt.MapClaims)["id"].(string)
		user := new(mccommon.UserModel)

		if err := UsersCollectionManager.FindModelByStringId(userTokenId, user); err != nil {
			req.Ws.SendErrorMsg("Incorrect token", authMsg.Method, 503, authMsg.Id)
			return nil
		}

		req.Ws.Authorized = true

		// Send success
		CreateAndSendWsAuthenticateActionRes(req.Ws, authMsg.Method, authMsg.Id)
		return nil
	})
	InitDeviceResource()
}

func userWsMiddleware(next mcws.HandlerFunc) mcws.HandlerFunc {
	return func(req *mcws.ReqSt) error {
		userId := req.Ws.Id

		userM := new(mccommon.UserModel)
		if err := UsersCollectionManager.FindModelByStringId(userId, userM); err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		req.Set("userM", userM)
		if err := next(req); err != nil {
			return err
		}

		return nil
	}
}

//func deviceWsMiddleware(next mcws.HandlerFunc) mcws.HandlerFunc {
//	return func(req *mcws.ReqSt) error {
//		userId := req.Ws.Id
//
//		userM := new(mccommon.UserModel)
//		if err := UsersCollectionManager.FindModelByStringId(userId, userM); err != nil {
//			return echo.NewHTTPError(http.StatusNotFound, "User not found")
//		}
//
//		req.Set("userM", userM)
//		if err := next(req); err != nil {
//			return err
//		}
//
//		return nil
//	}
//}

func InitDeviceResource() {
	deviceG := WSocketsResources.Group("device")
	deviceG.Use(userWsMiddleware)
	deviceG.AddHandler("list", func(req *mcws.ReqSt) error {
		return nil
	})
	deviceG.AddHandler("get", func(req *mcws.ReqSt) error {
		return nil
	})
	//deviceG.AddHandler("get", func(req *mcws.ReqSt) error {
	//	return nil
	//})
	//deviceG.AddHandler("get", func(req *mcws.ReqSt) error {
	//	return nil
	//})
}