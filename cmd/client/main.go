package main

import (
	models "base/models"
	"base/network"
	"base/redis"
	"base/serialization"
	"math/rand/v2"

	"capnproto.org/go/capnp/v3"
)

func createMsg(source string) *capnp.Message {
	chatMessage, err := serialization.SerializeChatMessage(models.ChatMessage{PlayerID: rand.Int32N(20), Text: source})
	if err != nil {
		return nil
	}
	return chatMessage
}

func main() {
	redis := redis.NewClient()
	shardAddress, ok := redis.GetRandomFromSet("shard_addresses").(string)
	if !ok {
		return
	}
	client := network.NewClientHandler(shardAddress)
	if client == nil {
		return
	}
	client.Connect()
	client.ListenReliable()

	select {}
}
