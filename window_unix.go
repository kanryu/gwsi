package gwsi

/*
#cgo linux,!android pkg-config: xcb
#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <string.h>

#include <xcb/xcb.h>
#include <xcb/xcb_keysyms.h>

*/
import "C"
import (
	"fmt"
	"gwsi/xcb"
	"gwsi/xcbimdkit"
	"unsafe"
)

type WsiWindow struct {
	Platform    *WsiPlatform
	XcbWindow   C.xcb_window_t
	XcbParent   C.xcb_window_t
	XcbColormap C.xcb_colormap_t
	UserWidth   int
	UserHeight  int

	UserData        unsafe.Pointer
	ConfigureWindow PFN_wsiConfigureWindow
	CloseWindow     PFN_wsiCloseWindow
}

// region XCB Events

func wsi_window_xcb_configure_notify(window *WsiWindow, event *xcb.XcbConfigureNotifyEventT) {

	info := WsiConfigureWindowEvent{
		Base: WsiEvent{
			Type:   WSI_EVENT_TYPE_CONFIGURE_WINDOW,
			Flags:  0,
			Serial: uint32(event.Sequence),
		},
		Window: window,
		Extent: WsiExtent{
			Width:  int32(event.Width),
			Height: int32(event.Height),
		},
	}

	window.ConfigureWindow(window.UserData, &info)
}

func wsi_window_xcb_client_message(window *WsiWindow, event *xcb.XcbClientMessageEventT) {
	data32 := (*uint32)(unsafe.Pointer(&event.Data))
	if C.uint(event.Type) == window.Platform.xcb_atom_wm_protocols &&
		*data32 == uint32(window.Platform.xcb_atom_wm_delete_window) {
		info := &WsiCloseWindowEvent{
			Base: WsiEvent{
				Type:   WSI_EVENT_TYPE_CLOSE_WINDOW,
				Flags:  0,
				Serial: uint32(event.Sequence),
				Time:   0,
			},
			Window: window,
		}

		window.CloseWindow(window.UserData, info)
	}
}

func wsiSetXicCallback2(im xcbimdkit.PtrXcbXim, new_ic xcbimdkit.XcbXicT, userData uintptr) {
	fmt.Println("wsi_set_xic_callback")
	p := (*WsiPlatform)(unsafe.Pointer(userData))
	p.xcb_xic = new_ic
}

func (p *WsiPlatform) CreateWindow(pCreateInfo *WsiWindowCreateInfo, title string) (*WsiWindow, WsiResult) {
	window := &WsiWindow{
		Platform:        p,
		XcbWindow:       C.xcb_generate_id(p.xcb_connection),
		UserWidth:       wsi_xcb_clamp(pCreateInfo.Extent.Width),
		UserHeight:      wsi_xcb_clamp(pCreateInfo.Extent.Height),
		UserData:        pCreateInfo.UserData,
		ConfigureWindow: pCreateInfo.ConfigureWindow,
		CloseWindow:     pCreateInfo.CloseWindow,
	}
	//C.xcb_window = window.XcbWindow

	if pCreateInfo.Parent != nil {
		window.XcbParent = pCreateInfo.Parent.XcbWindow
	} else {
		window.XcbParent = p.xcb_screen.root
	}

	value_mask := C.uint(C.XCB_CW_BACK_PIXEL | C.XCB_CW_EVENT_MASK)
	value_list := []C.uint32_t{
		p.xcb_screen.black_pixel,
		C.XCB_EVENT_MASK_EXPOSURE |
			// XCB_EVENT_MASK_RESIZE_REDIRECT |
			C.XCB_EVENT_MASK_STRUCTURE_NOTIFY |
			C.XCB_EVENT_MASK_KEY_PRESS |
			C.XCB_EVENT_MASK_KEY_RELEASE |
			C.XCB_EVENT_MASK_KEYMAP_STATE |
			C.XCB_EVENT_MASK_FOCUS_CHANGE |
			C.XCB_EVENT_MASK_PROPERTY_CHANGE |
			C.XCB_EVENT_MASK_BUTTON_PRESS |
			C.XCB_EVENT_MASK_BUTTON_RELEASE,
	}

	C.xcb_create_window_checked(
		p.xcb_connection,
		C.XCB_COPY_FROM_PARENT,
		window.XcbWindow,
		window.XcbParent,
		0, 0,
		C.ushort(window.UserWidth),
		C.ushort(window.UserHeight),
		10,
		C.XCB_WINDOW_CLASS_INPUT_OUTPUT,
		C.XCB_COPY_FROM_PARENT,
		value_mask,
		unsafe.Pointer(&value_list[0]),
	)

	properties := []C.xcb_atom_t{
		p.xcb_atom_wm_protocols,
		p.xcb_atom_wm_delete_window,
	}

	C.xcb_change_property(
		p.xcb_connection,
		C.XCB_PROP_MODE_REPLACE,
		window.XcbWindow,
		p.xcb_atom_wm_protocols,
		C.XCB_ATOM_ATOM,
		32,
		C.uint(len(properties)),
		unsafe.Pointer(&properties[0]),
	)

	C.xcb_map_window(p.xcb_connection, window.XcbWindow)
	C.xcb_flush(p.xcb_connection)

	if xcbimdkit.LibLoaded {
		// // Open connection to XIM server.
		// C.gwsi_xcb_xim_open(p.xcb_im, true, unsafe.Pointer(p))
		xcbimdkit.RegisterCreateIcCallback(wsiSetXicCallback2)
		xcbimdkit.XcbXimOpen(xcbimdkit.PtrXcbXim(unsafe.Pointer(p.xcb_im)), xcbimdkit.OpenCallback, true, uintptr(unsafe.Pointer(p)))
	}

	p.WindowList = append(p.WindowList, window)
	return window, WSI_SUCCESS
}

func WsiDestroyWindow(window *WsiWindow) {
	platform := window.Platform

	C.xcb_unmap_window(platform.xcb_connection, window.XcbWindow)
	C.xcb_destroy_window(platform.xcb_connection, window.XcbWindow)
	C.xcb_flush(platform.xcb_connection)

	winlist := []*WsiWindow{}
	for _, w := range platform.WindowList {
		if w != window {
			winlist = append(winlist, w)
		}
	}
	platform.WindowList = winlist
}

func WsiSetWindowParent(window *WsiWindow, parent *WsiWindow) WsiResult {
	platform := window.Platform
	if parent != nil {
		window.XcbParent = parent.XcbWindow
	} else {
		window.XcbParent = platform.xcb_screen.root
	}
	C.xcb_reparent_window(
		platform.xcb_connection,
		window.XcbWindow,
		parent.XcbWindow,
		0, 0)
	return WSI_SUCCESS
}

func wsiSetWindowTitle(window WsiWindow, pTitle string) WsiResult {
	if pTitle != "" {
		C.xcb_change_property(
			window.Platform.xcb_connection,
			C.XCB_PROP_MODE_REPLACE,
			window.XcbWindow,
			C.XCB_ATOM_WM_NAME,
			C.XCB_ATOM_STRING,
			8,
			C.uint(len(pTitle)),
			unsafe.Pointer(C.CString(pTitle)),
		)
	} else {
		C.xcb_delete_property(
			window.Platform.xcb_connection,
			window.XcbWindow,
			C.XCB_ATOM_WM_NAME)
	}

	return WSI_SUCCESS
}

func wsi_xcb_clamp(value int32) int {
	if value < 0 {
		return 0
	}
	if value > C.UINT16_MAX {
		return int(C.UINT16_MAX)
	}
	return int(value)
}
