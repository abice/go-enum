package cli

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/gommon/color"
	"github.com/mkideal/pkg/debug"
)

// RegisterHTTP init HTTPRouters for command
func (cmd *Command) RegisterHTTP(ctxs ...*Context) error {
	clr := color.Color{}
	clr.Disable()
	if len(ctxs) > 0 {
		clr = ctxs[0].color
	}
	if cmd.routersMap == nil {
		cmd.routersMap = make(map[string]string)
	}
	commands := []*Command{cmd}
	for len(commands) > 0 {
		c := commands[0]
		commands = commands[1:]
		if c.HTTPRouters != nil {
			for _, r := range c.HTTPRouters {
				if _, exists := c.routersMap[r]; exists {
					return throwRouterRepeat(clr.Yellow(r))
				}
				cmd.routersMap[r] = c.Path()
			}
		}
		if c.nochild() {
			continue
		}
		commands = append(commands, c.children...)
	}
	return nil
}

// ServeHTTP implements HTTP handler
func (cmd *Command) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	var (
		path  = r.URL.Path
		found = false
	)
	if cmd.routersMap != nil {
		path, found = cmd.routersMap[path]
	}
	if !found {
		path = strings.TrimPrefix(r.URL.Path, "/")
		path = strings.TrimSuffix(path, "/")
	}

	router := strings.Split(path, "/")
	args := make([]string, 0, len(r.Form)*2+len(router))
	for _, r := range router {
		args = append(args, r)
	}
	for key, values := range r.Form {
		if len(key) == 0 || len(values) == 0 {
			continue
		}
		if !strings.HasPrefix(key, dashOne) {
			if len(key) == 1 {
				key = dashOne + key
			} else {
				key = dashTwo + key
			}
		}
		args = append(args, key, values[len(values)-1])
	}
	debug.Debugf("agent: %s", r.UserAgent())
	debug.Debugf("path: %s", path)
	debug.Debugf("args: %q", args)

	buf := new(bytes.Buffer)
	statusCode := http.StatusOK
	if err := cmd.RunWith(args, buf, w, r.Method); err != nil {
		buf.Write([]byte(err.Error()))
		nativeError := err
		if werr, ok := err.(wrapError); ok {
			nativeError = werr.err
		}
		switch nativeError.(type) {
		case commandNotFoundError:
			statusCode = http.StatusNotFound

		case methodNotAllowedError:
			statusCode = http.StatusMethodNotAllowed

		default:
			statusCode = http.StatusInternalServerError
		}
	}
	debug.Debugf("resp: %s", buf.String())
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}

// ListenAndServeHTTP set IsServer flag with true and startup http service
func (cmd *Command) ListenAndServeHTTP(addr string) error {
	cmd.SetIsServer(true)
	return http.ListenAndServe(addr, cmd)
}

// Serve set IsServer with true and serve http with listeners
func (cmd *Command) Serve(listeners ...net.Listener) (err error) {
	cmd.SetIsServer(true)
	var g sync.WaitGroup
	for _, ln := range listeners {
		g.Add(1)
		go func(ln net.Listener) {
			if e := http.Serve(ln, cmd); e != nil {
				panic(e.Error())
			}
			g.Done()
		}(ln)
	}
	g.Wait()
	return
}

// RPC runs the command from remote
func (cmd *Command) RPC(httpc *http.Client, ctx *Context) error {
	addr := "http://rpc/" + ctx.Command().pathWithSep("/")
	method := "POST"
	if cmd.HTTPMethods != nil && len(cmd.HTTPMethods) > 0 {
		method = cmd.HTTPMethods[0]
	}
	var body io.Reader
	if values := ctx.FormValues(); values != nil {
		body = strings.NewReader(values.Encode())
	}
	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "cli-RPC")
	resp, err := httpc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(ctx, resp.Body)
	return err
}
