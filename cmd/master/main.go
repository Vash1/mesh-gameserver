package main

import (
	models "base/models"
	"base/network"
	"base/serialization"
	"math/rand/v2"

	"capnproto.org/go/capnp/v3"
)

func createMsg() *capnp.Message {
	chatMessage, err := serialization.SerializeChatMessage(models.ChatMessage{PlayerID: rand.Int32N(20), Text: "stream"})
	if err != nil {
		return nil
	}
	return chatMessage
}

func main() {
	master := network.NewMasterHandler()
	master.ResetRedisShards()
	go master.AcceptConnections()
	// redis := redis.NewClient()
	go master.AcceptStreams()
	// go master.BroadcastChat()
	select {}

}
