# Práctica guiada 1 – gRPC
## Integrantes
- Colman Lucas
- Keingeski Noam
- Martinez Alejandro
- Sanchez Rocio
- Schupiak Andres

## Instrucciones de ejecucion
1. **Antes de comenzar, asegurarse de tener**:
- Go instalado (go version)
- protoc instalado (Protocol Buffers Compiler)
- Plugin Go para protoc:
    ```
    go install google.golang.org/protobuf/cmd protoc-gen-go@latest

    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
De ser necesario agregar el directorio que contiene al compilador de Go al PATH.

2. **Definir el servicio en protobuf**
    - Archivo: proto/servicio.proto

3. **Generar el código gRPC y Go desde .proto**
    - Ejecutar desde la raíz del proyecto:
    ```
    protoc --go_out=. --go-grpc_out=. proto/servicio.proto
    ```


4. **Inicializar el Módulo Go**

    Desde la raíz del proyecto:
    - go mod init grpc-pg-1
    - go mod tidy

5. **Ejecutar**

    En dos terminales distintas:
    1. **Iniciar el servidor**: go run servidor/main.go

        - **Iniciar redirigiendo la salida a un archivo txt**: go run servidor/main.go > listado_servidor.txt 2>&1
    

    2. **Ejecutar el cliente**:
    go run cliente/main.go




