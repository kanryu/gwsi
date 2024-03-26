package main

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"

	"libwsi_test/internal/egl"
	"libwsi_test/libwsi"
)

const (
	TAU_F = 6.28318530717958647692
)

var (
	g_platform  libwsi.WsiPlatform
	g_window    libwsi.WsiWindow
	g_extent    libwsi.WsiExtent
	g_running   bool = true
	g_resized   bool = false
	g_display   egl.EGLDisplay
	g_config    egl.EGLConfig
	g_surface   egl.EGLSurface
	g_context   egl.EGLContext
	g_view_rotx float64 = 20.0
	g_view_roty float64 = 30.0
	g_view_rotz float64 = 0.0
	g_gear1     uint
	g_gear2     uint
	g_gear3     uint
	g_angle     float64 = 0.0

	g_config_attribs = []egl.EGLint{
		egl.egl.EGL_SURFACE_TYPE, egl.egl.EGL_WINDOW_BIT,
		egl.egl.EGL_RED_SIZE, 8,
		egl.egl.EGL_GREEN_SIZE, 8,
		egl.egl.EGL_BLUE_SIZE, 8,
		egl.egl.EGL_ALPHA_SIZE, 8,
		egl.egl.EGL_DEPTH_SIZE, 24,
		egl.egl.EGL_RENDERABLE_TYPE, egl.egl.EGL_OPENGL_BIT,
		egl.egl.EGL_NONE,
	}
	g_context_attribs = []egl.EGLint{
		egl.egl.EGL_CONTEXT_MAJOR_VERSION, 2,
		egl.egl.EGL_CONTEXT_MINOR_VERSION, 0,
		egl.egl.EGL_NONE,
	}
)

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

	gl.ShadeModel(gl.GL_FLAT)
	gl.Normal3f(0.0, 0.0, 1.0)

	gl.Begin(gl.GL_QUAD_STRIP)
	for i := 0; i <= teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(r0*math.Cos(angle), r0*math.Sin(angle), width*0.5)
		gl.Vertex3f(r1*math.Cos(angle), r1*math.Sin(angle), width*0.5)
		if i < teeth {
			gl.Vertex3f(r0*math.Cos(angle), r0*math.Sin(angle), width*0.5)
			gl.Vertex3f(r1*math.Cos(angle+3.0*da), r1*math.Sin(angle+3.0*da), width*0.5)
		}
	}
	gl.End()

	gl.Begin(gl.GL_QUADS)
	for i := 0; i < teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(r1*math.Cos(angle), r1*math.Sin(angle), width*0.5)
		gl.Vertex3f(r2*math.Cos(angle+da), r2*math.Sin(angle+da), width*0.5)
		gl.Vertex3f(r2*math.Cos(angle+2.0*da), r2*math.Sin(angle+2.0*da), width*0.5)
		gl.Vertex3f(r1*math.Cos(angle+3.0*da), r1*math.Sin(angle+3.0*da), width*0.5)
	}
	gl.End()

	gl.Normal3f(0.0, 0.0, -1.0)

	gl.Begin(gl.GL_QUAD_STRIP)
	for i := 0; i < teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(r1*math.Cos(angle), r1*math.Sin(angle), -width*0.5)
		gl.Vertex3f(r0*math.Cos(angle), r0*math.Sin(angle), -width*0.5)
		if i < teeth {
			gl.Vertex3f(r1*math.Cos(angle+3.0*da), r1*math.Sin(angle+3.0*da), -width*0.5)
			gl.Vertex3f(r0*math.Cos(angle), r0*math.Sin(angle), -width*0.5)
		}
	}
	gl.End()

	gl.Begin(gl.GL_QUADS)
	for i := 0; i < teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(r1*math.Cos(angle+3.0*da), r1*math.Sin(angle+3.0*da), -width*0.5)
		gl.Vertex3f(r2*math.Cos(angle+2.0*da), r2*math.Sin(angle+2.0*da), -width*0.5)
		gl.Vertex3f(r2*math.Cos(angle+da), r2*math.Sin(angle+da), -width*0.5)
		gl.Vertex3f(r1*math.Cos(angle), r1*math.Sin(angle), -width*0.5)
	}
	gl.nd()

	gl.Begin(gl.GL_QUAD_STRIP)
	for i := 0; i < teeth; i++ {
		angle := calc_angle(i, teeth)

		gl.Vertex3f(r1*math.Cos(angle), r1*math.Sin(angle), width*0.5)
		gl.Vertex3f(r1*math.Cos(angle), r1*math.Sin(angle), -width*0.5)
		u := r2*math.Cos(angle+da) - r1*math.Cos(angle)
		v := r2*math.Sin(angle+da) - r1*math.Sin(angle)
		leng := math.Sqrt(u*u + v*v)
		u /= leng
		v /= leng
		gl.Normal3f(v, -u, 0.0)
		gl.Vertex3f(r2*math.Cos(angle+da), r2*math.Sin(angle+da), width*0.5)
		gl.Vertex3f(r2*math.Cos(angle+da), r2*math.Sin(angle+da), -width*0.5)
		gl.Normal3f(math.Cos(angle), math.Sin(angle), 0.0)
		gl.Vertex3f(r2*math.Cos(angle+2.0*da), r2*math.Sin(angle+2.0*da), width*0.5)
		gl.Vertex3f(r2*math.Cos(angle+2.0*da), r2*math.Sin(angle+2.0*da), -width*0.5)
		u = r1*math.Cos(angle+3.0*da) - r2*math.Cos(angle+2.0*da)
		v = r1*math.Sin(angle+3.0*da) - r2*math.Sin(angle+2.0*da)
		gl.Normal3f(v, -u, 0.0)
		gl.Vertex3f(r1*math.Cos(angle+3.0*da), r1*math.Sin(angle+3.0*da), width*0.5)
		gl.Vertex3f(r1*math.Cos(angle+3.0*da), r1*math.Sin(angle+3.0*da), -width*0.5)
		gl.Normal3f(math.Cos(angle), math.Sin(angle), 0.0)
	}

	gl.Vertex3f(r1*math.Cos(0), r1*math.Sin(0), width*0.5)
	gl.Vertex3f(r1*math.Cos(0), r1*math.Sin(0), -width*0.5)

	gl.End()

	gl.ShadeModel(gl.GL_SMOOTH)

	gl.Begin(gl.GL_QUAD_STRIP)
	for i := 0; i <= teeth; i++ {
		angle := calc_angle(i, teeth)
		gl.Normal3f(-math.Cos(angle), -math.Sin(angle), 0.0)
		gl.Vertex3f(r0*math.Cos(angle), r0*math.Sin(angle), -width*0.5)
		gl.Vertex3f(r0*math.Cos(angle), r0*math.Sin(angle), width*0.5)
	}
	gl.End()
}

