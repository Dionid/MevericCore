package mcplantainer

import (
	"time"
	"mevericcore/mcmongo"
	"gopkg.in/mgo.v2"
	"mevericcore/mccommunication"
	"gopkg.in/mgo.v2/bson"
)

type RPCTaskSt struct {
	mcmongo.ModelBase `bson:",inline"`
	RPCMsg *mccommunication.RPCMsg

	InExec bool
	Done bool
	DoneError *map[string]interface{}

	ExecTime time.Time
	SecondsCheck bool
	HoursCheck bool
	DaysCheck bool
}

func NewRPCTask() *RPCTaskSt {
	return &RPCTaskSt{}
}

//easyjson:json
type RPCTasksListSt []RPCTaskSt

func NewRPCTasksList() *RPCTasksListSt {
	return &RPCTasksListSt{}
}

func (this *RPCTasksListSt) GetBaseQuery() *bson.M {
	return nil
}

type RPCTaskCollectionsManagerSt struct{
	mcmongo.CollectionManagerBaseSt
	Inited bool
}

func NewTaskCollectionsManager() *RPCTaskCollectionsManagerSt {
	return &RPCTaskCollectionsManagerSt{}
}

func (this *RPCTaskCollectionsManagerSt) Init(dbsession *mgo.Session, dbName string) {
	this.AddModel(&RPCTaskSt{})
	this.InitManager(dbsession, dbName, "tasks")
}

func (this *RPCTaskCollectionsManagerSt) FindSecondsTasks(tasks *RPCTasksListSt) {
	this.FindAllModels(&bson.M{"execTime": time.Now(), "inExec": false, "secondsCheck": true}, tasks)
}

func (this *RPCTaskCollectionsManagerSt) FindSecondsNotInExecTasks(tasks *RPCTasksListSt) {
	this.FindAllModels(&bson.M{"execTime": &bson.M{"$lt": time.Now()}, "inExec": false, "secondsCheck": true}, tasks)
}

type RPCTasksManagerSt struct {
	ColManager *RPCTaskCollectionsManagerSt
}

func NewTasksManager() *RPCTasksManagerSt {
	return &RPCTasksManagerSt{
		NewTaskCollectionsManager(),
	}
}

func (this *RPCTasksManagerSt) AddNewTask(rpcMsg *mccommunication.RPCMsg, execTime time.Time, ) {

}

func (this *RPCTasksManagerSt) StartMainLoop() {
	print("StartMainLoop")
	for range time.Tick(time.Second) {
		// 1. Go to DB and get tasksList for this second
		tasksList := NewRPCTasksList()
		this.ColManager.FindSecondsTasks(tasksList)
		// 2. Go through tasksList and exec them
		for _, task := range *tasksList {
			task.InExec = true
			if byteD, err := task.RPCMsg.MarshalJSON(); err != nil {
				errorMap := &map[string]interface{}{"message": err.Error()}
				task.RPCMsg.Error = errorMap
				task.Done = true
				task.DoneError = errorMap
			} else {
				innerRPCMan.Service.Publish(task.RPCMsg.Method, byteD)
			}
			this.ColManager.SaveModel(&task)
		}
	}
	print("DoneMainLoop")
}