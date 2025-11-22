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
    contrasena VARCHAR(100) NOT NULL
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

CREATE TABLE vendedor (
    id_vendedor INT PRIMARY KEY REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    area_responsable VARCHAR(100)
);

CREATE TABLE administrador (
    id_admin INT PRIMARY KEY REFERENCES usuario(id_usuario) ON DELETE CASCADE,
    permisos_especiales TEXT
);

CREATE TABLE producto (
    id_producto SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    precio NUMERIC(10,2) NOT NULL CHECK (precio >= 0),
    fecha_vencimiento DATE NOT NULL,
    stock INT DEFAULT 0 CHECK (stock >= 0),
    id_vendedor INT NOT NULL REFERENCES vendedor(id_vendedor) ON DELETE CASCADE
);

CREATE TABLE compra (
    id_compra SERIAL PRIMARY KEY,
    id_cliente INT NOT NULL REFERENCES cliente(id_cliente),
    id_vendedor INT REFERENCES vendedor(id_vendedor),
    fecha_compra DATE DEFAULT CURRENT_DATE NOT NULL
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

INSERT INTO rol (nombre) VALUES ('comprador');
INSERT INTO rol (nombre) VALUES ('vendedor');
INSERT INTO rol (nombre) VALUES ('admin');