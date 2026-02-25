-- Seed data for development/testing
-- Run: docker exec -i cqrs_postgres psql -U user -d cqrs_db < db/seeds/seed.sql

-- Clear existing data (optional, comment out if you want to preserve data)
TRUNCATE orders, products, users RESTART IDENTITY CASCADE;

-- Users
INSERT INTO users (name, email, age, created_at, updated_at) VALUES
  ('John Doe', 'john@example.com', 30, NOW(), NOW()),
  ('Jane Smith', 'jane@example.com', 25, NOW(), NOW()),
  ('Bob Wilson', 'bob@example.com', 35, NOW(), NOW());

-- Products
INSERT INTO products (name, price, created_at, updated_at) VALUES
  ('PS5', 499.99, NOW(), NOW()),
  ('Xbox Series X', 499.99, NOW(), NOW()),
  ('Nintendo Switch', 299.99, NOW(), NOW()),
  ('Gaming Headset', 79.99, NOW(), NOW());

-- Orders
INSERT INTO orders (user_id, product_id, quantity, created_at, updated_at) VALUES
  (1, 1, 1, NOW(), NOW()),  -- John buys PS5
  (1, 4, 2, NOW(), NOW()),  -- John buys 2 headsets
  (2, 2, 1, NOW(), NOW()),  -- Jane buys Xbox
  (3, 3, 1, NOW(), NOW());  -- Bob buys Switch
