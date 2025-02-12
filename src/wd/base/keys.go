package base

import "strings"

// KeyAction represents an activity involving a keyboard key.
type KeyAction map[string]interface{}

// Special keyboard keys, for SendKeys.
const (
	NullKey       = string('\ue000')
	CancelKey     = string('\ue001')
	HelpKey       = string('\ue002')
	BackspaceKey  = string('\ue003')
	TabKey        = string('\ue004')
	ClearKey      = string('\ue005')
	ReturnKey     = string('\ue006')
	EnterKey      = string('\ue007')
	ShiftKey      = string('\ue008')
	ControlKey    = string('\ue009')
	AltKey        = string('\ue00a')
	PauseKey      = string('\ue00b')
	EscapeKey     = string('\ue00c')
	SpaceKey      = string('\ue00d')
	PageUpKey     = string('\ue00e')
	PageDownKey   = string('\ue00f')
	EndKey        = string('\ue010')
	HomeKey       = string('\ue011')
	LeftArrowKey  = string('\ue012')
	UpArrowKey    = string('\ue013')
	RightArrowKey = string('\ue014')
	DownArrowKey  = string('\ue015')
	InsertKey     = string('\ue016')
	DeleteKey     = string('\ue017')
	SemicolonKey  = string('\ue018')
	EqualsKey     = string('\ue019')
	Numpad0Key    = string('\ue01a')
	Numpad1Key    = string('\ue01b')
	Numpad2Key    = string('\ue01c')
	Numpad3Key    = string('\ue01d')
	Numpad4Key    = string('\ue01e')
	Numpad5Key    = string('\ue01f')
	Numpad6Key    = string('\ue020')
	Numpad7Key    = string('\ue021')
	Numpad8Key    = string('\ue022')
	Numpad9Key    = string('\ue023')
	MultiplyKey   = string('\ue024')
	AddKey        = string('\ue025')
	SeparatorKey  = string('\ue026')
	SubtractKey   = string('\ue027')
	DecimalKey    = string('\ue028')
	DivideKey     = string('\ue029')
	F1Key         = string('\ue031')
	F2Key         = string('\ue032')
	F3Key         = string('\ue033')
	F4Key         = string('\ue034')
	F5Key         = string('\ue035')
	F6Key         = string('\ue036')
	F7Key         = string('\ue037')
	F8Key         = string('\ue038')
	F9Key         = string('\ue039')
	F10Key        = string('\ue03a')
	F11Key        = string('\ue03b')
	F12Key        = string('\ue03c')
	MetaKey       = string('\ue03d')
)

var _mapping = make(map[string]string)

func init() {
	var mapping = map[string]string{
		"Null":      NullKey,
		"Cancel":    CancelKey,
		"Help":      HelpKey,
		"Backspace": BackspaceKey,
		"Tab":       TabKey,
		"Clear":     ClearKey,
		"Return":    ReturnKey,
		"Enter":     EnterKey,
		"Shift":     ShiftKey,
		"Control":   ControlKey,
		"Alt":       AltKey,
		"Pause":     PauseKey,
		"Escape":    EscapeKey,
		"Space":     SpaceKey,
		"PageUp":    PageUpKey,
		"PageDown":  PageDownKey,
		"End":       EndKey,
		"Home":      HomeKey,
		"Left":      LeftArrowKey,
		"Up":        UpArrowKey,
		"Right":     RightArrowKey,
		"Down":      DownArrowKey,
		"Insert":    InsertKey,
		"Delete":    DeleteKey,
		"Semicolon": SemicolonKey,
		"Equals":    EqualsKey,
		"Numpad0":   Numpad0Key,
		"Numpad1":   Numpad1Key,
		"Numpad2":   Numpad2Key,
		"Numpad3":   Numpad3Key,
		"Numpad4":   Numpad4Key,
		"Numpad5":   Numpad5Key,
		"Numpad6":   Numpad6Key,
		"Numpad7":   Numpad7Key,
		"Numpad8":   Numpad8Key,
		"Numpad9":   Numpad9Key,
		"Multiply":  MultiplyKey,
		"Add":       AddKey,
		"Separator": SeparatorKey,
		"Subtract":  SubtractKey,
		"Decimal":   DecimalKey,
		"Divide":    DivideKey,
		"F1":        F1Key,
		"F2":        F2Key,
		"F3":        F3Key,
		"F4":        F4Key,
		"F5":        F5Key,
		"F6":        F6Key,
		"F7":        F7Key,
		"F8":        F8Key,
		"F9":        F9Key,
		"F10":       F10Key,
		"F11":       F11Key,
		"F12":       F12Key,
		"Meta":      MetaKey,
	}
	for k, v := range mapping {
		_mapping[strings.ToLower(k)] = v
	}
}

func KeyFromMapping(key string) string {
	if v, ok := _mapping[strings.ToLower(key)]; ok {
		return v
	}
	return key
}

// KeyPauseAction builds a KeyAction which pauses for the supplied duration.
func KeyPauseAction(ms uint) KeyAction {
	return KeyAction{
		"type":     "pause",
		"duration": ms,
	}
}

func KeyUpAction(key string) KeyAction {
	return KeyAction{
		"type":  "keyUp",
		"value": key,
	}
}

func KeyDownAction(key string) KeyAction {
	return KeyAction{
		"type":  "keyDown",
		"value": key,
	}
}
