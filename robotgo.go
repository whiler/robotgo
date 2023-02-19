// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

/*
Package robotgo Go native cross-platform system automation.
Please make sure Golang, GCC is installed correctly before installing RobotGo;
See Requirements:

	https://github.com/go-vgo/robotgo#requirements

Installation:
With Go module support (Go 1.11+), just import:

	import "github.com/go-vgo/robotgo"

Otherwise, to install the robotgo package, run the command:

	go get -u github.com/go-vgo/robotgo
*/
package robotgo

/*
#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit
#cgo darwin LDFLAGS: -framework Carbon -framework CoreFoundation

#cgo linux CFLAGS: -I/usr/src
#cgo linux LDFLAGS: -L/usr/src -lX11 -lXtst -lm

#cgo windows LDFLAGS: -lgdi32 -luser32

#include "screen/goScreen.h"
#include "mouse/goMouse.h"
#include "key/goKey.h"
#include "window/goWindow.h"
*/
import "C"

import (

	// "os"

	"time"
	"unsafe"

	// "syscall"
	"math/rand"
)

const (
	// Version get the robotgo version
	Version = "v0.100.10-kmactor, whiler"
)

var (
	MouseSleep = 0
	KeySleep   = 0

	Special = map[string]string{
		"~": "`",
		"!": "1",
		"@": "2",
		"#": "3",
		"$": "4",
		"%": "5",
		"^": "6",
		"&": "7",
		"*": "8",
		"(": "9",
		")": "0",
		"_": "-",
		"+": "=",
		"{": "[",
		"}": "]",
		"|": "\\",
		":": ";",
		`"`: "'",
		"<": ",",
		">": ".",
		"?": "/",
	}
)

var mouses = map[string]C.MMMouseButton{
	"left":       C.LEFT_BUTTON,
	"center":     C.CENTER_BUTTON,
	"right":      C.RIGHT_BUTTON,
	"wheelDown":  C.WheelDown,
	"wheelUp":    C.WheelUp,
	"wheelLeft":  C.WheelLeft,
	"wheelRight": C.WheelRight,
}

// GetVersion get the robotgo version
func GetVersion() string {
	return Version
}

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Millisecond)
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Second)
}

/*
      _______.  ______ .______       _______  _______ .__   __.
    /       | /      ||   _  \     |   ____||   ____||  \ |  |
   |   (----`|  ,----'|  |_)  |    |  |__   |  |__   |   \|  |
    \   \    |  |     |      /     |   __|  |   __|  |  . `  |
.----)   |   |  `----.|  |\  \----.|  |____ |  |____ |  |\   |
|_______/     \______|| _| `._____||_______||_______||__| \__|
*/

// SysScale get the sys scale
func SysScale() float64 {
	s := C.sys_scale()
	return float64(s)
}

// Scaled0 return int(x * f)
func Scaled0(x int, f float64) int {
	return int(float64(x) * f)
}

// GetScreenSize get the screen size
func GetScreenSize() (int, int) {
	size := C.get_screen_size()
	// fmt.Println("...", size, size.width)
	return int(size.w), int(size.h)
}

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

// CheckMouse check the mouse button
func CheckMouse(btn string) C.MMMouseButton {
	// button = args[0].(C.MMMouseButton)
	if v, ok := mouses[btn]; ok {
		return v
	}
	return C.LEFT_BUTTON
}

// Move move the mouse to (x, y)
//
// Examples:
//
//	robotgo.MouseSleep = 100  // 100 millisecond
//	robotgo.Move(10, 10)
func Move(x, y int) {
	// if runtime.GOOS == "windows" {
	// 	f := ScaleF()
	// 	x, y = Scaled0(x, f), Scaled0(y, f)
	// }

	cx := C.int32_t(x)
	cy := C.int32_t(y)
	C.move_mouse(cx, cy)

	MilliSleep(MouseSleep)
}

// Toggle toggle the mouse, support button:
//
//		"left", "center", "right",
//	 "wheelDown", "wheelUp", "wheelLeft", "wheelRight"
//
// Examples:
//
//	robotgo.Toggle("left") // default is down
//	robotgo.Toggle("left", "up")
func Toggle(key ...string) int {
	var button C.MMMouseButton = C.LEFT_BUTTON
	if len(key) > 0 {
		button = CheckMouse(key[0])
	}
	down := C.CString("down")
	if len(key) > 1 {
		down = C.CString(key[1])
	}

	i := C.mouse_toggle(down, button)
	C.free(unsafe.Pointer(down))
	MilliSleep(MouseSleep)
	return int(i)
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______/

*/

// KeyToggle toggle the keyboard, if there not have args default is "down"
//
// See keys:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
//
// Examples:
//
//	robotgo.KeyToggle("a")
//	robotgo.KeyToggle("a", "up")
//
//	robotgo.KeyToggle("a", "up", "alt", "cmd")
func KeyToggle(key string, args ...string) string {
	if len(args) <= 0 {
		args = append(args, "down")
	}

	if _, ok := Special[key]; ok {
		key = Special[key]
		if len(args) <= 1 {
			args = append(args, "shift")
		}
	}

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	ckeyArr := make([](*C.char), 0)
	if len(args) > 3 {
		num := len(args)
		for i := 0; i < num; i++ {
			ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(args[i]))))
		}

		str := C.key_Toggles(ckey, (**C.char)(unsafe.Pointer(&ckeyArr[0])), C.int(num))
		MilliSleep(KeySleep)
		return C.GoString(str)
	}

	// use key_toggle()
	var (
		down, mKey, mKeyT = "null", "null", "null"
		// keyDelay = 10
	)

	if len(args) > 0 {
		down = args[0]

		if len(args) > 1 {
			mKey = args[1]
			if len(args) > 2 {
				mKeyT = args[2]
			}
		}
	}

	cdown := C.CString(down)
	cmKey := C.CString(mKey)
	cmKeyT := C.CString(mKeyT)

	str := C.key_toggle(ckey, cdown, cmKey, cmKeyT)
	// str := C.key_Toggle(ckey, cdown, cmKey, cmKeyT, C.int(keyDelay))

	C.free(unsafe.Pointer(cdown))
	C.free(unsafe.Pointer(cmKey))
	C.free(unsafe.Pointer(cmKeyT))

	MilliSleep(KeySleep)
	return C.GoString(str)
}

// KeyPress press key string
func KeyPress(key string) {
	KeyDown(key)
	Sleep(1 + rand.Intn(3))
	KeyUp(key)
}

// KeyDown press down a key
func KeyDown(key string) {
	KeyToggle(key, "down")
}

// KeyUp press up a key
func KeyUp(key string) {
	KeyToggle(key, "up")
}
