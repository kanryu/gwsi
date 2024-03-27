//go:build linux || freebsd
// +build linux freebsd

package libwsi

import (
	"libwsi_test/egl"
	"unsafe"
)

/*
#cgo linux LDFLAGS: -L/home/ubuntu/work/libwsi/install/lib/aarch64-linux-gnu -lwsi

#include <EGL/egl.h>
#include <EGL/eglext.h>

#include "common.h"
#include "egl.h"
#include "input.h"
#include "output.h"
#include "platform.h"
#include "vulkan.h"
#include "window.h"

typedef struct InnerWindowCallbacks {
	PFN_wsiConfigureWindow ConfigureWindowCallback;
	PFN_wsiCloseWindow CloseWindowCallback;
} InnerWindowCallbacks;

void configure_window_callback(void *pUserData, const WsiConfigureWindowEvent *pConfig)
{
	InnerWindowCallbacks* inner;
	inner = (InnerWindowCallbacks*)pUserData;
	inner->ConfigureWindowCallback(pUserData, pConfig);
}
void close_window_callback(void *pUserData, const WsiCloseWindowEvent *pClose)
{
	InnerWindowCallbacks* inner;
	inner = (InnerWindowCallbacks*)pUserData;
	inner->CloseWindowCallback(pUserData, pClose);
}
*/
import "C"

type WsiPlatform C.WsiPlatform
type WsiSeat C.WsiSeat
type WsiWindow C.WsiWindow

type WsiPlatformCreateInfo C.WsiPlatformCreateInfo
type WsiWindowCreateInfo C.WsiWindowCreateInfo

type WsiConfigureWindowEvent C.WsiConfigureWindowEvent
type WsiCloseWindowEvent C.WsiCloseWindowEvent

type PFN_wsiConfigureWindow C.PFN_wsiConfigureWindow
type PFN_wsiCloseWindow C.PFN_wsiCloseWindow

type WsiExtent C.WsiExtent
type WsiEvent C.WsiEvent

type PFN_ConfigureWindowCallback func(unsafe.Pointer, *WsiConfigureWindowEvent)
type PFN_CloseWindowCallback func(unsafe.Pointer, *WsiCloseWindowEvent)

// platform

type InnerWindowCallbacks struct {
	ConfigureWindowCallback PFN_ConfigureWindowCallback
	CloseWindowCallback     PFN_CloseWindowCallback
}

func RegisterWindowCallbacks(info *WsiWindowCreateInfo, configure PFN_ConfigureWindowCallback, close PFN_CloseWindowCallback) {
	callbacks := &InnerWindowCallbacks{
		ConfigureWindowCallback: configure,
		CloseWindowCallback:     close,
	}
	info.PUserData = unsafe.Pointer(callbacks)
	info.PfnConfigureWindow = C.PFN_wsiConfigureWindow(C.configure_window_callback)
	info.PfnCloseWindow = C.PFN_wsiCloseWindow(C.close_window_callback)
}

func NewWsiPlatformCreateInfo(structureType WsiStructureType) WsiPlatformCreateInfo {
	platform_info := WsiPlatformCreateInfo{
		SType: C.int32_t(structureType),
		PNext: nil,
	}
	return platform_info
}
func NewWsiWindowCreateInfo(structureType WsiStructureType, title string, extent WsiExtent) WsiWindowCreateInfo {
	ptitle := C.CString("title")
	windowInfo := WsiWindowCreateInfo{
		SType:  C.int32_t(structureType),
		PTitle: ptitle,
		Extent: C.WsiExtent(extent),
	}
	return windowInfo
}

func WsiCreatePlatform(pCreateInfo *WsiPlatformCreateInfo, pPlatform *WsiPlatform) WsiResult {
	var pf C.WsiPlatform
	result := C.wsiCreatePlatform((*C.WsiPlatformCreateInfo)(unsafe.Pointer(pCreateInfo)), &pf)
	if WsiResult(result) == WSI_SUCCESS {
		*pPlatform = WsiPlatform(pf)
	}
	return WsiResult(result)
}

func WsiDestroyPlatform(platform WsiPlatform) {
	C.wsiDestroyPlatform(C.WsiPlatform(platform))
}

// EGL

func WsiGetEGLDisplay(platform WsiPlatform, pDisplay *egl.EGLDisplay) WsiResult {
	var disp C.EGLDisplay
	result := C.wsiGetEGLDisplay(C.WsiPlatform(platform), &disp)
	if WsiResult(result) == WSI_SUCCESS {
		*pDisplay = egl.EGLDisplay(disp)
	}
	return WsiResult(result)
}

func WsiCreateWindowEGLSurface(window WsiWindow, dpy egl.EGLDisplay, config egl.EGLConfig, pSurface *egl.EGLSurface) WsiResult {
	var surface C.EGLSurface
	result := C.wsiCreateWindowEGLSurface(C.WsiWindow(window), C.EGLDisplay(dpy), C.EGLConfig(config), &surface)
	if WsiResult(result) == WSI_SUCCESS {
		*pSurface = egl.EGLSurface(surface)
	}
	return WsiResult(result)
}

func WsiDestroyWindowEGLSurface(window WsiWindow, dpy egl.EGLDisplay, surface egl.EGLSurface) {
	C.wsiDestroyWindowEGLSurface(C.WsiWindow(window), C.EGLDisplay(dpy), C.EGLSurface(surface))
}

// window

func WsiCreateWindow(platform WsiPlatform, pCreateInfo *WsiWindowCreateInfo, pWindow *WsiWindow) WsiResult {
	var window C.WsiWindow
	result := C.wsiCreateWindow(C.WsiPlatform(platform), (*C.WsiWindowCreateInfo)(unsafe.Pointer(pCreateInfo)), &window)
	if WsiResult(result) == WSI_SUCCESS {
		*pWindow = WsiWindow(window)
	}
	return WsiResult(result)
}

func WsiDestroyWindow(window WsiWindow) {
	C.wsiDestroyWindow(C.WsiWindow((window)))
}

func WsiDispatchEvents(platform WsiPlatform, timeout int64) WsiResult {
	result := C.wsiDispatchEvents(C.WsiPlatform(platform), C.int64_t(timeout))
	return WsiResult(result)
}
