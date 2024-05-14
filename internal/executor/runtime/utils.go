package runtime

import "runtime"

type OS string

const (
	OSWindows OS = "windows"
	OSLinux   OS = "linux"
	OSMacOS   OS = "macos"
	OSUnknown OS = "unknown"
)

func GetHostOS() OS {
	os := runtime.GOOS
	switch os {
	case "windows":
		return OSWindows
	case "darwin":
		return OSMacOS
	case "linux":
		return OSLinux
	default:
		return OSUnknown
	}
}
