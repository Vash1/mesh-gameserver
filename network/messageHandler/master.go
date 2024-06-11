package messageHandler

import (
	capnp_schema "base/capnp"
	models "base/models"
	"base/serialization"
	"fmt"
	"log"

	"capnproto.org/go/capnp/v3"
)

var ShardJoinChan = make(chan *models.ClusterJoinRequest, 10)
var ChatMessageChan chan *models.ChatMessage

func ShardJoinRequest(msg *capnp.Message, sourceShard string) {
	fmt.Println("Handling ShardJoinRequest")
	rootMsg, err := capnp_schema.ReadRootClusterJoinRequest(msg)
	if err != nil {
		log.Println("Failed to read game message:", err)
	}
	joinMsg := serialization.DeserializeClusterJoinRequest(rootMsg)
	joinMsg.SourceID = sourceShard
	ShardJoinChan <- joinMsg
	log.Printf("Received cluster join request: %s from: %s", joinMsg.Address, sourceShard)
}

func GameMessage(msg *capnp.Message, sourceShard string) {
	gameMsg, err := capnp_schema.ReadRootGameMessage(msg)
	if err != nil {
		log.Println("Failed to read game message:", err)
	}

	switch gameMsg.Which() {
	case capnp_schema.GameMessage_Which_chatMessage:
		msgObj := serialization.DeserializeChatMessage(gameMsg)
		log.Printf("Received chat message: %d %s", msgObj.PlayerID, msgObj.Text)
	}
}
