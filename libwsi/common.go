package libwsi

type WsiOutput uint64

type WsiResult int
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

	WSI_EVENT_TYPE_CLOSE_WINDOW     = 1
	WSI_EVENT_TYPE_CONFIGURE_WINDOW = 2
	WSI_EVENT_TYPE_ENUM_MAX         = WSI_TYPE_MAX
	WSI_TYPE_MAX                    = 100
)

const (
	WSI_STRUCTURE_TYPE_PLATFORM_CREATE_INFO WsiStructureType = 0
	WSI_STRUCTURE_TYPE_ACQUIRE_SEAT_INFO    WsiStructureType = 1
	WSI_STRUCTURE_TYPE_WINDOW_CREATE_INFO   WsiStructureType = 2
	WSI_STRUCTURE_TYPE_ENUM_MAX             WsiStructureType = WSI_TYPE_MAX
)

// // egl.go

// type EGLDisplay uintptr
// type EGLConfig uintptr
// type EGLSurface uintptr

// type PFN_wsiGetEGLDisplay func(platform WsiPlatform, pDisplay *EGLDisplay) WsiResult
// type PFN_wsiCreateWindowEGLSurface func(window WsiWindow, dpy EGLDisplay, config EGLConfig, pSurface *EGLSurface) WsiResult
// type PFN_wsiDestroyWindowEGLSurface func(window WsiWindow, dpy EGLDisplay, surface EGLSurface)

// // func WsiGetEGLDisplay(platform WsiPlatform, pDisplay *EGLDisplay) WsiResult {
// // 	return WSI_SUCCESS
// // }

// // func WsiCreateWindowEGLSurface(window WsiWindow, dpy EGLDisplay, config EGLConfig, pSurface *EGLSurface) WsiResult {
// // 	return WSI_SUCCESS

// // }

// // func WsiDestroyWindowEGLSurface(window WsiWindow, dpy EGLDisplay, surface EGLSurface) {

// // }

// // input.go
// type WsiAcquireSeatInfo struct {
// 	Type  WsiStructureType
// 	Pnext uintptr
// 	ID    uint64
// }

// type PFN_wsiEnumerateSeats func(platform WsiPlatform, pIdCount *uint32, pIds *uint64) WsiResult
// type PFN_wsiAcquireSeat func(platform WsiPlatform, pCreateInfo *WsiAcquireSeatInfo, pSeat *WsiSeat) WsiResult
// type PFN_wsiReleaseSeat func(seat WsiSeat)

// func WsiEnumerateSeats(platform WsiPlatform, pIdCount *uint32, pIds *uint64) WsiResult {
// 	return WSI_SUCCESS
// }

// func WsiAcquireSeat(platform WsiPlatform, pAcquireInfo *WsiAcquireSeatInfo, pSeat *WsiSeat) WsiResult {
// 	return WSI_SUCCESS
// }

// func WsiReleaseSeat(seat WsiSeat) {

// }

// // output.go

// type PFN_wsiEnumerateOutputs func(platform WsiPlatform, pCount *uint32, pOutputs *WsiOutput) WsiResult

// func WsiEnumerateOutputs(platform WsiPlatform, pCount *uint32, pOutputs *WsiOutput) WsiResult {
// 	return WSI_SUCCESS
// }

// // platform.go

// type WsiPlatformCreateInfo struct {
// 	Type  WsiStructureType
// 	pNext uintptr
// }

// // type PFN_wsiCreatePlatform func(pCreateInfo WsiPlatformCreateInfo, pPlatform *WsiPlatform) WsiResult
// type PFN_wsiDestroyPlatform func(platform WsiPlatform)
// type PFN_wsiDispatchEvents func(platform WsiPlatform, timeout int64) WsiResult

// // func WsiCreatePlatform(pCreateInfo WsiPlatformCreateInfo, pPlatform *WsiPlatform) WsiResult {
// // 	return WSI_SUCCESS
// // }

// func WsiDestroyPlatform(platform WsiPlatform) {

// }

// func WsiDispatchEvents(platform WsiPlatform, timeout int64) WsiResult {
// 	return WSI_SUCCESS
// }

// // vulkan.go

// type VkInstance uintptr
// type VkPhysicalDevice uintptr
// type VkSurfaceKHR uintptr
// type VkBool32 int32
// type size_t int32
// type VkSystemAllocationScope int
// type VkInternalAllocationType int

