-- Crear la base de datos
CREATE DATABASE crm_users;

-- Conectar a la base de datos
\c crm_users;

-- Crear extensiones útiles
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- Para búsquedas por texto

-- GORM creará las tablas automáticamente con AutoMigrate
-- Pero si quieres crear manualmente con índices optimizados:

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(40) NOT NULL,
    dni VARCHAR(20) NOT NULL,
    gender VARCHAR(3) NOT NULL CHECK (gender IN ('M', 'F')),
    phone VARCHAR(15) NOT NULL,
    email VARCHAR(50) NOT NULL,
    birth_day TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- ✅ CRÍTICO: Índices únicos compuestos por tenant
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_tenant_dni 
ON users(tenant_id, dni) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_tenant_email 
ON users(tenant_id, email) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_tenant_phone 
ON users(tenant_id, phone) WHERE deleted_at IS NULL;

-- ✅ Índice para queries por tenant
CREATE INDEX IF NOT EXISTS idx_users_tenant_id 
ON users(tenant_id) WHERE deleted_at IS NULL;

-- ✅ Índice para soft deletes
CREATE INDEX IF NOT EXISTS idx_users_deleted_at 
ON users(deleted_at);

-- ✅ Índice para búsquedas por nombre (opcional)
CREATE INDEX IF NOT EXISTS idx_users_name_search 
ON users USING gin(to_tsvector('spanish', first_name || ' ' || last_name));

-- Trigger para updated_at automático
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at 
BEFORE UPDATE ON users 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();

-- Insertar datos de prueba (opcional)
-- Tenant 1
INSERT INTO users (tenant_id, first_name, last_name, dni, gender, phone, email, birth_day)
VALUES 
(1, 'Juan', 'Pérez', '12345678', 'M', '+573001234567', 'juan@example.com', '1990-01-15'),
(1, 'María', 'González', '87654321', 'F', '+573007654321', 'maria@example.com', '1992-05-20');

-- Tenant 2
INSERT INTO users (tenant_id, first_name, last_name, dni, gender, phone, email, birth_day)
VALUES 
(2, 'Carlos', 'Rodríguez', '11223344', 'M', '+573009876543', 'carlos@example.com', '1988-12-10'),
(2, 'Ana', 'Martínez', '44332211', 'F', '+573003456789', 'ana@example.com', '1995-08-25');

-- Verificar inserción
SELECT * FROM users;
