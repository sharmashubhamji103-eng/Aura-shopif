//go:build linux

package main

import (
	"log"
	"syscall"
)

// raiseNofileLimit increases the open file descriptor limit to 65535.
// Default on most Linux VPS is 1024 — too low for 50+ concurrent card checks.
func raiseNofileLimit() {
	var r syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &r); err != nil {
		log.Printf("[api] WARNING: Getrlimit: %v", err)
		return
	}
	want := uint64(65535)
	if r.Max > 0 && r.Max < want {
		want = r.Max
	}
	r.Cur = want
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &r); err != nil {
		log.Printf("[api] WARNING: could not raise NOFILE limit: %v (run as root or set LimitNOFILE in systemd)", err)
	} else {
		log.Printf("[api] NOFILE limit raised to %d", r.Cur)
	}
}
