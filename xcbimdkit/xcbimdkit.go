// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || freebsd || openbsd
// +build linux freebsd openbsd

package xcbimdkit

import (
	"fmt"
	"gwsi/xcb"
	"log"
	"runtime"
	"unsafe"

	"github.com/CannibalVox/cgoalloc"
	"github.com/ebitengine/purego"
)

type (
	xcb_keycode_t              uint8
	XcbXicT                    uint16
	xcb_window_t               uint32
	xcb_timestamp_t            uint32
	PtrXcbXim                  uintptr
	xcb_xim_encoding_t         uint32
	size_t                     uint32
	PtrXcbConnectionT          uintptr
	xcb_xim_create_ic_callback uintptr
	CreateIcCallbackType       func(im PtrXcbXim, new_ic XcbXicT, user_data uintptr)
	XcbXimOpenCallback         func(im PtrXcbXim, user_data uintptr)

	EGLenum = uint32
)

type xcb_key_press_event_t struct {
	response_type uint8
	detail        xcb_keycode_t
	sequence      uint16
	time          xcb_timestamp_t
	root          xcb_window_t
	event         xcb_window_t
	child         xcb_window_t
	root_x        int16
	root_y        int16
	event_x       int16
	event_y       int16
	state         uint16
	same_screen   uint8
	pad0          uint8
}

type xcb_point_t struct {
	X int16
	Y int16
}

type xcb_xim_nested_list struct {
	Data   uintptr
	Length size_t
}

type XcbXimImCallback struct {
	set_event_mask     uintptr
	forward_event      uintptr
	commit_string      uintptr
	geometry           uintptr
	preedit_start      uintptr
	preedit_draw       uintptr
	preedit_caret      uintptr
	preedit_done       uintptr
	status_start       uintptr
	status_draw_text   uintptr
	status_draw_bitmap uintptr
	status_done        uintptr
	sync               uintptr
	disconnected       uintptr
}

const (
	XCB_KEY_PRESS         uint8 = 2
	XCB_XIM_COMPOUND_TEXT       = 0
	XCB_XIM_UTF8_STRING         = 1

	XCB_IM_PreeditArea      uint32 = 0x0001
	XCB_IM_PreeditCallbacks uint32 = 0x0002
	XCB_IM_PreeditPosition  uint32 = 0x0004
	XCB_IM_PreeditNothing   uint32 = 0x0008
	XCB_IM_PreeditNone      uint32 = 0x0010
	XCB_IM_StatusArea       uint32 = 0x0100
	XCB_IM_StatusCallbacks  uint32 = 0x0200
	XCB_IM_StatusNothing    uint32 = 0x0400
	XCB_IM_StatusNone       uint32 = 0x0800

	XCB_XIM_XNQueryInputStyle       = "queryInputStyle"
	XCB_XIM_XNClientWindow          = "clientWindow"
	XCB_XIM_XNInputStyle            = "inputStyle"
	XCB_XIM_XNFocusWindow           = "focusWindow"
	XCB_XIM_XNFilterEvents          = "filterEvents"
	XCB_XIM_XNPreeditAttributes     = "preeditAttributes"
	XCB_XIM_XNStatusAttributes      = "statusAttributes"
	XCB_XIM_XNArea                  = "area"
	XCB_XIM_XNAreaNeeded            = "areaNeeded"
	XCB_XIM_XNSpotLocation          = "spotLocation"
	XCB_XIM_XNColormap              = "colorMap"
	XCB_XIM_XNStdColormap           = "stdColorMap"
	XCB_XIM_XNForeground            = "foreground"
	XCB_XIM_XNBackground            = "background"
	XCB_XIM_XNBackgroundPixmap      = "backgroundPixmap"
	XCB_XIM_XNFontSet               = "fontSet"
	XCB_XIM_XNLineSpace             = "lineSpace"
	XCB_XIM_XNSeparatorofNestedList = "separatorofNestedList"
)

var (
	LibLoaded = false

	XcbXimCreate           func(conn PtrXcbConnectionT, screen_id int32, imname *byte) PtrXcbXim
	XcbXimOpen             func(im PtrXcbXim, callback XcbXimOpenCallback, auto_connect bool, user_data uintptr) bool
	XcbXimSetImCallback    func(im PtrXcbXim, callbacks *XcbXimImCallback, user_data uintptr)
	XcbXimSetIcFocus       func(im PtrXcbXim, ic XcbXicT) bool
	XcbXimCreateNestedList func(im PtrXcbXim, style string, xcb_point uintptr, end *bool) xcb_xim_nested_list
	XcbXimCreateIc         func(im PtrXcbXim, callback uintptr, user_data uintptr,
		style1 string, input_style *uint32,
		style2 string, client_window *xcb_window_t,
		style3 string, focus_window *xcb_window_t,
		style4 string, nested *xcb_xim_nested_list,
		end *bool) bool

	XcbCompoundTextInit   func()
	XcbUtf8ToCompoundText func(utf8 string, length size_t, lenghtOut *size_t) string
	XcbCompoundTextToUtf8 func(compound_text string, length size_t, lenghtOut *size_t) string

	XcbXimSetUseCompoundText func(im PtrXcbXim, enable bool)
	XcbXimSetUseUtf8String   func(im PtrXcbXim, enable bool)
	XcbXimForwardEvent       func(im PtrXcbXim, ic XcbXicT, event *xcb.XcbGenericEventT) bool
	XcbXimFilterEvent        func(im PtrXcbXim, event *xcb.XcbGenericEventT) bool
	XcbXimGetEncoding        func(im PtrXcbXim) xcb_xim_encoding_t

	CreateIcCallback CreateIcCallbackType
)

func RegisterCreateIcCallback(callback CreateIcCallbackType) {
	CreateIcCallback = callback
}

