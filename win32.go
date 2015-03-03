package win32

import ( 
	"syscall"
	"unsafe"
	"fmt"
)

var (
	kernel32, _ = syscall.LoadLibrary("kernel32.dll")
	getModuleHandle, _ = syscall.GetProcAddress(kernel32, "GetModuleHandleW")
	user32, _ = syscall.LoadLibrary("user32.dll")
	getWindowDC, _ = syscall.GetProcAddress(user32, "GetWindowDC")
	getWindowRect, _ = syscall.GetProcAddress(user32, "GetWindowRect")
	findWindow, _ = syscall.GetProcAddress(user32, "FindWindowW")
	sendInput, _ = syscall.GetProcAddress(user32, "SendInput")
	mapVirtualKey, _ = syscall.GetProcAddress(user32, "MapVirtualKeyExW")
	vkKeyScan, vkError = syscall.GetProcAddress(user32, "VkKeyScanW")
	gdi32, _ = syscall.LoadLibrary("gdi32.dll")
	deleteDC, _ = syscall.GetProcAddress(gdi32, "DeleteDC")
	getPixel, _ = syscall.GetProcAddress(gdi32, "GetPixel")
)

type Hwnd uintptr

type ColorRef uint32

type HDC Hwnd 

type Rect struct {
	Left int32
	Top int32
	Right int32
	Bottom int32
}

type Input struct {
	Type int32
	Mi MouseInput
	Ki KeyboardInput
	Hi HardwareInput
}

type keyBoardInput struct {
	Type int32
	Ki KeyboardInput
}

type mouseInput struct {
	Type int32
	Mi MouseInput
}

type hardwareInput struct {
	Type int32
	Hi HardwareInput
}

type MouseInput struct {
	Dx int32
	Dy int32
	MouseData uint32
	DwFlags uint32
	Time uint32
	DwExtraInfo uintptr
}

type KeyboardInput struct {
	WVk uint16
	WScan uint16
	DwFlags uint32
	Time uint32
	DwExtraInfo uintptr
}

type HardwareInput struct {
	UMsg uint32
	wParamL uint16
	wParamH uint16
}

const (
	SRCCOPY = 0x00CC0020
	CLR_INVALD = 0xFFFFFFFF
	INPUT_MOUSE = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWARE = 2
	KEYEVENTF_KEYDOWN = 0
	KEYEVENTF_KEYUP = 0x0002
	KEYEVENTF_KEYUNICODE = 0x0004
	KEYEVENTF_SCANCODE = 0x0008
	MAPVK_VK_TO_CHAR = 2
	MAPVK_VK_TO_VSC = 0
	MAPVK_VSC_TO_VK = 1
	MAPVK_VSC_TO_VK_EX = 3
)

func GetRValue(c ColorRef) byte {
	return byte(c)
}

func GetGValue(c ColorRef) byte {
	return byte(c >> 8)
}

func GetBValue(c ColorRef) byte {
	return byte(c >> 16)
}

func GetWindowDC(handle Hwnd) (HDC, syscall.Errno) {
	var nargs uintptr = 1
	ret, _, callErr := syscall.Syscall(uintptr(getWindowDC), nargs, uintptr(handle), 0, 0)
	return HDC(ret), callErr
}

func GetPixel(hdc HDC, xpos int, ypos int) (ColorRef, syscall.Errno) {
	var nargs uintptr = 1
	ret, _, callErr := syscall.Syscall(uintptr(getPixel), nargs, uintptr(hdc), uintptr(xpos), uintptr(ypos))
	return ColorRef(ret), callErr
}

func GetWindowRect(handle Hwnd, rect *Rect) syscall.Errno {
	var nargs uintptr = 2
	_, _, callErr := syscall.Syscall(uintptr(getWindowRect), nargs, uintptr(handle), uintptr(unsafe.Pointer(rect)), 0)
	return callErr
}

func FindWindow(windowTitle string) (Hwnd, syscall.Errno) {
	var nargs uintptr = 1
	cstring, _ := syscall.UTF16PtrFromString(windowTitle)
	ret, _, callErr := syscall.Syscall(uintptr(findWindow), nargs, 
		uintptr(unsafe.Pointer(cstring)), 0, 0)
	return Hwnd(ret), callErr
}

func GetModuleHandle() (Hwnd, syscall.Errno) {
	var nargs uintptr = 0
	ret, _, callErr := syscall.Syscall(uintptr(getModuleHandle), nargs, 0, 0, 0)
	return Hwnd(ret), callErr
}

func SendInput(inputs []Input) (uintptr, syscall.Errno) {
	//mouseInput has the largest possible union type size, so we take the
	//size of this struct and use it as the possible size of all Input structs
	size := unsafe.Sizeof(mouseInput{})
	validInputs := make([]byte, int(size) * len(inputs))
	for i, input := range inputs {
		switch input.Type {
			case INPUT_KEYBOARD:
				var newInput = keyBoardInput{Type: input.Type, Ki: input.Ki}
				*(*keyBoardInput)(unsafe.Pointer(&validInputs[i * int(size)])) = newInput
			case INPUT_MOUSE:
				var newInput = mouseInput{Type: input.Type, Mi: input.Mi}
				*(*mouseInput)(unsafe.Pointer(&validInputs[i * int(size)])) = newInput
			case INPUT_HARDWARE:
				var newInput = hardwareInput{Type: input.Type, Hi: input.Hi}
				*(*hardwareInput)(unsafe.Pointer(&validInputs[i * int(size)])) = newInput
			default:
				panic("unknown type")
		}
	}
	var nargs uintptr = 3
	ret, _, callErr := syscall.Syscall(uintptr(sendInput), nargs, uintptr(1), uintptr(unsafe.Pointer(&validInputs[0])), uintptr(size))
	return ret, callErr
}

func MapVirtualKey(uCode uint, uMapType uint) (uint, syscall.Errno) {
	var nargs uintptr = 2
	ret, _, callErr := syscall.Syscall(uintptr(mapVirtualKey), nargs, uintptr(uCode), uintptr(uMapType), 0)
	return uint(ret), callErr
}

func VkKeyScan(char uint16) (int16, error) {
	var nargs uintptr = 1
	if vkError == nil {
		ret, _, callErr := syscall.Syscall(uintptr(vkKeyScan), nargs, uintptr(char), 0, 0)
		return int16(ret), callErr
	} else {
		fmt.Printf("%v\n", vkError)
		return 0, nil
	}
}
