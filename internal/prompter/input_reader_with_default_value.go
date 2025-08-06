package prompter

import (
	"os"
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
		return os.Stdin.Read(p)
	}

	// Read data from stdin
	n, err = os.Stdin.Read(p)
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
	return os.Stdin.Fd()
}
