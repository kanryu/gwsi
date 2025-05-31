// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || freebsd || openbsd
// +build linux freebsd openbsd

package egl

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

type (
	EGLBoolean           = uint32
	EGLenum              = uint32
	EGLint               = int32
	EGLDisplay           = uintptr
	EGLConfig            = uintptr
	EGLContext           = uintptr
	EGLSurface           = uintptr
	EGLAttrib            = uintptr
	EGLClientBuffer      = uintptr
	PtrXcbWindow         = unsafe.Pointer
	EGLNativePixmapType  = uintptr
	EGLNativeDisplayType = uintptr
	EGLNativeWindowType  = uintptr
)

var (
	_eglBindAPI              func(api EGLenum) EGLenum
	_eglGetConfigs           func(dpy EGLDisplay, configs *EGLConfig, config_size EGLint, num_config *EGLint) EGLBoolean
	_eglChooseConfig         func(dpy EGLDisplay, attrib_list *EGLint, configs *EGLConfig, config_size EGLint, num_config *EGLint) EGLBoolean
	_eglCreateContext        func(dpy EGLDisplay, config EGLConfig, share_context EGLContext, attrib_list *EGLint) EGLContext
	_eglCreatePbufferSurface func(dpy EGLDisplay, config EGLConfig, attrib_list *EGLint) EGLSurface
	_eglCreatePixmapSurface  func(dpy EGLDisplay, config EGLConfig, pixmap EGLNativePixmapType, attrib_list *EGLint) EGLSurface
	_eglCreateWindowSurface  func(dpy EGLDisplay, config EGLConfig, win EGLNativeWindowType, attrib_list *EGLint) EGLSurface
	_eglDestroyContext       func(dpy EGLDisplay, ctx EGLContext) EGLBoolean
	_eglDestroySurface       func(dpy EGLDisplay, surface EGLSurface) EGLBoolean
	_eglGetConfigAttrib      func(dpy EGLDisplay, config EGLConfig, attribute EGLint, value *EGLint) EGLBoolean
	_eglGetCurrentSurface    func(readdraw EGLint) EGLSurface
	_eglGetDisplay           func(display_id EGLNativeDisplayType) EGLDisplay
	_eglInitialize           func(dpy EGLDisplay, major *EGLint, minor *EGLint) EGLBoolean
	_eglMakeCurrent          func(dpy EGLDisplay, draw EGLSurface, read EGLSurface, ctx EGLContext) EGLBoolean
	_eglQueryContext         func(dpy EGLDisplay, ctx EGLContext, attribute EGLint, value *EGLint) EGLBoolean
	_eglQueryString          func(dpy EGLDisplay, name EGLint) string
	_eglQuerySurface         func(dpy EGLDisplay, surface EGLSurface, attribute EGLint, value *EGLint) EGLBoolean
	_eglSwapBuffers          func(dpy EGLDisplay, surface EGLSurface) EGLBoolean
	_eglTerminate            func(dpy EGLDisplay) EGLBoolean
	_eglWaitNative           func(engine EGLint) EGLBoolean
	_eglGetCurrentDisplay    func() EGLDisplay
	_eglGetProcAddress       func(procname string) uintptr
	_eglGetError             func() EGLint
	_eglWaitGL               func() EGLBoolean

	_eglQueryAPI                      func() EGLenum
	_eglCreatePbufferFromClientBuffer func(dpy EGLDisplay, buftype EGLenum, buffer EGLClientBuffer, config EGLConfig, attrib_list *EGLint) EGLSurface
	_eglReleaseThread                 func() EGLBoolean
	_eglWaitClient                    func() EGLBoolean

	_eglBindTexImage    func(dpy EGLDisplay, surface EGLSurface, buffer EGLint) EGLBoolean
	_eglReleaseTexImage func(dpy EGLDisplay, surface EGLSurface, buffer EGLint) EGLBoolean
	_eglSurfaceAttrib   func(dpy EGLDisplay, surface EGLSurface, attribute EGLint, value EGLint) EGLBoolean
	_eglSwapInterval    func(dpy EGLDisplay, interval EGLint) EGLBoolean

	EglGetPlatformDisplay          func(att uint32, conn unsafe.Pointer, attribs []EGLAttrib) EGLDisplay
	EglCreatePlatformWindowSurface func(disp EGLDisplay, conf EGLConfig, win PtrXcbWindow, attribs *EGLAttrib) EGLSurface
)

