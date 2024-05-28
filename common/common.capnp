using Go = import "/go.capnp";
@0xb81e0a7d72638bf1;
$Go.package("common");
$Go.import("base/common");

struct GameMessage {
  union {
    playerMove @0 :PlayerMove;
    playerAction @1 :PlayerAction;
    gameStateUpdate @2 :GameStateUpdate;
    chatMessage @3 :ChatMessage;
  }
}

struct Position {
    union {
        x @0 :Int32;
        y @1 :Int32;
    }
}

struct PlayerMove {
  playerId @0 :Int32;
  position @1 :Position;
}

struct PlayerAction {
  playerId @0 :Int32;
  action @1 :Text;
}

struct GameStateUpdate {
  playerStates @0 :List(PlayerState);
}

struct PlayerState {
  playerId @0 :Int32;
  position @1 :Position;
}

struct ChatMessage {
  playerId @0 :Int32;
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
  clusterId @0 :Int32;
  pos @1 :Position;
}