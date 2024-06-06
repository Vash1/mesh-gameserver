package network

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"sync"
	"time"

	"capnproto.org/go/capnp/v3"
	quic "github.com/quic-go/quic-go"
)

type ServerConfig struct {
	LocalAddr string
}

type NetServer struct {
	localAddr     net.Addr
	localConn     *net.UDPConn
	quicTransport *quic.Transport
	quicListener  *quic.Listener
}

type ClientConfig struct {
	RemoteAddr string
}

type NetClient struct {
	remoteAddr     *net.UDPAddr
	localConn      *net.UDPConn
	quicTransport  *quic.Transport
	quicConnection connection
	QuicStream     quicStream
}

type quicStream struct {
	quic.Stream
	mu *sync.Mutex
}

type connection struct {
	quic.Connection
}

func NewNetServer(config ServerConfig) (*NetServer, error) {
	addr, err := net.ResolveUDPAddr("udp4", config.LocalAddr)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, err
	}

	quicTransport := &quic.Transport{
		Conn: udpConn,
	}

	tlsConfig := generateTLSConfig()
	ln, err := quicTransport.Listen(tlsConfig, &quic.Config{EnableDatagrams: true, MaxIdleTimeout: 10 * time.Second, KeepAlivePeriod: 5 * time.Second})
	if err != nil {
		return nil, err
	}
	fmt.Println("QUIC server listening on", ln.Addr())
	return &NetServer{
		localAddr:     udpConn.LocalAddr(),
		localConn:     udpConn,
		quicTransport: quicTransport,
		quicListener:  ln,
	}, nil
}

func NewNetClient(config ClientConfig) (*NetClient, error) {
	localConn, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		return nil, err
	}

	remoteAddr, err := net.ResolveUDPAddr("udp4", config.RemoteAddr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	quicTransport := &quic.Transport{
		Conn: localConn,
	}
	quicConnection, err := quicTransport.Dial(ctx, remoteAddr, &tls.Config{
		InsecureSkipVerify: true, // For testing only; ensure proper TLS config in production
		NextProtos:         []string{"quic-echo-example"},
	}, &quic.Config{
		EnableDatagrams: true,
	})
	if err != nil {
		return nil, err
	}

	return &NetClient{
		remoteAddr:     remoteAddr,
		localConn:      localConn,
		quicTransport:  quicTransport,
		quicConnection: connection{quicConnection},
	}, nil

}

func (client *NetClient) OpenStream() error {
	stream, err := client.quicConnection.OpenStream()
	if err != nil {
		return err
	}

	log.Println("QUIC Stream opened successfully.")
	client.QuicStream = quicStream{stream, &sync.Mutex{}}
	return nil
}

func (client *NetClient) Listen() *capnp.Message {
	for {
		msg, ok := read(client.QuicStream)
		if !ok {
			fmt.Println("Stream closed")
			return nil
		}
		fmt.Println("Received message:", msg)
		// return msg
	}
}

func (stream *quicStream) SendMessage(msg *capnp.Message) error {
	if err := capnp.NewEncoder(stream).Encode(msg); err != nil {
		return err
	}
	return nil
}

func (conn *connection) SendDatagram(bytes []byte) error {
	err := conn.Connection.SendDatagram(bytes)
	if err != nil {
		return fmt.Errorf("failed to send datagram: %w", err)
	}
	return nil
}

func (server *NetServer) AcceptConnection() (*connection, error) {
	conn, err := server.quicListener.Accept(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("\nConnection accepted from", conn.RemoteAddr())
	return &connection{conn}, nil
}

func (serverConn *connection) AcceptStream() (quic.Stream, error) {
	stream, err := serverConn.Connection.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (serverConn *connection) receiveDatagram() ([]byte, error) {
	datagram, err := serverConn.ReceiveDatagram(context.Background())
	if err != nil {
		log.Println("Failed to receive datagram:", err)
		return nil, err
	}
	return datagram, nil
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Company, INC."},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}

func read(stream quic.Stream) (*capnp.Message, bool) {
	decoder := capnp.NewDecoder(stream)
	msg, err := decoder.Decode()
	if err != nil {
		if err == io.EOF {
			fmt.Println("Stream closed by client.")
			return nil, false
		}
		log.Println("Failed to decode message:", err)
		return nil, false
	}

	return msg, true
}
