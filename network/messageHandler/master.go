package messageHandler

import (
	"base/common"
	"base/message"
	"fmt"
	"log"

	"capnproto.org/go/capnp/v3"
)

var ShardJoinChan = make(chan *common.ClusterJoinRequestMsg, 10)
var ChatMessageChan chan *common.Message

func ShardJoinRequest(msg *capnp.Message, sourceShard string) {
	fmt.Println("Handling ShardJoinRequest")
	rootMsg, err := common.ReadRootClusterJoinRequest(msg)
	if err != nil {
		log.Println("Failed to read game message:", err)
	}
	joinMsg := message.ParseClusterJoinMessage(rootMsg)
	joinMsg.SourceID = sourceShard
	ShardJoinChan <- joinMsg
	log.Printf("Received cluster join request: %s from: %s", joinMsg.Address, sourceShard)
}

func GameMessage(msg *capnp.Message, sourceShard string) {
	gameMsg, err := common.ReadRootGameMessage(msg)
	if err != nil {
		log.Println("Failed to read game message:", err)
	}

	switch gameMsg.Which() {
	case common.GameMessage_Which_chatMessage:
		msgObj := message.ParseChatMessage(gameMsg)
		log.Printf("Received chat message: %d %s", msgObj.PlayerID, msgObj.Text)
	}
}
