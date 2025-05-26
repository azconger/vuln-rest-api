-- Database schema for the vulnerable REST API

-- Users table (vulnerable to SQL injection)
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL, -- Intentionally stored in plain text
    email TEXT,
    role TEXT DEFAULT 'user', -- User role (admin, user, guest)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table (vulnerable to session fixation)
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id),
    token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

-- Products table (vulnerable to SQL injection)
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price REAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Orders table (vulnerable to SQL injection)
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id),
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER,
    total_price REAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Comments table (vulnerable to XSS)
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id),
    content TEXT NOT NULL, -- Intentionally vulnerable to XSS
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Files table (vulnerable to path traversal)
CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    path TEXT NOT NULL, -- Intentionally vulnerable to path traversal
    size INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Commands table (vulnerable to command injection)
CREATE TABLE commands (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    command TEXT NOT NULL, -- Intentionally vulnerable to command injection
    output TEXT,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- XML data table (vulnerable to XXE)
CREATE TABLE xml_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL, -- Intentionally vulnerable to XXE
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default user (intentionally weak credentials)
INSERT INTO users (username, password, email, role) 
VALUES ('user', 'password', 'user@example.com', 'user');

-- Insert sample products
INSERT INTO products (name, description, price) VALUES
    ('Product 1', 'Description 1', 10.99),
    ('Product 2', 'Description 2', 20.99),
    ('Product 3', 'Description 3', 30.99);

-- Insert sample comments (intentionally vulnerable to XSS)
INSERT INTO comments (user_id, content) VALUES
    (1, '<script>alert("XSS")</script>'),
    (1, 'Normal comment');

-- Insert sample files (intentionally vulnerable to path traversal)
INSERT INTO files (name, path, size) VALUES
    ('test.txt', '/var/www/files/test.txt', 1024),
    ('config.json', '/etc/config.json', 2048); 