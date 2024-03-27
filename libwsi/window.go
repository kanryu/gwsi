package libwsi

/*
#include <stdlib.h>
#include <assert.h>
#include <string.h>

#include <xcb/xcb.h>

#include "utils.h"

#include "common_priv.h"
#include "platform_priv.h"
#include "window_priv.h"
*/
import "C"


type WsiWindow struct{
     platform *WsiPlatform
     link wsi_list
     xcb_window xcb_window_t
     xcb_parent xcb_window_t
     xcb_colormap xcb_colormap_t
     user_width uint16_t
     user_height uint16_t

    user_data uintptr
    pfn_configure PFN_wsiConfigureWindow
    pfn_close PFN_wsiCloseWindow
};

// region XCB Events

func WsiWindow_xcb_configure_notify(    window WsiWindow,    event *xcb_configure_notify_event_t) {
    assert(event.window == window.xcb_window);

    info := WsiConfigureWindowEvent{
        base.type = WSI_EVENT_TYPE_CONFIGURE_WINDOW,
        base.flags = 0,
        base.serial = event.sequence,
        window: window,
        extent.width = event.width,
        extent.height = event.height,
    };

    window.pfn_configure(window.user_data, &info);
}

func WsiWindow_xcb_client_message(    window *WsiWindow,    event *xcb_client_message_event_t){
    assert(event.window == window.xcb_window);

    if (event.type == window.platform.xcb_atom_wm_protocols &&
        event.data.data32[0] == window.platform.xcb_atom_wm_delete_window)
    {
        WsiCloseWindowEvent info = {
            .base.type = WSI_EVENT_TYPE_CLOSE_WINDOW,
            .base.flags = 0,
            .base.serial = event.sequence,
            .base.time = 0,
            .window = window,
        };

        window.pfn_close(window.user_data, &info);
    }
}

// endregion


func WsiCreateWindow(platform WsiPlatform, pCreateInfo *WsiWindowCreateInfo, pWindow *WsiWindow, title string) WsiResult {
    struct WsiWindow *window = calloc(1, sizeof(struct WsiWindow));
    if (!window) {
        return WSI_ERROR_OUT_OF_MEMORY;
    }

    window.platform = platform;
    window.xcb_window = xcb_generate_id(platform.xcb_connection);
    window.user_width = wsi_xcb_clamp(pCreateInfo.extent.width);
    window.user_height = wsi_xcb_clamp(pCreateInfo.extent.height);
    window.user_data = pCreateInfo.pUserData;
    window.pfn_configure = pCreateInfo.pfnConfigureWindow;
    window.pfn_close = pCreateInfo.pfnCloseWindow;

    if (pCreateInfo.parent) {
        window.xcb_parent = pCreateInfo.parent.xcb_window;
    } else {
        window.xcb_parent = platform.xcb_screen.root;
    }

    uint32_t value_list[2];
    uint32_t value_mask = XCB_CW_BACK_PIXEL | XCB_CW_EVENT_MASK ;
    value_list[0] = platform.xcb_screen.black_pixel;
    value_list[1] = XCB_EVENT_MASK_EXPOSURE
               // | XCB_EVENT_MASK_RESIZE_REDIRECT
                  | XCB_EVENT_MASK_STRUCTURE_NOTIFY
                  | XCB_EVENT_MASK_BUTTON_PRESS
                  | XCB_EVENT_MASK_BUTTON_RELEASE;

    xcb_create_window_checked(
        platform.xcb_connection,
        XCB_COPY_FROM_PARENT,
        window.xcb_window,
        window.xcb_parent,
        0, 0,
        window.user_width,
        window.user_height,
        10,
        XCB_WINDOW_CLASS_INPUT_OUTPUT,
        XCB_COPY_FROM_PARENT,
        value_mask,
        value_list);

    xcb_atom_t properties[] = {
        platform.xcb_atom_wm_protocols,
        platform.xcb_atom_wm_delete_window,
    };

    xcb_change_property(
        platform.xcb_connection,
        XCB_PROP_MODE_REPLACE,
        window.xcb_window,
        platform.xcb_atom_wm_protocols,
        XCB_ATOM_ATOM,
        32,
        wsi_array_length(properties),
        properties);

    xcb_map_window(platform.xcb_connection, window.xcb_window);
    xcb_flush(platform.xcb_connection);

    wsi_list_insert(&platform.window_list, &window.link);
    *pWindow = window;
    return WSI_SUCCESS;
}

func WsiDestroyWindow(window WsiWindow) {
    struct wsi_platform *platform = window.platform;

    xcb_unmap_window(platform.xcb_connection, window.xcb_window);
    xcb_destroy_window(platform.xcb_connection, window.xcb_window);
    xcb_flush(platform.xcb_connection);

    free(window);
}


func wsiSetWindowParent(window WsiWindow, parent WsiWindow) WsiResult{
    struct wsi_platform *platform = window.platform;
    if (parent) {
        window.xcb_parent = parent.xcb_window;
    } else {
        window.xcb_parent = platform.xcb_screen.root;
    }
    xcb_reparent_window(
        platform.xcb_connection,
        window.xcb_window,
        parent.xcb_window,
        0, 0);
    return WSI_SUCCESS;
}


func wsiSetWindowTitle(window WsiWindow, pTitle string) WsiResult{
    if (pTitle) {
        xcb_change_property(
            window.platform.xcb_connection,
            XCB_PROP_MODE_REPLACE,
            window.xcb_window,
            XCB_ATOM_WM_NAME,
            XCB_ATOM_STRING,
            8,
            strlen(pTitle),
            pTitle);
    } else {
        xcb_delete_property(
           window.platform.xcb_connection,
           window.xcb_window,
           XCB_ATOM_WM_NAME);
    }

    return WSI_SUCCESS;
}

func wsi_xcb_clamp(int32_t value) uint16_t{
    if (value < 0) {
        return 0;
    }
    if (value > UINT16_MAX) {
        return UINT16_MAX;
    }
    return (uint16_t)value;
}