func draw() {
    if g_resized {
        gl.Viewport(0, 0, (GLsizei) g_extent.Width, (GLsizei) g_extent.Height)

        gl.MatrixMode(gl.GL_PROJECTION)
        gl.LoadIdentity()

        hf := float64(g_extent.Height)
        wf := float64(g_extent.Width)

        if (hf > wf) {
        	 aspect := hf / wf
            gl.Frustum(-1.0, 1.0, -aspect, aspect, 5.0, 60.0)
        } else {
            aspect := wf / hf
            gl.Frustum(-aspect, aspect, -1.0, 1.0, 5.0, 60.0)
        }

        gl.MatrixMode(gl.GL_MODELVIEW)
        gl.LoadIdentity()
        gl.Translatef(0.0, 0.0, -40.0)
        g_resized = false
    }

    gl.ClearColor(0.0, 0.0, 0.0, 0.8f)
    gl.Clear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)

    gl.PushMatrix()
    gl.Rotatef(g_view_rotx, 1.0, 0.0, 0.0)
    gl.Rotatef(g_view_roty, 0.0, 1.0, 0.0)
    gl.Rotatef(g_view_rotz, 0.0, 0.0, 1.0)

    gl.PushMatrix()
    gl.Translatef(-3.0, -2.0, 0.0)
    gl.Rotatef(g_angle, 0.0, 0.0, 1.0)
    gl.CallList(g_gear1)
    gl.PopMatrix()

    gl.PushMatrix()
    gl.Translatef(3.1f, -2.0, 0.0)
    gl.Rotatef(-2.0 * g_angle - 9.0, 0.0, 0.0, 1.0)
    gl.CallList(g_gear2)
    gl.PopMatrix()

    gl.PushMatrix()
    gl.Translatef(-3.1f, 4.2f, 0.0)
    gl.Rotatef(-2.0 * g_angle - 25.0, 0.0, 0.0, 1.0)
    gl.CallList(g_gear3)
    gl.PopMatrix()

    gl.PopMatrix()
}

