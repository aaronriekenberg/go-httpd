//go:build !openbsd

package dropprivileges

import (
	"log"
	"syscall"
)

func internalSetGID(gid int) error {
	log.Printf("generic internalSetGID gid = %v", gid)

	return syscall.Setgid(gid)
}

func internalSetUID(uid int) error {
	log.Printf("generic internalSetGID uid = %v", uid)

	return syscall.Setuid(uid)
}
