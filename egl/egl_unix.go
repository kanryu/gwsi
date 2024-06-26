// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || freebsd || openbsd
// +build linux freebsd openbsd

package egl

/*
#cgo linux,!android  pkg-config: egl
#cgo freebsd openbsd android LDFLAGS: -lEGL
#cgo freebsd CFLAGS: -I/usr/local/include
#cgo freebsd LDFLAGS: -L/usr/local/lib
#cgo openbsd CFLAGS: -I/usr/X11R6/include
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib
#cgo CFLAGS: -DEGL_NO_X11

#include <EGL/egl.h>
#include <EGL/eglext.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

type (
	EGLenum           = C.EGLenum
	EGLint            = C.EGLint
	EGLDisplay        = C.EGLDisplay
	EGLConfig         = C.EGLConfig
	EGLContext        = C.EGLContext
	EGLSurface        = C.EGLSurface
	EGLAttrib         = C.long
	PtrXcbWindow      = unsafe.Pointer
	NativeDisplayType = C.EGLNativeDisplayType
	NativeWindowType  = C.EGLNativeWindowType
)

var (
	_eglBindAPI                    func(api EGLenum) EGLenum
	EglGetPlatformDisplay          func(att uint32, conn unsafe.Pointer, attribs []EGLAttrib) EGLDisplay
	EglCreatePlatformWindowSurface func(disp EGLDisplay, conf EGLConfig, win PtrXcbWindow, attribs *EGLAttrib) EGLSurface
)

func LoadEGL() error {
	lib, err := purego.Dlopen(getLibEGL(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	// var puts func(string)
	// purego.RegisterLibFunc(&puts, lib, "puts")
	// puts("Calling C from Go without Cgo!")
	purego.RegisterLibFunc(&_eglBindAPI, lib, "eglBindAPI")
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
	if ret := _eglBindAPI(EGLenum(api)); ret != C.EGL_TRUE {
		return false
	}
	return true
}

func EglGetConfigs(disp EGLDisplay, configs []EGLConfig, configSize int, numConfig *int) bool {
	var num_config EGLint
	config_size := EGLint(configSize)
	if configs == nil {
		if C.eglGetConfigs(disp, nil, 0, &num_config) != C.EGL_TRUE {
			return false
		}
		if numConfig != nil {
			*numConfig = int(num_config)
		}
		return true
	}
	if C.eglGetConfigs(disp, &configs[0], config_size, &num_config) != C.EGL_TRUE {
		return false
	}
	if numConfig != nil {
		*numConfig = int(num_config)
	}
	return true
}

func EglChooseConfig(disp EGLDisplay, attribs []EGLint) (EGLConfig, bool) {
	var cfg C.EGLConfig
	var ncfg C.EGLint
	if C.eglChooseConfig(disp, &attribs[0], &cfg, 1, &ncfg) != C.EGL_TRUE {
		return NilEGLConfig, false
	}
	return EGLConfig(cfg), true
}

func EglCreateContext(disp EGLDisplay, cfg EGLConfig, shareCtx EGLContext, attribs []EGLint) EGLContext {
	ctx := C.eglCreateContext(disp, cfg, shareCtx, &attribs[0])
	return EGLContext(ctx)
}

func EglDestroySurface(disp EGLDisplay, surf EGLSurface) bool {
	return C.eglDestroySurface(disp, surf) == C.EGL_TRUE
}

func EglDestroyContext(disp EGLDisplay, ctx EGLContext) bool {
	return C.eglDestroyContext(disp, ctx) == C.EGL_TRUE
}

func EglGetConfigAttrib(disp EGLDisplay, cfg EGLConfig, attr EGLint) (EGLint, bool) {
	var val EGLint
	ret := C.eglGetConfigAttrib(disp, cfg, attr, &val)
	return val, ret == C.EGL_TRUE
}

func EglGetError() EGLint {
	return C.eglGetError()
}

func EglInitialize(disp EGLDisplay) (EGLint, EGLint, bool) {
	var maj, min EGLint
	ret := C.eglInitialize(disp, &maj, &min)
	return maj, min, ret == C.EGL_TRUE
}

func EglMakeCurrent(disp EGLDisplay, draw, read EGLSurface, ctx EGLContext) bool {
	return C.eglMakeCurrent(disp, draw, read, ctx) == C.EGL_TRUE
}

func EglReleaseThread() bool {
	return C.eglReleaseThread() == C.EGL_TRUE
}

func EglSwapBuffers(disp EGLDisplay, surf EGLSurface) bool {
	return C.eglSwapBuffers(disp, surf) == C.EGL_TRUE
}

func EglSwapInterval(disp EGLDisplay, interval EGLint) bool {
	return C.eglSwapInterval(disp, interval) == C.EGL_TRUE
}

func EglTerminate(disp EGLDisplay) bool {
	return C.eglTerminate(disp) == C.EGL_TRUE
}

func EglQueryString(disp EGLDisplay, name EGLint) string {
	return C.GoString(C.eglQueryString(disp, name))
}

func EglGetDisplay(disp NativeDisplayType) EGLDisplay {
	return C.eglGetDisplay(disp)
}

func EglCreateWindowSurface(disp EGLDisplay, conf EGLConfig, win NativeWindowType, attribs []EGLint) EGLSurface {
	eglSurf := C.eglCreateWindowSurface(disp, conf, win, &attribs[0])
	return eglSurf
}

// func EglGetPlatformDisplay(att uint32, conn unsafe.Pointer, attribs []EGLAttrib) EGLDisplay {
// 	eglSurf := C.eglGetPlatformDisplay(C.uint(att), conn, &attribs[0])
// 	return eglSurf
// }

// func EglCreatePlatformWindowSurface(disp EGLDisplay, conf EGLConfig, win PtrXcbWindow, attribs *XCBAttrib) EGLSurface {
// 	eglSurf := C.libwsi_eglCreatePlatformWindowSurface(disp, conf, win, nil)
// 	return eglSurf
// }

func EglWaitClient() bool {
	return C.eglWaitClient() == C.EGL_TRUE
}
