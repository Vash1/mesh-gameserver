package messageHandler

import (
	capnp_schema "base/capnp"
	models "base/models"
	"base/serialization"
	"log"

	"capnproto.org/go/capnp/v3"
)

var ClientConnectionResponseChan = make(chan *models.ClientConnectionResponse, 10)

func ClientConnectionResponse(msg *capnp.Message, sourceClient string) {
	rootMsg, err := capnp_schema.ReadRootClientConnectionResponse(msg)
	if err != nil {
		log.Println("Failed to read ClientConnectionResponse:", err)
	}
	connMsg := serialization.DeserializeClientConnectionResponse(rootMsg)
	ClientConnectionResponseChan <- connMsg
	log.Printf("Received ClientConnectionResponse: %+v", connMsg)
}
