-- Insert 3 categories
INSERT INTO categories (code, name) VALUES
('CAT001', 'Clothing'),
('CAT002', 'Shoes'),
('CAT003', 'Accessories');

-- Assign each product to a single category using product/category codes to look up ids

-- Clothing: PROD001, PROD004, PROD007
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'CAT001')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

-- Shoes: PROD002, PROD006
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'CAT002')
WHERE code IN ('PROD002', 'PROD006');

-- Accessories: PROD003, PROD005, PROD008
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'CAT003')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');
