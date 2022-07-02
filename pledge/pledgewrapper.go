//go:build openbsd

package pledge

import "golang.org/x/sys/unix"

func pledgeWrapper(promises string) error {

	logger.Printf("pledgeWrapper openbsd promises = %q", promises)

	return unix.PledgePromises(promises)
}
