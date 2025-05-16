package main

import (
	"context"
	"log"
	"time"

	pb "grpc-pg-1/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	c := pb.NewServicioClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Hola(ctx, &pb.Requerimiento{Nombre: "Claudio"})
	if err != nil {
		log.Fatalf("Error al llamar al servidor: %v", err)
	}
	log.Printf("Respuesta: %s", r.Mensaje)
}
