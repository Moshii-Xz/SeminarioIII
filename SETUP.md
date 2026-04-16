# Guía de Configuración del Backend (VenciTrack)

Esta guía está diseñada para ayudar al equipo de frontend a levantar el servidor backend localmente de manera rápida y sencilla.

## 1. Prerrequisitos

Asegúrate de tener instalado lo siguiente en tu máquina:

*   **Go (Golang)**: Versión 1.24 o superior. [Descargar aquí](https://go.dev/dl/).
*   **PostgreSQL**: Base de datos relacional. [Descargar aquí](https://www.postgresql.org/download/).
*   **Git**: Para clonar el repositorio.

## 2. Configuración de la Base de Datos

El backend espera una base de datos PostgreSQL corriendo localmente con una configuración específica.

### Paso 2.1: Crear la Base de Datos y Usuario

Por defecto, el código está configurado para usar:
*   **Host**: `localhost`
*   **Puerto**: `5432`
*   **Usuario**: `postgres`
*   **Contraseña**: `1234`
*   **Nombre de BD**: `vencitrack`

Si tu configuración local de PostgreSQL es diferente, deberás ajustar el archivo `cmd/api/main.go` o configurar tu base de datos para que coincida.

**Comandos SQL para configurar la BD (ejecutar en pgAdmin o psql):**

```sql
-- 1. Crear la base de datos
CREATE DATABASE vencitrack;

-- 2. (Opcional) Si necesitas cambiar la contraseña del usuario postgres para que coincida con el código:
-- ALTER USER postgres WITH PASSWORD '1234';
```

### Paso 2.2: Crear las Tablas (Migración)

Debes ejecutar el script SQL que crea la estructura de la base de datos.

1.  Ubica el archivo `migrations/AS_BD.sql` en este proyecto.
2.  Ejecuta el contenido de ese archivo en tu base de datos `vencitrack` usando tu herramienta preferida (pgAdmin, DBeaver, TablePlus, o línea de comandos).

**Opción línea de comandos:**
```bash
psql -U postgres -d vencitrack -f migrations/AS_BD.sql
```

O descomenta las lineas 
```go
	//migrationsPath := filepath.Join(".", "migrations")
	//if err := database.RunMigrations(db, migrationsPath); err != nil {
	//	log.Fatalf("failed to run migrations: %v", err)
	//}
```
en el archivo main.go

## 3. Ejecutar el Servidor

Una vez que la base de datos está lista, puedes iniciar el servidor.

1.  Abre una terminal en la carpeta raíz del proyecto (`VenciTrack`).
2.  Descarga las dependencias:
    ```bash
    go mod download
    ```
3.  Inicia el servidor:
    ```bash
    go run cmd/api/main.go
    ```

Deberías ver un mensaje como:
> Starting server on :8080

## 4. Verificar que funciona

Abre tu navegador o Postman y visita:

`http://localhost:8080/health`

Deberías recibir una respuesta JSON:
```json
{
    "message": "dabase is ok",
    "status": "healthy"
}
```

## 5. Solución de Problemas Comunes

**Error: "database connection failed"**
*   Verifica que PostgreSQL esté corriendo.
*   Verifica que la base de datos `vencitrack` exista.
*   Verifica que el usuario `postgres` y la contraseña `1234` sean correctos. Si tu contraseña es diferente, edita el archivo `cmd/api/main.go` línea 16:
    ```go
    Password: "TU_CONTRASEÑA_REAL",
    ```

**Error: "bind: address already in use"**
*   Otro servicio está usando el puerto 8080. Cierra ese servicio o cambia el puerto en `internal/server/server_config.go`.
