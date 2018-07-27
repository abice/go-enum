package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mkideal/cli"
)

var sockFile = filepath.Join(os.Getenv("HOME"), ".rpc.sock")

func main() {
	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(daemon),
		cli.Tree(api,
			cli.Tree(ping),
		),
	).Run(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}

//------
// root
//------
var root = &cli.Command{
	Fn: func(ctx *cli.Context) error {
		ctx.WriteUsage()
		return nil
	},
}

//------
// help
//------
var help = &cli.Command{
	Name:        "help",
	Desc:        "display help",
	CanSubRoute: true,
	HTTPRouters: []string{"/v1/help"},
	HTTPMethods: []string{"GET"},

	Fn: cli.HelpCommandFn,
}

//--------
// daemon
//--------
type daemonT struct {
	cli.Helper
	Port uint16 `cli:"p,port" usage:"http port" dft:"8080"`
}

func (t *daemonT) Validate(ctx *cli.Context) error {
	if t.Port == 0 {
		return fmt.Errorf("please don't use 0 as http port")
	}
	return nil
}

var daemon = &cli.Command{
	Name: "daemon",
	Desc: "startup app as daemon",
	Argv: func() interface{} { return new(daemonT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*daemonT)
		addr := fmt.Sprintf(":%d", argv.Port)

		r := ctx.Command().Root()
		if err := r.RegisterHTTP(ctx); err != nil {
			return err
		}

		listeners := make([]net.Listener, 0, 2)

		// http listener
		if httpListener, err := net.Listen("tcp", addr); err != nil {
			return err
		} else if ln, ok := httpListener.(*net.TCPListener); ok {
			listeners = append(listeners, tcpKeepAliveListener{ln})
		}

		// unix listener
		if err := os.Remove(sockFile); err != nil {
			return err
		}
		if unixListener, err := net.Listen("unix", sockFile); err != nil {
			return err
		} else {
			listeners = append(listeners, unixListener)
		}

		ctx.String("listeners size: %d\n", len(listeners))

		return r.Serve(listeners...)
	},
}

// http client use unix sock
var httpc = &http.Client{
	Transport: &http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", sockFile)
		},
	},
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute)
	return tc, nil
}

//-----
// api
//-----
var api = &cli.Command{
	Name: "api",
	Desc: "display all api",
	Fn: func(ctx *cli.Context) error {
		cmd := ctx.Command().Root()
		if cmd.IsClient() {
			return cmd.RPC(httpc, ctx)
		}
		ctx.String("Commands:\n")
		ctx.String("    ping\n")
		return nil
	},
}

//------
// ping
//------
var ping = &cli.Command{
	Name: "ping",
	Desc: "ping server",
	Fn: func(ctx *cli.Context) error {
		cmd := ctx.Command().Root()
		if cmd.IsClient() {
			for {
				if err := cmd.RPC(httpc, ctx); err != nil {
					return err
				}
				time.Sleep(time.Millisecond * 1000)
			}
		}
		ctx.String(time.Now().Format(time.RFC3339) + " pong\n")
		return nil
	},
}
