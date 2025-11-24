CREATE TABLE rol (
    id_rol SERIAL PRIMARY KEY,
    nombre VARCHAR(50) UNIQUE NOT NULL
);
CREATE TABLE metodo_pago (
    id_metodo_pago SERIAL PRIMARY KEY,
    nombre VARCHAR(50) UNIQUE NOT NULL
);
CREATE TABLE usuario (
    id_usuario SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    correo VARCHAR(100) UNIQUE NOT NULL,
    contrasena VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE usuario_rol (
    id_usuario INT REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    id_rol INT REFERENCES rol(id_rol) ON DELETE CASCADE,
    PRIMARY KEY (id_usuario, id_rol)
);

CREATE TABLE cliente (
    id_cliente INT PRIMARY KEY REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    direccion VARCHAR(150),
    telefono VARCHAR(20)
);

CREATE TABLE tienda (
    id_tienda INT PRIMARY KEY REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    area_responsable VARCHAR(100),
    direccion VARCHAR(200),
    telefono VARCHAR(20)
);

CREATE TABLE administrador (
    id_admin INT PRIMARY KEY REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    permisos_especiales TEXT
);

CREATE TABLE categoria (
    id_categoria SERIAL PRIMARY KEY,
    nombre VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE producto (
    id_producto SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    imagen_url VARCHAR(500),
    precio NUMERIC(10,2) CHECK (precio >= 0),
    precio_original NUMERIC(10,2) CHECK (precio_original >= 0),
    precio_descuento NUMERIC(10,2) CHECK (precio_descuento >= 0),
    fecha_vencimiento DATE NOT NULL,
    stock INT DEFAULT 0 CHECK (stock >= 0),
    etiqueta VARCHAR(50) CHECK (etiqueta IN ('Oferta', 'Donación')),
    id_tienda INT NOT NULL REFERENCES tienda(id_tienda) ON DELETE CASCADE,
    id_categoria INT REFERENCES categoria(id_categoria) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE compra (
    id_compra SERIAL PRIMARY KEY,
    id_cliente INT NOT NULL REFERENCES cliente(id_cliente),
    id_tienda INT REFERENCES tienda(id_tienda),
    fecha_compra DATE DEFAULT CURRENT_DATE NOT NULL,
    estado VARCHAR(50) DEFAULT 'En preparación' CHECK (estado IN ('En preparación', 'Listo para recoger', 'Entregado'))
);

CREATE TABLE detalle_compra (
    id_detalle SERIAL PRIMARY KEY,
    id_compra INT NOT NULL REFERENCES compra(id_compra) ON DELETE CASCADE,
    id_producto INT NOT NULL REFERENCES producto(id_producto),
    cantidad INT NOT NULL CHECK (cantidad > 0),
    precio_unitario NUMERIC(10,2) NOT NULL CHECK (precio_unitario >= 0)
);

CREATE TABLE pago (
    id_pago SERIAL PRIMARY KEY,
    id_compra INT NOT NULL REFERENCES compra(id_compra) ON DELETE CASCADE,
    id_metodo_pago INT REFERENCES metodo_pago(id_metodo_pago),
    monto NUMERIC(10,2) NOT NULL CHECK (monto >= 0),
    fecha_pago DATE DEFAULT CURRENT_DATE NOT NULL
);

CREATE TABLE resena (
    id_resena SERIAL PRIMARY KEY,
    id_producto INT NOT NULL REFERENCES producto(id_producto) ON DELETE CASCADE,
    id_cliente INT NOT NULL REFERENCES cliente(id_cliente) ON DELETE CASCADE,
    calificacion INT NOT NULL CHECK (calificacion >= 1 AND calificacion <= 5),
    comentario TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

INSERT INTO rol (nombre) VALUES ('comprador');
INSERT INTO rol (nombre) VALUES ('tienda');
INSERT INTO rol (nombre) VALUES ('admin');

INSERT INTO categoria (nombre) VALUES ('Carnes y proteínas frescas');
INSERT INTO categoria (nombre) VALUES ('Lácteos y derivados');
INSERT INTO categoria (nombre) VALUES ('Frutas y verduras');
INSERT INTO categoria (nombre) VALUES ('Panadería y repostería fresca');
INSERT INTO categoria (nombre) VALUES ('Huevos');
INSERT INTO categoria (nombre) VALUES ('Congelados');
INSERT INTO categoria (nombre) VALUES ('Preparados listos para consumo');
INSERT INTO categoria (nombre) VALUES ('Bebidas');

-- Crear tienda de ejemplo con productos
-- Usuario tienda /Contraseña: password123
INSERT INTO usuario (nombre, correo, contrasena) VALUES 
('Supermercado Test', 'tienda@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy');

-- Asociar usuario con rol tienda (id_rol = 2)
INSERT INTO usuario_rol (id_usuario, id_rol) VALUES 
((SELECT id_usuario FROM usuario WHERE correo = 'tienda@test.com'), 2);

-- Crear perfil de tienda
INSERT INTO tienda (id_tienda, area_responsable, direccion, telefono) VALUES 
((SELECT id_usuario FROM usuario WHERE correo = 'tienda@test.com'), 'Alimentos', 'Calle Principal 123', '1234567890');

-- Crear productos: Leche, Pollo y Embutidos
-- Leche (categoría: Lácteos y derivados, id_categoria = 2) - Etiqueta: Oferta
INSERT INTO producto (nombre, descripcion, imagen_url, precio_original, precio_descuento, fecha_vencimiento, stock, etiqueta, id_tienda, id_categoria) VALUES 
('Leche', 'Leche entera fresca', '/uploads/images/ImagenLeche.jpg', 3500.00, 2500.00, CURRENT_DATE + INTERVAL '7 days', 50, 'Oferta', 
 (SELECT id_usuario FROM usuario WHERE correo = 'tienda@test.com'), 2);

-- Pollo (categoría: Carnes y proteínas frescas, id_categoria = 1) - Etiqueta: Donación
INSERT INTO producto (nombre, descripcion, imagen_url, fecha_vencimiento, stock, etiqueta, id_tienda, id_categoria) VALUES 
('Pollo', 'Pollo fresco entero', '/uploads/images/ImagenPollo.jpg', CURRENT_DATE + INTERVAL '3 days', 30, 'Donación', 
 (SELECT id_usuario FROM usuario WHERE correo = 'tienda@test.com'), 1);

-- Embutidos (categoría: Carnes y proteínas frescas, id_categoria = 1) - Etiqueta: Oferta
INSERT INTO producto (nombre, descripcion, imagen_url, precio_original, precio_descuento, fecha_vencimiento, stock, etiqueta, id_tienda, id_categoria) VALUES 
('Embutidos', 'Variedad de embutidos frescos', '/uploads/images/ImagenEmbutidos.jpg', 8500.00, 6000.00, CURRENT_DATE + INTERVAL '5 days', 25, 'Oferta', 
 (SELECT id_usuario FROM usuario WHERE correo = 'tienda@test.com'), 1);