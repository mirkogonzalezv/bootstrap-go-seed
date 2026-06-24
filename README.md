# xintecgo — Microservice Seed

Seed base para microservicios en Go con **Clean Architecture**.  
Diseñada para que te preocupes solo de los módulos de negocio en `internal/`.

## Requisitos

- Go 1.25+
- (Opcional) Air para hot-reload → `make install-air`

## Quick Start

```bash
cp .env.example .env.dev
# completar vars DB en .env.dev si tienes PostgreSQL
make start-dev          # hot-reload
# o
make start              # modo normal
```

Documentación Swagger: `http://localhost:3200/docs/index.html`

---

## Estructura del proyecto

```
cmd/                        # Bootstrap de la app (no tocar)
├── api/main.go             # Punto de entrada
├── app/                    # Init, Run, Shutdown (lifecycle)
├── config/                 # Variables de entorno
├── container/              # Inyección de dependencias (composition root)
└── routes/                 # Registro de rutas

internal/                   # ← Módulos de negocio (aquí trabajas)
└── <modulo>/
    ├── application/        # Lógica de negocio
    │   ├── dto/            # Objetos de transferencia
    │   └── usecase/        # Casos de uso
    ├── domain/             # Entidades e interfaces de repositorio
    ├── infrastructure/     # Implementaciones concretas (DB, APIs externas)
    └── presentation/       # Capa HTTP
        ├── handler/        # Controladores Gin
        └── router/         # Definición de rutas

migrations/                 # SQL de migraciones (ejecutar a mano)
pkg/                        # Utilidades compartidas
├── apperrors/              # Errores de dominio tipados
├── database/               # Database Provider (interfaz + drivers)
│   ├── database.go         # Interfaz Database (Ping, Close, Native)
│   └── postgres.go         # Driver PostgreSQL (PostgresConfig, NewPostgresDatabase)
├── env/                    # Tipos de entorno (dev, qa, prod)
├── httputil/               # Mapper errores de dominio → HTTP
├── logger/                 # Logger con Zap
├── middleware/             # Middlewares Gin (error handler, security)
└── pubsub/                 # Pub/Sub Provider (interfaz + drivers)
    ├── pubsub.go           # Interfaz PubSub (Publish, Subscribe, Close, Native)
    └── gcp.go              # Driver GCP Pub/Sub (NewGCPPubSub)
```

---

## Cómo crear un nuevo módulo

Ejemplo: módulo `tareas` con endpoint `GET /api/v1/tareas`.

### 1. Crear la estructura de carpetas

```bash
mkdir -p internal/tareas/application/{dto,usecase}
mkdir -p internal/tareas/domain
mkdir -p internal/tareas/infrastructure
mkdir -p internal/tareas/presentation/{handler,router}
```

### 2. Definir el DTO de respuesta

`internal/tareas/application/dto/tarea_response.go`:

```go
package dto

type TareaResponse struct {
    ID     string `json:"id"`
    Nombre string `json:"nombre"`
}

func NewTareaResponse(id, nombre string) *TareaResponse {
    return &TareaResponse{ID: id, Nombre: nombre}
}
```

### 3. Definir la interfaz del repositorio en domain

`internal/tareas/domain/repository.go`:

```go
package domain

import "context"

type TareaRepository interface {
    FindByID(ctx context.Context, id string) (*Tarea, error)
}
```

### 4. Crear el caso de uso

`internal/tareas/application/usecase/obtener_tarea.go`:

```go
package usecase

import (
    "context"
    "microservice/internal/tareas/application/dto"
    "microservice/internal/tareas/domain"
)

type ObtenerTareaUseCase struct {
    repo domain.TareaRepository
}

func NewObtenerTareaUseCase(repo domain.TareaRepository) *ObtenerTareaUseCase {
    return &ObtenerTareaUseCase{repo: repo}
}

func (uc *ObtenerTareaUseCase) Execute(ctx context.Context) (*dto.TareaResponse, error) {
    tarea, err := uc.repo.FindByID(ctx, "1")
    if err != nil {
        return nil, err
    }
    return dto.NewTareaResponse(tarea.ID, tarea.Nombre), nil
}
```

### 5. Crear la implementación del repositorio

`internal/tareas/infrastructure/postgres_repo.go`:

