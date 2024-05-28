package main

import (
	"base/common"
	"base/message"
	"base/network"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"capnproto.org/go/capnp/v3"
)

func createMsg(source string) *capnp.Message {
	chatMessage, err := message.CreateChatMessage(common.Message{PlayerId: rand.Int32N(20), Text: source})
	if err != nil {
		return nil
	}
	return chatMessage
}
func main() {
	// redis := redis.NewClient()
	// redis.Get("key")
	// redis.Set("key", "t123est")
	// redis.Get("key")

	// _ = rdb
	/*
		// for shard client
		datagramListener = server.listen(datagramConfig)
		  (sends constant game state updates to all clients)
		chatListener = server.listen(streamConfig)
		eventListener = server.listen(streamConfig)
		  (receives player actions and events)

		streams:


	*/
	// server := network.NewServer(true, true, true, true)
	// go server.BroadcastChat()
	// server.AcceptConnections()

	shard := network.NewShardHandler()
	fmt.Println("Shard network handler created", shard.GetAddress())
	master := network.NewMasterConn()
	if master == nil {
		return
	}
	master.OpenStream()
	msg, err := message.CreateClusterJoinMessage(common.ClusterJoinRequestMsg{Address: shard.GetAddress()})
	if err != nil {
		fmt.Println("Error creating cluster join message")
		return
	}
	master.SendReliable(msg)

	var wg = &sync.WaitGroup{}
	wg.Add(1)
	// chatMessage := createMsg()
	// client.SendReliable(chatMessage)
	// fmt.Println("Data1 sent successfully.")
	// time.Sleep(1 * time.Second)

	for i := 0; i < 1; i++ {
		go func() {
			chatMessage := createMsg("stream")
			unrealiableChatMessage := createMsg("datagram")
			for j := 0; j < 2; j++ {
				master.SendReliable(chatMessage)
				master.SendUnreliable(unrealiableChatMessage)
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// chatMessage = createMsg()
	// client.SendReliable(chatMessage)
	fmt.Println("Data sent successfully.")

	// wg.Wait()

	go shard.AcceptConnections()
	go shard.AcceptData()
	// go shard.AcceptStreams()
	// go shard.AcceptDatagrams()
	time.Sleep(20 * time.Second)

}
