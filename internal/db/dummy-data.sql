-- Insert dummy data into users table
INSERT INTO users (username, password) VALUES ('bob12', 'password123');
INSERT INTO users (username, password) VALUES ('alice9', 'password321');

-- Insert dummy data into expenses table
INSERT INTO expenses (userId, name, amount, category) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'bob12'),
        'wingdom',
        30.50,
        'sdo'
    );

INSERT INTO expenses (userId, name, amount, category) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'bob12'),
        'groceries',
        35.18,
        'sg'
    );

INSERT INTO expenses (userId, name, amount, category) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'alice9'),
        'salad n go',
        12.45,
        'pdo'
    );

INSERT INTO expenses (userId, name, amount, category) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'bob12'),
        'cleaning',
        5.55,
        'sc'
    );

INSERT INTO expenses (userId, name, amount, category) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'alice9'),
        'mouse food',
        15.60,
        'pg'
    );
