package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func getRemote(sock net.PacketConn, remote chan string) {
	var buf [1024]byte

	for {
		rlen, _, err := sock.ReadFrom(buf[:]) // addr not needed
		if err != nil {
			panic(err)
		}
		if rlen > 0 {
			data := string(buf[:rlen])
			remote <- data
		}
	}
}

func getInput(reader *bufio.Reader, input chan string) {
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		input <- text[:len(text)-1] // remove the newline
	}
}

func main() {
	fmt.Println("client started")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Get socket ready to "connect"
	addr := &net.UDPAddr{
		Port: 2000,
		IP:   net.ParseIP("127.0.0.1"),
	}
	sock, err := net.ListenPacket("udp", ":0")
	if err != nil {
		panic(err)
	}

	// Enter a name and join the server
	fmt.Println("Enter a name:")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	_, err = sock.WriteTo([]byte("join:"+text[:len(text)-1]), addr)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("join %d\n", r)

	// Remote data
	remote := make(chan string, 64)
	go getRemote(sock, remote)

	// Get inputs into a channel
	input := make(chan string, 1)
	go getInput(reader, input)

	// See if inputs are ready to go, if so, send a message or recv one
	done := false
	for !done {
		select {
		case c := <-remote:
			fmt.Println(c)
		case c := <-input:
			_, err := sock.WriteTo([]byte("send:"+c), addr)
			if err != nil {
				panic(err)
			}
			// fmt.Printf("sent %d\n", r)
		case signal := <-sigs:
			// Disconnect
			r, err := sock.WriteTo([]byte("disc:"), addr)
			if err != nil {
				panic(err)
			}
			fmt.Println("disc", r, signal)
			done = true
		}
	}

	fmt.Println("client stopped")

}
