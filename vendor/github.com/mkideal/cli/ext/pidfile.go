package ext

import (
	pidpkg "github.com/mkideal/pkg/osutil/pid"
)

// PidFile
type PidFile string

func (pid PidFile) String() string {
	return string(pid)
}

func (pid *PidFile) Decode(s string) error {
	*pid = PidFile(s)
	return nil
}

func (pid PidFile) New() error {
	return pidpkg.New(string(pid))
}

func (pid PidFile) Remove() error {
	return pidpkg.Remove(string(pid))
}
