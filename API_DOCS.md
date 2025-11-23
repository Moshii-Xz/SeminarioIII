# API Documentation

## Authentication
All protected endpoints require a valid JWT token in the `Authorization` header.
Format: `Authorization: Bearer <token>`

## Users Module

### POST /api/v1/users/
**Request:**
```json
{
  "nombre": "string (required, min=2, max=100)",
  "correo": "string (required, email)",
  "contrasena": "string (required, min=6)",
  "rol": "string (optional, values: 'comprador', 'tienda')"
}
```
**Response:**
```json
{
  "data": {
    "id_usuario": "uint",
    "nombre": "string",
    "correo": "string",
    "fecha_registro": "time"
  }
}
```

### GET /api/v1/users/:id
**Request:** None
**Response:**
```json
{
  "data": {
    "id_usuario": "uint",
    "nombre": "string",
    "correo": "string",
    "fecha_registro": "time"
  }
}
```

### PUT /api/v1/users/:id
**Request:**
```json
{
  "nombre": "string (optional)",
  "correo": "string (optional, email)"
}
```
**Response:**
```json
{
  "data": {
    "id_usuario": "uint",
    "nombre": "string",
    "correo": "string",
    "fecha_registro": "time"
  }
}
```

### DELETE /api/v1/users/:id
**Request:** None
**Response:**
```json
{
  "message": "user deleted"
}
```

### GET /api/v1/users/
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "data": [
    {
      "id_usuario": "uint",
      "nombre": "string",
      "correo": "string",
      "fecha_registro": "time"
    }
  ],
  "meta": {
    "total": "int64",
    "page": "int",
    "limit": "int"
  }
}
```

### POST /api/v1/users/login
**Request:**
```json
{
  "correo": "string (required, email)",
  "contrasena": "string (required)"
}
```
**Response:**
```json
{
  "data": {
    "token": "string",
    "usuario": {
      "id_usuario": "uint",
      "nombre": "string",
      "correo": "string",
      "fecha_registro": "time"
    }
  }
}
```

## Catalog Module (Products)

### POST /api/v1/products/upload-image (Protegido)
**Descripción:** Endpoint para subir una imagen desde el PC. Retorna la URL de la imagen que puede usarse en el campo `imagen_url` al crear o actualizar un producto.

**Request:** 
- Content-Type: `multipart/form-data`
- Body: Form data con campo `imagen` (archivo)
- Formatos permitidos: jpg, jpeg, png, gif, webp
- Tamaño máximo: 5MB

**Response:**
```json
{
  "data": {
    "imagen_url": "/uploads/images/1234567890_imagen.jpg",
    "filename": "1234567890_imagen.jpg"
  },
  "message": "image uploaded successfully"
}
```

**Ejemplo de uso:**
```javascript
const formData = new FormData();
formData.append('imagen', fileInput.files[0]);

fetch('/api/v1/products/upload-image', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer <token>'
  },
  body: formData
})
.then(res => res.json())
.then(data => {
  // Usar data.data.imagen_url en el campo imagen_url del producto
  console.log('URL de la imagen:', data.data.imagen_url);
});
```

### POST /api/v1/products/
**Request:**
```json
{
  "nombre": "string (required)",
  "descripcion": "string (optional)",
  "imagen_url": "string (optional, url, max=500)",
  "precio": "float64 (required, min=0)",
  "fecha_vencimiento": "time (required)",
  "stock": "int (optional, min=0)",
  "estado": "string (optional, values: 'En preparación', 'Listo para recoger', 'Entregado'. Default: 'En preparación')",
  "id_tienda": "uint (required)"
}
```
**Response:**
```json
{
  "data": {
    "id_producto": "uint",
    "nombre": "string",
    "descripcion": "string",
    "imagen_url": "string",
    "precio": "float64",
    "fecha_vencimiento": "time",
    "stock": "int",
    "estado": "string"
  }
}
```

### GET /api/v1/products/:id
**Request:** None
**Response:**
```json
{
  "data": {
    "id_producto": "uint",
    "nombre": "string",
    "descripcion": "string",
    "imagen_url": "string",
    "precio": "float64",
    "fecha_vencimiento": "time",
    "stock": "int",
    "estado": "string"
  }
}
```

### PUT /api/v1/products/:id
**Request:**
```json
{
  "nombre": "string (optional)",
  "descripcion": "string (optional)",
  "imagen_url": "string (optional, url, max=500)",
  "precio": "float64 (optional)",
  "fecha_vencimiento": "time (optional)",
  "stock": "int (optional)",
  "estado": "string (optional, values: 'En preparación', 'Listo para recoger', 'Entregado')"
}
```
**Response:**
```json
{
  "data": {
    "id_producto": "uint",
    "nombre": "string",
    "descripcion": "string",
    "imagen_url": "string",
    "precio": "float64",
    "fecha_vencimiento": "time",
    "stock": "int",
    "estado": "string"
  }
}
```

### DELETE /api/v1/products/:id
**Request:** None
**Response:**
```json
{
  "message": "product deleted"
}
```

### GET /api/v1/products/
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "productos": [
    {
      "id_producto": "uint",
      "nombre": "string",
      "descripcion": "string",
      "imagen_url": "string",
      "precio": "float64",
      "fecha_vencimiento": "time",
      "stock": "int",
      "estado": "string"
    }
  ],
  "total": "int64",
  "pagina": "int",
  "limite": "int"
}
```

