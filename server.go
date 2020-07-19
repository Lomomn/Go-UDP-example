package main

import (
	"fmt"
	"net"
	"strings"
)

func serve() {

}

func main() {
	fmt.Println("server started")

	sock, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 2000,
		IP:   net.ParseIP("127.0.0.1"),
	})
	if err != nil {
		panic(err)
	}
	defer sock.Close()

	connections := make(map[string]string)

	var buf [1024]byte
	for {
		rlen, addr, err := sock.ReadFromUDP(buf[:])
		if err != nil {
			panic(err)
		}
		data := string(buf[:rlen])
		fmt.Println(data, rlen, addr)

		split := strings.Split(data, ":")
		switch split[0] {
		case "join":
			fmt.Println("player join", split)
			connections[addr.String()] = split[1] // save name against addr
			fmt.Println(connections, addr, split[1], connections[addr.String()])
		case "send":
			fmt.Println(connections[addr.String()], "says:", split[1:])

			// Send message to everybody else
			from := connections[addr.String()]
			fmt.Println("\nsending to others:")
			for caddr, name := range connections {
				ip, err := net.ResolveUDPAddr("udp", caddr)
				if err != nil {
					panic(err)
				}

				// Don't send the message back to sender
				if caddr != addr.String() {
					msg := from + ": " + strings.Join(split[1:], "")
					fmt.Println(caddr, name, ip, msg)
					sock.WriteTo([]byte(msg), ip) // No command added, maybe add
				}
			}
			fmt.Println()
		case "disc":
			fmt.Println(connections[addr.String()], "disconnected")
			delete(connections, addr.String())
		}
	}

	fmt.Println("server stopped")
}
