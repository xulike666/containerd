// +build linux

package console

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func tcget(fd uintptr, p *syscall.Termios) error {
	return ioctl(fd, syscall.TCGETS, uintptr(unsafe.Pointer(p)))
}

func tcset(fd uintptr, p *syscall.Termios) error {
	return ioctl(fd, syscall.TCSETS, uintptr(unsafe.Pointer(p)))
}

func ioctl(fd, flag, data uintptr) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, flag, data); err != 0 {
		return err
	}
	return nil
}

// unlockpt unlocks the slave pseudoterminal device corresponding to the master pseudoterminal referred to by f.
// unlockpt should be called before opening the slave side of a pty.
func unlockpt(f *os.File) error {
	var u int32
	return ioctl(f.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
}

// ptsname retrieves the name of the first available pts for the given master.
func ptsname(f *os.File) (string, error) {
	var n int32
	if err := ioctl(f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n))); err != nil {
		return "", err
	}
	return fmt.Sprintf("/dev/pts/%d", n), nil
}

func saneTerminal(f *os.File) error {
	// Go doesn't have a wrapper for any of the termios ioctls.
	var termios syscall.Termios
	if err := tcget(f.Fd(), &termios); err != nil {
		return err
	}
	// Set -onlcr so we don't have to deal with \r.
	termios.Oflag &^= syscall.ONLCR
	return tcset(f.Fd(), &termios)
}