// type PFN_wsiEnumerateRequiredInstanceExtensions func(platform WsiPlatform, pExtensionCount *uint32, ppExtensions *string) WsiResult
// type PFN_wsiEnumerateRequiredDeviceExtensions func(platform WsiPlatform, pExtensionCount *uint32, ppExtensions *string) WsiResult
// type PFN_wsiCreateWindowSurface func(window WsiWindow, instance VkInstance, pAllocator VkAllocationCallbacks, pSurface VkSurfaceKHR) WsiResult
// type PFN_wsiGetPhysicalDevicePresentationSupport func(platform WsiPlatform, physicalDevice VkPhysicalDevice, queueFamilyIndex uint32) VkBool32

// type PFN_vkAllocationFunction func(
// 	pUserData uintptr,
// 	size size_t,
// 	alignment size_t,
// 	allocationScope VkSystemAllocationScope,
// )

// type PFN_vkReallocationFunction func(
// 	pUserData uintptr,
// 	pOriginal uintptr,
// 	size size_t,
// 	alignment size_t,
// 	allocationScope VkSystemAllocationScope,
// )

// type PFN_vkFreeFunction func(
// 	pUserData uintptr,
// 	pMemory uintptr,
// )

// type PFN_vkInternalAllocationNotification func(
// 	pUserData uintptr,
// 	size size_t,
// 	allocationType VkInternalAllocationType,
// 	allocationScope VkSystemAllocationScope,
// )

// type PFN_vkInternalFreeNotification func(
// 	pUserData uintptr,
// 	size size_t,
// 	allocationType VkInternalAllocationType,
// 	allocationScope VkSystemAllocationScope,
// )
// type VkAllocationCallbacks struct {
// 	pUserData             uintptr
// 	pfnAllocation         PFN_vkAllocationFunction
// 	pfnReallocation       PFN_vkReallocationFunction
// 	pfnFree               PFN_vkFreeFunction
// 	pfnInternalAllocation PFN_vkInternalAllocationNotification
// 	pfnInternalFree       PFN_vkInternalFreeNotification
// }

// func WsiEnumerateRequiredInstanceExtensions(
// 	platform WsiPlatform,
// 	pExtensionCount *uint32,
// 	ppExtensions *string) WsiResult {
// 	return WSI_SUCCESS
// }

// func WsiEnumerateRequiredDeviceExtensions(
// 	platform WsiPlatform,
// 	pExtensionCount *uint32,
// 	ppExtensions *string) WsiResult {
// 	return WSI_SUCCESS
// }

// func WsiCreateWindowSurface(
// 	window WsiWindow,
// 	instance VkInstance,
// 	pAllocator *VkAllocationCallbacks,
// 	pSurface VkSurfaceKHR,
// ) WsiResult {
// 	return WSI_SUCCESS

// }

// func WsiGetPhysicalDevicePresentationSupport(
// 	platform WsiPlatform,
// 	physicalDevice VkPhysicalDevice,
// 	queueFamilyIndex uint32,
// ) VkBool32 {
// 	return 1
// }

// // window.go

// type WsiConfigureWindowEvent struct {
// 	base   WsiEvent
// 	window WsiWindow
// 	extent WsiExtent
// }

// type WsiCloseWindowEvent struct {
// 	base   WsiEvent
// 	window WsiWindow
// }

// type PFN_wsiConfigureWindow func(pUserData uintptr, pConfig *WsiConfigureWindowEvent)
// type PFN_wsiCloseWindow func(pUserData uintptr, pClose *WsiCloseWindowEvent)

// type WsiWindowCreateInfo struct {
// 	sType              WsiStructureType
// 	pNext              uintptr
// 	parent             WsiWindow
// 	extent             WsiExtent
// 	pTitle             string
// 	pUserData          uintptr
// 	pfnConfigureWindow PFN_wsiConfigureWindow
// 	pfnCloseWindow     PFN_wsiCloseWindow
// }

// type PFN_wsiCreateWindow func(platform WsiPlatform, pConfig *WsiConfigureWindowEvent, pWindow *WsiWindow) WsiResult
// type PFN_wsiDestroyWindow func(window WsiWindow)
// type PFN_wsiSetWindowParent func(window WsiWindow, parent WsiWindow) WsiResult
// type PFN_wsiSetWindowTitle func(window WsiWindow, title string) WsiResult

// func WsiCreateWindow(platform WsiPlatform, pConfig *WsiConfigureWindowEvent, pWindow *WsiWindow) WsiResult {
// 	return WSI_SUCCESS
// }

// func WsiDestroyWindow(window WsiWindow) {

// }

// func WsiSetWindowParent(window WsiWindow, parent WsiWindow) WsiResult {
// 	return WSI_SUCCESS
// }

// func WsiSetWindowTitle(window WsiWindow, title string) WsiResult {
// 	return WSI_SUCCESS
// }
