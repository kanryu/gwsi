//go:build linux || freebsd
// +build linux freebsd

package libwsi

/*
#cgo linux LDFLAGS: -lwsi

#include "common.h"
#include "input.h"
#include "output.h"
#include "platform.h"
#include "vulkan.h"
#include "window.h"
*/
import "C"
import "unsafe"

// platform

type WsiPlatformCreateInfo struct {
	Type  WsiStructureType
	pNext uintptr
}

func WsiCreatePlatform(pCreateInfo *WsiPlatformCreateInfo, pPlatform *WsiPlatform) WsiResult {
	pf C.WsiPlatform
	result := C.wsiCreatePlatform((*C.WsiPlatformCreateInfo)(unsafe.Pointer(pCreateInfo)), &pf)
	if WsiResult(result) == WSI_SUCCESS {
		*pPlatform = uintptr(pf)
	}
	return WsiResult(result)
}

func WsiDestroyPlatform(platform WsiPlatform) {
	C.wsiDestroyPlatform((*C.WsiPlatform)(unsafe.Pointer(platform)))
}

// EGL

func WsiGetEGLDisplay(platform WsiPlatform, pDisplay *EGLDisplay) WsiResult {
	disp C.EGLDisplay
	result := C.wsiGetEGLDisplay((*C.WsiPlatform)(unsafe.Pointer(platform)), &disp)
	if WsiResult(result) == WSI_SUCCESS {
		*pDisplay = uintptr(disp)
	}
	return WsiResult(result)
}

func WsiCreateWindowEGLSurface(window WsiWindow, dpy EGLDisplay, config EGLConfig, pSurface *EGLSurface) WsiResult {
	surface C.EGLSurface
	result := C.wsiCreateWindowEGLSurface((*C.WsiWindow)(unsafe.Pointer(window)), C.EGLDisplay(dpy), C.EGLConfig(config), &surface)
	if WsiResult(result) == WSI_SUCCESS {
		*pSurface = uintptr(surface)
	}
	return WsiResult(result)
}

func WsiDestroyWindowEGLSurface(window WsiWindow, dpy EGLDisplay, surface EGLSurface) {
	C.wsiDestroyWindowEGLSurface((*C.WsiWindow)(unsafe.Pointer(window)), C.EGLDisplay(dpy), C.EGLSurface(surface))
}

// window

func WsiCreateWindow(platform WsiPlatform, pConfig *WsiConfigureWindowEvent, pWindow *WsiWindow) WsiResult {
	window C.WsiWindow
	result := C.wsiCreateWindow((*C.WsiPlatform)(unsafe.Pointer(platform)), (*C.WsiConfigureWindowEvent)(unsafe.Pointer(pConfig)), &window)
	if WsiResult(result) == WSI_SUCCESS {
		*pWindow = uintptr(window)
	}
	return WsiResult(result)
}

func WsiDestroyWindow(window WsiWindow) {
	C.wsiDestroyWindow((*C.WsiWindow)(unsafe.Pointer(window)))
}

func WsiDispatchEvents(platform WsiPlatform, timeout int64) WsiResult {
	result := C.WsiDispatchEvents((*C.WsiPlatform)(unsafe.Pointer(platform)), C.int64(timeout))
	return WsiResult(result)
}
