package base

//15.6.4 Wheel actions

type WheelAction map[string]interface{}

func CreateWheelAction(x int, y int, deltaX int, deltaY int) WheelAction {
	return WheelAction{
		"type":   "scroll",
		"x":      x,
		"y":      y,
		"deltaX": deltaX,
		"deltaY": deltaY,
	}
}

func CreateWheelPauseAction(ms int) WheelAction {
	return WheelAction{
		"type":     "pause",
		"duration": ms,
	}
}
