package robotgo

import (
	"fmt"
	"testing"
)

func TestBasicFunctions(t *testing.T) {
	w, h := GetScreenSize()
	fmt.Printf("screen: %dx%d\r\n", w, h)
	KeyDown("a")
	KeyUp("a")
	KeyPress("x")
	Move(w/2, h/2)
	Toggle("right", "down")
	Toggle("right", "up")
}
