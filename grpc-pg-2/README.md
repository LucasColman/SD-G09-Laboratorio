# Práctica Guiada 2 – gRPC Detector de Fallos

## Integrantes
- Keingeski Noam
- Lucas Colman
- Rocío Sanchez
- Alejandro Martinez
- Andrés Schupiak

## Cómo compilar y ejecutar
descomprimir la carpeta "grpc-pg-2" en el escritorio.

```bash
# Generar código protobuf
protoc --go_out=. --go-grpc_out=. proto/monitor.proto

#Ubicarse en la carpeta en los 3 CMD 
cd C:\Users\usuario\Desktop\grpc-pg-2

# Iniciar servidor
go run servidor\main.go

# Iniciar clientes (en otras ventanas CMD)
go run cliente/main.go nodo1
go run cliente/main.go nodo2

#Matar nodo2
Ctrl+C en el CMD

#Efecto de la frecuencia de Heartbeats
Aumentar la frecuencia (por ejemplo, cada 2 s):
    Detecta fallos más rápido.
    Más tráfico de red y CPU.

Disminuir frecuencia (p. ej. cada 10 s):
    Menos carga de red/CPU.
    Detección de fallos más lenta.
