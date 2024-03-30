#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <string.h>

#include <xcb/xcb.h>
#include <xcb/xcb_keysyms.h>
#include <xcb-imdkit/encoding.h>
#include <xcb-imdkit/encoding.h>
#include <xcb-imdkit/ximproto.h>
#include <xcb-imdkit/imclient.h>
#include "_cgo_export.h"
xcb_window_t xcb_window;
extern void wsiSetXicCallback(xcb_xic_t new_ic, void *user_data);

void create_ic_callback(xcb_xim_t *im, xcb_xic_t new_ic, void *user_data) {
	puts("create_ic_callback");
	wsiSetXicCallback(new_ic, user_data);
    if (new_ic) {
        fprintf(stderr, "icid:%u\n", new_ic);
        xcb_xim_set_ic_focus(im, new_ic);
    }
}


void open_callback(xcb_xim_t *im, void *user_data) {
	puts("open_callback");
    uint32_t input_style = XCB_IM_PreeditPosition | XCB_IM_StatusArea;
    xcb_point_t spot;
    spot.x = 0;
    spot.y = 0;
    xcb_xim_nested_list nested =
        xcb_xim_create_nested_list(im, XCB_XIM_XNSpotLocation, &spot, NULL);
    xcb_xim_create_ic(im, create_ic_callback, user_data, XCB_XIM_XNInputStyle,
                      &input_style, XCB_XIM_XNClientWindow, &xcb_window,
                      XCB_XIM_XNFocusWindow, &xcb_window, XCB_XIM_XNPreeditAttributes,
                      &nested, NULL);
    free(nested.data);
}
bool gwsi_xcb_xim_open(xcb_xim_t *im,
	bool auto_connect, void *user_data)
{
	puts("gwsi_xcb_xim_open");
	return xcb_xim_open(im, open_callback, auto_connect, user_data);
}