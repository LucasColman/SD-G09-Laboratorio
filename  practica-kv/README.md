# 🗃️ Sistema Clave-Valor Multi-Master con Coordinador y Réplicas (Go + gRPC)

Este proyecto implementa un sistema distribuido de almacenamiento **clave-valor** con replicación asíncrona, detección de conflictos mediante **relojes vectoriales**, y un **coordinador** central para gestionar las operaciones de los clientes.

### Integrantes
    - Colman Lucas
    - Keingeski Noam
    - Martinez Alejandro
    - Sanchez Rocio
    - Schupiak Andres

---

    

## 📚 Estructura del Proyecto
```
practica-kv/
├── cliente/
│ └── cliente_ejemplo.go
├── coordinador/
│ └── servidor_coordinador.go
├── proto/
│ └── kv.proto
├── replica/
│ └── servidor_replica.go
├── go.mod
```
---

## 🔧 Requisitos

- Go 1.16 o superior
- Protocol Buffers (`protoc`)
- Plugins Go:
  - `protoc-gen-go`
  - `protoc-gen-go-grpc`

  Instalación de plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Asegurarse de tener $GOPATH/bin en tu PATH.

## 🛠️ Compilación y Generación de Código
1. Desde la raíz del proyecto:
``` bash
go mod tidy
```

2. Generar el código gRPC:
``` bash
protoc --go_out=. --go-grpc_out=. proto/kv.proto
```

## 🚀 Ejecución del Sistema
🧱 1. Levantar las 3 réplicas (en terminales separadas)
``` bash
go run replica/servidor_replica.go 0 :50051 :50052 :50053
go run replica/servidor_replica.go 1 :50052 :50051 :50053
go run replica/servidor_replica.go 2 :50053 :50051 :50052
```

🧭 2. Levantar el Coordinador
``` bash
go run coordinador/servidor_coordinador.go -listen :6000 :50051 :50052 :50053
```

🤖 3. Ejecutar el Cliente
``` bash
go run cliente/cliente_ejemplo.go
```

## 🧪 Prueba básica
El cliente realiza lo siguiente:

    - Guarda la clave "usuario123" con el valor "datosImportantes".

    - Recupera y muestra el valor y el reloj vectorial.

    - Elimina la clave.

    - Verifica que haya sido eliminada.



