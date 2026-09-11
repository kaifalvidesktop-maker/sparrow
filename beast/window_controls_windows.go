//go:build windows

package main

import (
	"syscall"
)

const (
	gwlStyle        = 0xfffffff0
	wsCaption       = 0x00C00000
	wsThickFrame    = 0x00040000
	wsMinimizeBox   = 0x00020000
	wsMaximizeBox   = 0x00010000
	wsSysMenu       = 0x00080000
	wmClose         = 0x0010
	wmNCLButtonDown = 0x00A1
	htCaption       = 0x0002
	swMinimize      = 6
	swMaximize      = 3
	swRestore       = 9
	swpNoMove       = 0x0002
	swpNoSize       = 0x0001
	swpNoZOrder     = 0x0004
	swpFrameChanged = 0x0020
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getWindowLongPtr = user32.NewProc("GetWindowLongPtrW")
	setWindowLongPtr = user32.NewProc("SetWindowLongPtrW")
	setWindowPos     = user32.NewProc("SetWindowPos")
	showWindow       = user32.NewProc("ShowWindow")
	isZoomed         = user32.NewProc("IsZoomed")
	postMessage      = user32.NewProc("PostMessageW")
	releaseCapture   = user32.NewProc("ReleaseCapture")
	sendMessage      = user32.NewProc("SendMessageW")
)

func configureNativeWindow(hwnd uintptr) {
	style, _, _ := getWindowLongPtr.Call(hwnd, uintptr(gwlStyle))
	style &^= wsCaption
	style |= wsThickFrame | wsMinimizeBox | wsMaximizeBox | wsSysMenu
	setWindowLongPtr.Call(hwnd, uintptr(gwlStyle), style)
	setWindowPos.Call(hwnd, 0, 0, 0, 0, 0, swpNoMove|swpNoSize|swpNoZOrder|swpFrameChanged)
}

func closeNativeWindow(hwnd uintptr) {
	postMessage.Call(hwnd, wmClose, 0, 0)
}

func minimizeNativeWindow(hwnd uintptr) {
	showWindow.Call(hwnd, swMinimize)
}

func toggleMaximizeNativeWindow(hwnd uintptr) {
	zoomed, _, _ := isZoomed.Call(hwnd)
	if zoomed != 0 {
		showWindow.Call(hwnd, swRestore)
		return
	}
	showWindow.Call(hwnd, swMaximize)
}

func toggleFullscreenNativeWindow(hwnd uintptr) {
	toggleMaximizeNativeWindow(hwnd)
}

func beginNativeWindowDrag(hwnd uintptr) {
	releaseCapture.Call()
	sendMessage.Call(hwnd, wmNCLButtonDown, htCaption, 0)
}
