package main

import (
	"base/common"
	"base/message"
	"base/network"
	"base/redis"
	"math/rand/v2"

	"capnproto.org/go/capnp/v3"
)

func createMsg(source string) *capnp.Message {
	chatMessage, err := message.CreateChatMessage(common.Message{PlayerID: rand.Int32N(20), Text: source})
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
	// client.ListenReliable()

	// for i := 0; i < 5; i++ {
	// 	go func() {
	// 		chatMessage := createMsg("stream")
	// 		unrealiableChatMessage := createMsg("datagram")
	// 		for j := 0; j < 10; j++ {
	// 			client.SendReliable(chatMessage)
	// 			client.SendUnreliable(unrealiableChatMessage)
	// 			time.Sleep(1 * time.Second)
	// 		}
	// 	}()
	// }

	select {}
}
