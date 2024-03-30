package gwsi

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <xcb/xcb.h>
#include <xcb/xcb_keysyms.h>
#include <xcb-imdkit/encoding.h>
#include <xcb-imdkit/encoding.h>
#include <xcb-imdkit/ximproto.h>
#include <xcb-imdkit/imclient.h>

void forward_event(xcb_xim_t *im, xcb_xic_t ic, xcb_key_press_event_t *event,
                   void *user_data) {
    fprintf(stderr, "Key %s Keycode %u, State %u\n",
            event->response_type == XCB_KEY_PRESS ? "press" : "release",
            event->detail, event->state);
}

void commit_string(xcb_xim_t *im, xcb_xic_t ic, uint32_t flag, char *str,
                   uint32_t length, uint32_t *keysym, size_t nKeySym,
                   void *user_data) {
    if (xcb_xim_get_encoding(im) == XCB_XIM_UTF8_STRING) {
        fprintf(stderr, "key commit utf8: %.*s\n", length, str);
    } else if (xcb_xim_get_encoding(im) == XCB_XIM_COMPOUND_TEXT) {
        size_t newLength = 0;
        char *utf8 = xcb_compound_text_to_utf8(str, length, &newLength);
        if (utf8) {
            int l = newLength;
            fprintf(stderr, "key commit: %.*s\n", l, utf8);
        }
    }
}

void disconnected(xcb_xim_t *im, void *user_data) {
    fprintf(stderr, "Disconnected from input method server.\n");
    //ic = 0;
}

xcb_xim_im_callback callback = {
    .forward_event = forward_event,
    .commit_string = commit_string,
    .disconnected = disconnected,
};
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
	xcb_key_symbols           *C.xcb_key_symbols_t
	xcb_im                    *C.xcb_xim_t
	xcb_xic                   C.xcb_xic_t
	xim_callback              C.xcb_xim_im_callback
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
	// Init global state for compound text encoding.
	C.xcb_compound_text_init()
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

	platform.xcb_im = C.xcb_xim_create(platform.xcb_connection, C.int(platform.xcb_screen_id), nil)

	C.xcb_xim_set_im_callback(platform.xcb_im, &C.callback, nil)
	//C.xcb_xim_set_log_handler(platform.xcb_im, logger)
	C.xcb_xim_set_use_compound_text(platform.xcb_im, true)
	C.xcb_xim_set_use_utf8_string(platform.xcb_im, true)

	// // Open connection to XIM server.
	// C.xcb_xim_open(platform.xcb_im, open_callback, true, nil)

	result = WSI_SUCCESS
	return platform, nil
}

func (p *WsiPlatform) Destroy() {
	C.xcb_disconnect(p.xcb_connection)
}

func (p *WsiPlatform) DispatchEvents(timeout int64) WsiResult {
	for {
		event := C.xcb_poll_for_event(p.xcb_connection)
		if event == nil {
			break
		}
		evtp := event.response_type & 0x7f
		//fmt.Println("DispatchEvents:", event.response_type)
		if p.xcb_im == nil || !C.xcb_xim_filter_event(p.xcb_im, event) {
			switch evtp {
			case C.XCB_CONFIGURE_NOTIFY:
				fmt.Println("DispatchEvents:", "XCB_CONFIGURE_NOTIFY", event)
				notify := (*C.xcb_configure_notify_event_t)(unsafe.Pointer(event))
				window := p.FindWindow(notify.window)
				if window != nil {
					wsi_window_xcb_configure_notify(window, notify)
				}
			case C.XCB_CLIENT_MESSAGE:
				message := (*C.xcb_client_message_event_t)(unsafe.Pointer(event))

				window := p.FindWindow(message.window)
				if window != nil {
					wsi_window_xcb_client_message(window, message)
				}
			case C.XCB_EXPOSE:
				fmt.Println("DispatchEvents:", "XCB_EXPOSE", event)
			case C.XCB_BUTTON_PRESS:
				fmt.Println("DispatchEvents:", "XCB_BUTTON_PRESS", event)
			case C.XCB_BUTTON_RELEASE:
				fmt.Println("DispatchEvents:", "XCB_BUTTON_RELEASE", event)
			case C.XCB_MAP_NOTIFY:
				fmt.Println("DispatchEvents:", "XCB_MAP_NOTIFY", event)
			case C.XCB_REPARENT_NOTIFY:
				fmt.Println("DispatchEvents:", "XCB_REPARENT_NOTIFY", event)
			case C.XCB_KEY_PRESS, C.XCB_KEY_RELEASE:
				fmt.Println("DispatchEvents:", "XCB_KEY_PRESS", event)
				// Forward event to input method if IC is created.
				if p.xcb_xic != 0 {
					C.xcb_xim_forward_event(p.xcb_im, p.xcb_xic, (*C.xcb_key_press_event_t)(unsafe.Pointer(event)))
				}
			case C.XCB_KEYMAP_NOTIFY:
				fmt.Println("DispatchEvents:", "XCB_KEYMAP_NOTIFY", event)
			default:
				fmt.Println("DispatchEvents:", event.response_type, event)
			}
		}

		C.free(unsafe.Pointer(event))
	}

	return WSI_SUCCESS
}
