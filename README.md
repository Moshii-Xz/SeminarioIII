# VenciTrack

VenciTrack es una API REST en Go para gestion de caducidades, catalogo de productos y flujo de pedidos, con autenticacion JWT, carga de imagenes y modulos separados por dominio.

## Resumen

- Backend en Go 1.24 con Gin y GORM.
- Persistencia en PostgreSQL.
- Modulos para usuarios, catalogo, pedidos, pagos, reseñas, tiendas y reportes.
- Servido de imagenes subidas desde `/uploads`.

## Stack

- Go 1.24
- Gin
- GORM
- PostgreSQL
- JWT

## Estructura

- `cmd/api/`: punto de entrada de la aplicacion.
- `internal/domain/`: entidades de dominio.
- `internal/modules/`: logica por modulo, con handler, service y repository.
- `internal/platform/`: utilidades de infraestructura.
- `internal/server/`: configuracion del servidor y rutas.
- `migrations/`: script SQL de la base de datos.

## Requisitos

- Go 1.24 o superior.
- PostgreSQL 15 o superior.
- Git.

## Configuracion local

1. Clona el repositorio.
2. Descarga dependencias con `go mod download`.
3. Configura estas variables de entorno si quieres cambiar los valores por defecto:
   - `DB_HOST` (por defecto `localhost`)
   - `DB_PORT` (por defecto `5432`)
   - `DB_USER` (por defecto `postgres`)
   - `DB_PASSWORD` (por defecto `1234`)
   - `DB_NAME` (por defecto `vencitrack`)
   - `DB_SSLMODE` (por defecto `disable`)
   - `DB_TIMEZONE` (por defecto `America/Bogota`)
4. Levanta PostgreSQL y ejecuta las migraciones en `migrations/AS_BD.sql`.
5. Inicia la API con `go run cmd/api/main.go`.

## Docker

La forma mas simple de levantar todo el entorno es con Docker Compose:

```bash
docker-compose up --build
```

La API queda expuesta en `http://localhost:8081` y el healthcheck en `http://localhost:8081/health`.

## Endpoints principales

- `GET /health`: validacion de conexion a la base de datos.
- `POST /api/v1/users/`: registro de usuarios.
- `POST /api/v1/users/login`: inicio de sesion.
- `GET /api/v1/products/`: listado publico de productos.
- `POST /api/v1/products/`: creacion de productos autenticada.
- `POST /api/v1/products/upload-image`: subida de imagenes autenticada.
- `GET /api/v1/orders/`: gestion de pedidos autenticada.
- `GET /api/v1/payments/`: gestion de pagos autenticada.
- `GET /api/v1/reviews/`: listado publico de reseñas.
- `GET /api/v1/reports/sales`: reportes autenticados.
- `GET /api/v1/stores/`: gestion de tiendas autenticada.

La documentacion completa de la API esta en [API_DOCS.md](API_DOCS.md).

## Documentacion adicional

- [API_DOCS.md](API_DOCS.md): referencia detallada de endpoints y payloads.
- [ARCHITECTURE.md](ARCHITECTURE.md): resumen de la arquitectura, modulos y modelo de dominio.
- [SETUP.md](SETUP.md): guia de configuracion local.
- [DOCKER_SETUP.md](DOCKER_SETUP.md): guia de despliegue con Docker.

## Notas

- El servidor escucha en `:8080` dentro del contenedor.
- Las imagenes subidas se sirven desde `/uploads`.
- Si compilas en Windows, evita subir binarios generados como `main.exe`.
