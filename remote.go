package gopwntools

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"

	"os/signal"
)

type RemoteConf struct {
	Protocol string
}

func remote[V number](host string, port V, protocol string) *Conn {
	info = connInfo{host: host, port: fmt.Sprintf("%v", port), isRemote: true}
	p := Progress(fmt.Sprintf("Opening connection to %s on port %s", info.host, info.port))

	netConn, err := net.Dial(protocol, net.JoinHostPort(info.host, info.port))
	p.Status(fmt.Sprintf("Trying %s:%s", info.host, info.port))

	if err != nil {
		panic(err)
	}

	stdin := io.WriteCloser(netConn)
	stdout := io.ReadCloser(netConn)
	p.Success("Done")

	conn := Conn{stdin: stdin, stdout: stdout, errChan: make(chan error, 1), conn: netConn}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		conn.errChan <- fmt.Errorf("Control-C")
	}()

	return &conn
}

func Remote[V number](host string, port V) *Conn {
	return remote(host, port, "tcp")
}

func RemoteWithConf[V number](host string, port V, conf RemoteConf) *Conn {
	return remote(host, port, conf.Protocol)
}
