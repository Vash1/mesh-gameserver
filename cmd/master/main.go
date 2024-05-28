package main

import (
	"base/common"
	"base/message"
	"base/network"
	"fmt"
	"math/rand/v2"
	"time"

	"capnproto.org/go/capnp/v3"
)

func createMsg() *capnp.Message {
	chatMessage, err := message.CreateChatMessage(common.Message{PlayerId: rand.Int32N(20), Text: "stream"})
	if err != nil {
		return nil
	}
	return chatMessage
}

func main() {

	shard := network.NewMasterHandler()
	fmt.Println("Shard network handler created", shard)
	go shard.AcceptConnections()
	go shard.AcceptStreams()
	time.Sleep(20 * time.Second)

}
