package tcp

import (
	"fmt"
	"net"
	"time"

	"github.com/byteweap/wukong/pkg/knet"
)

type server struct {
	opts *Options
}

var _ knet.Server = (*server)(nil)

func NewServer(opts ...Option) knet.Server {

	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &server{
		opts: options,
	}
}

func (s *server) Addr() string {
	return s.opts.Addr
}

func (s *server) Protocol() string {
	return "tcp"
}

func (s *server) Start() {

	addr, err := net.ResolveTCPAddr("tcp", s.opts.Addr)
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("tcp accept err:%v\n", err)
			continue
		}

		go func() {
			for {
				buf := make([]byte, 1024*2)
				cnt, err := conn.Read(buf)
				if err != nil {
					fmt.Printf("tcp read err:%v\n", err)
					continue
				}
				fmt.Printf("tcp read len: %d, content: %s \n", cnt, string(buf[:cnt]))
				if _, err := conn.Write(buf[:cnt]); err != nil {
					fmt.Printf("tcp write err:%v\n", err)
				}

				time.Sleep(time.Second * 3)
			}
		}()
	}

}

func (s *server) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s *server) OnStart(handler knet.StartHandler) {
	//TODO implement me
	panic("implement me")
}

func (s *server) OnStop(handler knet.StopHandler) {
	//TODO implement me
	panic("implement me")
}

func (s *server) OnConnect(handler knet.ConnectHandler) {
	//TODO implement me
	panic("implement me")
}

func (s *server) OnDisconnect(handler knet.ConnectHandler) {
	//TODO implement me
	panic("implement me")
}

func (s *server) OnTextMessage(handler knet.ConnMessageHandler) {
	//TODO implement me
	panic("implement me")
}

func (s *server) OnBinaryMessage(handler knet.ConnMessageHandler) {
	//TODO implement me
	panic("implement me")
}

func (s *server) OnError(handler knet.ErrorHandler) {
	//TODO implement me
	panic("implement me")
}
