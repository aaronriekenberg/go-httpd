//go:build !openbsd

package pledge

import _ "golang.org/x/sys/unix"

func pledgeWrapper(promises string) error {

	logger.Printf("pledgeWrapper noop promises = %q", promises)

	return nil
}
