// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || freebsd || openbsd
// +build linux freebsd openbsd

package xcb

import "C"
import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

type (
	XcbAtom       uint32
	XcbButton     uint8
	XcbKeycode    uint8
	XcbScreenId   int32
	XcbTimestamp  uint32
	XcbVoidCookie uint32
	XcbWindow     uint32

	PtrXcbConnection           uintptr
	PtrXcbKeyPressEvent        uintptr
	PtrXcbKeySymbols           uintptr
	PtrXcbScreen               uintptr
	PtrXcbSetup                uintptr
	PtrXcbXim                  uintptr
	PtrXcbAtom                 uintptr
	PtrXcbClientMessageEvent   uintptr
	PtrXcbConfigureNotifyEvent uintptr
)

type XcbGenericEventT struct {
	ResponseType uint8     /* Type of the response */
	pad0         uint8     /* Padding */
	Sequence     uint16    /* Sequence number */
	pad          [7]uint32 /* Padding */
	FullSequence uint32    /* full sequence */
}

type XcbConfigureNotifyEventT struct {
	ResponseType     uint8
	pad0             uint8
	Sequence         uint16
	Event            XcbWindow
	Window           XcbWindow
	AboveSibling     XcbWindow
	X                int16
	Y                int16
	Width            uint16
	Height           uint16
	BorderWidth      uint16
	OverrideRedirect uint8
	pad1             uint8
}

type XcbClientMessageEventT struct {
	ResponseType uint8
	Format       uint8
	Sequence     uint16
	Window       XcbWindow
	Type         XcbAtom
	Data         XcbClientMessageData
}

type XcbExposeEventT struct {
	ResponseType uint8
	pad0         uint8
	Sequence     uint16
	Window       XcbWindow
	X            uint16
	Y            uint16
	Width        uint16
	Height       uint16
	Count        uint16
	pad1         [2]uint8
}

type XcbButtonPressEventT struct {
	ResponseType uint8 /* The type of the event, here it is xcb_button_press_event_t or xcb_button_release_event_t */
	Detail       XcbButton
	Sequence     uint16
	Time         XcbTimestamp /* Time, in milliseconds the event took place in */
	Root         XcbWindow
	Event        XcbWindow
	Child        XcbWindow
	RootX        int16
	RootY        int16
	EventX       int16  /* The x coordinate where the mouse has been pressed in the window */
	EventY       int16  /* The y coordinate where the mouse has been pressed in the window */
	State        uint16 /* A mask of the buttons (or keys) during the event */
	SameScreen   uint8
}
type XcbMapNotifyEventT struct {
	ResponseType     uint8
	pad0             uint8
	Sequence         uint16
	Event            XcbWindow
	Window           XcbWindow
	OverrideRedirect uint8
	pad1             [3]uint8
}
type XcbReparentNotifyEventT struct {
	ResponseType     uint8 /* The type of the event, here it is xcb_button_press_event_t or xcb_button_release_event_t */
	pad0             uint8
	Sequence         uint16
	Event            XcbWindow
	Window           XcbWindow
	Parent           XcbWindow
	X                int16
	Y                int16
	OverrideRedirect uint8
	pad1             [3]uint8
}
type XcbKeyPressEventT struct {
	ResponseType uint8 /* The type of the event, here it is xcb_button_press_event_t or xcb_button_release_event_t */
	Detail       XcbKeycode
	Sequence     uint16
	Time         XcbTimestamp /* Time, in milliseconds the event took place in */
	Root         XcbWindow
	Event        XcbWindow
	Child        XcbWindow
	RootX        int16
	RootY        int16
	EventX       int16  /* The x coordinate where the mouse has been pressed in the window */
	EventY       int16  /* The y coordinate where the mouse has been pressed in the window */
	State        uint16 /* A mask of the buttons (or keys) during the event */
	SameScreen   uint8
	pad0         uint8
}
type XcbKeymapNotifyEventT struct {
	ResponseType uint8 /* The type of the event, here it is xcb_button_press_event_t or xcb_button_release_event_t */
	Keys         [31]uint8
}

type XcbClientMessageData struct {
	Data8 [20]byte
}

func (e *XcbGenericEventT) AsXcbConfigureNotifyEventT() *XcbConfigureNotifyEventT {
	return (*XcbConfigureNotifyEventT)(unsafe.Pointer(e))
}

func (e *XcbGenericEventT) AsXcbClientMessageEventT() *XcbClientMessageEventT {
	return (*XcbClientMessageEventT)(unsafe.Pointer(e))
}

