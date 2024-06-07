package network

import (
	models "base/models"
	"base/network/messageHandler"
	"base/serialization"
	"fmt"

	"capnproto.org/go/capnp/v3"
)

type ClientHandler struct {
	*NetClient
}

func NewClientHandler(addr string) *ClientHandler {
	client, err := NewNetClient(ClientConfig{RemoteAddr: addr})
	if err != nil {
		fmt.Printf("Error creating client: %s\n", err)
		return nil
	}
	return &ClientHandler{
		NetClient: client,
	}
}

func (client *ClientHandler) Connect() {
	client.OpenStream()
	msg, err := serialization.SerializeClientConnectionRequest(models.ClientConnectionRequest{})
	if err != nil {
		fmt.Println("Error creating message")
		return
	}
	client.SendReliable(msg)
}

func ListenUnreliable() {
}

func (client *ClientHandler) SendUnreliable(msg *capnp.Message) {
	bytes, _ := serialization.MsgToBytes(msg)
	client.quicConnection.SendDatagram(bytes)

}

func (client *ClientHandler) ListenReliable() {
	msgHandler := messageHandler.NewMessageHandler()
	msgHandler.AddHandler(messageHandler.ClientConnectionResponse)
	// go client.Listen()
	for {
		msg, ok := read(client.QuicStream)
		if !ok {
			fmt.Println("Stream closed")
			return
		}
		msgHandler.HandleMessage(msg, "shard")
	}
}

func (client *ClientHandler) SendReliable(msg *capnp.Message) {
	client.QuicStream.mu.Lock()
	defer client.QuicStream.mu.Unlock()

	err := client.QuicStream.SendMessage(msg)
	if err != nil {
		fmt.Println("couldn't write to stream")
	}
}
