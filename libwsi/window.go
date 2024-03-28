package libwsi

/*
#include <stdlib.h>
#include <assert.h>
#include <string.h>

#include <xcb/xcb.h>
*/
import "C"

type WsiWindow struct {
	Platform    *WsiPlatform
	XcbWindow   C.xcb_window_t
	XcbParent   C.xcb_window_t
	XcbColormap C.xcb_colormap_t
	UserWidth   int
	UserHeight  int

	UserData        uintptr
	ConfigureWindow PFN_wsiConfigureWindow
	CloseWindow     PFN_wsiCloseWindow
}

// region XCB Events

func wsi_window_xcb_configure_notify(window *WsiWindow, event *C.xcb_configure_notify_event_t) {
	assert(event.window == window.xcb_window)

	info := WsiConfigureWindowEvent{
		Base: WsiEvent{
			Type:   WSI_EVENT_TYPE_CONFIGURE_WINDOW,
			Flags:  0,
			Serial: event.sequence,
		},
		Window: window,
		Extent: WsiExtent{
			Width:  event.width,
			Height: event.height,
		},
	}

	window.pfn_configure(window.user_data, &info)
}

func wsi_window_xcb_client_message(window *WsiWindow, event *C.xcb_client_message_event_t) {
	if event.Type == window.platform.xcb_atom_wm_protocols &&
		event.data.data32[0] == window.platform.xcb_atom_wm_delete_window {
		info := &WsiCloseWindowEvent{
			Base: WsiEvent{
				Type:   WSI_EVENT_TYPE_CLOSE_WINDOW,
				Flags:  0,
				Serial: event.sequence,
				Time:   0,
			},
			Window: window,
		}

		window.pfn_close(window.user_data, info)
	}
}

// endregion

func (p *WsiPlatform) CreateWindow(pCreateInfo *WsiWindowCreateInfo, title string) (*WsiWindow, WsiResult) {
	window := &WsiWindow{
		Platform:        p,
		XcbWindow:       xcb_generate_id(p.xcb_connection),
		UserWidth:       wsi_xcb_clamp(pCreateInfo.extent.width),
		UserHeight:      wsi_xcb_clamp(pCreateInfo.extent.height),
		UserData:        pCreateInfo.pUserData,
		ConfigureWindow: pCreateInfo.ConfigureWindow,
		CloseWindow:     pCreateInfo.CloseWindow,
	}

	if pCreateInfo.parent != nil {
		window.xcb_parent = pCreateInfo.parent.xcb_window
	} else {
		window.xcb_parent = p.xcb_screen.root
	}

	value_mask := XCB_CW_BACK_PIXEL | XCB_CW_EVENT_MASK
	value_list := []C.uint32_t{
		p.xcb_screen.black_pixel,
		XCB_EVENT_MASK_EXPOSURE |
			// XCB_EVENT_MASK_RESIZE_REDIRECT |
			XCB_EVENT_MASK_STRUCTURE_NOTIFY |
			XCB_EVENT_MASK_BUTTON_PRESS |
			XCB_EVENT_MASK_BUTTON_RELEASE,
	}

	C.xcb_create_window_checked(
		p.xcb_connection,
		C.XCB_COPY_FROM_PARENT,
		window.XcbWindow,
		window.XcbParent,
		0, 0,
		window.UserWidth,
		window.UserHeight,
		10,
		C.XCB_WINDOW_CLASS_INPUT_OUTPUT,
		C.XCB_COPY_FROM_PARENT,
		value_mask,
		value_list)

	properties := []xcb_atom_t{
		p.xcb_atom_wm_protocols,
		p.xcb_atom_wm_delete_window,
	}

	C.xcb_change_property(
		p.xcb_connection,
		C.XCB_PROP_MODE_REPLACE,
		window.xcb_window,
		p.xcb_atom_wm_protocols,
		C.XCB_ATOM_ATOM,
		32,
		wsi_array_length(properties),
		properties)

	C.xcb_map_window(p.xcb_connection, window.xcb_window)
	C.xcb_flush(p.xcb_connection)

	wsi_list_insert(&p.WindowList, &window.link)
	return window, WSI_SUCCESS
}

func WsiDestroyWindow(window *WsiWindow) {
	platform := window.platform

	C.xcb_unmap_window(platform.xcb_connection, window.xcb_window)
	C.xcb_destroy_window(platform.xcb_connection, window.xcb_window)
	C.xcb_flush(platform.xcb_connection)
}

func WsiSetWindowParent(window *WsiWindow, parent *WsiWindow) WsiResult {
	platform := window.platform
	if parent {
		window.xcb_parent = parent.xcb_window
	} else {
		window.xcb_parent = platform.xcb_screen.root
	}
	C.xcb_reparent_window(
		platform.xcb_connection,
		window.xcb_window,
		parent.xcb_window,
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
			len(pTitle),
			C.CString(pTitle),
		)
	} else {
		C.xcb_delete_property(
			window.Platform.xcb_connection,
			window.XcbWindow,
			C.XCB_ATOM_WM_NAME)
	}

	return WSI_SUCCESS
}

func wsi_xcb_clamp(value uint32) uint16 {
	if value < 0 {
		return 0
	}
	if value > C.UINT16_MAX {
		return uint16(C.UINT16_MAX)
	}
	return uint16(value)
}
