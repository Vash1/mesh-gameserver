package messageHandler

import (
	"base/common"
	"base/message"
	"fmt"
	"log"

	"capnproto.org/go/capnp/v3"
)

var ShardJoinResponseChan = make(chan *common.ClusterJoinResponseMsg, 10)

func ShardJoinResponse(msg *capnp.Message, source string) {
	fmt.Println("Handling ShardJoinResponse")
	rootMsg, err := common.ReadRootClusterJoinResponse(msg)
	if err != nil {
		log.Println("Failed to read ShardJoinResponse:", err)
	}
	joinMsg := message.ParseClusterJoinResponseMessage(rootMsg)
	log.Printf("shard %s is located at: %s", joinMsg.ShardID, fmt.Sprint(joinMsg.Pos))
	ShardJoinResponseChan <- joinMsg
}

func ClientGameMessage(msg *capnp.Message, source string) {
	gameMsg, err := common.ReadRootGameMessage(msg)
	if err != nil {
		log.Println("Failed to read game message:", err)
	}

	switch gameMsg.Which() {
	// case common.GameMessage_Which_playerMove:
	// 	...
	// case common.GameMessage_Which_playerAction:
	// 	...
	// case common.GameMessage_Which_gameStateUpdate:
	// 	...
	case common.GameMessage_Which_chatMessage:
		msgObj := message.ParseChatMessage(gameMsg)
		log.Printf("Received chat message: %d %s", msgObj.PlayerId, msgObj.Text)
	}
}
