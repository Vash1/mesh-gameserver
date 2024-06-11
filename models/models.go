package models

type Vector struct {
	X float32
	Y float32
}

type Dimensions struct {
	Width  int32
	Height int32
}

var EmptyPos = Vector{X: 0, Y: 0}
var ZeroPos = Vector{}

func (a Vector) Add(b Vector) Vector {
	return Vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

type Player struct {
	Id int
}

type ServerData struct {
	Coords  Vector
	Players map[int]Player
}

type ChatMessage struct {
	PlayerID int32
	Text     string
	NetworkMessage
}

type NetworkMessage struct {
	SourceID string
}

type ClusterJoinRequest struct {
	Address string
	NetworkMessage
}
type ClusterJoinResponse struct {
	ShardID string
	Pos     Vector
}

type ClientConnectionRequest struct {
	NetworkMessage
}

type ClientConnectionResponse struct {
	ClientID string
	Position Vector
	MapData  MapData
}

type MapData struct {
	Size Dimensions
}