func createIcCallback(im PtrXcbXim, new_ic XcbXicT, user_data uintptr) {
	fmt.Print("create_ic_callback")
	if CreateIcCallback == nil {
		log.Fatal(fmt.Errorf("CreateIcCallback must be registered"))
	}
	CreateIcCallback(im, new_ic, user_data)
	if new_ic > 0 {
		fmt.Printf("icid:%u\n", new_ic)
		XcbXimSetIcFocus(im, new_ic)
	}
}

// Checks whether the key event is trapped by IME.
// Requires libxcb-imdkit.so to work, and returns false if it does not exist.
func PluginXcbXimFilterEvent(im PtrXcbXim, event *xcb.XcbGenericEventT) bool {
	if !LibLoaded {
		return false
	}
	return XcbXimFilterEvent(im, event)
}

func OpenCallback(im PtrXcbXim, user_data uintptr) {
	fmt.Print("open_callback")
	input_style := XCB_IM_PreeditPosition | XCB_IM_StatusArea
	var xcb_window xcb_window_t
	spot := xcb_point_t{X: 0, Y: 0}

	nested := XcbXimCreateNestedList(im, XCB_XIM_XNSpotLocation, uintptr(unsafe.Pointer(&spot)), nil)

	XcbXimCreateIc(im, purego.NewCallback(createIcCallback), user_data,
		XCB_XIM_XNInputStyle, &input_style,
		XCB_XIM_XNClientWindow, &xcb_window,
		XCB_XIM_XNFocusWindow, &xcb_window,
		XCB_XIM_XNPreeditAttributes, &nested,
		nil)

	allocator := cgoalloc.DefaultAllocator{}
	allocator.Free(unsafe.Pointer(nested.Data))
}

var ImCallback XcbXimImCallback = XcbXimImCallback{
	forward_event: purego.NewCallback(callback_forward_event),
	commit_string: purego.NewCallback(func(im PtrXcbXim, ic XcbXicT, flag uint32, str *byte,
		length size_t, keysym *uint32, nKeySym size_t, user_data uintptr) {
		sstr := unsafe.String(str, length)
		callback_commit_string(im, ic, flag, sstr, length, keysym, nKeySym, user_data)
	}),
	disconnected: purego.NewCallback(callback_disconnected),
}

func callback_forward_event(im PtrXcbXim, ic XcbXicT, event *xcb_key_press_event_t,
	user_data uintptr) {
	key := "release"
	if event.response_type == XCB_KEY_PRESS {
		key = "press"
	}
	fmt.Printf("Key %s Keycode %u, State %u\n",
		key,
		event.detail, event.state)
}

func callback_commit_string(im PtrXcbXim, ic XcbXicT, flag uint32, str string,
	length size_t, keysym *uint32, nKeySym size_t, user_data uintptr) {
	if XcbXimGetEncoding(im) == XCB_XIM_UTF8_STRING {
		fmt.Printf("key commit utf8: %.*s\n", length, str)
	} else if XcbXimGetEncoding(im) == XCB_XIM_COMPOUND_TEXT {
		var newLength size_t
		utf8 := XcbCompoundTextToUtf8(str, length, &newLength)
		if utf8 != "" {
			fmt.Printf("key commit: %.*s\n", newLength, utf8)
		}
	}
}

func callback_disconnected(im PtrXcbXim, user_data uintptr) {
	fmt.Printf("Disconnected from input method server.\n")
	//ic = 0;
}

func LoadEGLXcbImdkit() error {
	if LibLoaded {
		return nil
	}
	var lib uintptr
	var err error
	for _, libname := range getLibXcbImdkit() {
		lib, err = purego.Dlopen(libname, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if lib != 0 {
			break
		}
	}
	if err != nil {
		//panic(err)
		fmt.Printf("libxcb-imdkit is required to support IME, but it cannot be found\nError: %s\n", err.Error())
		return err
	}
	purego.RegisterLibFunc(&XcbXimCreate, lib, "xcb_xim_create")
	purego.RegisterLibFunc(&XcbXimOpen, lib, "xcb_xim_open")

	purego.RegisterLibFunc(&XcbCompoundTextInit, lib, "xcb_compound_text_init")
	purego.RegisterLibFunc(&XcbUtf8ToCompoundText, lib, "xcb_utf8_to_compound_text")
	purego.RegisterLibFunc(&XcbCompoundTextToUtf8, lib, "xcb_compound_text_to_utf8")

	purego.RegisterLibFunc(&XcbXimSetImCallback, lib, "xcb_xim_set_im_callback")
	purego.RegisterLibFunc(&XcbXimSetIcFocus, lib, "xcb_xim_set_ic_focus")
	purego.RegisterLibFunc(&XcbXimCreateNestedList, lib, "xcb_xim_create_nested_list")
	purego.RegisterLibFunc(&XcbXimCreateIc, lib, "xcb_xim_create_ic")

	purego.RegisterLibFunc(&XcbXimSetUseCompoundText, lib, "xcb_xim_set_use_compound_text")
	purego.RegisterLibFunc(&XcbXimSetUseUtf8String, lib, "xcb_xim_set_use_utf8_string")
	purego.RegisterLibFunc(&XcbXimForwardEvent, lib, "xcb_xim_forward_event")
	purego.RegisterLibFunc(&XcbXimFilterEvent, lib, "xcb_xim_filter_event")
	purego.RegisterLibFunc(&XcbXimGetEncoding, lib, "xcb_xim_get_encoding")

	LibLoaded = true
	return nil
}

func getLibXcbImdkit() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"/usr/lib/libSystem.B.dylib"}
	case "linux":
		return []string{"libxcb-imdkit.so.1", "libxcb-imdkit.so"}
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}
