package wangxsrpc

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

func startServerDebug(addrCh chan string) {
	var foo Foo
	l, _ := net.Listen("tcp", ":9999")
	log.Println("debug page: http://localhost:9999/debug/wangxsrpc")
	_ = Register(&foo)
	HandleHTTP()
	addrCh <- l.Addr().String()
	_ = http.Serve(l, nil)
}

func call(addrCh chan string) {
	client, _ := DialHTTP("tcp", <-addrCh)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func TestDebugHTTP_ServeHTTP(t *testing.T) {
	t.Run("test_debug", func(t *testing.T) {
		log.SetFlags(0)
		ch := make(chan string)
		go call(ch)
		startServerDebug(ch)
	})
}


