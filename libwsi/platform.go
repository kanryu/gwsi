package libwsi

/*
#include <stdlib.h>
#include <string.h>

#include <xcb/xcb.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type WsiPlatform struct {
	xcb_connection            *C.xcb_connection_t
	xcb_screen                *C.xcb_screen_t
	xcb_screen_id             C.int
	xcb_atom_wm_protocols     C.xcb_atom_t
	xcb_atom_wm_delete_window C.xcb_atom_t
	WindowList                []*WsiWindow
}

var (
	mapAtoms     map[string]*C.xcb_atom_t
	platformList []*WsiPlatform
)

func (p *WsiPlatform) QueryExtension(name string) bool {
	cookie := C.xcb_query_extension(
		p.xcb_connection,
		C.ushort(len(name)),
		C.CString(name),
	)

	reply := C.xcb_query_extension_reply(
		p.xcb_connection,
		cookie,
		nil)

	present := false
	if reply != nil {
		present = reply.present != 0
		C.free(unsafe.Pointer(reply))
	}

	return present
}

func (p *WsiPlatform) InitAtoms() {
	mapAtoms["WM_PROTOCOLS"] = &p.xcb_atom_wm_protocols
	mapAtoms["WM_DELETE_WINDOW"] = &p.xcb_atom_wm_delete_window

	for name, atom := range mapAtoms {
		cookie := C.xcb_intern_atom(
			p.xcb_connection, 0, C.ushort(len(name)), C.CString(name))
		reply := C.xcb_intern_atom_reply(
			p.xcb_connection, cookie, nil)

		if reply != nil {
			*atom = reply.atom
			C.free(unsafe.Pointer(reply))
		} else {
			*atom = C.XCB_ATOM_NONE
		}
	}
}

func wsi_xcb_get_screen(setup *C.xcb_setup_t, screen int) *C.xcb_screen_t {
	iter := C.xcb_setup_roots_iterator(setup)
	for ; iter.rem != 0; screen-- {
		if screen == 0 {
			return iter.data
		}
		C.xcb_screen_next(&iter)
	}

	return nil
}

func (p *WsiPlatform) FindWindow(window C.xcb_window_t) *WsiWindow {
	for _, w := range p.WindowList {
		if w.XcbWindow == window {
			return w
		}
	}

	return nil
}

func wsi_list_init(windowList []*WsiWindow) {

}

func WsiCreatePlatform(pCreateInfo *WsiPlatformCreateInfo) (*WsiPlatform, error) {
	platform := &WsiPlatform{}
	result := WSI_ERROR_PLATFORM

	wsi_list_init(platform.WindowList)
	platform.xcb_connection = C.xcb_connect(nil, &platform.xcb_screen_id)
	defer func() {
		if result != WSI_SUCCESS {
			C.xcb_disconnect(platform.xcb_connection)
			platform.xcb_connection = nil
		}
	}()
	if err := C.xcb_connection_has_error(platform.xcb_connection); err > 0 {
		return nil, fmt.Errorf("xcb_connect failed (%d)", err)
	}

	setup := C.xcb_get_setup(platform.xcb_connection)

	platform.xcb_screen = wsi_xcb_get_screen(setup, platform.xcb_screen_id)
	if !platform.xcb_screen {
		return nil, fmt.Errorf("xcb_get_screen failed")
	}

	if !platform.QueryExtension("RANDR") {
		return nil, fmt.Errorf("extension RANDR not found")
	}

	if !platform.QueryExtension("XInputExtension") {
		return nil, fmt.Errorf("extension XInputExtension not found")
	}

	platform.InitAtoms()
	result = WSI_SUCCESS
	return platform, nil
}

func (p *WsiPlatform) Destroy() {
	C.xcb_disconnect(p.xcb_connection)
}

func (p *WsiPlatform) DispatchEvents(timeout int64) WsiResult {
	for {
		event := C.xcb_poll_for_event(platform.xcb_connection)
		if event == nil {
			break
		}

		switch event.response_type & ~0x80 {
		case C.XCB_CONFIGURE_NOTIFY:
			notify := (*C.xcb_configure_notify_event_t)(event)
			window := p.FindWindow(notify.window)
			if window != nil {
				wsi_window_xcb_configure_notify(window, notify)
			}
			break
		case C.XCB_CLIENT_MESSAGE:
			message := (*C.xcb_client_message_event_t)(event)

			window := p.FindWindow(message.window)
			if window != nil {
				wsi_window_xcb_client_message(window, message)
			}
			break
		}

		C.free(event)
	}

	return WSI_SUCCESS
}
