package network

import (
	models "base/models"
	"base/network/messageHandler"
	"base/serialization"
	"fmt"
	"math/rand/v2"
	"sync"

	"capnproto.org/go/capnp/v3"
	"github.com/google/uuid"
)

type ShardHandler struct {
	server               *NetServer
	MasterClient         *NetClient
	clientConnectionChan chan *connection
	clientConnections    map[string]*ClientConnection
	PoolJoined           chan struct{}
}

type ClientConnection struct {
	id         string
	stream     *quicStream
	connection *connection
	isActive   bool
}

func (c *NetClient) SendUnreliable(msg *capnp.Message) {
	bytes, _ := serialization.MsgToBytes(msg)
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
	c.QuicStream.SendReliable(msg)
}

func (s *quicStream) SendReliable(msg *capnp.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.SendMessage(msg)
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
		server:               server,
		clientConnectionChan: make(chan *connection),
		clientConnections:    make(map[string]*ClientConnection),
		PoolJoined:           make(chan struct{}),
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
	go func() {
		for newClient := range messageHandler.ClientConnectionRequestChan {
			client := server.clientConnections[newClient.SourceID]
			client.isActive = true
			response := models.ClientConnectionResponse{ClientID: client.id, Position: models.Vector{X: rand.Float32() * 100, Y: rand.Float32() * 100}, MapData: models.MapData{Size: models.Dimensions{Width: 10, Height: 10}}}
			msg, err := serialization.SerializeClientConnectionResponse(&response)
			if err != nil {
				fmt.Println("Error serializing client connection response")
				return
			}
			client.stream.SendReliable(msg)
		}
	}()

	for conn := range server.clientConnectionChan {
		go func(conn *connection) {
			msgHandler := messageHandler.NewMessageHandler()
			msgHandler.AddHandler(messageHandler.ClientConnectionRequest)
			stream, err := conn.AcceptStream()
			if err != nil {
				fmt.Printf("Error accepting stream: %s\n", err)
				return
			}
			quicStream := quicStream{stream, &sync.Mutex{}}

			clientConn := ClientConnection{
				id:         uuid.New().String(),
				stream:     &quicStream,
				connection: conn,
			}
			server.clientConnections[clientConn.id] = &clientConn

			for {
				msg, ok := read(quicStream)
				if !ok {
					fmt.Println("Stream closed")
					delete(server.clientConnections, clientConn.id)
					return
				}
				msgHandler.HandleMessage(msg, clientConn.id)
			}
		}(conn)

		go func(conn *connection) {
			for {
				bytes, err := conn.receiveDatagram()
				if err != nil {
					// fmt.Printf("Error reading datagram: %s\n", err)
					return
				}
				msg, _ := serialization.BytesToMsg(bytes)
				_ = msg // TODO: handle datagrams
			}
		}(conn)
	}
}

func Respond(stream quicStream) {
	chatMessage, err := serialization.SerializeChatMessage(models.ChatMessage{PlayerID: rand.Int32N(20), Text: "stream"})
	if err != nil {
		return
	}
	fmt.Println("responding")
	stream.SendMessage(chatMessage)

}