func (e *XcbGenericEventT) AsXcbExposeEventT() *XcbExposeEventT {
	return (*XcbExposeEventT)(unsafe.Pointer(e))
}

func (e *XcbGenericEventT) AsXcbButtonPressEventT() *XcbButtonPressEventT {
	return (*XcbButtonPressEventT)(unsafe.Pointer(e))
}

func (e *XcbGenericEventT) AsXcbMapNotifyEventT() *XcbMapNotifyEventT {
	return (*XcbMapNotifyEventT)(unsafe.Pointer(e))
}

func (e *XcbGenericEventT) AsXcbReparentNotifyEventT() *XcbReparentNotifyEventT {
	return (*XcbReparentNotifyEventT)(unsafe.Pointer(e))
}

func (e *XcbGenericEventT) AsXcbKeyPressEventT() *XcbKeyPressEventT {
	return (*XcbKeyPressEventT)(unsafe.Pointer(e))
}

func (e *XcbGenericEventT) AsXcbKeymapNotifyEventT() *XcbKeymapNotifyEventT {
	return (*XcbKeymapNotifyEventT)(unsafe.Pointer(e))
}

// XcbAtom
// XcbColormap
// XcbParent
// XcbWindow
// XcbXic
// XcbXimImCallbac

var (
	// xcb_change_property func(xcb_connection_t *conn, uint8_t mode, xcb_window_t window, xcb_atom_t property, xcb_atom_t atom, uint8_t format, uint32_t data_len, data uintptr) XcbVoidCookie
	// xcb_change_window_attributes_checked func(xcb_connection_t *conn, xcb_window_t window, uint32_t value_mask, uint32_t *values) XcbVoidCookie
	// xcb_change_window_attributes func(xcb_connection_t *conn, xcb_window_t window, uint32_t value_mask, const void *value_list) XcbVoidCookie
	// xcb_connect func(displayname string, int *screenp *int32) PtrXcbConnection
	// xcb_connection_has_error func(xcb_connection_t *c) int32
	// xcb_create_colormap_checked
	// xcb_create_window_checked
	// xcb_delete_property
	// xcb_destroy_window xcb_void_cookie_t xcb_destroy_window 	( 	xcb_connection_t *  	c,
	// 		xcb_window_t  	window
	// 	) XcbVoidCookie
	// xcb_disconnect func(xcb_connection_t *c)
	// xcb_flush
	// xcb_free_colormap
	// xcb_generate_id
	// xcb_get_setup
	// xcb_intern_atom_reply
	// xcb_intern_atom
	// xcb_map_window
	// xcb_poll_for_event
	// xcb_query_extension_reply
	// xcb_query_extension
	// xcb_reparent_window
	// xcb_screen_next
	// xcb_setup_roots_iterator
	// xcb_unmap_window
	// xcb_visualid_t

	// xcb_compound_text_init
	// xcb_xim_create
	// xcb_xim_filter_event
	// xcb_xim_forward_event
	// xcb_xim_set_im_callback
	// xcb_xim_set_use_compound_text
	// xcb_xim_set_use_utf8_string

	// XcbChangeProperty
	// XcbChangeWindowAttributesChecked
	// XcbChangeWindowAttributes
	// XcbCompoundTextInit
	// XcbConnect
	// XcbConnectionHasError
	// XcbCreateColormapChecked
	// XcbCreateWindowChecked
	// XcbDeleteProperty
	// XcbDestroyWindow
	// XcbDisconnect
	// XcbFlush
	// XcbFreeColormap
	// XcbGenerateId func (conn PtrXcbConnection) uint32
	// XcbGetSetup
	// XcbInternAtomReply
	// XcbInternAtom
	// XcbMapWindow
	XcbPollForEvent func(conn PtrXcbConnection) *XcbGenericEventT

// XcbQueryExtensionReply
// XcbQueryExtension
// XcbReparentWindow
// XcbScreenNext
// XcbSetupRootsIterator
// XcbUnmapWindow func (conn PtrXcbConnection, window XcbWindow) XcbVoidCookie
// XcbVisualidT
// XcbXimCreate
// XcbXimFilterEvent
// XcbXimForwardEvent
// XcbXimSetImCallback
// XcbXimSetUseCompoundText
// XcbXimSetUseUtf8String
)

func LoadXcb() error {
	lib, err := purego.Dlopen(getLibXcb(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	// puts("Calling C from Go without Cgo!")
	purego.RegisterLibFunc(&XcbPollForEvent, lib, "xcb_poll_for_event")
	return nil
}

func getLibXcb() string {
	switch runtime.GOOS {
	case "linux":
		return "libxcb.so"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}
