package network

import (
	"base/message"
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
	// ResolveLocalUDP -> localaddr net.UDPAddr
	// * requirement from local udp conn
	// Listen Local UDP -> conn net.UDPConn
	// * requirement from transport
	// Setup QUIC Transport -> quic.Transport
	// * requirement from Dial
	// ResolveRemoteUDP -> remoteaddr net.UDPAddr
	// * requirement from Dial
	// NewConnection{Dial()} -> quic.Connection
	// * requirement from send/receive datagrams
	// * requirement from OpenStream
	// OpenStream -> quic.stream
}

// clienthandler.SendReliable()
func ListenUnreliable() {
	// quic.Connection.ReceiveDatagram
	// Unmarshal
	// filter if timestamp is older than most recent message
	// HandleGameMessage(channels)
}

func (client *ClientHandler) SendUnreliable(msg *capnp.Message) {
	bytes, _ := message.MsgToBytes(msg)
	client.quicConnection.SendDatagram(bytes)

}

func (client *ClientHandler) ListenReliable() {
	go client.Listen()
	// decode(quic.Stream) -> capnp.Message
	// HandleGameMessage(capnp.Message, channels)
}

func (client *ClientHandler) SendReliable(msg *capnp.Message) {
	client.QuicStream.mu.Lock()
	defer client.QuicStream.mu.Unlock()

	err := client.QuicStream.SendMessage(msg)
	if err != nil {
		fmt.Println("couldn't write to stream")
	}

	fmt.Println("msg sent")
}

// logic splits here for handling client side messages via channels