## Orders Module

### POST /api/v1/orders/
**Request:**
```json
{
  "client_id": "uint (required)",
  "store_id": "uint (optional)",
  "items": [
    {
      "product_id": "uint (required)",
      "quantity": "int (required, min=1)",
      "unit_price": "float64 (required, min=0)"
    }
  ]
}
```
**Response:**
```json
{
  "data": {
    "id": "uint",
    "client_id": "uint",
    "store_id": "uint",
    "date": "time",
    "items": [],
    "total": "float64"
  }
}
```

### GET /api/v1/orders/:id
**Request:** None
**Response:**
```json
{
  "data": {
    "id": "uint",
    "client_id": "uint",
    "store_id": "uint",
    "date": "time",
    "items": [],
    "total": "float64"
  }
}
```

### PUT /api/v1/orders/:id
**Request:**
```json
{
  "store_id": "uint (optional)"
}
```
**Response:**
```json
{
  "data": {
    "id": "uint",
    "client_id": "uint",
    "store_id": "uint",
    "date": "time",
    "items": [],
    "total": "float64"
  }
}
```

### DELETE /api/v1/orders/:id
**Request:** None
**Response:**
```json
{
  "message": "order deleted"
}
```

### GET /api/v1/orders/
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "orders": [],
  "total": "int64",
  "page": "int",
  "limit": "int"
}
```

### POST /api/v1/orders/:id/items
**Request:**
```json
{
  "product_id": "uint (required)",
  "quantity": "int (required, min=1)",
  "unit_price": "float64 (required, min=0)"
}
```
**Response:**
```json
{
  "data": {
    "id": "uint",
    "product_id": "uint",
    "quantity": "int",
    "unit_price": "float64",
    "subtotal": "float64"
  }
}
```

### PUT /api/v1/orders/:id/items/:itemId
**Request:**
```json
{
  "quantity": "int (optional)",
  "unit_price": "float64 (optional)"
}
```
**Response:**
```json
{
  "data": {
    "id": "uint",
    "product_id": "uint",
    "quantity": "int",
    "unit_price": "float64",
    "subtotal": "float64"
  }
}
```

### DELETE /api/v1/orders/:id/items/:itemId
**Request:** None
**Response:**
```json
{
  "message": "item removed"
}
```

## Payments Module

### POST /api/v1/payments/
**Request:**
```json
{
  "order_id": "uint (required)",
  "payment_method_id": "uint (optional)",
  "amount": "float64 (required, min=0)"
}
```
**Response:**
```json
{
  "data": {
    "id": "uint",
    "order_id": "uint",
    "payment_method_id": "uint",
    "amount": "float64",
    "date": "time"
  }
}
```

### GET /api/v1/payments/:id
**Request:** None
**Response:**
```json
{
  "data": {
    "id": "uint",
    "order_id": "uint",
    "payment_method_id": "uint",
    "amount": "float64",
    "date": "time"
  }
}
```

### PUT /api/v1/payments/:id
**Request:**
```json
{
  "payment_method_id": "uint (optional)",
  "amount": "float64 (optional)"
}
```
**Response:**
```json
{
  "data": {
    "id": "uint",
    "order_id": "uint",
    "payment_method_id": "uint",
    "amount": "float64",
    "date": "time"
  }
}
```

### DELETE /api/v1/payments/:id
**Request:** None
**Response:**
```json
{
  "message": "payment deleted"
}
```

### GET /api/v1/payments/
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "payments": [],
  "total": "int64",
  "page": "int",
  "limit": "int"
}
```