func LoadEGL() error {
	lib, err := purego.Dlopen(getLibEGL(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	purego.RegisterLibFunc(&_eglBindAPI, lib, "eglBindAPI")
	purego.RegisterLibFunc(&_eglGetConfigs, lib, "eglGetConfigs")
	purego.RegisterLibFunc(&_eglChooseConfig, lib, "eglChooseConfig")
	purego.RegisterLibFunc(&_eglCreateContext, lib, "eglCreateContext")
	purego.RegisterLibFunc(&_eglCreatePbufferSurface, lib, "eglCreatePbufferSurface")
	purego.RegisterLibFunc(&_eglCreatePixmapSurface, lib, "eglCreatePixmapSurface")
	purego.RegisterLibFunc(&_eglCreateWindowSurface, lib, "eglCreateWindowSurface")
	purego.RegisterLibFunc(&_eglDestroyContext, lib, "eglDestroyContext")
	purego.RegisterLibFunc(&_eglDestroySurface, lib, "eglDestroySurface")
	purego.RegisterLibFunc(&_eglGetConfigAttrib, lib, "eglGetConfigAttrib")
	purego.RegisterLibFunc(&_eglGetCurrentSurface, lib, "eglGetCurrentSurface")
	purego.RegisterLibFunc(&_eglGetDisplay, lib, "eglGetDisplay")
	purego.RegisterLibFunc(&_eglInitialize, lib, "eglInitialize")
	purego.RegisterLibFunc(&_eglMakeCurrent, lib, "eglMakeCurrent")
	purego.RegisterLibFunc(&_eglQueryContext, lib, "eglQueryContext")
	purego.RegisterLibFunc(&_eglQueryString, lib, "eglQueryString")
	purego.RegisterLibFunc(&_eglQuerySurface, lib, "eglQuerySurface")
	purego.RegisterLibFunc(&_eglSwapBuffers, lib, "eglSwapBuffers")
	purego.RegisterLibFunc(&_eglTerminate, lib, "eglTerminate")
	purego.RegisterLibFunc(&_eglWaitNative, lib, "eglWaitNative")
	purego.RegisterLibFunc(&_eglGetCurrentDisplay, lib, "eglGetCurrentDisplay")
	purego.RegisterLibFunc(&_eglGetProcAddress, lib, "eglGetProcAddress")
	purego.RegisterLibFunc(&_eglGetError, lib, "eglGetError")
	purego.RegisterLibFunc(&_eglWaitGL, lib, "eglWaitGL")
	purego.RegisterLibFunc(&_eglQueryAPI, lib, "eglQueryAPI")
	purego.RegisterLibFunc(&_eglCreatePbufferFromClientBuffer, lib, "eglCreatePbufferFromClientBuffer")
	purego.RegisterLibFunc(&_eglReleaseThread, lib, "eglReleaseThread")
	purego.RegisterLibFunc(&_eglWaitClient, lib, "eglWaitClient")
	purego.RegisterLibFunc(&_eglBindTexImage, lib, "eglBindTexImage")
	purego.RegisterLibFunc(&_eglReleaseTexImage, lib, "eglReleaseTexImage")
	purego.RegisterLibFunc(&_eglSurfaceAttrib, lib, "eglSurfaceAttrib")
	purego.RegisterLibFunc(&_eglSwapInterval, lib, "eglSwapInterval")

	purego.RegisterLibFunc(&EglGetPlatformDisplay, lib, "eglGetPlatformDisplay")
	purego.RegisterLibFunc(&EglCreatePlatformWindowSurface, lib, "eglCreatePlatformWindowSurface")
	return nil
}

func getLibEGL() string {
	switch runtime.GOOS {
	case "darwin":
		return "/usr/lib/libSystem.B.dylib"
	case "linux":
		return "libEGL.so"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func EglBindAPI(api uint) bool {
	if ret := _eglBindAPI(EGLenum(api)); ret != EGL_TRUE {
		return false
	}
	return true
}

func EglGetConfigs(disp EGLDisplay, configs []EGLConfig, configSize int, numConfig *int) bool {
	var num_config EGLint
	config_size := EGLint(configSize)
	if configs == nil {
		if _eglGetConfigs(disp, nil, 0, &num_config) != EGL_TRUE {
			return false
		}
		if numConfig != nil {
			*numConfig = int(num_config)
		}
		return true
	}
	if _eglGetConfigs(disp, &configs[0], config_size, &num_config) != EGL_TRUE {
		return false
	}
	if numConfig != nil {
		*numConfig = int(num_config)
	}
	return true
}

func EglChooseConfig(disp EGLDisplay, attribs []EGLint) (EGLConfig, bool) {
	var cfg EGLConfig
	var ncfg EGLint
	if _eglChooseConfig(disp, &attribs[0], &cfg, 1, &ncfg) != EGL_TRUE {
		return NilEGLConfig, false
	}
	return EGLConfig(cfg), true
}

func EglCreateContext(disp EGLDisplay, cfg EGLConfig, shareCtx EGLContext, attribs []EGLint) EGLContext {
	ctx := _eglCreateContext(disp, cfg, shareCtx, &attribs[0])
	return EGLContext(ctx)
}

func EglDestroySurface(disp EGLDisplay, surf EGLSurface) bool {
	return _eglDestroySurface(disp, surf) == EGL_TRUE
}

func EglDestroyContext(disp EGLDisplay, ctx EGLContext) bool {
	return _eglDestroyContext(disp, ctx) == EGL_TRUE
}

func EglGetConfigAttrib(disp EGLDisplay, cfg EGLConfig, attr EGLint) (EGLint, bool) {
	var val EGLint
	ret := _eglGetConfigAttrib(disp, cfg, attr, &val)
	return val, ret == EGL_TRUE
}

func EglGetError() EGLint {
	return _eglGetError()
}

func EglInitialize(disp EGLDisplay) (EGLint, EGLint, bool) {
	var maj, min EGLint
	ret := _eglInitialize(disp, &maj, &min)
	return maj, min, ret == EGL_TRUE
}

func EglMakeCurrent(disp EGLDisplay, draw, read EGLSurface, ctx EGLContext) bool {
	return _eglMakeCurrent(disp, draw, read, ctx) == EGL_TRUE
}

func EglReleaseThread() bool {
	return _eglReleaseThread() == EGL_TRUE
}

func EglSwapBuffers(disp EGLDisplay, surf EGLSurface) bool {
	return _eglSwapBuffers(disp, surf) == EGL_TRUE
}

func EglSwapInterval(disp EGLDisplay, interval EGLint) bool {
	return _eglSwapInterval(disp, interval) == EGL_TRUE
}

func EglTerminate(disp EGLDisplay) bool {
	return _eglTerminate(disp) == EGL_TRUE
}

func EglQueryString(disp EGLDisplay, name EGLint) string {
	return _eglQueryString(disp, name)
}

func EglGetDisplay(disp EGLNativeDisplayType) EGLDisplay {
	return _eglGetDisplay(disp)
}

func EglCreateWindowSurface(disp EGLDisplay, conf EGLConfig, win EGLNativeWindowType, attribs []EGLint) EGLSurface {
	eglSurf := _eglCreateWindowSurface(disp, conf, win, &attribs[0])
	return eglSurf
}

// func EglGetPlatformDisplay(att uint32, conn unsafe.Pointer, attribs []EGLAttrib) EGLDisplay {
// 	eglSurf := _eglGetPlatformDisplay(_uint(att), conn, &attribs[0])
// 	return eglSurf
// }

// func EglCreatePlatformWindowSurface(disp EGLDisplay, conf EGLConfig, win PtrXcbWindow, attribs *XCBAttrib) EGLSurface {
// 	eglSurf := _libwsi_eglCreatePlatformWindowSurface(disp, conf, win, nil)
// 	return eglSurf
// }

func EglWaitClient() bool {
	return _eglWaitClient() == EGL_TRUE
}
