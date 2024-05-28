package network

import (
	"base/message"
	"fmt"
	"sync"
)

type MasterHandler struct {
	*Server
	shardConnectionChan chan *connection
}

// Shard:
func NewMasterHandler() *MasterHandler {
	server, err := NewServer(ServerConfig{LocalAddr: "127.0.0.1:1234"})
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		return nil
	}
	return &MasterHandler{
		Server:              server,
		shardConnectionChan: make(chan *connection),
	}

}

func (server *MasterHandler) AcceptConnections() {
	for {
		conn, err := server.AcceptConnection()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			return
		}
		server.shardConnectionChan <- conn
	}
}

func (server *MasterHandler) AcceptStreams() {
	for conn := range server.shardConnectionChan {
		go func(conn *connection) {
			fmt.Println("starting conn goroutine")
			stream, err := conn.AcceptStream()
			if err != nil {
				fmt.Printf("Error accepting stream: %s\n", err)
				return
			}
			fmt.Println("accepted stream:")

			quicStream := quicStream{stream, &sync.Mutex{}}
			for {
				message.HandleGameMessage(read(stream))
				Respond(quicStream)
				// go server.handleStream(quicStream)
			}
		}(conn)

	}
}
