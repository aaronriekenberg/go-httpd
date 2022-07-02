package pledge

import (
	"github.com/aaronriekenberg/go-httpd/logging"
)

var logger = logging.GetLogger()

func FinalPledge() {
	const promises = "stdio rpath wpath inet unix"

	if err := pledgeWrapper(promises); err != nil {
		logger.Fatalf("FinalPledge pledgeWrapper err = %v", err)
	}

}
