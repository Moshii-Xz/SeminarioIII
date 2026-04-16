# Arquitectura de Expirapp

Este proyecto sigue una arquitectura modular por dominio. La idea central es mantener cada funcionalidad con su propia ruta de ejecucion: handler, service y repository.

## Flujo general

1. El router de Gin recibe la solicitud en `internal/server/gin_server.go`.
2. El handler del modulo valida y transforma la entrada.
3. El service aplica reglas de negocio.
4. El repository interactua con PostgreSQL mediante GORM.
5. La respuesta vuelve al cliente usando el formato JSON definido por la API.

## Modulos principales

- Usuarios: registro, login, consulta y administracion de perfiles.
- Catalogo: creacion, listado, edicion y carga de imagenes de productos.
- Orders: gestion de compras y sus detalles.
- Payments: registro de pagos y metodos de pago.
- Reviews: reseñas por producto y moderacion basica.
- Reports: reportes de ventas, stock y productos por vencer.
- Stores: registro y gestion de tiendas.

## Entidades de dominio

- User: usuario autenticado del sistema.
- Store: tienda que publica productos.
- Product: producto con precio, stock, imagen y fecha de vencimiento.
- Category: categoria de producto.
- Order y OrderDetail: compra y sus items.
- Payment y PaymentMethod: pago asociado a una compra.
- Review: reseña de un producto.

## Autenticacion y permisos

- Los endpoints protegidos usan JWT en el header `Authorization: Bearer <token>`.
- La creacion de productos, pedidos, pagos, reportes y gestion de tiendas requieren autenticacion.
- El catalogo publico permite listar productos y ver detalle sin token.

## Base de datos

- El proyecto usa PostgreSQL como persistencia principal.
- Los modelos de dominio ya declaran nombres de tabla y relaciones GORM.
- Las migraciones base estan en `migrations/AS_BD.sql`.

## Archivos clave

- `cmd/api/main.go`: arranque de la aplicacion y conexion a base de datos.
- `internal/server/gin_server.go`: configuracion del router y grupos de rutas.
- `internal/platform/database/postgres.go`: inicializacion de PostgreSQL.
- `internal/platform/http/response.go`: formato de respuestas.

## Convenciones utiles

- Los modulos se organizan en archivos separados por responsabilidad.
- Los nombres de campo en los modelos usan mapeo GORM hacia columnas en español.
- Las imagenes subidas se publican desde la ruta `/uploads`.
