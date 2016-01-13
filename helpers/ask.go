package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func Ask(ask string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(ask)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(string(text))
}

func Password(ask string) string {
	echo(false)
	defer echo(true)
	ret := Ask(ask)
	fmt.Println("")
	return ret
}

// Enable or disable echoing terminal input.
func echo(show bool) {
	var termios = &syscall.Termios{}
	var fd = os.Stdout.Fd()

	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios))); err != 0 {
		return
	}

	if show {
		termios.Lflag |= syscall.ECHO
	} else {
		termios.Lflag &^= syscall.ECHO
	}

	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(termios))); err != 0 {
		return
	}
}
