package network

import (
	models "base/models"
	"base/network/messageHandler"
	"base/redis"
	"base/serialization"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type MasterHandler struct {
	*NetServer
	shardConnectionChan chan *connection
	shardConnections    map[string]*ShardConnection
}

type ShardConnection struct {
	address string
	id      string
	stream  *quicStream
}

var redisClient *redis.Redis = redis.NewClient()

// Shard:
func NewMasterHandler() *MasterHandler {
	server, err := NewNetServer(ServerConfig{LocalAddr: "127.0.0.1:1234"})
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		return nil
	}
	return &MasterHandler{
		NetServer:           server,
		shardConnectionChan: make(chan *connection),
		shardConnections:    make(map[string]*ShardConnection),
	}
}

func (server *MasterHandler) ResetRedisShards() {
	redisClient.DeleteAll(REDIS_SHARD_KEY)
}

func (server *MasterHandler) AcceptConnections() {
	for {
		conn, err := server.AcceptConnection()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			return
		}
		server.shardConnectionChan <- conn
	}
}

func (server *MasterHandler) addToPool(shard *ShardConnection) {
	redisClient.AddToSet(REDIS_SHARD_KEY, shard.address)
}

func (server *MasterHandler) removeFromPool(shard *ShardConnection) {
	redisClient.RemoveFromSet(REDIS_SHARD_KEY, shard.address)
	delete(server.shardConnections, shard.id)
}

func (server *MasterHandler) BroadcastChat() {
	for {
		for _, shard := range server.shardConnections {
			Respond(*shard.stream)
		}
	}
}

func (server *MasterHandler) AcceptStreams() {
	go func() {
		for shardJoined := range messageHandler.ShardJoinChan {
			sourceShard := server.shardConnections[shardJoined.SourceID]
			sourceShard.address = shardJoined.Address
			server.addToPool(sourceShard)
			shardJoinResponse, _ := serialization.SerializeClusterJoinResponse(&models.ClusterJoinResponse{ShardID: sourceShard.id, Pos: models.Vector{X: 3, Y: 5}})
			sourceShard.stream.SendMessage(shardJoinResponse)
		}
	}()

	for conn := range server.shardConnectionChan {
		go func(conn *connection) {
			msgHandler := initMessageHandlers()
			stream, err := conn.AcceptStream()
			if err != nil {
				fmt.Printf("Error accepting stream: %s\n", err)
				return
			}
			quicStream := quicStream{stream, &sync.Mutex{}}

			shardConn := ShardConnection{
				id:     uuid.New().String(),
				stream: &quicStream,
			}
			server.shardConnections[shardConn.id] = &shardConn

			for {
				msg, ok := read(quicStream)
				if !ok {
					server.removeFromPool(&shardConn)
					fmt.Println("Stream closed")
					return
				}
				msgHandler.HandleMessage(msg, shardConn.id)
				// Respond(quicStream)
			}
		}(conn)
	}
}

func initMessageHandlers() *messageHandler.MessageHandler {
	handler := messageHandler.NewMessageHandler()
	handler.AddHandler(messageHandler.ShardJoinRequest)
	handler.AddHandler(messageHandler.GameMessage)
	return handler
}
