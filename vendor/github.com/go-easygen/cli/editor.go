package cli

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const DefaultEditor = "vim"

// GetEditor sets callback to get editor program
var GetEditor func() (string, error)

func getEditor() (string, error) {
	if GetEditor != nil {
		return GetEditor()
	}
	return exec.LookPath(DefaultEditor)
}

func randomFilename() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "CLI_EDIT_FILE"
	}
	return fmt.Sprintf(".%x", buf)
}

func LaunchEditor(editor string) (content []byte, err error) {
	return launchEditorWithFilename(editor, randomFilename())
}

func launchEditorWithFilename(editor, filename string) (content []byte, err error) {
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	defer os.Remove(filename)
	err = cmd.Run()
	if err != nil {
		if _, isExitError := err.(*exec.ExitError); !isExitError {
			return
		}
	}
	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, nil
	}
	return
}
