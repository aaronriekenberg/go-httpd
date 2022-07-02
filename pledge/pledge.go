package pledge

import (
	"github.com/aaronriekenberg/go-httpd/logging"
)

var logger = logging.GetLogger()

func Pledge() {
	const promises = "stdio rpath wpath cpath inet unix"

	if err := pledgeWrapper(promises); err != nil {
		logger.Fatalf("Pledge pledgeWrapper err = %v", err)
	}

}
