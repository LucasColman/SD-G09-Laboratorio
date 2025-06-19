package main

/*
Implementacion del cliente de ejemplo, que interactúa con el Coordinador para probar el sistema completo: guardar, obtener, eliminar y verificar una clave.
*/


//  Paso 1: Importaciones necesarias
import(
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "practica-kv/proto"
)

// Paso 2: Función principal
func main() {
	//Conectarse al Coordinador (Puerto 6000)
	conn, err := grpc.Dial(":6000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar con el coordinador: %v",err)
	}
	defer conn.Close()

	cliente:= pb.NewCoordinadorClient(conn)

	ctx, cancel:= context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	//  Paso 3: Guardar una clave
	fmt.Println("👉 Guardando clave 'usuario123' con valor 'datosImportantes'...")
	respGuardar, err := cliente.Guardar(ctx, &pb.SolicitudGuardar{
		Clave:       "usuario123",
		Valor:       []byte("datosImportantes"),
		RelojVector: nil, // El cliente no necesita enviar vector inicialmente
	})
	if err != nil {
		log.Fatalf("Error al guardar: %v", err)
	}
	fmt.Printf("✔️ Guardado exitoso. Nuevo reloj: %v\n", respGuardar.NuevoRelojVector)

	// Paso 4: Obtener la clave
	fmt.Println("🔎 Obteniendo clave 'usuario123'...")
	respObtener,err:= cliente.Obtener(ctx,&pb.SolicitudObtener{
		Clave: "usuario123",
	})
	if err != nil {
		log.Fatalf("Error al obtener: %v", err)
	}
	if respObtener.Existe {
		fmt.Printf("✅ Valor obtenido: %s\n", string(respObtener.Valor))
		fmt.Printf("⏱️ Reloj vectorial: %v\n", respObtener.RelojVector)

	}else{
		fmt.Println("⚠️ Clave no encontrada")
	}


	// Paso 5: Eliminar la clave
	fmt.Println("🗑️ Eliminando clave 'usuario123'...")
	respEliminar,err:= cliente.Eliminar(ctx,&pb.SolicitudEliminar{
		Clave: "Usuario123",
		RelojVector: respObtener.RelojVector, // Usamos el vector recibido
	})

	if err != nil {
		log.Fatalf("Error al eliminar: %v", err)
	}

	fmt.Printf("✔️ Eliminación exitosa. Nuevo reloj: %v\n", respEliminar.NuevoRelojVector)

	// Paso 6: Verificar que la clave fue eliminada
	fmt.Println("🔄 Verificando que la clave fue eliminada...")
	respVerificacion, err := cliente.Obtener(ctx, &pb.SolicitudObtener{
		Clave: "usuario123",
	})
	if err != nil {
		log.Fatalf("Error al verificar: %v", err)
	}

	if(!respVerificacion.Existe){
		fmt.Println("✅ La clave ya no existe.")
	}else{
		fmt.Printf("❌ La clave aún existe con valor: %s\n", string(respVerificacion.Valor))

	}

	/*
	Cliente completo
	Este cliente:
		Guarda una clave con un valor.
		La recupera y muestra su reloj vectorial.
		La elimina.
		Verifica que haya sido eliminada.
	*/




}