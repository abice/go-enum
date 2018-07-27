package ext_test

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/mkideal/cli"
	"github.com/mkideal/cli/ext"
)

// This example demonstrates usage of Time decoder
func ExampleTime() {
	type argT struct {
		When ext.Time `cli:"w"`
	}
	for _, args := range [][]string{
		[]string{"app", "-w", "2016-01-02 12:12:22"},
		[]string{"app", "-w2016-01-02"},
	} {
		cli.RunWithArgs(new(argT), args, func(ctx *cli.Context) error {
			argv := ctx.Argv().(*argT)
			ctx.String("time=%s\n", argv.When.Time.Format(time.ANSIC))
			return nil
		})
	}
	// Output:
	// time=Sat Jan  2 12:12:22 2016
	// time=Sat Jan  2 00:00:00 2016
}

// This example demonstrates uage of Duration decoder
func ExampleDuration() {
	type argT struct {
		Duration ext.Duration `cli:"d"`
	}
	for _, args := range [][]string{
		[]string{"app", "-d10s"},
		[]string{"app", "-d10ms"},
	} {
		cli.RunWithArgs(new(argT), args, func(ctx *cli.Context) error {
			argv := ctx.Argv().(*argT)
			ctx.String("duration=%v\n", argv.Duration.Duration)
			return nil
		})
	}
	// Output:
	// duration=10s
	// duration=10ms
}

// This example demonstrates usage of File decoder
func ExampleFile() {
	type argT struct {
		Data ext.File `cli:"f"`
	}
	filename := "test.txt"
	ioutil.WriteFile(filename, []byte("hello, File decoder"), 0644)
	defer os.Remove(filename)

	args := []string{"app", "-f", filename}
	cli.RunWithArgs(new(argT), args, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("%s\n", argv.Data)
		return nil
	})
	// Output:
	// hello, File decoder
}

// This example demonstrates usage of PidFile decoder
//func ExamplePidFile() {
//	type argT struct {
//		Pid ext.PidFile `cli:"pid"`
//	}
//
//	args := []string{"app", "--pid=test.pid"}
//	cli.RunWithArgs(new(argT), args, func(ctx *cli.Context) error {
//		argv := ctx.Argv().(*argT)
//
//		if err := argv.Pid.New(); err != nil {
//			return err
//		}
//		defer argv.Pid.Remove()
//
//		return nil
//	})
//}
