package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
	pb "practica-kv/proto"
)

// ----  Paso 3: Implementar el Reloj Vectorial -------
// VectorReloj representa un reloj vectorial de longitud 3 (tres réplicas).
type VectorReloj [3]uint64

// Incrementar aumenta en 1 el componente correspondiente a la réplica que llama.
func (vr *VectorReloj) Incrementar(idReplica int) {
	vr[idReplica]++
}

// Fusionar toma el máximo elemento a elemento entre dos vectores.
func (vr *VectorReloj) Fusionar(otro VectorReloj) {
	for i := 0; i < 3; i++ {
		if otro[i] > (vr)[i] {
			vr[i] = otro[i]
		}
	}
}

// AntesDe devuelve true si vr < otro3 el sentido estricto (strictly less).
func (vr VectorReloj) AntesDe(otro VectorReloj) bool {
	menor := false
	for i := 0; i < 3; i++ {
		if(vr[i] > otro[i]){
			return false
		}
		if vr[i] < otro[i] {
			menor = true
		}
	}
	return menor
}

// Serialización  encodeVector. Serializa el VectorReloj a []byte para enviarlo por gRPC.
func encodeVector(vr VectorReloj) []byte {
	buf := make([]byte, 8*3) // 3 uint64 = 3 * 8 bytes
	for i := 0; i < 3; i++ {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], vr[i])
	}

	return buf
}

// Deserialización decodeVector. Convierte []byte a VectorReloj.
func decodeVector(b []byte) VectorReloj {
	var vr VectorReloj
	for i := 0; i < 3; i++ {
		vr[i] = binary.BigEndian.Uint64(b[i*8 : (i+1)*8])
	}
	return vr
}




// ---------- Paso 4: Implementar el Servidor Replica -----------
// ValorConVersion guarda el valor y su reloj vectorial asociado.
type ValorConVersion struct {
	Valor []byte
	RelojVector VectorReloj
}


// ServidorReplica implementa pb.ReplicaServer
type ServidorReplica struct {
	pb.UnimplementedReplicaServer
	
	mu 		sync.Mutex
	almacen map[string]ValorConVersion // map[clave]ValorConVersion
	relojVector VectorReloj
	idReplica int  // 0, 1 o 2
	clientesPeer []pb.ReplicaClient // stubs gRPC a las otras réplicas
}

// NewServidorReplica crea una instancia de ServidorReplica
// idReplica: 0, 1 o 2
// peerAddrs: direcciones gRPC de los otros dos peers (ej.: []string{":50052", ":50053"})

func NewServidorReplica(idReplica int, peerAddrs []string)*ServidorReplica {
	clientes := make([]pb.ReplicaClient, 0)


	for _, addr := range peerAddrs {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("No se pudo conectar con peer %s: %v", addr, err)
		}
		clientes = append(clientes, pb.NewReplicaClient(conn))
	}

	return &ServidorReplica{
		almacen:      make(map[string]ValorConVersion),
		idReplica:    idReplica,
		clientesPeer: clientes,
	}
}

// Meteoddo GuardarLocal
func (r *ServidorReplica) GuardarLocal(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Incrementar reloj local
	r.relojVector.Incrementar(r.idReplica)

	// 2. Guardar valor localmente
	r.almacen[req.Clave] = ValorConVersion{
		Valor:       req.Valor,
		RelojVector: r.relojVector,
	}

	// 3. Construir mutación
	mutacion := &pb.Mutacion{
		Tipo:        pb.Mutacion_GUARDAR,
		Clave:       req.Clave,
		Valor:       req.Valor,
		RelojVector: encodeVector(r.relojVector),
	}

	// 4. Replicar asíncronamente
	for _, cliente := range r.clientesPeer {
		go func(c pb.ReplicaClient) {
			_, _ = c.ReplicarMutacion(context.Background(), mutacion)
		}(cliente)
	}

	// 5. Responder
	return &pb.RespuestaGuardar{
		Exito:           true,
		NuevoRelojVector: encodeVector(r.relojVector),
	}, nil
}


//Método EliminarLocal
func (r *ServidorReplica) EliminarLocal(ctx context.Context, req *pb.SolicitudEliminar) (*pb.RespuestaEliminar, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.relojVector.Incrementar(r.idReplica)

	delete(r.almacen, req.Clave)

	mutacion := &pb.Mutacion{
		Tipo:        pb.Mutacion_ELIMINAR,
		Clave:       req.Clave,
		RelojVector: encodeVector(r.relojVector),
	}

	for _, cliente := range r.clientesPeer {
		go func(c pb.ReplicaClient) {
			_, _ = c.ReplicarMutacion(context.Background(), mutacion)
		}(cliente)
	}

	return &pb.RespuestaEliminar{
		Exito:           true,
		NuevoRelojVector: encodeVector(r.relojVector),
	}, nil
}

// Método ObtenerLocal
func (r *ServidorReplica) ObtenerLocal(ctx context.Context, req *pb.SolicitudObtener) (*pb.RespuestaObtener, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	valor, ok := r.almacen[req.Clave]
	return &pb.RespuestaObtener{
		Valor:       valor.Valor,
		RelojVector: encodeVector(valor.RelojVector),
		Existe:      ok,
	}, nil
}

// Método ReplicarMutacion
func (r *ServidorReplica) ReplicarMutacion(ctx context.Context, m *pb.Mutacion) (*pb.Reconocimiento, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	remoto := decodeVector(m.RelojVector)
	local, existe := r.almacen[m.Clave]

	// Si no existe localmente o el remoto es más nuevo
	if !existe || local.RelojVector.AntesDe(remoto) {
		if m.Tipo == pb.Mutacion_GUARDAR {
			r.almacen[m.Clave] = ValorConVersion{
				Valor:       m.Valor,
				RelojVector: remoto,
			}
		} else if m.Tipo == pb.Mutacion_ELIMINAR {
			delete(r.almacen, m.Clave)
		}
	}

	r.relojVector.Fusionar(remoto)

	return &pb.Reconocimiento{
		Ok:             true,
		RelojVectorAck: encodeVector(r.relojVector),
	}, nil
}


// Función main
func main() {
	if len(os.Args) != 5 {
		log.Fatalf("Uso: %s <idReplica> <direccionEscucha> <peer1> <peer2>", os.Args[0])
	}

	idReplica := os.Args[1]
	listenAddr := os.Args[2]
	peer1 := os.Args[3]
	peer2 := os.Args[4]

	// Convertir idReplica a int
	var id int
	fmt.Sscanf(idReplica, "%d", &id)

	servidor := grpc.NewServer()
	replica := NewServidorReplica(id, []string{peer1, peer2})
	pb.RegisterReplicaServer(servidor, replica)

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	log.Printf("Réplica %d escuchando en %s", id, listenAddr)
	if err := servidor.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}




