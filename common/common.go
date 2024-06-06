package common

type Pos struct {
	X int32
	Y int32
}

var EmptyPos = Pos{X: 0, Y: 0}
var ZeroPos = Pos{}

func (a Pos) Add(b Pos) Pos {
	return Pos{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

type Player struct {
	Id int
}

type ServerData struct {
	Coords  Pos
	Players map[int]Player
}

type Message struct {
	PlayerID int32
	Text     string
	NetworkMessage
}

type NetworkMessage struct {
	SourceID string
}

type ClusterJoinRequestMsg struct {
	Address string
	NetworkMessage
}
type ClusterJoinResponseMsg struct {
	ShardID string
	Pos     Pos
}

type ClientConnectionRequestMsg struct {
}
