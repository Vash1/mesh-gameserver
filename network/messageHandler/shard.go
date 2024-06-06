package messageHandler

import (
	capnp_schema "base/capnp"
	models "base/models"
	"base/serialization"
	"fmt"
	"log"

	"capnproto.org/go/capnp/v3"
)

var ShardJoinResponseChan = make(chan *models.ClusterJoinResponse, 10)
var ClientConnectionRequestChan = make(chan *models.ClientConnectionRequest, 10)

func ShardJoinResponse(msg *capnp.Message, source string) {
	fmt.Println("Handling ShardJoinResponse")
	rootMsg, err := capnp_schema.ReadRootClusterJoinResponse(msg)
	if err != nil {
		log.Println("Failed to read ShardJoinResponse:", err)
	}
	joinMsg := serialization.DeserializeClusterJoinResponse(rootMsg)
	log.Printf("shard %s is located at: %s", joinMsg.ShardID, fmt.Sprint(joinMsg.Pos))
	ShardJoinResponseChan <- joinMsg
}

func ClientGameMessage(msg *capnp.Message, source string) {
	gameMsg, err := capnp_schema.ReadRootGameMessage(msg)
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
	case capnp_schema.GameMessage_Which_chatMessage:
		msgObj := serialization.DeserializeChatMessage(gameMsg)
		log.Printf("Received chat message: %d %s", msgObj.PlayerID, msgObj.Text)
	}
}

func ClientConnectionRequest(msg *capnp.Message, source string) {
	fmt.Println("Handling ClientConnectionRequest")
	rootMsg, err := capnp_schema.ReadRootClientConnectionRequest(msg)
	if err != nil {
		log.Println("Failed to read ClientConnectionRequest:", err)
	}
	connMsg := serialization.DeserializeClientConnectionRequest(rootMsg)
	ClientConnectionRequestChan <- connMsg
	log.Printf("Received client connection request: %s", source)
}
