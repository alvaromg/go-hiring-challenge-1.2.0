CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32),
    name VARCHAR(256) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Ensure no duplicated product codes

CREATE UNIQUE INDEX IF NOT EXISTS categories_code_idx ON categories (code);

-- Ensure no duplicated category names

CREATE UNIQUE INDEX IF NOT EXISTS categories_name_idx ON categories (name);

-- Categories/products relationship as many to one (a product belongs to at most one category)

ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS products_category_id_idx ON products (category_id);
