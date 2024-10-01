//go:build windows
// +build windows

package memory

import (
	"fmt"
	"regexp"
	"strings"
	"syscall"
	"unsafe"

	windows "github.com/elastic/go-windows"
	xsyscall "golang.org/x/sys/windows"
)

var (
	kernel32 = xsyscall.NewLazySystemDLL("kernel32.dll")
	user32   = xsyscall.NewLazySystemDLL("user32.dll")

	procEnumWindows          = user32.NewProc("EnumWindows")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	getWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")

	procVirtualQueryEx         = kernel32.NewProc("VirtualQueryEx")
	queryFullProcessImageNameW = kernel32.NewProc("QueryFullProcessImageNameW")
)

func enumWindows(enumFunc uintptr, lparam uintptr) error {
	r1, _, e1 := syscall.SyscallN(procEnumWindows.Addr(), enumFunc, lparam, 0)
	if r1 == 0 {
		if e1 != 0 {
			return e1
		}

		return syscall.EINVAL
	}
	return nil
}

func getWindowText(handle syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(
		procGetWindowTextW.Addr(),
		3,
		uintptr(handle),
		uintptr(unsafe.Pointer(str)),
		uintptr(maxCount),
	)

	if r0 == 0 {
		if e1 != 0 {
			return 0, e1
		}
		return 0, syscall.EINVAL
	}

	return int32(r0), nil
}

func GetWindowThreadProcessID(hwnd syscall.Handle) (int32, error) {
	var processID int32

	if _, _, err := getWindowThreadProcessID.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&processID)),
	); err != nil {
		return 0, err
	}

	return processID, nil
}

func FindWindow(title string) (syscall.Handle, error) {
	var handle syscall.Handle

	cb := syscall.NewCallback(func(cHandle syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)

		if _, err := getWindowText(cHandle, &b[0], int32(len(b))); err != nil {
			return 1
		}

		if strings.Contains(syscall.UTF16ToString(b), title) {
			handle = cHandle
			return 0
		}

		return 1
	})

	_ = enumWindows(cb, 0)
	if handle == 0 {
		return 0, fmt.Errorf("no window with title '%s' found", title)
	}

	return handle, nil
}

func virtualQueryEx(handle syscall.Handle, off int64) (region, error) {
	var reg region

	// Syscall6 is deprecated
	r1, _, e1 := syscall.SyscallN(
		procVirtualQueryEx.Addr(),
		uintptr(handle),
		uintptr(off),
		uintptr(unsafe.Pointer(&reg)),
		unsafe.Sizeof(reg),
	)

	if r1 == 0 {
		if e1 != 0 {
			return region{}, e1
		} else {
			return region{}, syscall.EINVAL
		}
	}

	return reg, nil
}

func queryFullProcessImageName(hProcess syscall.Handle) (string, error) {
	var buf [syscall.MAX_PATH]uint16
	n := uint32(len(buf))

	r1, _, e1 := queryFullProcessImageNameW.Call(
		uintptr(hProcess),
		uintptr(0),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&n)),
	)

	if r1 == 0 {
		if e1 != nil {
			return "", e1
		} else {
			return "", syscall.EINVAL
		}
	}

	return syscall.UTF16ToString(buf[:n]), nil
}

func FindProcess(re *regexp.Regexp, blacklistedTitles ...string) ([]Process, error) {
	var processes []Process

	pids, err := windows.EnumProcesses()
	if err != nil {
		return nil, err
	}

	for _, pid := range pids {
		handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pid)
		if err != nil {
			continue
		}

		name, err := windows.GetProcessImageFileName(handle)
		if err != nil {
			_ = syscall.CloseHandle(handle)
			continue
		}

		if re.MatchString(name) {
			var bannedHandles []syscall.Handle

			for _, title := range blacklistedTitles {
				h, _ := FindWindow(title)
				if h != 0 {
					bannedHandles = append(bannedHandles, h)
				}
			}

			var isBanned = false
			for _, bHandle := range bannedHandles {
				pid2, err := GetWindowThreadProcessID(bHandle)
				if err != nil {
					return nil, err
				}

				if int32(pid) == pid2 {
					isBanned = true
					break
				}
			}

			if !isBanned {
				processes = append(processes, process{pid, handle})
			}
		}
	}

	if len(processes) < 1 {
		return nil, ErrNoProcess
	}

	return processes, nil
}

type process struct {
	pid uint32
	h   syscall.Handle
}

func (p process) HandleFromTitle() (string, error) {
	return queryFullProcessImageName(p.h)
}

func (p process) ExecutablePath() (string, error) {
	return queryFullProcessImageName(p.h)
}

func (p process) Close() error {
	return syscall.CloseHandle(p.h)
}

func (p process) Pid() int {
	return int(p.pid)
}

func (p process) ReadAt(b []byte, off int64) (n int, err error) {
	un, err := windows.ReadProcessMemory(p.h, uintptr(off), b)
	return int(un), err
}

func (p process) Maps() ([]Map, error) {
	lastAddr := int64(0)
	var maps []Map
	for {
		reg, err := virtualQueryEx(p.h, lastAddr)
		if err != nil {
			if lastAddr == 0 {
				return nil, err
			}
			break
		}
		maps = append(maps, reg)
		lastAddr = reg.Start() + reg.Size()
	}
	return maps, nil
}

type region struct {
	baseAddress       uintptr
	allocationBase    uintptr
	allocationProtect int32
	regionSize        int
	state             int32
	protect           int32
	type_             int32
}

func (r region) Start() int64 {
	return int64(r.baseAddress)
}

func (r region) Size() int64 {
	return int64(r.regionSize)
}
