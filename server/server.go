package main

import (
	"fmt"
	"log"
	"net"
)

/**
This will maintain
*/
type Message struct {
	sender  string
	payload []byte
}

/**
Server Struct
ln : the listener for the server
quitChan: the channel we use to kill the server (empty struct uses no memory)


// we can create a map of connections to save info about the user
*/
type Server struct {
	listenerAddr string
	ln           net.Listener
	quitChan     chan struct{}
	msgChan      chan Message
}

/**
We need to buffer the msgChan because if don't then that can introduce errors if
no one is actively listening to this channel
*/
func NewServer(listenAddr string) *Server {
	return &Server{
		listenerAddr: listenAddr,
		quitChan:     make(chan struct{}),
		msgChan:      make(chan Message, 10),
	}
}

// Start up our TCP server and wait for connections
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

		fmt.Printf("Server accepted new connection: %s\n", conn.RemoteAddr())
		// start a goroutine to avoid blocking incomming connections
		go s.readLoop(conn)
	}
}

// Every time we get a new connection, we need to spin up another loop
// so that we can read and write
func (s *Server) readLoop(conn net.Conn) {
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("Read Error: %s\n", err)
			// we could check for EOF here...
			continue
		}

		// send the message to the msgChan
		s.msgChan <- Message{
			sender:  conn.RemoteAddr().String(),
			payload: buff[:n],
		}

		conn.Write([]byte("message sent!"))
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
