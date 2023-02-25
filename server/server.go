package main

import (
	"fmt"
	"log"
	"net"
)

/*
This will maintain contain info about each message
*/
type Message struct {
	sender  string
	payload []byte
}

/*
Server Struct
This will represent a "chat room" for now
*/
type Server struct {
	listenerAddr string        //this is the address we would like to spin up the server on
	ln           net.Listener  // the listener for the server
	quitChan     chan struct{} // the channel we use to kill the server (empty struct uses no memory)
	msgChan      chan Message  // we will send messages through this channel
}

/*
Constructor for the server
We need to buffer the msgChan because if we don't then that can
introduce errors if no one is actively listening to this channel
*/
func NewServer(listenAddr string) *Server {
	return &Server{
		listenerAddr: listenAddr,
		quitChan:     make(chan struct{}),
		msgChan:      make(chan Message, 10),
	}
}

// Function to start up our TCP server and wait for connections
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenerAddr)
	if err != nil {
		return err
	}

	defer ln.Close() // let's close the listener when done
	s.ln = ln

	go s.acceptLoop()

	<-s.quitChan     // wait for the quit channel here
	close(s.msgChan) // let liteners know we are done with this channel
	return nil
}

// Function to accept connections
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()

		// this err should no trigger an exit
		if err != nil {
			fmt.Printf("Accept error: %s\n", err)
			continue
		}

		fmt.Printf("New connection joined the server: %s\n", conn.RemoteAddr())
		// start a goroutine to avoid blocking incomming connections
		go s.readLoop(conn)
	}
}

/*
Every time we get a new connection, we need to spin up another loop
so that we can read and write. This will prevent blocking others
*/
func (s *Server) readLoop(conn net.Conn) {
	buff := make([]byte, 2048)
read:
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("-Error-: %s\n", err)
			// close the connection if we have an error
			break read
		}

		// send the message to the msgChan
		s.msgChan <- Message{
			sender:  conn.RemoteAddr().String(),
			payload: buff[:n],
		}

		conn.Write([]byte("Message Received...\n"))
	}
}

func main() {
	fmt.Println("Go Server!")
	server := NewServer(":3000")

	// we can have our main func print all the messages that come into this channel
	go func() {
		for msg := range server.msgChan {
			fmt.Printf("(%s): %s\n", string(msg.sender), string(msg.payload))
		}
	}()
	log.Fatal(server.Start())
}
