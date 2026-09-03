
-- Ensure no duplicated product codes

CREATE UNIQUE INDEX IF NOT EXISTS products_code_idx ON products (code);

-- Ensure no duplicated product variant SKUs

CREATE UNIQUE INDEX IF NOT EXISTS product_variant_sku_idx ON product_variants (sku);

-- Speed up filtering products by price

CREATE INDEX IF NOT EXISTS products_price_idx ON products (price);

-- Speed up looking up variants by product (preload/cascade delete)

CREATE INDEX IF NOT EXISTS product_variants_product_id_idx ON product_variants (product_id);
