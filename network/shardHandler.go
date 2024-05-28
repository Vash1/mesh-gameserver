package network

import (
	"base/common"
	"base/message"
	"fmt"
	"math/rand/v2"
	"sync"

	"capnproto.org/go/capnp/v3"
)

type ShardHandler struct {
	*Server
	clientConnectionChan chan *connection
}

type MasterConn struct {
	*Client
}

func NewMasterConn() *MasterConn {
	client, err := NewClient(ClientConfig{RemoteAddr: "127.0.0.1:1234"})
	if err != nil {
		fmt.Printf("Error creating client: %s\n", err)
		return nil
	}
	return &MasterConn{
		Client: client,
	}
}

func (client *MasterConn) SendUnreliable(msg *capnp.Message) {
	bytes, _ := message.MsgToBytes(msg)
	client.quicConnection.SendDatagram(bytes)

}

func (client *MasterConn) ListenReliable() {
	go client.Listen()
	// decode(quic.Stream) -> capnp.Message
	// HandleGameMessage(capnp.Message, channels)
}

func (client *MasterConn) SendReliable(msg *capnp.Message) {
	client.QuicStream.mu.Lock()
	defer client.QuicStream.mu.Unlock()

	err := client.QuicStream.SendMessage(msg)
	if err != nil {
		fmt.Println("couldn't write to stream")
	}

	fmt.Println("msg sent")
}

// Shard:
func NewShardHandler() *ShardHandler {
	server, err := NewServer(ServerConfig{LocalAddr: "127.0.0.1:0"})
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		return nil
	}
	return &ShardHandler{
		Server:               server,
		clientConnectionChan: make(chan *connection),
	}

}

func (server *ShardHandler) GetAddress() string {
	return server.quicListener.Addr().String()
}

func (server *ShardHandler) AcceptConnections() {
	for {
		conn, err := server.AcceptConnection()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			return
		}
		server.clientConnectionChan <- conn
	}
}
func (server *ShardHandler) AcceptData() {
	for conn := range server.clientConnectionChan {
		go func(conn *connection) {
			fmt.Println("starting stream goroutine")
			stream, err := conn.AcceptStream()
			if err != nil {
				fmt.Printf("Error accepting stream: %s\n", err)
				return
			}

			quicStream := quicStream{stream, &sync.Mutex{}}
			for {
				message.HandleGameMessage(read(quicStream))
				Respond(quicStream)
				// go server.handleStream(quicStream)
			}
		}(conn)

		go func(conn *connection) {
			fmt.Println("starting datagram goroutine")
			for {
				bytes, err := conn.receiveDatagram()
				if err != nil {
					fmt.Printf("Error reading datagram: %s\n", err)
					return
				}
				msg, _ := message.BytesToMsg(bytes)
				message.HandleGameMessage(msg)
			}
		}(conn)
	}
}

func (server *ShardHandler) AcceptDatagrams() {
	for conn := range server.clientConnectionChan {
		go func(conn *connection) {
			fmt.Println("starting datagram goroutine")

		}(conn)
	}
}

func Respond(stream quicStream) {
	chatMessage, err := message.CreateChatMessage(common.Message{PlayerId: rand.Int32N(20), Text: "stream"})
	if err != nil {
		return
	}
	stream.SendMessage(chatMessage)

}

func Setup2() {
	// generateTLSConfig
	// ResolveLocalUDP -> localaddr net.UDPAddr
	// * requirement from local udp conn
	// ListenUDP -> conn net.UDPConn
	// * requirement from transport
	// Setup QUIC Transport -> quic.Transport
	// * requirement from transport.Listen
	// Listen -> quic.Listener
	// * requirement from Accept
	// Accept -> quic.Connection
	// * requirement for connection.AcceptStream
	// * requirement for connection.ReceiveDatagram
	// * requirement for connection.SendDatagram
	// AcceptStream -> quic.Stream
	// * requirement for decode(quic.Stream)
	// * requirement for stream.WriteStream
}
