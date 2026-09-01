
-- Ensure no duplicated product codes

CREATE UNIQUE INDEX IF NOT EXISTS products_code_idx ON products (code);

-- Ensure no duplicated product variant SKUs

CREATE UNIQUE INDEX IF NOT EXISTS product_variant_sku_idx ON product_variants (sku);