func create_gears(){
    pos = []float64{ 5.0, 5.0, 10.0, 0.0 }
    red = []float64{ 0.8, 0.1, 0.0, 1.0 }
    green = []float64{ 0.0, 0.8, 0.2, 1.0 }
    blue = []float64{ 0.2, 0.2, 1.0, 1.0 }

    gl.Lightfv(gl.GL_LIGHT0, gl.GL_POSITION, pos)
    gl.Enable(gl.GL_CULL_FACE)
    gl.Enable(gl.GL_LIGHTING)
    gl.Enable(gl.GL_LIGHT0)
    gl.Enable(gl.GL_DEPTH_TEST)

    g_gear1 = gl.GenLists(1)
    gl.NewList(g_gear1, gl.GL_COMPILE)
    gl.Materialfv(gl.GL_FRONT, gl.GL_AMBIENT_AND_DIFFUSE, red)
    gear(1.0, 4.0, 1.0, 20, 0.7)
    gl.EndList()

    g_gear2 = gl.GenLists(1)
    gl.NewList(g_gear2, gl.GL_COMPILE)
    gl.Materialfv(gl.GL_FRONT, gl.GL_AMBIENT_AND_DIFFUSE, green)
    gear(0.5, 2.0, 2.0, 10, 0.7)
    gl.EndList()

    g_gear3 = gl.GenLists(1)
    gl.NewList(g_gear3, gl.GL_COMPILE)
    gl.Materialfv(gl.GL_FRONT, gl.GL_AMBIENT_AND_DIFFUSE, blue)
    gear(1.3, 2.0, 0.5, 10, 0.7)
    gl.EndList()

    gl.Enable(GL_NORMALIZE)
}

