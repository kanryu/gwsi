package gwsi

import "unsafe"

type WsiOutput uint64

type WsiResult int
type WsiEventType int
type WsiStructureType int32

const (
	WSI_SUCCESS             WsiResult = 0
	WSI_INCOMPLETE          WsiResult = 1
	WSI_TIMEOUT             WsiResult = 2
	WSI_EVENT_SET           WsiResult = 3
	WSI_EVENT_UNSET         WsiResult = 4
	WSI_SKIPPED             WsiResult = 5
	WSI_ERROR_UNKNOWN       WsiResult = -1
	WSI_ERROR_EGL           WsiResult = -2
	WSI_ERROR_VULKAN        WsiResult = -3
	WSI_ERROR_OUT_OF_RANGE  WsiResult = -4
	WSI_ERROR_OUT_OF_MEMORY WsiResult = -5
	WSI_ERROR_UNSUPPORTED   WsiResult = -6
	WSI_ERROR_DISCONNECTED  WsiResult = -7
	WSI_ERROR_PLATFORM      WsiResult = -8
	WSI_ERROR_PLATFORM_LOST WsiResult = -9
	WSI_ERROR_SEAT_LOST     WsiResult = -10
	WSI_ERROR_SEAT_IN_USE   WsiResult = -11
	WSI_ERROR_UNINITIALIZED WsiResult = -12
	WSI_ERROR_WINDOW_IN_USE WsiResult = -13
	WSI_RESULT_ENUM_MAX     WsiResult = WSI_TYPE_MAX

	WSI_EVENT_TYPE_CLOSE_WINDOW     WsiEventType = 1
	WSI_EVENT_TYPE_CONFIGURE_WINDOW WsiEventType = 2
	WSI_EVENT_TYPE_ENUM_MAX         WsiEventType = WSI_TYPE_MAX

	WSI_TYPE_MAX = 100
)

const (
	WSI_STRUCTURE_TYPE_PLATFORM_CREATE_INFO WsiStructureType = 0
	WSI_STRUCTURE_TYPE_ACQUIRE_SEAT_INFO    WsiStructureType = 1
	WSI_STRUCTURE_TYPE_WINDOW_CREATE_INFO   WsiStructureType = 2
	WSI_STRUCTURE_TYPE_ENUM_MAX             WsiStructureType = WSI_TYPE_MAX
)

type WsiPlatformCreateInfo struct {
	Type WsiStructureType
	Next string
}

type WsiWindowCreateInfo struct {
	Type            WsiStructureType
	Parent          *WsiWindow
	Extent          WsiExtent
	Title           string
	UserData        unsafe.Pointer
	ConfigureWindow PFN_wsiConfigureWindow
	CloseWindow     PFN_wsiCloseWindow
}

type WsiConfigureWindowEvent struct {
	Base   WsiEvent
	Window *WsiWindow
	Extent WsiExtent
}

type WsiCloseWindowEvent struct {
	Base   WsiEvent
	Window *WsiWindow
}

type PFN_wsiConfigureWindow func(pUserData unsafe.Pointer, pConfig *WsiConfigureWindowEvent)
type PFN_wsiCloseWindow func(pUserData unsafe.Pointer, pClose *WsiCloseWindowEvent)

type WsiExtent struct {
	Width  int32
	Height int32
}

type WsiEvent struct {
	Type   WsiEventType
	Flags  uint32
	Serial uint32
	Time   int64
}

type PFN_ConfigureWindowCallback func(unsafe.Pointer, *WsiConfigureWindowEvent)
type PFN_CloseWindowCallback func(unsafe.Pointer, *WsiCloseWindowEvent)
