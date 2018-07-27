package pid

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Exsit(pid int) bool {
	return pidIsExist(pid)
}

func New(filename string) error {
	dir, _ := filepath.Split(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("PidFile: %v", err)
		}
	}
	if content, err := ioutil.ReadFile(filename); err == nil {
		pidStr := strings.TrimSpace(string(content))
		if pidStr == "" {
			return nil
		}
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return nil
		}
		if pidIsExist(pid) {
			return fmt.Errorf("pid file found, ensoure %s is not running", os.Args[0])
		}
	}
	if err := ioutil.WriteFile(filename, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return err
	}
	return nil
}

func Remove(filename string) error {
	if filename != "" {
		return os.Remove(filename)
	}
	return nil
}
