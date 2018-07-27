package netutil

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type TCPKeepAliveListener struct {
	*net.TCPListener
	duration time.Duration
}

func NewTCPKeepAliveListener(ln *net.TCPListener, d time.Duration) *TCPKeepAliveListener {
	return &TCPKeepAliveListener{
		TCPListener: ln,
		duration:    d,
	}
}

func (ln TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	if ln.duration == 0 {
		ln.duration = 3 * time.Minute
	}
	tc.SetKeepAlivePeriod(ln.duration)
	return tc, nil
}

func KeepAliveTCPConn(conn net.Conn, d time.Duration) {
	tc, ok := conn.(*net.TCPConn)
	if ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(d)
	}
}

func listenTCP(addrStr string) (net.Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp4", addrStr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	return listener, err
}

func listenUDP(addrStr string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenUDP("udp", addr)
	return listener, err
}

func serve(listener net.Listener, handler func(net.Conn), async bool) error {
	serveFunc := func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handler(conn)
		}
	}
	if async {
		go serveFunc()
	} else {
		serveFunc()
	}
	return nil
}

func ListenAndServeTCP(addrStr string, handler func(net.Conn), async bool, certs ...tls.Certificate) error {
	var (
		listener net.Listener
		err      error
	)
	if len(certs) > 0 {
		config := &tls.Config{Certificates: certs}
		listener, err = tls.Listen("tcp", addrStr, config)
	} else {
		listener, err = listenTCP(addrStr)
	}
	if err == nil {
		err = serve(listener, handler, async)
	}
	return err
}

func ListenAndServeWebsocket(addrStr, path string, handler func(net.Conn), async bool) error {
	mux := http.NewServeMux()
	mux.Handle(path, websocket.Handler(func(conn *websocket.Conn) {
		handler(conn)
	}))
	httpServer := &http.Server{
		Addr:           addrStr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if async {
		ln, err := net.Listen("tcp", httpServer.Addr)
		if err != nil {
			return err
		}
		go httpServer.Serve(NewTCPKeepAliveListener(ln.(*net.TCPListener), time.Minute*3))
		return nil
	}
	return httpServer.ListenAndServe()
}

func ListenAndServeUDP(addrStr string, handler func(*net.UDPConn), async bool) error {
	udpconn, err := listenUDP(addrStr)
	if err != nil {
		return err
	}

	serveFunc := func() {
		defer udpconn.Close()
		for {
			handler(udpconn)
		}
	}
	if async {
		go serveFunc()
	} else {
		serveFunc()
	}
	return nil
}
