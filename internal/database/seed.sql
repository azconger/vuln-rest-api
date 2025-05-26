-- Additional test data for the vulnerable REST API

-- Insert more test users with weak credentials
INSERT INTO users (username, password, email, role) VALUES
    ('admin', 'admin123', 'admin@example.com', 'admin'),
    ('test', 'test123', 'test@example.com', 'user'),
    ('guest', 'guest123', 'guest@example.com', 'guest');

-- Insert more test products
INSERT INTO products (name, description, price) VALUES
    ('Vulnerable Product 1', 'Contains SQL injection vulnerability', 99.99),
    ('Vulnerable Product 2', 'Contains XSS vulnerability', 149.99),
    ('Vulnerable Product 3', 'Contains path traversal vulnerability', 199.99);

-- Insert more test comments with XSS payloads
INSERT INTO comments (user_id, content) VALUES
    (1, '<img src=x onerror=alert("XSS1")>'),
    (2, '<svg/onload=alert("XSS2")>'),
    (3, '"><script>alert("XSS3")</script>');

-- Insert more test files with path traversal examples
INSERT INTO files (name, path, size) VALUES
    ('passwd', '/etc/passwd', 2048),
    ('shadow', '/etc/shadow', 1024),
    ('config', '/etc/config.json', 4096);

-- Insert test commands
INSERT INTO commands (command, output) VALUES
    ('ls -la', 'total 1234\ndrwxr-xr-x ...'),
    ('pwd', '/app'),
    ('whoami', 'root');

-- Insert test XML data
INSERT INTO xml_data (content) VALUES
    ('<?xml version="1.0"?><!DOCTYPE test [<!ENTITY xxe SYSTEM "file:///etc/passwd">]><test>&xxe;</test>'),
    ('<?xml version="1.0"?><root><user>admin</user><password>admin123</password></root>'); 