```go
package infrastructure

import (
    "context"
    "microservice/internal/tareas/domain"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTareaRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresTareaRepository(pool *pgxpool.Pool) *PostgresTareaRepository {
    return &PostgresTareaRepository{pool: pool}
}

func (r *PostgresTareaRepository) FindByID(ctx context.Context, id string) (*domain.Tarea, error) {
    // query a la DB...
}
```

### 6. Crear el handler HTTP

`internal/tareas/presentation/handler/tarea.go`:

```go
package handler

import (
    "net/http"
    "microservice/internal/tareas/application/usecase"
    "github.com/gin-gonic/gin"
)

type TareaHandler struct {
    uc *usecase.ObtenerTareaUseCase
}

func NewTareaHandler(uc *usecase.ObtenerTareaUseCase) *TareaHandler {
    return &TareaHandler{uc: uc}
}

func (h *TareaHandler) GetTarea(c *gin.Context) {
    res, err := h.uc.Execute(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, res)
}
```

### 7. Definir las rutas

`internal/tareas/presentation/router/tarea.go`:

```go
package router

import (
    "microservice/internal/tareas/presentation/handler"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.TareaHandler) {
    v1 := api.Group("/v1/tareas")
    v1.GET("", h.GetTarea)
}
```

### 8. Registrar en el container de dependencias

`cmd/container/container.go`:

```go
import (
    tareaDomain "microservice/internal/tareas/domain"
    tareaInfra "microservice/internal/tareas/infrastructure"
    tareaUsecase "microservice/internal/tareas/application/usecase"
    tareaHandler "microservice/internal/tareas/presentation/handler"
)

type Container struct {
    DB            *pgxpool.Pool
    HealthHandler *handler.HealthHandler
    TareaHandler  *tareaHandler.TareaHandler   // ← agregar
    Log           *zap.Logger
}

func NewContainer(cfg *config.Configuration, log *zap.Logger) (*Container, error) {
    // ... pool, health ...

    var tareaRepo tareaDomain.TareaRepository = tareaInfra.NewPostgresTareaRepository(pool)
    tareaUC := tareaUsecase.NewObtenerTareaUseCase(tareaRepo)
    tareaHdlr := tareaHandler.NewTareaHandler(tareaUC)

    return &Container{
        DB:            pool,
        HealthHandler: healthHdlr,
        TareaHandler:  tareaHdlr,              // ← agregar
        Log:           log,
    }, nil
}
```

### 9. Registrar la ruta

`cmd/routes/register.go`:

```go
import tareaRouter "microservice/internal/tareas/presentation/router"

func ConfigureModule(router *gin.Engine, c *container.Container) {
    api := router.Group("/api")
    health.RegisterRoutes(api, c.HealthHandler)
    tareaRouter.RegisterRoutes(api, c.TareaHandler)   // ← agregar
}
```

✅ **Listo** — tu endpoint ya responde en `GET /api/v1/tareas`.

---

## Base de datos con Database Provider

El proyecto separa la inicialización de la base de datos en una capa independiente (`pkg/database/`) para que el **container** no dependa del motor concreto. Usa el patrón **Provider**:

```
App.Init()
  → database.NewPostgresDatabase(cfg)   ← acá se crea la DB
  → container.NewContainer(..., db)     ← container recibe Database sin saber el motor

Handler (HTTP)
  → UseCase (lógica de negocio)
    → Repository (interfaz en domain/)
      → PostgresRepository (implementación en infrastructure/)
        → Database.Native().(*pgxpool.Pool)
```

### Interfaz `Database`

Definida en `pkg/database/database.go`:

```go
type Database interface {
    Ping(ctx context.Context) error
    Close()
    Native() interface{}
}
```

- `Ping` / `Close` — ciclo de vida estándar
- `Native()` — expone el pool/conexión nativa para que los repositorios hagan type assertion al tipo concreto (`*pgxpool.Pool`, `*sql.DB`, etc.)

Cada driver implementa esta interfaz con su propio struct privado y factory: `PostgresConfig` + `NewPostgresDatabase` en `postgres.go`, etc.

### Migraciones

Las migraciones son **SQL puro**, ejecutadas a mano. No hay ORM, no hay auto-migrate:

```bash
psql -d microservice -f migrations/001_create_health_check.up.sql
```

### Pool de conexiones

Se crea en `cmd/app/app.go` usando la factory del driver correspondiente:

```go
db, err := database.NewPostgresDatabase(ctx, database.PostgresConfig{
    Host: cfg.DBHost, Port: cfg.DBPort, User: cfg.DBUser,
    Password: cfg.DBPassword, DBName: cfg.DBName, SSLMode: cfg.DBSSLMode,
})
```

