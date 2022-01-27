package main

import (
	"bytes"
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {

	// Create an inotify instance represent flag arguements Fd is the file descriptor that we'll use to read events
	// We need to specify what files/directories we want to monitor

	fd, err := unix.InotifyInit1(0)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	defer unix.Close(fd)

	// We are monitoring the pwd for 6 events. Added a watch to the inotify instance here

	watch_descriptor, err := unix.InotifyAddWatch(fd, ".", unix.IN_CREATE|
		unix.IN_DELETE|
		unix.IN_MOVED_TO|
		unix.IN_MOVED_FROM|
		unix.IN_MOVE_SELF|
		unix.IN_MODIFY,
	)

	fmt.Printf("print the watch descriptor: %d", watch_descriptor)

	if err != nil {
		log.Fatalf("err: %v\n", err)

	}

	var buff [(unix.SizeofInotifyEvent + unix.NAME_MAX + 1) * 20]byte

	for {
		offset := 0
		n, err := unix.Read(fd, buff[:])
		if err != nil {
			log.Fatalf("err: %v\n", err)
		}

		for offset < n {
			e := (*unix.InotifyEvent)(unsafe.Pointer(&buff[offset]))

			nameBs := buff[offset+unix.SizeofInotifyEvent : offset+unix.SizeofInotifyEvent+int(e.Len)]
			name := string(bytes.TrimRight(nameBs, "\x00"))
			if len(name) > 0 && e.Mask&unix.IN_ISDIR == unix.IN_ISDIR {
				name += " (dir)"
			}

			switch {
			case e.Mask&unix.IN_CREATE == unix.IN_CREATE:
				fmt.Printf("CREATE %v\n", name)
			case e.Mask&unix.IN_DELETE == unix.IN_DELETE:
				fmt.Printf("DELETE %v\n", name)
			case e.Mask&unix.IN_CLOSE_WRITE == unix.IN_CLOSE_WRITE:
				fmt.Printf("CLOSE_WRITE %v\n", name)
			case e.Mask&unix.IN_MOVED_TO == unix.IN_MOVED_TO:
				fmt.Printf("IN_MOVED_TO [%v] %v\n", e.Cookie, name)
			case e.Mask&unix.IN_MOVED_FROM == unix.IN_MOVED_FROM:
				fmt.Printf("IN_MOVED_FROM [%v] %v\n", e.Cookie, name)
			case e.Mask&unix.IN_MOVE_SELF == unix.IN_MOVE_SELF:
				fmt.Printf("IN_MOVE_SELF %v\n", name)
			}

			offset += int(unix.SizeofInotifyEvent + e.Len)
		}
	}
}
