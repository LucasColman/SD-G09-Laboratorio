package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"grpc-pg-2/proto"

	"google.golang.org/grpc"
)

type servidor struct {
	proto.UnimplementedMonitorServer
	mu          sync.Mutex
	ultimaVista map[string]time.Time
}

func (s *servidor) EnviarHeartbeat(stream proto.Monitor_EnviarHeartbeatServer) error {
	for {
		hb, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.Ack{Mensaje: "Stream cerrado"})
		}
		if err != nil {
			log.Printf("Error en stream: %v", err)
			return err
		}
		s.mu.Lock()
		s.ultimaVista[hb.NodoId] = time.Unix(hb.MarcaTiempo, 0)
		s.mu.Unlock()
	}
}

func (s *servidor) detectorFallas(intervalo time.Duration) {
	for {
		time.Sleep(intervalo)
		s.mu.Lock()
		ahora := time.Now()
		vivos := []string{}
		for nodo, ultimo := range s.ultimaVista {
			if ahora.Sub(ultimo) <= 3*intervalo {
				vivos = append(vivos, nodo)
			}
		}
		s.mu.Unlock()

		//log.Printf("Nodos vivos: %v\n", vivos)
		//fmt.Println("Nodos vivos:", vivos)
		fmt.Println("Nodos vivos:", vivos)

	}
}

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}
	s := grpc.NewServer()
	srv := &servidor{ultimaVista: make(map[string]time.Time)}
	proto.RegisterMonitorServer(s, srv)

	// lanza el detector cada 5s
	go srv.detectorFallas(5 * time.Second)

	fmt.Println("Servidor escuchando en :8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
