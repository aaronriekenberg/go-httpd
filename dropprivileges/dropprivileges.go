package dropprivileges

import (
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/aaronriekenberg/go-httpd/config"
	"github.com/aaronriekenberg/go-httpd/logging"
)

var logger = logging.GetLogger()

func DropPrivileges(
	config *config.DropPrivileges,
) {

	if config == nil {
		return
	}

	userObject, err := user.Lookup(config.UserName)
	if err != nil {
		logger.Fatalf("user.Lookup failed config.UserName = %q error: %v", config.UserName, err)
	}

	uidInt, err := strconv.Atoi(userObject.Uid)
	if err != nil {
		logger.Fatalf("strconv.Atoi userObject.Uid = %q error: %v", userObject.Uid, err)
	}

	groupObject, err := user.LookupGroup(config.GroupName)
	if err != nil {
		logger.Fatalf("user.LookupGroup failed config.GroupName = %q error: %v", config.GroupName, err)
	}

	gidInt, err := strconv.Atoi(groupObject.Gid)
	if err != nil {
		logger.Fatalf("strconv.Atoi groupObject.Gid = %q error: %v", groupObject.Gid, err)
	}

	if config.ChrootEnabled {
		if err := syscall.Chroot(config.ChrootDirectory); err != nil {
			logger.Fatalf("Chroot failed error: %v", err)
		}

		if err := os.Chdir("/"); err != nil {
			logger.Fatalf("Chdir / failed error: %v", err)
		}
	}

	if err := syscall.Setgroups([]int{gidInt}); err != nil {
		logger.Fatalf("syscall.Setgroups %v error: %v", []int{gidInt}, err)
	}

	if err := syscall.Setgid(gidInt); err != nil {
		logger.Fatalf("syscall.Setgid gitInt = %v error: %v", gidInt, err)
	}

	if err := syscall.Setuid(uidInt); err != nil {
		logger.Fatalf("syscall.Setuid uidInt = %v error: %v", uidInt, err)
	}

}
