package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"

	"gwsi"
	"gwsi/egl"
	"gwsi/xcb"
	"gwsi/xcbimdkit"
)

const (
	TAU_F = 6.28318530717958647692
)

var (
	g_platform  *gwsi.WsiPlatform
	g_window    *gwsi.WsiWindow
	g_extent    gwsi.WsiExtent
	g_running   bool = true
	g_resized   bool = false
	g_display   egl.EGLDisplay
	g_config    egl.EGLConfig
	g_surface   egl.EGLSurface
	g_context   egl.EGLContext
	g_view_rotx float64 = 20.0
	g_view_roty float64 = 30.0
	g_view_rotz float64 = 0.0
	g_gear1     uint32
	g_gear2     uint32
	g_gear3     uint32
	g_angle     float64 = 0.0

	g_config_attribs = []egl.EGLint{
		egl.EGL_SURFACE_TYPE, egl.EGL_WINDOW_BIT,
		egl.EGL_RED_SIZE, 8,
		egl.EGL_GREEN_SIZE, 8,
		egl.EGL_BLUE_SIZE, 8,
		egl.EGL_ALPHA_SIZE, 8,
		egl.EGL_DEPTH_SIZE, 24,
		egl.EGL_RENDERABLE_TYPE, egl.EGL_OPENGL_BIT,
		egl.EGL_NONE,
	}
	g_context_attribs = []egl.EGLint{
		egl.EGL_CONTEXT_MAJOR_VERSION_KHR, 2,
		egl.EGL_CONTEXT_MINOR_VERSION_KHR, 0,
		egl.EGL_NONE,
	}
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func calc_angle(i int, teeth int) float64 {
	return float64(i) * TAU_F / float64(teeth)
}

func gear(
	inner_radius float64,
	outer_radius float64,
	width float64,
	teeth int,
	tooth_depth float64,
) {

	r0 := inner_radius
	r1 := outer_radius - tooth_depth/2.0
	r2 := outer_radius + tooth_depth/2.0

	da := TAU_F / float64(teeth) / 4.0

	gl.ShadeModel(gl.FLAT)
	gl.Normal3f(0.0, 0.0, 1.0)

	gl.Begin(gl.QUAD_STRIP)
	for i := 0; i <= teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(width*0.5))
		if i < teeth {
			gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(width*0.5))
			gl.Vertex3f(float32(r1*math.Cos(angle+3.0*da)), float32(r1*math.Sin(angle+3.0*da)), float32(width*0.5))
		}
	}
	gl.End()

	gl.Begin(gl.QUADS)
	for i := 0; i < teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+da)), float32(r2*math.Sin(angle+da)), float32(width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+2.0*da)), float32(r2*math.Sin(angle+2.0*da)), float32(width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle+3.0*da)), float32(r1*math.Sin(angle+3.0*da)), float32(width*0.5))
	}
	gl.End()

	gl.Normal3f(0.0, 0.0, -1.0)

	gl.Begin(gl.QUAD_STRIP)
	for i := 0; i < teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(-width*0.5))
		gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(-width*0.5))
		if i < teeth {
			gl.Vertex3f(float32(r1*math.Cos(angle+3.0*da)), float32(r1*math.Sin(angle+3.0*da)), float32(-width*0.5))
			gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(-width*0.5))
		}
	}
	gl.End()

	gl.Begin(gl.QUADS)
	for i := 0; i < teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(float32(r1*math.Cos(angle+3.0*da)), float32(r1*math.Sin(angle+3.0*da)), float32(-width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+2.0*da)), float32(r2*math.Sin(angle+2.0*da)), float32(-width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+da)), float32(r2*math.Sin(angle+da)), float32(-width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(-width*0.5))
	}
	gl.End()

	gl.Begin(gl.QUAD_STRIP)
	for i := 0; i < teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(-width*0.5))
		u := r2*math.Cos(angle+da) - r1*math.Cos(angle)
		v := r2*math.Sin(angle+da) - r1*math.Sin(angle)
		leng := math.Sqrt(u*u + v*v)
		u /= leng
		v /= leng
		gl.Normal3f(float32(v), float32(-u), 0.0)
		gl.Vertex3f(float32(r2*math.Cos(angle+da)), float32(r2*math.Sin(angle+da)), float32(width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+da)), float32(r2*math.Sin(angle+da)), float32(-width*0.5))
		gl.Normal3f(float32(math.Cos(angle)), float32(math.Sin(angle)), float32(0.0))
		gl.Vertex3f(float32(r2*math.Cos(angle+2.0*da)), float32(r2*math.Sin(angle+2.0*da)), float32(width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+2.0*da)), float32(r2*math.Sin(angle+2.0*da)), float32(-width*0.5))
		u = r1*math.Cos(angle+3.0*da) - r2*math.Cos(angle+2.0*da)
		v = r1*math.Sin(angle+3.0*da) - r2*math.Sin(angle+2.0*da)
		gl.Normal3f(float32(v), float32(-u), 0.0)
		gl.Vertex3f(float32(r1*math.Cos(angle+3.0*da)), float32(r1*math.Sin(angle+3.0*da)), float32(width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle+3.0*da)), float32(r1*math.Sin(angle+3.0*da)), float32(-width*0.5))
		gl.Normal3f(float32(math.Cos(angle)), float32(math.Sin(angle)), float32(0.0))
	}

	gl.Vertex3f(float32(r1*math.Cos(0)), float32(r1*math.Sin(0)), float32(width*0.5))
	gl.Vertex3f(float32(r1*math.Cos(0)), float32(r1*math.Sin(0)), float32(-width*0.5))

	gl.End()

	gl.ShadeModel(gl.SMOOTH)

	gl.Begin(gl.QUAD_STRIP)
	for i := 0; i <= teeth; i++ {
		angle := calc_angle(i, teeth)
		gl.Normal3f(float32(-math.Cos(angle)), float32(-math.Sin(angle)), float32(0.0))
		gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(-width*0.5))
		gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(width*0.5))
	}
	gl.End()
}

