package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-pg-1/proto"       // alias “pb” para abreviar
	"google.golang.org/grpc"
)

type servidor struct {
	pb.UnimplementedServicioServer
}

func (s *servidor) Hola(ctx context.Context, req *pb.Requerimiento) (*pb.Respuesta, error) {
	log.Printf("Recibido: %s", req.Nombre)
	return &pb.Respuesta{Mensaje: "Hola " + req.Nombre}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterServicioServer(s, &servidor{})

	fmt.Println("Servidor escuchando en :8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
