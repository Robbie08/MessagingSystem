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
	dialAddr string
	dialer   net.Dialer
	quitChan chan struct{}
	msgChan  chan Message
}

/*
Constructor for our client
*/
func NewClient(dialAddr string) *Client {
	return &Client{
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
		c.sendMessage(conn)
		go c.readMessage(conn)
	}
}

func (c *Client) sendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> ")
	text, _ := reader.ReadString('\n')
	conn.Write([]byte(text))
}

func (c *Client) readMessage(conn net.Conn) {
	msg, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("\n->: %s\n>> ", msg)
}

func main() {
	fmt.Println("Go Client!")
	client := NewClient(":3000")
	log.Fatal(client.Start())
}
