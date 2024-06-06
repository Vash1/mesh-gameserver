using Go = import "/go.capnp";
@0xb81e0a7d72638bf1;
$Go.package("capnp");
$Go.import("base/capnp");

struct GameMessage {
  union {
    playerMove @0 :PlayerMove;
    playerAction @1 :PlayerAction;
    gameStateUpdate @2 :GameStateUpdate;
    chatMessage @3 :ChatMessage;
  }
}

struct Vector {
  x @0 :Int32;
  y @1 :Int32;
}

struct PlayerMove {
  playerID @0 :Int32;
  position @1 :Vector;
}

struct PlayerAction {
  playerID @0 :Int32;
  action @1 :Text;
}

struct GameStateUpdate {
  playerStates @0 :List(PlayerState);
}

struct PlayerState {
  playerID @0 :Int32;
  position @1 :Vector;
}

struct ChatMessage {
  playerID @0 :Int32;
  text @1 :Text;
}
#
#struct ServerMessage {
#  union {
#    clusterJoin @0 :ClusterJoin;
#  }
#}

struct ClusterJoinRequest {
  address @0 :Text;
}

struct ClusterJoinResponse {
  shardID @0 :Text;
  pos @1 :Vector;
}

struct MapData {
  size @0 :Vector;
}

struct ClientConnectionRequest {
}

struct ClientConnectionResponse {
  clientID @0 :Text;
  position @1 :Vector;
  mapData @2 :MapData;
}