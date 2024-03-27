#ifndef WSI_INCLUDE_WINDOW_H
#define WSI_INCLUDE_WINDOW_H

#include "common.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct WsiConfigureWindowEvent {
    WsiEvent Base;
    WsiWindow Window;
    WsiExtent Extent;
} WsiConfigureWindowEvent;

typedef struct WsiCloseWindowEvent {
    WsiEvent Base;
    WsiWindow Window;
} WsiCloseWindowEvent;

typedef void (*PFN_wsiConfigureWindow)(void *pUserData, const WsiConfigureWindowEvent *pConfig);
typedef void (*PFN_wsiCloseWindow)(void *pUserData, const WsiCloseWindowEvent *pClose);

typedef struct WsiWindowCreateInfo {
    int32_t SType;
    const void *pNext;
    WsiWindow Parent;
    WsiExtent Extent;
    const char *PTitle;
    void *PUserData;
    PFN_wsiConfigureWindow PfnConfigureWindow;
    PFN_wsiCloseWindow PfnCloseWindow;
} WsiWindowCreateInfo;

typedef WsiResult (*PFN_wsiCreateWindow)(WsiPlatform platform, const WsiWindowCreateInfo *pCreateInfo, WsiWindow *pWindow);
typedef void (*PFN_wsiDestroyWindow)(WsiWindow window);
typedef WsiResult (*PFN_wsiSetWindowParent)(WsiWindow window, WsiWindow parent);
typedef WsiResult (*PFN_wsiSetWindowTitle)(WsiWindow window, const char *pTitle);

WsiResult
wsiCreateWindow(WsiPlatform platform, const WsiWindowCreateInfo *pCreateInfo, WsiWindow *pWindow);

void
wsiDestroyWindow(WsiWindow window);

WsiResult
wsiSetWindowParent(WsiWindow window, WsiWindow parent);

WsiResult
wsiSetWindowTitle(WsiWindow window, const char *pTitle);

#ifdef __cplusplus
}
#endif

#endif
