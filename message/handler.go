package message

import (
	"base/common"
	"log"

	"capnproto.org/go/capnp/v3"
)

func newCapnpMessage() (*capnp.Message, *capnp.Segment, error) {
	data, segment, err := capnp.NewMessage(capnp.SingleSegment(nil))
	return data, segment, err
}

func CreateChatMessage(msg common.Message) (*capnp.Message, error) {
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

	chatMsg.SetPlayerID(msg.PlayerID)
	chatMsg.SetText(msg.Text)
	return data, nil
}

func ParseChatMessage(msg common.GameMessage) *common.Message {
	chatMessage, _ := msg.ChatMessage()
	id := chatMessage.PlayerID()
	text, _ := chatMessage.Text()

	return &common.Message{
		PlayerID: id,
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

	message.SetShardID(msg.ShardID)
	pos, err := message.NewPos()
	if err != nil {
		return nil, err
	}

	pos.SetX(msg.Pos.X)
	pos.SetY(msg.Pos.Y)
	return data, nil
}

func ParseClusterJoinResponseMessage(msg common.ClusterJoinResponse) *common.ClusterJoinResponseMsg {
	shardID, _ := msg.ShardID()
	pos, _ := msg.Pos()

	return &common.ClusterJoinResponseMsg{
		ShardID: shardID,
		Pos:     common.Pos{X: pos.X(), Y: pos.Y()},
	}
}

func CreateClientConnectionMsg(msg common.ClientConnectionRequestMsg) (*capnp.Message, error) {
	data, segment, err := newCapnpMessage()
	if err != nil {
		return nil, err
	}
	_, err = common.NewRootClientConnectionRequest(segment)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ParseClientConnectionRequestMsg(msg common.ClientConnectionRequest) *common.ClientConnectionRequestMsg {
	return &common.ClientConnectionRequestMsg{}
}

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
