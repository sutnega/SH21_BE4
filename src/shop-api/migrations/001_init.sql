CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE gender_type AS ENUM ('male', 'female', 'other');

CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country VARCHAR(100) NOT NULL,
    city    VARCHAR(100) NOT NULL,
    street  VARCHAR(200) NOT NULL
);

CREATE TABLE categories (
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE images (
    id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image BYTEA NOT NULL
);

CREATE TABLE suppliers (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(200) NOT NULL,
    address_id   UUID NOT NULL REFERENCES addresses(id) ON DELETE RESTRICT,
    phone_number VARCHAR(20)  NOT NULL
);

CREATE TABLE clients (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_name       VARCHAR(100) NOT NULL,
    client_surname    VARCHAR(100) NOT NULL,
    birthday          DATE NOT NULL,
    gender            gender_type NOT NULL,
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    address_id        UUID NOT NULL REFERENCES addresses(id) ON DELETE RESTRICT
);

CREATE TABLE products (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(200) NOT NULL,
    category_id       UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    price             NUMERIC(12,2) NOT NULL CHECK (price >= 0),
    available_stock   INT NOT NULL CHECK (available_stock >= 0),
    last_update_date  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    supplier_id       UUID NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,
    image_id          UUID REFERENCES images(id) ON DELETE SET NULL
);

CREATE INDEX idx_clients_name_surname ON clients (client_name, client_surname);
CREATE INDEX idx_products_stock       ON products (available_stock);
