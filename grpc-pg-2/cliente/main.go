package main

import (
	"context"
	"log"
	"os"
	"time"

	"grpc-pg-2/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Uso: cliente <idNodo>")
	}
	nodo := os.Args[1]

	conn, err := grpc.Dial("localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	c := proto.NewMonitorClient(conn)
	stream, err := c.EnviarHeartbeat(context.Background())
	if err != nil {
		log.Fatalf("No se pudo abrir stream: %v", err)
	}

	for {
		hb := &proto.Heartbeat{
			NodoId:      nodo,
			MarcaTiempo: time.Now().Unix(),
		}
		if err := stream.Send(hb); err != nil {
			log.Fatalf("Error enviando heartbeat: %v", err)
		}
		log.Printf("[%s] Heartbeat enviado", nodo)
		time.Sleep(5 * time.Second)
	}
}
