// emulate_right_click translates CONTROL+left_click to a right_click mouse event.
//
// The combination CONTROL+left_click is hardcoded for now, but changing it is relatively
// simple.
//
// Requirements:
//
//   * Go compiler: http://golang.org/
//   * gcc compiler for windows:
//     http://mingw-w64.sourceforge.net/
//   * Golang package, installed with
//     go get github.com/AllenDang/w32
//     go install github.com/AllenDang/w32
//
// Documentation reference:
//   * godoc for package github.com/AllenDang/w32, and looking at the source code.
//   * Example code in http://play.golang.org/p/kwfYDhhiqk
//   * C++ example of mouse emulation (oudated but it was helpful to get an idea):
//     http://www.codeproject.com/Articles/194265/Mouse-emulating-software
//   * Windows API Index:
//     https://msdn.microsoft.com/en-US/library/windows/desktop/ff818516(v=vs.85).aspx
package main

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/AllenDang/w32"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procExitProcess = modkernel32.NewProc("ExitProcess")
)

const (
	MOUSEEVENTF_RIGHTDOWN = 0x0008
	MOUSEEVENTF_RIGHTUP   = 0x0010
)

const (
	HC_ACTION   = 0
	HC_NOREMOVE = 3
)

const (
	VK_A_KEY uint16 = 0x41
	VK_S_KEY uint16 = 0x53
)

var hHook w32.HHOOK

func main() {
	hinst := w32.GetModuleHandle("")
	//hHook = w32.SetWindowsHookEx(w32.WH_KEYBOARD_LL, KeyboardHook, hinst, 0)
	hHook = w32.SetWindowsHookEx(w32.WH_MOUSE_LL, MouseHook, hinst, 0)
	for w32.GetMessage(nil, 0, w32.WM_QUIT, w32.WM_QUIT+1) != 0 {
		// NOP while not WM_QUIT
	}
	w32.UnhookWindowsHookEx(hHook)
}

func ExitProcess(uExitCode uint32) {
	procExitProcess.Call(uintptr(uExitCode))
}

func MouseHook(nCode int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	if nCode == HC_ACTION && (wParam == w32.WM_LBUTTONDOWN || wParam == w32.WM_LBUTTONUP) {
		if checkClick(wParam) {
			// We processed and want to consume the event.
			return 1
		}
	}

	// Pass-through event.
	return w32.CallNextHookEx(hHook, nCode, wParam, lParam)
}

func checkClick(wParam w32.WPARAM) bool {
	// Check that CONTROL is clicked.
	if w32.GetAsyncKeyState(w32.VK_CONTROL)&uint16(0x8000) == 0 {
		return false
	}

	// Secret exit: Control + A + Left click
	if w32.GetAsyncKeyState(int(VK_A_KEY))&uint16(0x8000) != 0 {
		w32.UnhookWindowsHookEx(hHook)
		log.Printf("Finishing ...\n")
		ExitProcess(0)
		return false
	}

	if wParam == w32.WM_LBUTTONUP {
		rightButtonUp()
	} else {
		rightButtonDown()
	}
	return true
}

func rightButtonUp() {
	var inputs []w32.INPUT
	inputs = append(inputs, w32.INPUT{
		Type: w32.INPUT_MOUSE,
		Mi: w32.MOUSEINPUT{
			DwFlags: MOUSEEVENTF_RIGHTUP,
		},
	})
	w32.SendInput(inputs)
}

func rightButtonDown() {
	var inputs []w32.INPUT
	inputs = append(inputs, w32.INPUT{
		Type: w32.INPUT_MOUSE,
		Mi: w32.MOUSEINPUT{
			DwFlags: MOUSEEVENTF_RIGHTDOWN,
		},
	})
	w32.SendInput(inputs)
}

func KeyboardHook(nCode int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	if nCode == HC_ACTION && (wParam == w32.WM_KEYDOWN || wParam == w32.WM_KEYUP) {
		if checkKey(nCode, wParam, lParam) {
			return 1
		}
	}
	return w32.CallNextHookEx(hHook, nCode, wParam, lParam)
}

func checkKey(nCode int, wParam w32.WPARAM, lParam w32.LPARAM) bool {
	hookStruct := (*w32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	if hookStruct.VkCode != w32.DWORD(VK_S_KEY) || w32.GetAsyncKeyState(w32.VK_CONTROL)&uint16(0x8000) == 0 {
		return false
	}
	if wParam == w32.WM_KEYDOWN {
		rightButtonUp()
	} else {
		rightButtonDown()
	}
	return true
}
