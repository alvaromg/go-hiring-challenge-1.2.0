-- Insert 3 categories
INSERT INTO categories (code, name) VALUES
('CAT001', 'Clothing'),
('CAT002', 'Shoes'),
('CAT003', 'Accessories');

-- Assign products to categories using product/category codes to look up ids

-- Clothing: PROD001, PROD004, PROD007
INSERT INTO product_categories (product_id, category_id) VALUES
((SELECT id FROM products WHERE code = 'PROD001'), (SELECT id FROM categories WHERE code = 'CAT001')),
((SELECT id FROM products WHERE code = 'PROD004'), (SELECT id FROM categories WHERE code = 'CAT001')),
((SELECT id FROM products WHERE code = 'PROD007'), (SELECT id FROM categories WHERE code = 'CAT001'));

-- Shoes: PROD002, PROD006
INSERT INTO product_categories (product_id, category_id) VALUES
((SELECT id FROM products WHERE code = 'PROD002'), (SELECT id FROM categories WHERE code = 'CAT002')),
((SELECT id FROM products WHERE code = 'PROD006'), (SELECT id FROM categories WHERE code = 'CAT002'));

-- Accessories: PROD003, PROD005, PROD008
INSERT INTO product_categories (product_id, category_id) VALUES
((SELECT id FROM products WHERE code = 'PROD003'), (SELECT id FROM categories WHERE code = 'CAT003')),
((SELECT id FROM products WHERE code = 'PROD005'), (SELECT id FROM categories WHERE code = 'CAT003')),
((SELECT id FROM products WHERE code = 'PROD008'), (SELECT id FROM categories WHERE code = 'CAT003'));