### GET /api/v1/payments/order/:orderId
**Request:** None
**Response:**
```json
{
  "data": {
    "order_id": "uint",
    "order_total": "float64",
    "total_paid": "float64",
    "pending": "float64",
    "payments": []
  }
}
```

### GET /api/v1/payments/order/:orderId/status
**Request:** None
**Response:**
```json
{
  "data": {
    "order_id": "uint",
    "order_total": "float64",
    "total_paid": "float64",
    "pending": "float64",
    "payments": []
  }
}
```

### POST /api/v1/payments/methods
**Request:**
```json
{
  "name": "string (required)"
}
```
**Response:**
```json
{
  "data": {
    "id": "uint",
    "name": "string"
  }
}
```

### GET /api/v1/payments/methods
**Request:** None
**Response:**
```json
{
  "methods": [
    {
      "id": "uint",
      "name": "string"
    }
  ],
  "total": "int64"
}
```

### GET /api/v1/payments/methods/:id
**Request:** None
**Response:**
```json
{
  "data": {
    "id": "uint",
    "name": "string"
  }
}
```

### PUT /api/v1/payments/methods/:id
**Request:**
```json
{
  "name": "string (required)"
}
```
**Response:**
```json
{
  "data": {
    "id": "uint",
    "name": "string"
  }
}
```

### DELETE /api/v1/payments/methods/:id
**Request:** None
**Response:**
```json
{
  "message": "payment method deleted"
}
```

## Reviews Module

### POST /api/v1/reviews/
**Request:**
```json
{
  "id_producto": "uint (required)",
  "id_cliente": "uint (required)",
  "calificacion": "int (required, 1-5)",
  "comentario": "string (optional)"
}
```
**Response:**
```json
{
  "data": {
    "id_resena": "uint",
    "id_producto": "uint",
    "id_cliente": "uint",
    "calificacion": "int",
    "comentario": "string",
    "fecha_creacion": "time",
    "fecha_actualizacion": "time"
  }
}
```

### GET /api/v1/reviews/:id
**Request:** None
**Response:**
```json
{
  "data": {
    "id_resena": "uint",
    "id_producto": "uint",
    "id_cliente": "uint",
    "calificacion": "int",
    "comentario": "string",
    "fecha_creacion": "time",
    "fecha_actualizacion": "time"
  }
}
```

### PUT /api/v1/reviews/:id
**Request:**
```json
{
  "calificacion": "int (optional)",
  "comentario": "string (optional)"
}
```
**Response:**
```json
{
  "data": {
    "id_resena": "uint",
    "id_producto": "uint",
    "id_cliente": "uint",
    "calificacion": "int",
    "comentario": "string",
    "fecha_creacion": "time",
    "fecha_actualizacion": "time"
  }
}
```

### DELETE /api/v1/reviews/:id
**Request:** None
**Response:**
```json
{
  "message": "review deleted"
}
```

### GET /api/v1/reviews/
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "data": [],
  "meta": {
    "total": "int64",
    "page": "int",
    "limit": "int"
  }
}
```

### GET /api/v1/reviews/product/:productId
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "data": [],
  "meta": {
    "total": "int64",
    "page": "int",
    "limit": "int"
  }
}
```

## Reports Module

### GET /api/v1/reports/sales
**Request:** Query Params: `start_date` (YYYY-MM-DD), `end_date` (YYYY-MM-DD)
**Response:**
```json
{
  "data": {
    "total_ordenes": "int64",
    "total_ingresos": "float64",
    "ticket_promedio": "float64",
    "total_items_vendidos": "int64"
  }
}
```

