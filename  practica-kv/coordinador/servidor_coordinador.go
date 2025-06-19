package main

/*
 Implementacion del servidor Coordinador, que actúa como punto de entrada para los clientes, redirigiendo las solicitudes a las réplicas de manera balanceada.
*/

//Paso 1: Importaciones necesarias
import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	pb "practica-kv/proto"
)

// Paso 2: Estructura del Coordinador. ServidorCoordinador implementa pb.CoordinadorServer.
type ServidorCoordinador struct {
	pb.UnimplementedCoordinadorServer

	listaReplicas []string // ej: []string{":50051", ":50052",":50053"}
	stubs         []pb.ReplicaClient // conexiones gRPC a las réplicas 
	indiceRR uint64 // contador atómico para round-robin

}


// Paso 3: Constructor. NewServidorCoordinador crea un Coordinador con direcciones de réplica.
func NewServidorCoordinador(replicas []string) *ServidorCoordinador {
	stubs := make([]pb.ReplicaClient, 0)
	for _, addr := range replicas {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("No se pudo conectar con réplica %s: %v", addr, err)
		}
		stubs = append(stubs, pb.NewReplicaClient(conn))
	}
	
	
	return &ServidorCoordinador{
		listaReplicas: replicas,
		stubs:         stubs,
		indiceRR:      0,
	}
}

// Paso 4: Selección de réplicas (round-robin)

func (c *ServidorCoordinador) elegirReplica() int {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return int(idx) % len(c.stubs)
}



// elegirReplicaParaEscritura: round-robin simple (ignora la clave).
func (c *ServidorCoordinador) elegirReplicaParaEscritura(clave string) string {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return c.listaReplicas[int(idx)%len(c.listaReplicas)]
}

// elegirReplicaParaLectura: también round-robin.
func (c *ServidorCoordinador) elegirReplicaParaLectura() string {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return c.listaReplicas[int(idx)%len(c.listaReplicas)]
}

// Paso 5: Operaciones del Coordinador. 
// Obtener: redirige petición de lectura a una réplica.
func (c *ServidorCoordinador) Obtener(ctx context.Context, req *pb.SolicitudObtener) (*pb.RespuestaObtener, error) {
	idx := c.elegirReplica()
	return c.stubs[idx].ObtenerLocal(ctx, req)
}



// Guardar: redirige petición de escritura a una réplica elegida.
func (c *ServidorCoordinador) Guardar(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
	idx := c.elegirReplica()
	return c.stubs[idx].GuardarLocal(ctx, req)
}



// Eliminar: redirige petición de eliminación a una réplica elegida.
func (c *ServidorCoordinador) Eliminar(ctx context.Context, req *pb.SolicitudEliminar) (*pb.RespuestaEliminar, error) {
	idx := c.elegirReplica()
	return c.stubs[idx].EliminarLocal(ctx, req)
}


func main() {

	// Definir bandera para la dirección de escucha del Coordinador.
	listen := flag.String("listen", ":6000", "dirección para que escuche el Coordinador (p.ej., :6000)")
	flag.Parse()
	replicas := flag.Args()
	if len(replicas) < 3 {
		log.Fatalf("Debe proveer al menos 3 direcciones de réplicas, p.ej.: go run servidor_coordinador.go -listen :6000 :50051 :50052 :50053")
	}

	lis, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}


	grpcServer := grpc.NewServer()
	coordinador := NewServidorCoordinador(replicas)
	pb.RegisterCoordinadorServer(grpcServer, coordinador)

	fmt.Printf("Coordinador escuchando en %s\n", *listen)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}


/*
✅ Coordinador listo
Con esto:

El cliente puede conectarse al Coordinador sin preocuparse de las réplicas.

El Coordinador redirige las peticiones balanceadamente con round-robin.
*/
