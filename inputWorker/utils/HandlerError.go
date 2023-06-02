package utils

import (
	"github.com/Temctl/E-Notification/util/elog"
)

// -------------------------------------------------------
// PUBLIC VAR --------------------------------------------
// -------------------------------------------------------
func HandleError(err error, msg string) {
	if err != nil {
		elog.Error(msg, err)
	}
}
