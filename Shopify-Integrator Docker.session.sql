-- name: UpsertProduct :one
INSERT INTO products(
    id,
    product_code,
    active,
    title,
    body_html,
    category,
    vendor,
    product_type,
    created_at,
    updated_at
) VALUES (uuid_generate_v4(),'GenImp-SongPiner', '1', 'Song of Broken Pines', '<p>A greatsword as light as the sigh of grass in the breeze, yet as merciless to the corrupt as a typhoon.</p>\n<ul>\n<li>Claymore weapon</li>\n<li>5 star</li>\n<li>Base stat - Phys. damage</li>\n</ul>', 'Genshin Impact', 'Mihoyo', 'Weapon', '2024-01-24 13:05:34.558834', '2024-01-24 13:05:34.558834')
ON CONFLICT(product_code)
DO UPDATE 
SET
    active = COALESCE('1', products.active),
    title = COALESCE('Song of Broken Pines', products.title),
    body_html = COALESCE('<p>A greatsword as light as the sigh of grass in the breeze, yet as merciless to the corrupt as a typhoon.</p>\n<ul>\n<li>Claymore weapon</li>\n<li>5 star</li>\n<li>Base stat - Phys. damage</li>\n</ul>', products.body_html),
    category = COALESCE('Genshin Impact', products.category),
    vendor = COALESCE('Mihoyo', products.vendor),
    product_type = COALESCE('Weapon', products.product_type),
    updated_at = '2024-01-24 13:05:34.558834'
RETURNING *, (xmax = 0) AS inserted;

select * from products where product_code = 'GenImp-SongPines'