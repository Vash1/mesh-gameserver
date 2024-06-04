package network

import (
	"base/common"
	"base/message"
	"base/network/messageHandler"
	"fmt"
	"math/rand/v2"
	"sync"

	"capnproto.org/go/capnp/v3"
)

type ShardHandler struct {
	server               *NetServer
	MasterClient         *NetClient
	clientConnectionChan chan *connection
	PoolJoined           chan struct{}
}

func (c *NetClient) SendUnreliable(msg *capnp.Message) {
	bytes, _ := message.MsgToBytes(msg)
	c.quicConnection.SendDatagram(bytes)

}

func (h *NetClient) ListenReliable(PoolJoined chan struct{}) {
	handler := messageHandler.NewMessageHandler()
	handler.AddHandler(messageHandler.ShardJoinResponse)
	go func() {
		defer close(PoolJoined)
		<-messageHandler.ShardJoinResponseChan
		fmt.Println("Joined pool")
	}()
	for {
		msg, ok := read(h.QuicStream)
		if !ok {
			fmt.Println("Stream closed")
			return
		}
		handler.HandleMessage(msg, "master")
	}
}

func (c *NetClient) SendReliable(msg *capnp.Message) {
	c.QuicStream.mu.Lock()
	defer c.QuicStream.mu.Unlock()

	err := c.QuicStream.SendMessage(msg)
	if err != nil {
		fmt.Println("couldn't write to stream")
	}
}

func NewShardHandler() *ShardHandler {
	// Listen Locally
	server, err := NewNetServer(ServerConfig{LocalAddr: "127.0.0.1:0"})
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		return nil
	}
	return &ShardHandler{
		server:     server,
		PoolJoined: make(chan struct{}),
	}
}

func (h *ShardHandler) JoinCluster() error {
	client, err := NewNetClient(ClientConfig{RemoteAddr: "127.0.0.1:1234"})
	if err != nil {
		fmt.Printf("Error creating client: %s\n", err)
		return err
	}
	h.MasterClient = client
	return nil
}

func (h *ShardHandler) GetAddress() string {
	return h.server.quicListener.Addr().String()
}

func (h *ShardHandler) AcceptConnections() {
	h.clientConnectionChan = make(chan *connection)

	for {
		conn, err := h.server.AcceptConnection()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			return
		}
		h.clientConnectionChan <- conn
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
				msg, ok := read(quicStream)
				if !ok {
					fmt.Println("Stream closed")
					return
				}
				message.HandleGameMessage(msg)
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
	fmt.Println("responding")
	stream.SendMessage(chatMessage)

}

func initShardMessageHandlers() *messageHandler.MessageHandler {
	handler := messageHandler.NewMessageHandler()
	handler.AddHandler(messageHandler.ShardJoinResponse)
	handler.AddHandler(messageHandler.ClientGameMessage)
	return handler
}
