package proxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

type PostgresProxy struct {
	namespace string
}

func NewPostgresProxy(namespace string) *PostgresProxy {
	return &PostgresProxy{namespace: namespace}
}

func (p *PostgresProxy) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	log.Printf("TCP Database Proxy listening on :%s", port)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept TCP connection: %v", err)
			continue
		}
		go p.handleConnection(clientConn)
	}
}

func (p *PostgresProxy) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	// 1. Read the initial Postgres Startup Message
	buffer := make([]byte, 1024)
	n, err := clientConn.Read(buffer)
	if err != nil {
		return
	}

	username := extractPostgresUsername(buffer[:n])
	if username == "" {
		log.Println("Proxy: Connection rejected (No username found)")
		return
	}

	tenantID := username

	log.Printf("Proxy: Routing valid connection for Tenant [%s]", tenantID)

	// 4. Build the internal K8s DNS target
	internalK8sTarget := fmt.Sprintf("%s-db-svc.%s.svc.cluster.local:5432", tenantID, p.namespace)

	// 5. Connect to the internal database pod
	dbConn, err := net.Dial("tcp", internalK8sTarget)
	if err != nil {
		log.Printf("Proxy: Failed to reach internal DB %s: %v", internalK8sTarget, err)
		return
	}
	defer dbConn.Close()

	// 6. Send the intercepted startup message to the DB
	dbConn.Write(buffer[:n])

	// 7. Bridge the connections
	go io.Copy(dbConn, clientConn)
	io.Copy(clientConn, dbConn)
}

func extractPostgresUsername(data []byte) string {
	if len(data) < 8 {
		return ""
	}
	payload := data[8:]
	parts := bytes.Split(payload, []byte{0})
	for i := 0; i < len(parts)-1; i++ {
		if string(parts[i]) == "user" {
			return string(parts[i+1])
		}
	}
	return ""
}
