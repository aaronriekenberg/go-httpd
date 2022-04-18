package dropprivileges

import (
	"log"
	"syscall"
)

func internalSetGID(gid int) error {
	log.Printf("darwin internalSetGID gid = %v", gid)

	return syscall.Setregid(gid, gid)
}

func internalSetUID(uid int) error {
	log.Printf("darwin internalSetGID uid = %v", uid)

	return syscall.Setregid(uid, uid)
}
