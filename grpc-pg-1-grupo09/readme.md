# Práctica guiada 1 – gRPC
## Integrantes
- Colman Lucas
- Keingeski Noam
- Martinez Alejandro
- Sanchez Rocio
- Schupiak Andres

## Instrucciones
1. **Antes de comenzar, asegurarse de tener**:
- Go instalado (go version)
- protoc instalado (Protocol Buffers Compiler)
- Plugin Go para protoc:
```
go install google.golang.org/protobuf/cmd protoc-gen-go@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
De ser necesario agregar el directorio que contiene al compilador de Go al PATH.


2. **Inicializar el Módulo Go**

Desde la raíz del proyecto:
- go mod init grpc-pg-1
- go mod tidy

3. **Ejecutar**

En dos terminales distintas:
1. **Iniciar el servidor**:
go run servidor/main.go

2. **Ejecutar el cliente**:
go run cliente/main.go

Deberías ver en el servidor que se recibió un nombre, y en el cliente un mensaje como:
- Respuesta: Hola Claudio


