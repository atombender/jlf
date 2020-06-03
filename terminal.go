package main

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func getTerminalSize() size {
	if size, err := getTerminalWidthFromFd(os.Stdout.Fd()); err == nil {
		return size
	}
	if size, err := getTerminalWidthFromFd(os.Stderr.Fd()); err == nil {
		return size
	}
	if size, err := getTerminalWidthFromFd(os.Stdin.Fd()); err == nil {
		return size
	}

	f, err := os.Open("/dev/tty")
	if err == nil {
		defer f.Close()
		if size, err := getTerminalWidthFromFd(f.Fd()); err == nil {
			return size
		}
	}

	if w, ok := getEnvAsInt("COLUMNS"); ok {
		if h, ok := getEnvAsInt("LINES"); ok {
			return size{w, h}
		}
	}

	return size{w: 80, h: 24}
}

func getEnvAsInt(env string) (int, bool) {
	s := os.Getenv(env)
	if s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			return v, true
		}
	}
	return 0, false
}

func getTerminalWidthFromFd(fd uintptr) (size, error) {
	var dims [4]uint16
	if _, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dims)), 0, 0, 0); err != 0 {
		return size{}, err
	}
	return size{
		w: int(dims[1]),
		h: int(dims[0]),
	}, nil
}

type size struct {
	w, h int
}
