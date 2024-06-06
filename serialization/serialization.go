package serialization

import (
	capnp_schema "base/capnp"
	models "base/models"
	"log"

	"capnproto.org/go/capnp/v3"
)

func newCapnpMessage() (*capnp.Message, *capnp.Segment, error) {
	data, segment, err := capnp.NewMessage(capnp.SingleSegment(nil))
	return data, segment, err
}

func SerializeChatMessage(msg models.ChatMessage) (*capnp.Message, error) {
	data, segment, err := newCapnpMessage()
	if err != nil {
		return nil, err
	}
	message, err := capnp_schema.NewRootGameMessage(segment)
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

func DeserializeChatMessage(msg capnp_schema.GameMessage) *models.ChatMessage {
	chatMessage, _ := msg.ChatMessage()
	id := chatMessage.PlayerID()
	text, _ := chatMessage.Text()

	return &models.ChatMessage{
		PlayerID: id,
		Text:     text,
	}
}
func SerializeClusterJoinRequest(msg models.ClusterJoinRequest) (*capnp.Message, error) {
	data, segment, err := newCapnpMessage()
	if err != nil {
		return nil, err
	}
	message, err := capnp_schema.NewRootClusterJoinRequest(segment)
	if err != nil {
		return nil, err
	}

	message.SetAddress(msg.Address)
	return data, nil
}

func DeserializeClusterJoinRequest(msg capnp_schema.ClusterJoinRequest) *models.ClusterJoinRequest {
	address, _ := msg.Address()

	return &models.ClusterJoinRequest{
		Address: address,
	}
}

func SerializeClusterJoinResponse(msg *models.ClusterJoinResponse) (*capnp.Message, error) {
	data, segment, err := newCapnpMessage()
	if err != nil {
		return nil, err
	}
	message, err := capnp_schema.NewRootClusterJoinResponse(segment)
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

func DeserializeClusterJoinResponse(msg capnp_schema.ClusterJoinResponse) *models.ClusterJoinResponse {
	shardID, _ := msg.ShardID()
	pos, _ := msg.Pos()

	return &models.ClusterJoinResponse{
		ShardID: shardID,
		Pos:     models.Vector{X: pos.X(), Y: pos.Y()},
	}
}

func SerializeClientConnectionRequest(msg models.ClientConnectionRequest) (*capnp.Message, error) {
	data, segment, err := newCapnpMessage()
	if err != nil {
		return nil, err
	}
	_, err = capnp_schema.NewRootClientConnectionRequest(segment)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DeserializeClientConnectionRequest(msg capnp_schema.ClientConnectionRequest) *models.ClientConnectionRequest {
	return &models.ClientConnectionRequest{}
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
