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
	client := network.NewClientHandler("127.0.0.1:32985")
	if client == nil {
		return
	}
	client.Connect()
	// client.ListenReliable()
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	// chatMessage := createMsg()
	// client.SendReliable(chatMessage)
	// fmt.Println("Data1 sent successfully.")
	// time.Sleep(1 * time.Second)

	for i := 0; i < 5; i++ {
		go func() {
			chatMessage := createMsg("stream")
			unrealiableChatMessage := createMsg("datagram")
			for j := 0; j < 10; j++ {
				client.SendReliable(chatMessage)
				client.SendUnreliable(unrealiableChatMessage)
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// chatMessage = createMsg()
	// client.SendReliable(chatMessage)
	fmt.Println("Data sent successfully.")

	wg.Wait()
}
