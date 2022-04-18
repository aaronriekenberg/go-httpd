package dropprivileges

import (
	"log"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/aaronriekenberg/go-httpd/config"
)

func DropPrivileges(
	config *config.DropPrivileges,
) {

	log.Printf("begin DropPrivileges")

	if config == nil {
		log.Printf("DropPrivileges config is nil returning")
		return
	}

	userObject, err := user.Lookup(config.UserName)
	if err != nil {
		log.Fatalf("Lookup failed: %q", config.UserName)
	}
	log.Printf("got userObject %+v", userObject)

	uidInt, err := strconv.Atoi(userObject.Uid)
	if err != nil {
		log.Fatalf("strconv.Atoi userObject.Uid = %q error: %v", userObject.Uid, err)
	}
	log.Printf("uidInt = %v", uidInt)

	groupObject, err := user.LookupGroup(config.GroupName)
	if err != nil {
		log.Fatalf("LookupGroup failed: %q", config.GroupName)
	}
	log.Printf("got groupObject %+v", groupObject)

	gidInt, err := strconv.Atoi(groupObject.Gid)
	if err != nil {
		log.Fatalf("strconv.Atoi groupObject.Gid = %q error: %v", groupObject.Gid, err)
	}
	log.Printf("gidInt = %v", gidInt)

	if config.ChrootEnabled {
		log.Printf("Chroot to %q", config.ChrootDirectory)
		err := syscall.Chroot(config.ChrootDirectory)
		if err != nil {
			log.Fatalf("Chroot failed: %v", err)
		}

		err = os.Chdir("/")
		if err != nil {
			log.Fatalf("Chdir / failed: %v", err)
		}
	}

	err = syscall.Setgroups([]int{gidInt})
	if err != nil {
		log.Fatalf("syscall.Setgroups %v error: %v", []int{gidInt}, err)
	}

	err = internalSetGID(gidInt)
	if err != nil {
		log.Fatalf("internalSetGID gitInt = %v error: %v", gidInt, err)
	}

	err = internalSetUID(uidInt)
	if err != nil {
		log.Fatalf("internalSetUID uidInt = %v error: %v", uidInt, err)
	}

	log.Printf("end DropPrivileges")

}
