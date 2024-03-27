//go:build linux
// +build linux

package linux

/*
#include <time.h>
static unsigned long long get_nsecs(void)
{
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (unsigned long long)ts.tv_sec * 1000000000UL + ts.tv_nsec;
}*/
import "C"

func GetTimeNs() int64 {
	ts := C.get_nsecs()
	return int64(ts)
}
