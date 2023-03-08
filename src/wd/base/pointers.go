package base

import "time"

type MouseButton int

const (
	LeftButton MouseButton = iota
	MiddleButton
	RightButton
)

// PointerType is the type of pointer used by StorePointerActions.
// There are 3 different types according to the WC3 implementation.
type PointerType string

// PointerAction represents an activity involving a pointer.
type PointerAction map[string]interface{}

const (
	MousePointer PointerType = "mouse"
	PenPointer               = "pen"
	TouchPointer             = "touch"
)

// PointerMoveOrigin controls how the offset for
// the pointer move action is calculated.
type PointerMoveOrigin string

const (
	// FromViewport calculates the offset from the viewport at 0,0.
	FromViewport PointerMoveOrigin = "viewport"
	// FromPointer calculates the offset from the current pointer position.
	FromPointer = "pointer"
)

func PointerPauseAction(duration time.Duration) PointerAction {
	return PointerAction{
		"type":     "pause",
		"duration": uint(duration / time.Millisecond),
	}
}

func PointerMoveAction(duration time.Duration, offset Point, origin PointerMoveOrigin) PointerAction {
	return PointerAction{
		"type":     "pointerMove",
		"duration": uint(duration / time.Millisecond),
		"origin":   origin,
		"x":        offset.X,
		"y":        offset.Y,
	}
}

func PointerUpAction(button MouseButton) PointerAction {
	return PointerAction{
		"type":   "pointerUp",
		"button": button,
	}
}

func PointerDownAction(button MouseButton) PointerAction {
	return PointerAction{
		"type":   "pointerDown",
		"button": button,
	}
}
