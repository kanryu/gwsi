package gwsi

/*
#include <stdio.h>
#include <stdlib.h>
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

type WsiPlatform struct {
	xcb_connection            *C.xcb_connection_t
	xcb_screen                *C.xcb_screen_t
	xcb_screen_id             C.int
	xcb_key_symbols           *C.xcb_key_symbols_t
	xcb_atom_wm_protocols     C.xcb_atom_t
	xcb_atom_wm_delete_window C.xcb_atom_t

	//xcb_im *C.xcb_xim_t
	//xcb_xic         C.xcb_xic_t
	//xim_callback              C.xcb_xim_im_callback
	xcb_im  xcbimdkit.PtrXcbXim
	xcb_xic xcbimdkit.XcbXicT

	WindowList []*WsiWindow
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

	mapAtoms = map[string]*C.xcb_atom_t{
		"WM_PROTOCOLS":     &p.xcb_atom_wm_protocols,
		"WM_DELETE_WINDOW": &p.xcb_atom_wm_delete_window,
	}
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

func (p *WsiPlatform) FindWindow(window xcb.XcbWindow) *WsiWindow {
	for _, w := range p.WindowList {
		if w.XcbWindow == C.uint(window) {
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
	// Init global state for compound text encoding.
	//C.xcb_compound_text_init()
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

	platform.xcb_screen = wsi_xcb_get_screen(setup, int(platform.xcb_screen_id))
	if platform.xcb_screen == nil {
		return nil, fmt.Errorf("xcb_get_screen failed")
	}

	if !platform.QueryExtension("RANDR") {
		return nil, fmt.Errorf("extension RANDR not found")
	}

	if !platform.QueryExtension("XInputExtension") {
		return nil, fmt.Errorf("extension XInputExtension not found")
	}

	platform.InitAtoms()

	if xcbimdkit.LibLoaded {
		xcbimdkit.XcbCompoundTextInit()
		platform.xcb_im = xcbimdkit.XcbXimCreate(xcbimdkit.PtrXcbConnectionT(unsafe.Pointer(platform.xcb_connection)), int32(platform.xcb_screen_id), nil)
		xcbimdkit.XcbXimSetImCallback(platform.xcb_im, &xcbimdkit.ImCallback, 0)
		xcbimdkit.XcbXimSetUseCompoundText(platform.xcb_im, true)
		xcbimdkit.XcbXimSetUseUtf8String(platform.xcb_im, true)

		// Open connection to XIM server.
		xcbimdkit.XcbXimOpen(platform.xcb_im, xcbimdkit.OpenCallback, true, 0)
	}

	result = WSI_SUCCESS
	return platform, nil
}

func (p *WsiPlatform) Destroy() {
	C.xcb_disconnect(p.xcb_connection)
}

func (p *WsiPlatform) DispatchEvents(timeout int64) WsiResult {
	for {
		var conn xcb.PtrXcbConnection
		conn = xcb.PtrXcbConnection(unsafe.Pointer(p.xcb_connection))
		event := xcb.XcbPollForEvent(conn)
		if event == nil {
			break
		}
		evtp := event.ResponseType & 0x7f
		xim_filtered := xcbimdkit.PluginXcbXimFilterEvent(p.xcb_im, event)
		if p.xcb_im == 0 || !xim_filtered {
			switch evtp {
			case C.XCB_CONFIGURE_NOTIFY:
				notify := event.AsXcbConfigureNotifyEventT()
				fmt.Println("DispatchEvents:", "XCB_CONFIGURE_NOTIFY", notify)
				window := p.FindWindow(notify.Window)
				if window != nil {
					wsi_window_xcb_configure_notify(window, notify)
				}
			case C.XCB_CLIENT_MESSAGE:
				message := event.AsXcbClientMessageEventT()

				window := p.FindWindow(message.Window)
				if window != nil {
					wsi_window_xcb_client_message(window, message)
				}
			case C.XCB_EXPOSE:
				expose := event.AsXcbExposeEventT()
				fmt.Println("DispatchEvents:", "XCB_EXPOSE", expose)
			case C.XCB_BUTTON_PRESS:
				button := event.AsXcbButtonPressEventT()
				fmt.Println("DispatchEvents:", "XCB_BUTTON_PRESS", button)
			case C.XCB_BUTTON_RELEASE:
				button := event.AsXcbButtonPressEventT()
				fmt.Println("DispatchEvents:", "XCB_BUTTON_RELEASE", button)
			case C.XCB_MAP_NOTIFY:
				notify := event.AsXcbMapNotifyEventT()
				fmt.Println("DispatchEvents:", "XCB_MAP_NOTIFY", notify)
			case C.XCB_REPARENT_NOTIFY:
				notify := event.AsXcbReparentNotifyEventT()
				fmt.Println("DispatchEvents:", "XCB_REPARENT_NOTIFY", notify)
			case C.XCB_KEY_PRESS:
				notify := event.AsXcbKeyPressEventT()
				fmt.Println("DispatchEvents:", "XCB_KEY_PRESS", notify)
				// Forward event to input method if IC is created.
				if p.xcb_xic != 0 {
					xcbimdkit.XcbXimForwardEvent(xcbimdkit.PtrXcbXim(unsafe.Pointer(p.xcb_im)), p.xcb_xic, event)
				}
			case C.XCB_KEY_RELEASE:
				notify := event.AsXcbKeyPressEventT()
				fmt.Println("DispatchEvents:", "XCB_KEY_RELEASE", notify)
				// Forward event to input method if IC is created.
				if p.xcb_xic != 0 {
					xcbimdkit.XcbXimForwardEvent(xcbimdkit.PtrXcbXim(unsafe.Pointer(p.xcb_im)), p.xcb_xic, event)
				}
			case C.XCB_KEYMAP_NOTIFY:
				notify := event.AsXcbKeymapNotifyEventT()
				fmt.Println("DispatchEvents:", "XCB_KEYMAP_NOTIFY", notify)
			default:
				//fmt.Println("DispatchEvents:", event.response_type, event)
				fmt.Println("DispatchEvents:", event.ResponseType, event)
			}
		}

		C.free(unsafe.Pointer(event))
	}

	return WSI_SUCCESS
}
