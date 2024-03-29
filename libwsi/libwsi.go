package libwsi

import "unsafe"

type WsiPlatformCreateInfo struct {
	Type WsiStructureType
	Next string
}

type WsiWindowCreateInfo struct {
	Type            WsiStructureType
	Parent          *WsiWindow
	Extent          WsiExtent
	Next            uintptr
	Title           string
	UserData        unsafe.Pointer
	ConfigureWindow PFN_wsiConfigureWindow
	CloseWindow     PFN_wsiCloseWindow
}

type WsiConfigureWindowEvent struct {
	Base   WsiEvent
	Window *WsiWindow
	Extent WsiExtent
}

type WsiCloseWindowEvent struct {
	Base   WsiEvent
	Window *WsiWindow
}

type PFN_wsiConfigureWindow func(pUserData unsafe.Pointer, pConfig *WsiConfigureWindowEvent)
type PFN_wsiCloseWindow func(pUserData unsafe.Pointer, pClose *WsiCloseWindowEvent)

type WsiExtent struct {
	Width  int32
	Height int32
}

type WsiEvent struct {
	Type   WsiEventType
	Flags  uint32
	Serial uint32
	Time   int64
}

type PFN_ConfigureWindowCallback func(unsafe.Pointer, *WsiConfigureWindowEvent)
type PFN_CloseWindowCallback func(unsafe.Pointer, *WsiCloseWindowEvent)

// platform

// type InnerWindowCallbacks struct {
// 	ConfigureWindowCallback PFN_ConfigureWindowCallback
// 	CloseWindowCallback     PFN_CloseWindowCallback
// }

// var callbacks InnerWindowCallbacks

// func RegisterWindowCallbacks(info *WsiWindowCreateInfo, configure PFN_ConfigureWindowCallback, close PFN_CloseWindowCallback) {
// 	callbacks = InnerWindowCallbacks{
// 		ConfigureWindowCallback: configure,
// 		CloseWindowCallback:     close,
// 	}
// 	info.PUserData = unsafe.Pointer(&callbacks)
// 	info.PfnConfigureWindow = C.PFN_wsiConfigureWindow(C.configure_window_callback)
// 	info.PfnCloseWindow = C.PFN_wsiCloseWindow(C.close_window_callback)
// }

// func NewWsiPlatformCreateInfo(structureType WsiStructureType) WsiPlatformCreateInfo {
// 	platform_info := WsiPlatformCreateInfo{
// 		SType: C.int32_t(structureType),
// 		PNext: nil,
// 	}
// 	return platform_info
// }
// func NewWsiWindowCreateInfo(structureType WsiStructureType, extent WsiExtent) WsiWindowCreateInfo {
// 	windowInfo := WsiWindowCreateInfo{
// 		SType:  C.int32_t(structureType),
// 		Extent: C.WsiExtent(extent),
// 	}
// 	return windowInfo
// }

// func WsiCreatePlatform(pCreateInfo *WsiPlatformCreateInfo, pPlatform *WsiPlatform) WsiResult {
// 	var pf C.WsiPlatform
// 	result := C.wsiCreatePlatform((*C.WsiPlatformCreateInfo)(unsafe.Pointer(pCreateInfo)), &pf)
// 	if WsiResult(result) == WSI_SUCCESS {
// 		*pPlatform = WsiPlatform(pf)
// 	}
// 	return WsiResult(result)
// }

// func WsiDestroyPlatform(platform WsiPlatform) {
// 	C.wsiDestroyPlatform(C.WsiPlatform(platform))
// }

// EGL

// func WsiGetEGLDisplay(platform WsiPlatform, pDisplay *egl.EGLDisplay) WsiResult {
// 	var disp C.EGLDisplay
// 	result := C.wsiGetEGLDisplay(C.WsiPlatform(platform), &disp)
// 	if WsiResult(result) == WSI_SUCCESS {
// 		*pDisplay = egl.EGLDisplay(disp)
// 	}
// 	return WsiResult(result)
// }

// func WsiCreateWindowEGLSurface(window WsiWindow, dpy egl.EGLDisplay, config egl.EGLConfig, pSurface *egl.EGLSurface) WsiResult {
// 	var surface C.EGLSurface
// 	result := C.wsiCreateWindowEGLSurface(C.WsiWindow(window), C.EGLDisplay(dpy), C.EGLConfig(config), &surface)
// 	if WsiResult(result) == WSI_SUCCESS {
// 		*pSurface = egl.EGLSurface(surface)
// 	}
// 	return WsiResult(result)
// }

// func WsiDestroyWindowEGLSurface(window WsiWindow, dpy egl.EGLDisplay, surface egl.EGLSurface) {
// 	C.wsiDestroyWindowEGLSurface(C.WsiWindow(window), C.EGLDisplay(dpy), C.EGLSurface(surface))
// }

// window

// func WsiCreateWindow(platform WsiPlatform, pCreateInfo *WsiWindowCreateInfo, pWindow *WsiWindow, title string) WsiResult {
// 	ttl := C.CString(title)
// 	p := runtime.Pinner{}
// 	p.Pin(ttl)
// 	defer C.free(unsafe.Pointer(ttl))
// 	defer p.Unpin()
// 	pCreateInfo.PTitle = ttl
// 	var window C.WsiWindow
// 	result := C.wsiCreateWindow(C.WsiPlatform(platform), (*C.WsiWindowCreateInfo)(unsafe.Pointer(pCreateInfo)), &window)
// 	pCreateInfo.PTitle = nil
// 	if WsiResult(result) == WSI_SUCCESS {
// 		*pWindow = WsiWindow(window)
// 	}
// 	return WsiResult(result)
// }

// func WsiDestroyWindow(window WsiWindow) {
// 	C.wsiDestroyWindow(C.WsiWindow((window)))
// }

// func WsiDispatchEvents(platform WsiPlatform, timeout int64) WsiResult {
// 	result := C.wsiDispatchEvents(C.WsiPlatform(platform), C.int64_t(timeout))
// 	return WsiResult(result)
// }