La DB se cierra en el shutdown de la app (`cmd/app/shutdown.go`), **no** desde el container.

---

## Cómo agregar un nuevo motor de base de datos

Ejemplo: agregar **MySQL** como segundo motor.

### 1. Crear el driver en `pkg/database/mysql.go`

```go
package database

import (
    "context"
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
}

func (c MySQLConfig) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
        c.User, c.Password, c.Host, c.Port, c.DBName)
}

type mysqlDatabase struct {
    db *sql.DB
}

func (d *mysqlDatabase) Ping(ctx context.Context) error { return d.db.PingContext(ctx) }
func (d *mysqlDatabase) Close()                         { d.db.Close() }
func (d *mysqlDatabase) Native() interface{}            { return d.db }

func NewMySQLDatabase(ctx context.Context, cfg MySQLConfig) (Database, error) {
    db, err := sql.Open("mysql", cfg.DSN())
    if err != nil {
        return nil, fmt.Errorf("database mysql: %w", err)
    }
    return &mysqlDatabase{db: db}, nil
}
```

### 2. Crear la implementación del repositorio para MySQL

`internal/<módulo>/infrastructure/mysql_repo.go`:

```go
package infrastructure

import (
    "context"
    "database/sql"
    "microservice/pkg/database"
)

type MySQLHealthRepository struct {
    db database.Database
}

func NewMySQLHealthRepository(db database.Database) *MySQLHealthRepository {
    return &MySQLHealthRepository{db: db}
}

func (r *MySQLHealthRepository) Ping(ctx context.Context) error {
    sqlDB := r.db.Native().(*sql.DB)
    return sqlDB.PingContext(ctx)
}
```

### 3. Conectar en `cmd/app/app.go`

```go
db, err := database.NewMySQLDatabase(context.Background(), database.MySQLConfig{
    Host: cfg.DBHost, Port: cfg.DBPort, User: cfg.DBUser,
    Password: cfg.DBPassword, DBName: cfg.DBName,
})
```

### 4. Inyectar el repositorio en `cmd/container/container.go`

```go
repo := healthInfra.NewMySQLHealthRepository(db)
```

### Resumen por driver nuevo

| Qué | Archivos |
|---|---|
| Config + factory + struct privado que implementa `Database` | 1 archivo: `pkg/database/<driver>.go` |
| Repositorio por módulo (type assertion con `Native()`) | 1 archivo por módulo en `internal/<módulo>/infrastructure/` |
| Cableado | `app.go` (crear DB) + `container.go` (inyectar repo) |

**El container nunca sabe qué motor es.** Solo recibe `database.Database` y llama al constructor del repo que corresponda.

---

## Pub/Sub con PubSub Provider

El proyecto sigue el mismo patrón que Database para mensajería asíncrona. La interfaz `PubSub` está en `pkg/pubsub/` y los drivers concretos implementan el contrato.

### Flujo

```
App.Init()
  → pubsub.NewGCPPubSub(ctx, projectID)   ← acá se crea el cliente
  → container.NewContainer(..., ps)        ← container recibe PubSub sin saber el proveedor

Handler (HTTP)
  → UseCase (lógica de negocio)
    → Publisher (implementación en infrastructure/)
      → PubSub.Publish(topic, data, attrs)
```

### Interfaz `PubSub`

Definida en `pkg/pubsub/pubsub.go`:

```go
type PubSub interface {
    Publish(ctx context.Context, topic string, data []byte, attrs map[string]string) error
    Subscribe(ctx context.Context, subscription string, handler MessageHandler) error
    Close() error
    Native() interface{}
}

type Message struct {
    ID         string
    Data       []byte
    Attributes map[string]string
    AckFunc    func()
    NackFunc   func()
}

type MessageHandler func(ctx context.Context, msg *Message) error
```

### Driver GCP

`pkg/pubsub/gcp.go` implementa `PubSub` usando `cloud.google.com/go/pubsub/v2`:

```go
ps, err := pubsub.NewGCPPubSub(context.Background(), cfg.PubSubProjectID)
```

- Cachea `Publisher` por topic para reutilizar el batching interno del SDK
- `Subscribe` recibe mensajes vía `Subscriber.Receive` y mapea a `Message` propio
- `Close()` detiene todos los publishers y cierra el cliente gRPC

### Cableado

Se crea en `cmd/app/app.go` y se cierra en `cmd/app/shutdown.go`:

```go
// app.go
ps, err := pubsub.NewGCPPubSub(context.Background(), cfg.PubSubProjectID)
a.ps = ps

c, err := container.NewContainer(a.config, db, ps, a.log)

// shutdown.go
if a.ps != nil {
    a.ps.Close()
}
```

---

## Cómo agregar un nuevo proveedor Pub/Sub

Ejemplo: agregar **RabbitMQ** como mensajería.

### 1. Crear el driver en `pkg/pubsub/rabbitmq.go`

```go
package pubsub

import (
    "context"
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
    URL string
}

type rabbitMQPubSub struct {
    conn *amqp.Connection
}

func NewRabbitMQPubSub(ctx context.Context, cfg RabbitMQConfig) (PubSub, error) {
    conn, err := amqp.Dial(cfg.URL)
    if err != nil {
        return nil, fmt.Errorf("pubsub rabbitmq: %w", err)
    }
    return &rabbitMQPubSub{conn: conn}, nil
}

func (p *rabbitMQPubSub) Publish(ctx context.Context, topic string, data []byte, attrs map[string]string) error {
    ch, _ := p.conn.Channel()
    defer ch.Close()
    return ch.PublishWithContext(ctx, topic, "", false, false, amqp.Publishing{
        ContentType: "application/json",
        Body:        data,
    })
}

func (p *rabbitMQPubSub) Subscribe(ctx context.Context, subscription string, handler MessageHandler) error {
    ch, _ := p.conn.Channel()
    msgs, _ := ch.Consume(subscription, "", false, false, false, false, nil)
    for msg := range msgs {
        m := &Message{
            ID:   msg.MessageId,
            Data: msg.Body,
            AckFunc: func() { msg.Ack(false) },
            NackFunc: func() { msg.Nack(false, true) },
        }
        if err := handler(ctx, m); err != nil {
            msg.Nack(false, true)
            continue
        }
        msg.Ack(false)
    }
    return nil
}

func (p *rabbitMQPubSub) Close() error { return p.conn.Close() }
func (p *rabbitMQPubSub) Native() interface{} { return p.conn }
```

### 2. Conectar en `cmd/app/app.go`

```go
ps, err := pubsub.NewRabbitMQPubSub(context.Background(), pubsub.RabbitMQConfig{
    URL: "amqp://guest:guest@localhost:5672/",
})
```

### Resumen por proveedor nuevo

| Qué | Archivos |
|---|---|
| Config + factory + struct privado que implementa `PubSub` | 1 archivo: `pkg/pubsub/<driver>.go` |
| Servicio por módulo (publisher/subscriber) | 1 archivo por módulo en `internal/<módulo>/infrastructure/` |
| Cableado | `app.go` (crear PubSub) + `container.go` (inyectar al servicio) |

---

## Tests

Los tests de use cases usan **mocks manuales** (sin frameworks externos):

```go
type mockHealthRepository struct{}

func (m *mockHealthRepository) Ping(ctx context.Context) error {
    return nil
}

func TestHealthUseCase(t *testing.T) {
    mockRepo := &mockHealthRepository{}
    uc := NewHealthUseCase(mockRepo)

    res, err := uc.Execute(context.Background())

    if err != nil {
        t.Fatalf("Execute() returned unexpected error: %v", err)
    }

    if res.Database != "connected" {
        t.Errorf("Execute().Database = %q, want %q", res.Database, "connected")
    }
}
```

**Patrón**: mock manual → inyectar vía constructor → testear solo lógica de negocio, sin DB real.

Ejecutar:

```bash
make test
```

---

## Logger

Usa `go.uber.org/zap`. Dos formas de usarlo:

**En módulos de negocio (`internal/`)** — helpers simples, sin acoplamiento:

```go
import "microservice/pkg/logger"

logger.Success("todo ok")
logger.Error("algo falló: %s", err)
logger.General("mensaje informativo")
```

**En bootstrap e infraestructura** — logger inyectado vía DI:

```go
func ErrorHandlerMiddleware(log *zap.Logger) gin.HandlerFunc {
    // usa log.Error(), log.Warn(), etc.
}
```

---

## Manejo de errores

Usa errores de dominio tipados en `pkg/apperrors/`:

```go
import "microservice/pkg/apperrors"

apperrors.NewNotFoundError("USER_001", "Usuario no encontrado")
apperrors.NewValidationError("EMAIL_001", "Email inválido")
apperrors.NewUnauthorizedError("AUTH_001", "Token expirado")
apperrors.NewInternalError("DB_001", "Error de conexión")
```

