package models

type Vector struct {
	X int32
	Y int32
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
}

type ClientConnectionResponse struct {
	clientID string
	position Vector
	mapData  MapData
}

type MapData struct {
	size Vector
}