func main() {
	int ret := -1

    WsiPlatformCreateInfo platform_info = {
        .sType = libwgi.WSI_STRUCTURE_TYPE_PLATFORM_CREATE_INFO,
        .pNext = NULL,
    }

    WsiResult res = liwsi.WsiCreatePlatform(&platform_info, &g_platform)
    if (res != libwgi.WSI_SUCCESS) {
        fprintf(stderr, "wsiCreatePlatform failed: %d\n", res)
        goto err_wsi_platform
    }

    res = liwsi.WsiGetEGLDisplay(g_platform, &g_display)
    if (res != libwgi.WSI_SUCCESS) {
        fprintf(stderr, "wsiGetEGLDisplay failed: %d", res)
        if (res == libwgi.WSI_ERROR_EGL) {
            fprintf(stderr, " 0x%08x", egl.EglGetError())
        }
        fprintf(stderr, "\n")
        goto err_wsi_display
    }

    var major, minor egl.EGLint
    EGLBoolean ok = egl.EglInitialize(g_display, &major, &minor)
    if (ok == egl.EGL_FALSE) {
        fprintf(stderr, "eglInitialize failed: 0x%08x\n", egl.EglGetError())
        goto err_egl_init
    }

    if (major < 1 || (major == 1 && minor < 4)) {
        fprintf(stderr, "EGL version %d.%d is too old\n", major, minor)
        goto err_egl_version
    }

    ok = egl.EglBindAPI(egl.EGL_OPENGL_API)
    if (ok == egl.EGL_FALSE) {
        fprintf(stderr, "eglBindAPI failed: 0x%08x\n", egl.EglGetError())
        goto err_egl_bind
    }

     num_configs  := egl.EGLint(1)
    ok = egl.EglChooseConfig(g_display, g_config_attribs, &g_config, 1, &num_configs)
    if (ok == egl.EGL_FALSE) {
        fprintf(stderr, "eglChooseConfig failed: 0x%08x\n", egl.EglGetError())
        goto err_egl_config
    }

    if (num_configs == 0) {
        fprintf(stderr, "eglChooseConfig failed: no configs found\n")
        goto err_egl_config
    }

    g_context = egl.EglCreateContext(
        g_display,
        g_config,
        egl.EGL_NO_CONTEXT,
        g_context_attribs)
    if (g_context == egl.EGL_NO_CONTEXT) {
        fprintf(stderr, "eglCreateContext failed: 0x%08x\n", egl.EglGetError())
        goto err_egl_context
    }

    WsiWindowCreateInfo info = {
        .sType = libwgi.WSI_STRUCTURE_TYPE_WINDOW_CREATE_INFO,
        .pNext = NULL,
        .extent.width = 300,
        .extent.height = 300,
        .pTitle = "Gears",
        .pUserData = NULL,
        .pfnCloseWindow = close_window,
        .pfnConfigureWindow = configure_window,
    }

    res = liwsi.WsiCreateWindow(g_platform, &info, &g_window)
    if (res != libwgi.WSI_SUCCESS) {
        fprintf(stderr, "wsiCreateWindow failed: %d\n", res)
        goto err_wsi_window
    }

    while (true) {
        res = liwsi.WsiDispatchEvents(g_platform, -1)
        if (res != libwgi.WSI_SUCCESS || g_resized) {
            break
        }
    }
    if (res != libwgi.WSI_SUCCESS) {
        goto err_wsi_dispatch
    }

    res = liwsi.WsiCreateWindowEGLSurface(g_window, g_display, g_config, &g_surface)
    if (res != libwgi.WSI_SUCCESS) {
        fprintf(stderr, "wsiCreateWindowEGLSurface failed: %d", res)
        if (res == libwgi.WSI_ERROR_EGL) {
            fprintf(stderr, " 0x%08x", egl.EglGetError())
        }
        fprintf(stderr, "\n")
        goto err_wsi_surface
    }

    ok = egl.EglMakeCurrent(g_display, g_surface, g_surface, g_context)
    if (ok == egl.EGL_FALSE) {
        fprintf(stderr, "eglMakeCurrent failed: 0x%08x\n", egl.EglGetError())
        goto err_egl_current
    }

    ok = egl.EglSwapInterval(g_display, 0)
    if (ok == egl.EGL_FALSE) {
        fprintf(stderr, "eglSwapInterval failed: 0x%08x\n", egl.EglGetError())
        goto err_egl_interval
    }

    create_gears()

    int64_t last_time = get_time_ns()

    while(true) {
        draw()

        ok = egl.EglSwapBuffers(g_display, g_surface)
        if (ok == egl.EGL_FALSE) {
            printf("eglSwapBuffers failed: %d\n", egl.EglGetError())
            break
        }

        int64_t now = get_time_ns()
        int64_t dt = now - last_time
        last_time = now

        float time = (float)dt / 1000000000.0
        g_angle += time * 70.0
        g_angle = fmodf(g_angle, 360.0)

        res = liwsi.WsiDispatchEvents(g_platform, 0)
        if (res != libwgi.WSI_SUCCESS || !g_running) {
            break
        }
    }

    ret = EXIT_SUCCESS
err_egl_interval:
err_egl_current:
	liwsi.WsiDestroyWindowEGLSurface(g_window, g_display, g_surface)
err_wsi_surface:
err_wsi_dispatch:
	liwsi.WsiDestroyWindow(g_window)
err_wsi_window:
    egl.EglDestroyContext(g_display, g_context)
err_egl_context:
err_egl_config:
err_egl_bind:
err_egl_version:
    egl.EglTerminate(g_display)
err_egl_init:
err_wsi_display:
	liwsi.WsiDestroyPlatform(g_platform)
err_wsi_platform:
    return ret
}
