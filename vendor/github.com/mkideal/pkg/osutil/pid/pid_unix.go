// +build !windows,!darwin

package pid

import (
	"os"
	"path/filepath"
	"strconv"
)

func pidIsExist(pid int) bool {
	_, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid)))
	return err == nil
}
