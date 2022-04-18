package dropprivileges

import (
	"log"

	"golang.org/x/sys/unix"
)

func internalSetGID(gid int) error {
	log.Printf("darwin internalSetGID gid = %v", gid)

	return unix.Setregid(gid, gid)
}

func internalSetUID(uid int) error {
	log.Printf("darwin internalSetGID uid = %v", uid)

	return unix.Setregid(uid, uid)
}