El middleware global los mapea automáticamente al HTTP status code correcto:

| Error | HTTP Status |
|---|---|
| `Validation` | 400 |
| `NotFound` | 404 |
| `Conflict` | 409 |
| `Unauthorized` | 401 |
| `Forbidden` | 403 |
| `Internal` | 500 |

---

## Security Headers

Usa [`github.com/goddtriffin/helmet`](https://github.com/goddtriffin/helmet), el port a Go de HelmetJS.

Se activa automáticamente en todas las rutas vía middleware:

```go
router.Use(middleware.SecurityHeaders())
```

**Cabeceras incluidas por defecto:**

| Header | Valor |
|---|---|
| `X-Frame-Options` | `SAMEORIGIN` |
| `X-Content-Type-Options` | `nosniff` |
| `X-XSS-Protection` | `1; mode=block` |
| `X-DNS-Prefetch-Control` | `off` |
| `X-Download-Options` | `noopen` |
| `Strict-Transport-Security` | `max-age=5184000; includeSubDomains` |
| `X-Powered-By` | removido |

---

## Configuración

Las variables de entorno se cargan automáticamente desde `.env.{environment}` (dev/qa) o del sistema (prod).

| Variable | Default | Descripción |
|---|---|---|
| `PORT` | `3200` | Puerto del servidor |
| `ENVIRONMENT` | `dev` | Entorno: `dev`, `qa` o `prod` |
| `LOG_LEVEL` | `info` | Nivel de log |
| `SHUTDOWN_TIMEOUT` | `30` | Timeout de graceful shutdown (seg) |
| `DB_HOST` | `localhost` | Host de PostgreSQL |
| `DB_PORT` | `5432` | Puerto de PostgreSQL |
| `DB_USER` | `postgres` | Usuario de PostgreSQL |
| `DB_PASSWORD` | `postgres` | Contraseña de PostgreSQL |
| `DB_NAME` | `microservice` | Nombre de la base de datos |
| `DB_SSL_MODE` | `disable` | Modo SSL (`disable`, `require`, `verify-full`) |
| `PUBSUB_PROJECT_ID` | `""` | Project ID de GCP para Pub/Sub |

```bash
cp .env.example .env.dev
# editar .env.dev con tus valores
```

---

## Comandos disponibles

```bash
make help              # Lista todos los comandos
make start-dev         # Iniciar con hot-reload (Air)
make start             # Iniciar en modo normal
make test              # Ejecutar tests de internal/
make test-coverage     # Tests con reporte de cobertura
make docs-generate     # Generar documentación Swagger
make docs-serve        # Iniciar con documentación ReDoc
make audit             # Auditoría de dependencias (govulncheck)
make security-scan     # Escaneo de secretos (gitleaks)
make deps-check        # Verificar y limpiar dependencias
```

---

## Documentación API (Swagger)

La documentación se genera desde las anotaciones en los handlers:

```go
// @Summary Obtener tarea
// @Description Retorna una tarea por ID
// @Tags Tareas
// @Success 200 {object} dto.TareaResponse
// @Router /tareas [get]
```

Generar y visualizar:

```bash
make docs-generate
# Abrir http://localhost:3200/docs/index.html (solo dev/qa)
```

---

## Docker

```bash
docker build -t microservice .
docker run -p 3200:3200 microservice
```

**Características de la imagen:**
- Base: `gcr.io/distroless/static-debian12:nonroot`
- Tamaño: ~15 MB
- Corre como non-root (mejores prácticas)
- Certificados CA incluidos (HTTPS listo)
- Sin shell, sin package manager (mínima superficie de ataque)

Compatible con: **Cloud Run, AWS ECS/Fargate, Railway, Azure Container Apps**.

---

## Convenciones Go usadas

| Regla | Ejemplo |
|---|---|
| Sin shadow de stdlib | `pkg/apperrors`, `pkg/httputil` |
| Packages en singular | `middleware`, `usecase`, `handler` |
| Sin prefijo `Get` en getters | `Version()` en vez de `GetVersion()` |
| Archivos en snake_case | `error_handler.go`, `health_response.go` |
| Package name = nombre del directorio | `package dto` dentro de `dto/` |
| Interfaz en domain, implementación en infrastructure | `domain.HealthRepository` → `infrastructure.PostgresHealthRepository` |
