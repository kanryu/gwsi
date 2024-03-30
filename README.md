# gowsi

This is an experimental implementation of [libwsi](markbolhuis/libwsi) ported to the Go language.

This project is not currently intended to be reused in other projects.

## Significance of this project

This project was carried out to develop mado. A significant portion of this project's output will eventually be incorporated as part of mado.

A minimal implementation of OpenGL rendering that combines X11/EGL/xcb.

There are very few implementation examples of programming that combines EGL and xcb. I had a very hard time getting the program to work properly. This is because GLFW does not use xcb.

The first thing I understood after reading the source of this library is that Windows created with xcb requires its own EGLDisplay and EGLSurface. This is handled differently from a Window created with XCreateWindow. Other parts were the same.

libwsi solves this problem at a high level.
This library also has an interesting implementation for Wayland, which I would like to work on at a later date.

Some C APIs that I used for the first time, such as eglCreatePlatformWindowSurface(), were very difficult to execute from cgo, so I reimplemented them in purego. I would like to actively use it because it works fine and the build is lighter compared to cgo.

# Liscense

MIT