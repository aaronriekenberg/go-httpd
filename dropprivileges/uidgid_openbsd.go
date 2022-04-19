//go:build openbsd

package dropprivileges

import (
	"log"
	"syscall"
)

func internalSetGID(gid int) error {
	log.Printf("openbsd internalSetGID gid = %v", gid)

	return syscall.Setresgid(gid, gid, gid)
}

func internalSetUID(uid int) error {
	log.Printf("openbsd internalSetUID uid = %v", uid)

	return syscall.Setresuid(uid, uid, uid)
}