### GET /api/v1/reports/stock
**Request:** Query Params: `threshold` (int, default 10)
**Response:**
```json
{
  "data": [
    {
      "id_producto": "uint",
      "nombre_producto": "string",
      "stock": "int"
    }
  ]
}
```

### GET /api/v1/reports/expiring
**Request:** Query Params: `days` (int, default 30)
**Response:**
```json
{
  "data": [
    {
      "product_id": "uint",
      "product_name": "string",
      "expiration_date": "time",
      "stock": "int"
    }
  ]
}
```

## Stores Module

### POST /api/v1/stores/ (Público - Registro de Tiendas)
**Descripción:** Endpoint público para registrar una nueva tienda. Crea automáticamente un usuario con rol "tienda" y su perfil de tienda asociado.

**Request:**
```json
{
  "nombre": "string (required, min=1, max=100)",
  "correo": "string (required, email)",
  "contrasena": "string (required, min=6)",
  "area_responsable": "string (optional, max=100)",
  "direccion": "string (optional, max=200)",
  "telefono": "string (optional, max=20)"
}
```
**Response:**
```json
{
  "data": {
    "id_tienda": "uint",
    "area_responsable": "string",
    "direccion": "string",
    "telefono": "string",
    "usuario": {
      "id_usuario": "uint",
      "nombre": "string",
      "correo": "string",
      "fecha_registro": "time"
    },
    "fecha_creacion": "time"
  }
}
```

### GET /api/v1/stores/
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "tiendas": [
    {
      "id_tienda": "uint",
      "area_responsable": "string",
      "direccion": "string",
      "telefono": "string",
      "usuario": {
        "id_usuario": "uint",
        "nombre": "string",
        "correo": "string",
        "fecha_registro": "time"
      },
      "fecha_creacion": "time"
    }
  ],
  "total": "int64",
  "pagina": "int",
  "limite": "int"
}
```

### GET /api/v1/stores/:id
**Request:** None
**Response:**
```json
{
  "data": {
    "id_tienda": "uint",
    "area_responsable": "string",
    "direccion": "string",
    "telefono": "string",
    "usuario": {
      "id_usuario": "uint",
      "nombre": "string",
      "correo": "string",
      "fecha_registro": "time"
    },
    "fecha_creacion": "time"
  }
}
```

### PUT /api/v1/stores/:id
**Request:**
```json
{
  "area_responsable": "string (optional, max=100)",
  "direccion": "string (optional, max=200)",
  "telefono": "string (optional, max=20)"
}
```
**Response:**
```json
{
  "data": {
    "id_tienda": "uint",
    "area_responsable": "string",
    "direccion": "string",
    "telefono": "string",
    "usuario": {
      "id_usuario": "uint",
      "nombre": "string",
      "correo": "string",
      "fecha_registro": "time"
    },
    "fecha_creacion": "time"
  }
}
```

### GET /api/v1/stores/:id/products
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "data": {
    "tienda": {
      "id_tienda": "uint",
      "area_responsable": "string",
      "direccion": "string",
      "telefono": "string",
      "usuario": {
        "id_usuario": "uint",
        "nombre": "string",
        "correo": "string",
        "fecha_registro": "time"
      },
      "fecha_creacion": "time"
    },
    "productos": [
      {
        "id_producto": "uint",
        "nombre": "string",
        "descripcion": "string",
        "imagen_url": "string",
        "precio": "float64",
        "fecha_vencimiento": "time",
        "stock": "int"
      }
    ],
    "total": "int64"
  }
}
```

### GET /api/v1/stores/:id/orders
**Request:** Query Params: `page` (int), `limit` (int)
**Response:**
```json
{
  "data": {
    "tienda": {
      "id_tienda": "uint",
      "area_responsable": "string",
      "direccion": "string",
      "telefono": "string",
      "usuario": {
        "id_usuario": "uint",
        "nombre": "string",
        "correo": "string",
        "fecha_registro": "time"
      },
      "fecha_creacion": "time"
    },
    "ordenes": [
      {
        "id": "uint",
        "client_id": "uint",
        "date": "time",
        "total": "float64"
      }
    ],
    "total": "int64"
  }
}
```
