package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Message struct {
	Payload string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			return
		}

		message := Message{text}

		encoder := json.NewEncoder(conn)
		err := encoder.Encode(message)
		if err != nil {
			return
		}

		decoder := json.NewDecoder(conn)
		err = decoder.Decode(&message)
		if err != nil {
			continue
		}

		_, err = fmt.Fprintf(os.Stdout, "%s\n", message.Payload)
		if err != nil {
			return
		}
	}
}
