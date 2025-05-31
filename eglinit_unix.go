// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || freebsd || openbsd
// +build linux freebsd openbsd

package gwsi

import (
	"gwsi/egl"
	"unsafe"
)

/*
#include <assert.h>

#include <xcb/xcb.h>

#include <EGL/egl.h>
#include <EGL/eglext.h>
*/
import "C"

func WsiEglInit() {
	egl.LoadEGL()
}

func WsiGetEGLDisplay(platform *WsiPlatform, pDisplay *egl.EGLDisplay) WsiResult {
	attrib := []egl.EGLAttrib{
		egl.EGL_PLATFORM_XCB_SCREEN_EXT, egl.EGLAttrib(platform.xcb_screen_id),
		egl.EGL_NONE,
	}

	*pDisplay = egl.EglGetPlatformDisplay(egl.EGL_PLATFORM_XCB_EXT, unsafe.Pointer(platform.xcb_connection), attrib)
	if *pDisplay == egl.NilEGLDisplay {
		return WSI_ERROR_EGL
	}

	return WSI_SUCCESS
}

func WsiCreateWindowEGLSurface(window *WsiWindow, dpy egl.EGLDisplay, config egl.EGLConfig, pSurface *egl.EGLSurface) WsiResult {
	visualid := egl.EGLint(0)
	if id, ok := egl.EglGetConfigAttrib(dpy, config, egl.EGL_NATIVE_VISUAL_ID); ok == false {
		return WSI_ERROR_EGL
	} else {
		visualid = id
	}
	platform := window.Platform

	window.XcbColormap = C.xcb_generate_id(platform.xcb_connection)
	if window.XcbColormap == C.UINT32_MAX {
		return WSI_ERROR_PLATFORM
	}

	C.xcb_create_colormap_checked(
		platform.xcb_connection,
		C.XCB_COLORMAP_ALLOC_NONE,
		window.XcbColormap,
		platform.xcb_screen.root,
		C.xcb_visualid_t(visualid))

	attribs := []C.uint32_t{window.XcbColormap}

	C.xcb_change_window_attributes_checked(
		platform.xcb_connection,
		window.XcbWindow,
		C.XCB_CW_COLORMAP,
		unsafe.Pointer(&attribs[0]),
	)

	*pSurface = egl.EglCreatePlatformWindowSurface(dpy, config, egl.PtrXcbWindow(&window.XcbWindow), nil)
	if *pSurface == egl.NilEGLSurface {
		return WSI_ERROR_EGL
	}

	return WSI_SUCCESS
}

func WsiDestroyWindowEGLSurface(window *WsiWindow, dpy egl.EGLDisplay, surface egl.EGLSurface) {

	egl.EglDestroySurface(dpy, surface)
	platform := window.Platform

	value_list := []C.uint32_t{
		platform.xcb_screen.default_colormap,
	}

	C.xcb_change_window_attributes(
		platform.xcb_connection,
		window.XcbWindow,
		C.XCB_CW_COLORMAP,
		unsafe.Pointer(&value_list[0]))

	C.xcb_free_colormap(platform.xcb_connection, window.XcbColormap)
	window.XcbColormap = C.XCB_NONE
}
