package mcplantainer


import (
	"gopkg.in/mgo.v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"mevericcore/mcinnerrpc"
	"mevericcore/mccommunication"
	"fmt"
)

var (
	isDBDrop = false
	mainDBName = "tztatom"
)

func initMongoDbConnection() *mgo.Session {
	session, err := mgo.Dial("tzta:qweqweqwe@localhost")
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	if isDBDrop {
		err = session.DB("tztatom").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	return session
}

func initEcho() *echo.Echo {
	e := echo.New()

	// Debug
	e.Debug = true
	e.Logger.SetLevel(1)

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "")

	return e
}

var (
	innerRPCMan = mcinnerrpc.New()
)

func initInnerRPCMan() {
	innerRPCMan.Init()
	innerRPCMan.Service.Subscribe("Devices.Plantainer.RPC.Send", func(msg *mcinnerrpc.Msg) {
		return
	})
	innerRPCMan.Service.Subscribe("User.RPC.Devices.Plantainer.>", func(req *mcinnerrpc.Msg) {
		msg := &mccommunication.ClientToServerRPCReqSt{}
		if err := msg.UnmarshalJSON(req.Data); err != nil {
			fmt.Println("msg.UnmarshalJSON error: " + err.Error())
			return
		}

		respChan := make(mccommunication.ClientToServerHandleResultChannel)

		go func() {
			if err := UserRPCManager.Handle(respChan, msg); err != nil {
				data := &mccommunication.RPCMsg{
					Method: msg.RPCMsg.Method,
					Id: msg.RPCMsg.Id,
					Src: msg.RPCMsg.Dst,
					Dst: msg.RPCMsg.Src,
					Error: &map[string]interface{}{
						"message": err.Error(),
						"code": 500,
					},
				}
				if bData, err := data.MarshalJSON(); err != nil {
					//data := &mccommunication.RPCMsg{
					//	Method: msg.RPCMsg.Method,
					//	Id: msg.RPCMsg.Id,
					//	Src: msg.RPCMsg.Dst,
					//	Dst: msg.RPCMsg.Src,
					//	Error: &map[string]interface{}{
					//		"message": "Marshaling error problem",
					//		"code": 500,
					//	},
					//}
					//ebData, _ := data.MarshalJSON()
					//userWS.SendMsg(ebData)
					print(bData)
					return
				} else {
					return
					//userWS.SendMsg(bData)
				}
			}
		}()

		for resultSt := range respChan {
			if resultSt.Error != nil {
				if bData, err := resultSt.Error.MarshalJSON(); err != nil {
					//data := &mccommunication.RPCMsg{
					//	Method: msg.Method,
					//	Id: msg.Id,
					//	Src: msg.Dst,
					//	Dst: msg.Src,
					//	Error: &map[string]interface{}{
					//		"message": "Marshaling error problem",
					//		"code": 500,
					//	},
					//}
					//ebData, _ := data.MarshalJSON()
					//userWS.SendMsg(ebData)
				} else {
					innerRPCMan.Service.Publish("User.RPC.Send", bData)
				}
			}
			if resultSt.Res != nil {
				if bData, err := resultSt.Res.MarshalJSON(); err != nil {
					//data := &mccommunication.RPCMsg{
					//	Method: msg.Method,
					//	Id: msg.Id,
					//	Src: msg.Dst,
					//	Dst: msg.Src,
					//	Error: &map[string]interface{}{
					//		"message": "Marshaling error problem",
					//		"code": 500,
					//	},
					//}
					//ebData, _ := data.MarshalJSON()
					//userWS.SendMsg(ebData)
				} else {
					innerRPCMan.Service.Publish("User.RPC.Send", bData)
				}
			}
		}
	})
}

var (
	plantainerDataCollectionManager = NewPlantainerDataCollectionManager()
	plantainerCollectionManager = NewPlantainerCollectionManager(plantainerDataCollectionManager)
)

func initCollections(session *mgo.Session) {
	plantainerDataCollectionManager.Init(session, mainDBName)
	plantainerCollectionManager.Init(session, mainDBName)
}

func initRoutes(e *echo.Echo) {

}

func Init() {
	// 1. Init MongoDB session
	session := initMongoDbConnection()
	defer session.Close()

	// 2. Init Echo server for Devices and Users
	e := initEcho()

	// 3. Init Collections
	initCollections(session)

	// 4. Init Routes
	initRoutes(e)

	initInnerRPCMan()

	InitUserRPCManager()

	e.Logger.Fatal(e.Start("localhost:3001"))
}