func draw() {
	if g_resized {
		gl.Viewport(0, 0, int32(g_extent.Width), int32(g_extent.Height))

		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()

		hf := float64(g_extent.Height)
		wf := float64(g_extent.Width)

		if hf > wf {
			aspect := hf / wf
			gl.Frustum(-1.0, 1.0, -aspect, aspect, 5.0, 60.0)
		} else {
			aspect := wf / hf
			gl.Frustum(-aspect, aspect, -1.0, 1.0, 5.0, 60.0)
		}

		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Translatef(0.0, 0.0, -40.0)
		g_resized = false
	}

	gl.ClearColor(0.0, 0.0, 0.0, 0.8)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.PushMatrix()
	gl.Rotatef(float32(g_view_rotx), 1.0, 0.0, 0.0)
	gl.Rotatef(float32(g_view_roty), 0.0, 1.0, 0.0)
	gl.Rotatef(float32(g_view_rotz), 0.0, 0.0, 1.0)

	gl.PushMatrix()
	gl.Translatef(-3.0, -2.0, 0.0)
	gl.Rotatef(float32(g_angle), 0.0, 0.0, 1.0)
	gl.CallList(uint32(g_gear1))
	gl.PopMatrix()

	gl.PushMatrix()
	gl.Translatef(3.1, -2.0, 0.0)
	gl.Rotatef(float32(-2.0*g_angle-9.0), 0.0, 0.0, 1.0)
	gl.CallList(uint32(g_gear2))
	gl.PopMatrix()

	gl.PushMatrix()
	gl.Translatef(-3.1, 4.2, 0.0)
	gl.Rotatef(float32(-2.0*g_angle-25.0), 0.0, 0.0, 1.0)
	gl.CallList(uint32(g_gear3))
	gl.PopMatrix()

	gl.PopMatrix()
}

func create_gears() {
	pos := []float32{5.0, 5.0, 10.0, 0.0}
	red := []float32{0.8, 0.1, 0.0, 1.0}
	green := []float32{0.0, 0.8, 0.2, 1.0}
	blue := []float32{0.2, 0.2, 1.0, 1.0}

	gl.Lightfv(gl.LIGHT0, gl.POSITION, &pos[0])
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)
	gl.Enable(gl.DEPTH_TEST)

	g_gear1 = gl.GenLists(1)
	gl.NewList(g_gear1, gl.COMPILE)
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &red[0])
	gear(1.0, 4.0, 1.0, 20, 0.7)
	gl.EndList()

	g_gear2 = gl.GenLists(1)
	gl.NewList(g_gear2, gl.COMPILE)
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &green[0])
	gear(0.5, 2.0, 2.0, 10, 0.7)
	gl.EndList()

	g_gear3 = gl.GenLists(1)
	gl.NewList(g_gear3, gl.COMPILE)
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &blue[0])
	gear(1.3, 2.0, 0.5, 10, 0.7)
	gl.EndList()

	gl.Enable(gl.NORMALIZE)
}

func configure_window(pUserData unsafe.Pointer, pConfig *gwsi.WsiConfigureWindowEvent) {
	g_extent = gwsi.WsiExtent(pConfig.Extent)
	g_resized = true
}

func close_window(pUserData unsafe.Pointer, pInfo *gwsi.WsiCloseWindowEvent) {
	g_running = false
}

