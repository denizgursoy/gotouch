package prompter

import (
	"sync"
	"syscall"
)

type defaultMessageReader struct {
	prepended    bool
	initialValue string
	mu           sync.Mutex
}

func (d *defaultMessageReader) Read(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.prepended {
		return syscall.Read(int(d.Fd()), p)
	}

	// Read data from stdin
	n, err = syscall.Read(int(d.Fd()), p)
	if err != nil {
		return n, err
	}

	totalLength := len(d.initialValue) + n

	if len(p) < totalLength {
		return 0, syscall.EINVAL
	}

	// Shift data to the right
	copy(p[len(d.initialValue):], p[:n])
	copy(p[:len(d.initialValue)], d.initialValue)

	d.prepended = false // Mark as already prepended
	return totalLength, nil
}

func (d *defaultMessageReader) Fd() uintptr {
	return uintptr(syscall.Stdin)
}
