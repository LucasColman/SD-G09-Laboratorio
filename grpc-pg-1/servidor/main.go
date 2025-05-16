package main

import (
	"context"
	"fmt"
	"grpc-pg-1/proto"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type servidor struct {
	proto.UnimplementedServicioServer
	personas []string
	mu       sync.Mutex
}

// Método Hola
func (s *servidor) Hola(ctx context.Context, req *proto.Requerimiento) (*proto.Respuesta, error) {
	log.Printf("Recibido: %s", req.Nombre)

	s.mu.Lock()
	s.personas = append(s.personas, req.Nombre) // Agregar persona a la lista
	s.mu.Unlock()

	return &proto.Respuesta{Mensaje: "Hola " + req.Nombre}, nil
}

// Método ListadoPersonas
func (s *servidor) ListadoPersonas(ctx context.Context, req *proto.Vacio) (*proto.Lista, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &proto.Lista{Personas: s.personas}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterServicioServer(s, &servidor{})
	fmt.Println("Servidor escuchando en :8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
