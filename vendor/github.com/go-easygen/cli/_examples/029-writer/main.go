package main

import (
	"bytes"
	"os"

	"github.com/mkideal/cli"
	clix "github.com/mkideal/cli/ext"
)

type argT struct {
	Writer *clix.Writer `cli:"w,writer" usage:"write to file, stdout or any io.Writer"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		argv.Writer.Close()
		n, err := argv.Writer.Write([]byte("hello,writer\n"))
		argv.Writer.Close()
		if err != nil {
			return err
		}
		ctx.String("writes %d byte(s) to file or stdout\n", n)
		ctx.String("filename: %s, isStdout: %v\n", argv.Writer.Name(), argv.Writer.IsStdout())

		// Replace the writer
		w := bytes.NewBufferString("")
		argv.Writer.SetWriter(w)
		n, err = argv.Writer.Write([]byte("hello,bytes.Writer"))
		if err != nil {
			return err
		}
		ctx.String("writes %d bytes to bytes.Writer: %s\n", n, w.String())
		ctx.String("filename: %s, isStdout: %v\n", argv.Writer.Name(), argv.Writer.IsStdout())
		return nil
	}))
}
