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

func createMsg(source string) *capnp.Message {
	chatMessage, err := message.CreateChatMessage(common.Message{PlayerId: rand.Int32N(20), Text: source})
	if err != nil {
		return nil
	}
	return chatMessage
}
func main() {
	// redis := redis.NewClient()

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
	shard.JoinCluster()
	if shard.MasterClient == nil {
		return
	}
	shard.MasterClient.OpenStream()
	go shard.MasterClient.ListenReliable()
	msg, err := message.CreateClusterJoinMessage(common.ClusterJoinRequestMsg{Address: shard.GetAddress()})
	if err != nil {
		fmt.Println("Error creating cluster join message")
		return
	}
	shard.MasterClient.SendReliable(msg)

	for i := 0; i < 5; i++ {
		go func() {
			chatMessage := createMsg("stream")
			unrealiableChatMessage := createMsg("datagram")
			for j := 0; j < 2; j++ {
				shard.MasterClient.SendReliable(chatMessage)
				shard.MasterClient.SendUnreliable(unrealiableChatMessage)
				fmt.Println("Data sent successfully.")
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// chatMessage = createMsg()
	// client.SendReliable(chatMessage)

	// wg.Wait()

	go shard.AcceptConnections()
	go shard.AcceptData()
	// go shard.AcceptStreams()
	// go shard.AcceptDatagrams()
	select {}

}
