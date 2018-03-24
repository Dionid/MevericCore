package mcplantainer

type TimersManagerSt struct {
	TimersByDeviceShadowId map[string]interface{}
}

func NewTimersManager() *TimersManagerSt{
	return &TimersManagerSt{
		map[string]interface{}{},
	}
}

var timersManager = NewTimersManager()
