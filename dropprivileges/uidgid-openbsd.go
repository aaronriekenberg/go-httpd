//go:build openbsd

package dropprivileges

import (
	"log"

	"golang.org/x/sys/unix"
)

func internalSetGID(gid int) error {
	log.Printf("openbsd internalSetGID gid = %v", gid)

	return unix.Setresgid(gid, gid, gid)
}

func internalSetUID(uid int) error {
	log.Printf("openbsd internalSetUID uid = %v", uid)

	return unix.Setresuid(uid, uid, uid)
}