func main() {
	var info gwsi.WsiWindowCreateInfo
	var num_configs egl.EGLint
	var major, minor egl.EGLint
	var ok bool

	//	last_time := linux.GetTimeNs()
	last_time := time.Now()

	gwsi.WsiEglInit()
	xcb.LoadXcb()
	xcbimdkit.LoadEGLXcbImdkit()
	platform_info := gwsi.WsiPlatformCreateInfo{
		Type: gwsi.WSI_STRUCTURE_TYPE_PLATFORM_CREATE_INFO,
	}
	if platform, err := gwsi.WsiCreatePlatform(&platform_info); err != nil {
		panic(err)
	} else {
		g_platform = platform
	}

	res := gwsi.WsiGetEGLDisplay(g_platform, &g_display)
	if res != gwsi.WSI_SUCCESS {
		fmt.Printf("wsiGetEGLDisplay failed: %d", res)
		if res == gwsi.WSI_ERROR_EGL {
			fmt.Printf(" 0x%08x", egl.EglGetError())
		}
		fmt.Printf("\n")
		goto err_wsi_display
	}

	major, minor, ok = egl.EglInitialize(g_display)
	if !ok {
		fmt.Printf("eglInitialize failed: 0x%08x\n", egl.EglGetError())
		goto err_egl_init
	}

	if major < 1 || (major == 1 && minor < 4) {
		fmt.Printf("EGL version %d.%d is too old\n", major, minor)
		goto err_egl_version
	}

	ok = egl.EglBindAPI(egl.EGL_OPENGL_API)
	if !ok {
		fmt.Printf("eglBindAPI failed: 0x%08x\n", egl.EglGetError())
		goto err_egl_bind
	}

	num_configs = egl.EGLint(1)
	g_config, ok = egl.EglChooseConfig(g_display, g_config_attribs)
	if !ok {
		fmt.Printf("eglChooseConfig failed: 0x%08x\n", egl.EglGetError())
		goto err_egl_config
	}

	if num_configs == 0 {
		fmt.Printf("eglChooseConfig failed: no configs found\n")
		goto err_egl_config
	}

	g_context = egl.EglCreateContext(
		g_display,
		g_config,
		egl.NilEGLContext,
		g_context_attribs)
	if g_context == egl.NilEGLContext {
		fmt.Printf("eglCreateContext failed: 0x%08x\n", egl.EglGetError())
		goto err_egl_context
	}

	info = gwsi.WsiWindowCreateInfo{
		Type:            gwsi.WSI_STRUCTURE_TYPE_WINDOW_CREATE_INFO,
		Extent:          gwsi.WsiExtent{Width: 300, Height: 300},
		ConfigureWindow: configure_window,
		CloseWindow:     close_window,
	}

	g_window, res = g_platform.CreateWindow(&info, "gears")
	if res != gwsi.WSI_SUCCESS {
		fmt.Printf("wsiCreateWindow failed: %d\n", res)
		goto err_wsi_window
	}

	for {
		res = g_platform.DispatchEvents(-1)
		if res != gwsi.WSI_SUCCESS || g_resized {
			break
		}
	}
	if res != gwsi.WSI_SUCCESS {
		goto err_wsi_dispatch
	}

	res = gwsi.WsiCreateWindowEGLSurface(g_window, g_display, g_config, &g_surface)
	if res != gwsi.WSI_SUCCESS {
		fmt.Printf("wsiCreateWindowEGLSurface failed: %d", res)
		if res == gwsi.WSI_ERROR_EGL {
			fmt.Printf(" 0x%08x", egl.EglGetError())
		}
		fmt.Printf("\n")
		goto err_wsi_surface
	}

	ok = egl.EglMakeCurrent(g_display, g_surface, g_surface, g_context)
	if !ok {
		fmt.Printf("eglMakeCurrent failed: 0x%08x\n", egl.EglGetError())
		goto err_egl_current
	}

	ok = egl.EglSwapInterval(g_display, 0)
	if !ok {
		fmt.Printf("eglSwapInterval failed: 0x%08x\n", egl.EglGetError())
		goto err_egl_interval
	}
	gl.Init()

	create_gears()

	//last_time = linux.GetTimeNs()
	last_time = time.Now()

	for {
		draw()

		ok = egl.EglSwapBuffers(g_display, g_surface)
		if !ok {
			fmt.Printf("eglSwapBuffers failed: %d\n", egl.EglGetError())
			break
		}

		// now := linux.GetTimeNs()
		// dt := now - last_time
		// last_time = now
		now := time.Now()
		dt := now.Sub(last_time)
		last_time = now

		time := float64(dt) / 1000000000.0
		g_angle += time * 70.0
		g_angle = math.Mod(float64(g_angle), 360.0)

		res = g_platform.DispatchEvents(0)
		if res != gwsi.WSI_SUCCESS || !g_running {
			break
		}
	}

err_egl_interval:
err_egl_current:
	gwsi.WsiDestroyWindowEGLSurface(g_window, g_display, g_surface)
err_wsi_surface:
err_wsi_dispatch:
	gwsi.WsiDestroyWindow(g_window)
err_wsi_window:
	egl.EglDestroyContext(g_display, g_context)
err_egl_context:
err_egl_config:
err_egl_bind:
err_egl_version:
	egl.EglTerminate(g_display)
err_egl_init:
err_wsi_display:
	g_platform.Destroy()
}
