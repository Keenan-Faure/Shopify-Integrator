delete from users where email = 'test@test.com';
delete from products where product_code = 'product_code';

delete from products;
delete from customers;
delete from orders;
delete from register_tokens;
delete from shopify_location;
delete from shopify_inventory;
delete from inventory_location;
delete from shopify_vid;
delete from shopify_pid;
delete from queue_items;

select * from orders;
select * from customer_orders;
select * from customers;
select * from products;
select * from variants;
select * from product_options;
select * from customers;
select * from inventory_location;
select * from shopify_inventory;
select * from shopify_location;

select * from fetch_stats;

select * from queue_items;

select * from runtime_flags;
update runtime_flags
set flag_value = TRUE
where id = '23f6b38c-acf4-451e-b071-d143734c0889';

SELECT id, active, product_code, title, category, vendor, product_type, updated_at FROM products WHERE product_type = 'simple' AND category = 'test';

select * from orders where id = 'bf3e2161-0f9a-4ddf-b84c-2426ae00b508';
