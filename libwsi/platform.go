package libwsi

/*
#include <stdlib.h>
#include <string.h>

#include <xcb/xcb.h>

#include "utils.h"

#include "common_priv.h"
#include "platform_priv.h"
#include "window_priv.h"
*/
import "C"


type WsiPlatform struct{
     xcb_connection *xcb_connection_t
         xcb_screen *xcb_screen_t
                  xcb_screen_id int
           xcb_atom_wm_protocols xcb_atom_t
           xcb_atom_wm_delete_window xcb_atom_t
      window_list wsi_list
};

var (
	platform WsiPlatform
	protocolTable map[string]*xcb_atom_t
)

func wsi_query_extension(platform *WsiPlatform,  name string) bool{
    cookie := C.xcb_query_extension(
        platform.xcb_connection,
        strlen(name),
        name);

    reply := C.xcb_query_extension_reply(
        platform.xcb_connection,
        cookie,
        NULL);

    present := false;
    if reply {
        present = reply.present;
        free(reply);
    }

    return present;
}

func wsi_init_atoms(platform *WsiPlatform){
	protocolTable["WM_PROTOCOLS"] = &platform.xcb_atom_wm_protocols
	protocolTable["WM_DELETE_WINDOW"] = &platform.xcb_atom_wm_delete_window

    const size_t count = wsi_array_length(table);
    xcb_intern_atom_cookie_t cookies[count];

    for (size_t i = 0; i < count; ++i) {
        cookies[i] = xcb_intern_atom(
            platform.xcb_connection, 0, strlen(table[i].name), table[i].name);
    }

    for (size_t i = 0; i < count; ++i) {
        xcb_intern_atom_reply_t *reply = xcb_intern_atom_reply(
            platform.xcb_connection, cookies[i], NULL);

        if (reply) {
            *table[i].atom = reply.atom;
            free(reply);
        } else {
            *table[i].atom = XCB_ATOM_NONE;
        }
    }
}

func wsi_xcb_get_screen(setup *xcb_setup_t, int screen) *xcb_screen_t{
    xcb_screen_iterator_t iter = xcb_setup_roots_iterator(setup);
    for (; iter.rem; --screen, xcb_screen_next(&iter))
    {
        if (screen == 0) {
            return iter.data;
        }
    }

    return NULL;
}

func wsi_find_window(platform *WsiPlatform, window xcb_window_t) *WsiWindow{
    struct WsiWindow *WsiWindow = NULL;
    wsi_list_for_each(WsiWindow, &platform.window_list, link) {
        if (WsiWindow.xcb_window == window) {
            return WsiWindow;
        }
    }

    return NULL;
}

func wsi_platform_init(pCreateInfo *WsiPlatformCreateInfo, platform *WsiPlatform) WsiResult{
    WsiResult result;

    wsi_list_init(&platform.window_list);

    platform.xcb_connection = xcb_connect(NULL, &platform.xcb_screen_id);
    int err = xcb_connection_has_error(platform.xcb_connection);
    if (err > 0) {
        result = WSI_ERROR_PLATFORM;
        goto err_connect;
    }

    const xcb_setup_t *setup = xcb_get_setup(platform.xcb_connection);

    platform.xcb_screen = wsi_xcb_get_screen(setup, platform.xcb_screen_id);
    if (!platform.xcb_screen) {
        result = WSI_ERROR_PLATFORM;
        goto err_screen;
    }

    if (!wsi_query_extension(platform, "RANDR")) {
        result = WSI_ERROR_PLATFORM;
        goto err_extension;
    }

    if (!wsi_query_extension(platform, "XInputExtension")) {
        result = WSI_ERROR_PLATFORM;
        goto err_extension;
    }

    wsi_init_atoms(platform);

    return WSI_SUCCESS;

err_extension:
err_screen:
err_connect:
    xcb_disconnect(platform.xcb_connection);
    return result;
}

func wsi_platform_uninit(platform *WsiPlatform){
    xcb_disconnect(platform.xcb_connection);
}

func WsiCreatePlatform(pCreateInfo *WsiPlatformCreateInfo, pPlatform *WsiPlatform) WsiResult {
    struct WsiPlatform *p = calloc(1, sizeof(struct WsiPlatform));
    if (!p) {
        return WSI_ERROR_OUT_OF_MEMORY;
    }

    WsiResult result = wsi_platform_init(pCreateInfo, p);
    if (result != WSI_SUCCESS) {
        free(p);
        return result;
    }

    *pPlatform = p;
    return WSI_SUCCESS;
}

func WsiDestroyPlatform(platform WsiPlatform) {
    wsi_platform_uninit(platform);
    free(platform);
}

func WsiDispatchEvents(platform WsiPlatform, timeout int64) WsiResult {
    for {
        xcb_generic_event_t *event = xcb_poll_for_event(platform.xcb_connection);
        if (!event) {
            break;
        }

        switch (event.response_type & ~0x80) {
            case XCB_CONFIGURE_NOTIFY: {
                xcb_configure_notify_event_t *notify
                    = (xcb_configure_notify_event_t *)event;

                struct WsiWindow *window
                    = wsi_find_window(platform, notify.window);
                if (window) {
                    WsiWindow_xcb_configure_notify(window, notify);
                }
                break;
            }
            case XCB_CLIENT_MESSAGE: {
                xcb_client_message_event_t *message
                    = (xcb_client_message_event_t *)event;

                struct WsiWindow *window
                    = wsi_find_window(platform, message.window);
                if (window) {
                    WsiWindow_xcb_client_message(window, message);
                }
                break;
            }
        }

        free(event);
    }

    return WSI_SUCCESS;
}
