package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"grpc-pg-1/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	
	defer conn.Close()

	c := proto.NewServicioClient(conn)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			nombre := fmt.Sprintf("Persona %d", i)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err := c.Hola(ctx, &proto.Requerimiento{Nombre: nombre})
			if err != nil {
				log.Printf("Error al saludar a %s: %v", nombre, err)
			}
		}(i)
	}

	wg.Wait()

	// Luego de saludar a todos, pedimos el listado
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	respuesta, err := c.ListadoPersonas(ctx, &proto.Vacio{})
	if err != nil {
		log.Fatalf("Error al obtener el listado: %v", err)
	}

	fmt.Println("Listado de personas saludadas:")
	for _, nombre := range respuesta.Personas {
		fmt.Println("-", nombre)
	}
}
