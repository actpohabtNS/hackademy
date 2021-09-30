package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"sync"
	"time"
)

type Message struct {
	Payload string
}

type FibonacciNum struct {
	value *big.Int
	prev  *big.Int
}

type Server struct {
	Cached map[int]*FibonacciNum
	*sync.Mutex
}

func ServerInit() *Server {
	server := &Server{}

	server.Cached = make(map[int]*FibonacciNum)
	server.Cached[1] = &FibonacciNum{big.NewInt(1), big.NewInt(0)}
	server.Cached[2] = &FibonacciNum{big.NewInt(1), big.NewInt(1)}
	server.Mutex = &sync.Mutex{}

	return server
}

func (server *Server) Run(port string) {
	listener, _ := net.Listen("tcp", port)
	if listener == nil {
		panic("ERROR: can not start listening " + port)
	}
	fmt.Println("listening port", port)

	for {
		connection, _ := listener.Accept()
		go server.handleConnection(connection)
	}
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func (server *Server) cacheFibonacci(num int) big.Int {
	fib, ok := server.Cached[num]
	if ok {
		fmt.Println("Getting num from cache")
		return *fib.value
	}

	closest := 1
	for i := range server.Cached {
		if abs(num-i) < abs(num-closest) {
			closest = i
		}
	}

	cur := new(big.Int).Set(server.Cached[closest].value)
	prev := new(big.Int).Set(server.Cached[closest].prev)
	temp := new(big.Int)

	if closest < num {
		for closest < num {
			temp.Set(cur)
			cur.Add(cur, prev)
			prev.Set(temp)
			closest += 1
		}
	} else {
		for closest > num {
			temp.Set(prev)
			prev.Sub(cur, prev)
			cur.Set(temp)
			closest -= 1
		}
	}

	server.Lock()
	server.Cached[num] = &FibonacciNum{cur, prev}
	server.Unlock()

	return *cur
}

func (server *Server) handleConnection(client net.Conn) {
	for {
		message := Message{}
		decoder := json.NewDecoder(client)
		encoder := json.NewEncoder(client)

		err := decoder.Decode(&message)
		if err != nil {
			continue
		}

		var num int
		_, err = fmt.Sscan(message.Payload, &num)
		if err != nil {
			message.Payload = "Wrong Argument"
			encErr := encoder.Encode(message)
			if encErr != nil {
				return
			}
			continue
		}

		start := time.Now()
		fib := server.cacheFibonacci(num)
		elapsed := time.Since(start)

		message.Payload = fmt.Sprintf("%s ", elapsed)
		message.Payload += fmt.Sprintf("%s", fib.String())
		encErr := encoder.Encode(message)
		if encErr != nil {
			return
		}
		fmt.Println("Message sent")
	}
}

func main() {
	go ServerInit().Run(":8081")
	fmt.Println("Press any key to stop")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
