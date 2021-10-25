package codec_test

import (
	"encoding/json"
	"fmt"
	"geerpc"
	"geerpc/codec"
	"log"
	"net"
	"testing"
	"time"
)

type Bar int

func (b Bar) Timeout(argv int, reply *int) error {
	time.Sleep(2 * time.Second)
	return nil
}

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assetion failed: "+msg, v...))
	}
}

func TestNewJsonCodec(t *testing.T) {
	addr := make(chan string)
	go startServer(addr)
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { conn.Close() }()

	time.Sleep(time.Second)
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewJsonCodec(conn)
	defer cc.Close()
	h := &codec.Header{
		ServiceMethod: "Bar.Timeout",
		Seq:           uint64(1),
	}
	_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
	_ = cc.ReadHeader(h)
	var reply string
	_ = cc.ReadBody(&reply)
	log.Println("reply:", reply)
	_assert(reply != "", "reply is not empty")
}
