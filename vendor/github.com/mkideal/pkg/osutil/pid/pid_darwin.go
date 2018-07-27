// +build darwin

package pid

import (
	"syscall"
)

func pidIsExist(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}
