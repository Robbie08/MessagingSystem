package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Message struct {
	sender  string
	target  string
	payload []byte
}

type Client struct {
	name     string
	dialAddr string
	dialer   net.Dialer
	quitChan chan struct{}
	msgChan  chan Message
}

/*
Constructor for our client
*/
func NewClient(name, dialAddr string) *Client {
	return &Client{
		name:     name,
		dialAddr: dialAddr,
		quitChan: make(chan struct{}),
		msgChan:  make(chan Message, 10),
	}
}

/*
This will start up our TCP client
*/
func (c *Client) Start() error {
	conn, err := net.Dial("tcp", c.dialAddr)
	if err != nil {
		fmt.Println("Failed to connect...")
		os.Exit(1)
	}

	defer conn.Close()

	go c.connectionLoop(conn)
	<-c.quitChan
	close(c.msgChan)
	return nil
}

/*
This loop maintains the communication between target and self
*/
func (c *Client) connectionLoop(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		conn.Write([]byte(text))
	}
}

func main() {
	fmt.Println("Go Client!")
	client := NewClient("TheBoys", ":3000")
	log.Fatal(client.Start())
}
