package message

import (
	"base/common"
	"log"

	"capnproto.org/go/capnp/v3"
)

// type networkChatMessage struct {
// 	playerId int32
// 	common.Message
// }

func newCapnpMessage() (*capnp.Message, *capnp.Segment, error) {
	data, segment, err := capnp.NewMessage(capnp.SingleSegment(nil))
	return data, segment, err
}

// func (c *Client) CreateSimpleChatMessage(msg common.Message) (*capnp.Message, error) {
// 	// networkMsg := &networkChatMessage{
// 	// 	playerId: c.id,
// 	// 	Message:  msg,
// 	// }

// 	data, segment, err := newCapnpMessage()
// 	if err != nil {
// 		return nil, err
// 	}

// 	chatMsg, err := common.NewRootChatMessage(segment)
// 	if err != nil {
// 		return nil, err
// 	}

// 	chatMsg.SetPlayerId(c.id)
// 	chatMsg.SetText(msg.Text)
// 	return data, nil
// }

func CreateChatMessage(msg common.Message) (*capnp.Message, error) {
	// networkMsg := &networkChatMessage{
	// 	playerId: c.id,
	// 	Message:  msg,
	// }

	data, segment, err := newCapnpMessage()
	if err != nil {
		return nil, err
	}
	message, err := common.NewRootGameMessage(segment)
	if err != nil {
		return nil, err
	}

	chatMsg, err := message.NewChatMessage()
	if err != nil {
		return nil, err
	}

	chatMsg.SetPlayerId(msg.PlayerId)
	chatMsg.SetText(msg.Text)
	return data, nil
}

// func HandleGameMessage(msg *capnp.Message, channels *Event, source string) {
// 	gameMsg, err := common.ReadRootGameMessage(msg)
// 	if err != nil {
// 		log.Println("Failed to read game message:", err)
// 	}

// 	switch gameMsg.Which() {
// 	// case common.GameMessage_Which_playerMove:
// 	// 	...
// 	// case common.GameMessage_Which_playerAction:
// 	// 	...
// 	// case common.GameMessage_Which_gameStateUpdate:
// 	// 	...
// 	case common.GameMessage_Which_chatMessage:
// 		msg := ParseChatMessage(gameMsg)
// 		log.Printf("Received %s chat message: %d %s", source, msg.PlayerId, msg.Text)
// 		channels.chat <- msg
// 	}
// }

func HandleGameMessage(msg *capnp.Message) {
	log.Printf("Received chat message:")
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
		msg := ParseChatMessage(gameMsg)
		log.Printf("Received chat message: %d %s", msg.PlayerId, msg.Text)
	}
}

func ParseChatMessage(msg common.GameMessage) *common.Message {
	chatMessage, _ := msg.ChatMessage()
	id := chatMessage.PlayerId()
	text, _ := chatMessage.Text()

	return &common.Message{
		PlayerId: id,
		Text:     text,
	}
}
func CreateClusterJoinMessage(msg common.ClusterJoinRequestMsg) (*capnp.Message, error) {
	data, segment, err := newCapnpMessage()
	if err != nil {
		return nil, err
	}
	message, err := common.NewRootClusterJoinRequest(segment)
	if err != nil {
		return nil, err
	}

	message.SetAddress(msg.Address)
	return data, nil
}

func ParseClusterJoinMessage(msg common.ClusterJoinRequest) *common.ClusterJoinRequestMsg {
	address, _ := msg.Address()

	return &common.ClusterJoinRequestMsg{
		Address: address,
	}
}

func CreateClusterJoinResponseMessage(msg *common.ClusterJoinResponseMsg) (*capnp.Message, error) {
	data, segment, err := newCapnpMessage()
	if err != nil {
		return nil, err
	}
	message, err := common.NewRootClusterJoinResponse(segment)
	if err != nil {
		return nil, err
	}

	message.SetShardId(msg.ShardID)
	pos, err := message.NewPos()
	if err != nil {
		return nil, err
	}

	pos.SetX(msg.Pos.X)
	pos.SetY(msg.Pos.Y)
	return data, nil
}

func ParseClusterJoinResponseMessage(msg common.ClusterJoinResponse) *common.ClusterJoinResponseMsg {
	shardID, _ := msg.ShardId()
	pos, _ := msg.Pos()

	return &common.ClusterJoinResponseMsg{
		ShardID: shardID,
		Pos:     common.Pos{X: pos.X(), Y: pos.Y()},
	}
}

// 	data, segment, err := newCapnpMessage()
// 	if err != nil {
// 		return nil, err
// 	}
// 	joinMsg, err := common.New(segment)
// 	if err != nil {
// 		return nil, err
// 	}

// 	joinMsg.SetAddress("1")
// 	return data, nil
// }

func MsgToBytes(msg *capnp.Message) ([]byte, bool) {
	bytes, err := msg.Marshal()
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return nil, false
	}
	return bytes, true
}

func BytesToMsg(bytes []byte) (*capnp.Message, bool) {
	msg, err := capnp.Unmarshal(bytes)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return nil, false
	}
	return msg, true
}
