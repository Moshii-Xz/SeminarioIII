# Diagrama de Base de Datos

El archivo [database-diagram.puml](database-diagram.puml) contiene el diagrama ER de VenciTrack basado en las clases de dominio y en `migrations/AS_BD.sql`.

Incluye:

- `usuario`, `rol` y la tabla puente `usuario_rol`.
- Los perfiles derivados `cliente`, `tienda` y `administrador`.
- Catálogo de productos con `categoria`.
- `compra` y `detalle_compra`.
- `pago` y `metodo_pago`.
- `resena`.

Si quieres renderizarlo en VS Code, puedes usar una extensión de PlantUML y abrir el archivo `.puml` directamente.
