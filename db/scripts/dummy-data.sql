-- Insert dummy data into users table
-- INSERT INTO users (username, password) VALUES ('bob12', 'password123');
-- INSERT INTO users (username, password) VALUES ('alice9', 'password321');

-- Insert dummy data into categories table --
INSERT INTO categories (categoryName, categoryLimit)
    VALUES(
        "Seattle Dine Out",
        150
    );

INSERT INTO categories (categoryName, categoryLimit)
    VALUES(
        "Seattle Groceries",
        150
    );

INSERT INTO categories (categoryName, categoryLimit)
    VALUES(
        "Phx Dine Out",
        150
    );

INSERT INTO categories (categoryName, categoryLimit)
    VALUES(
        "Seattle Cleaning",
        50
    );
INSERT INTO categories (categoryName, categoryLimit)
    VALUES(
        "Phx Groceries",
        300
    );

-- Insert dummy data into expenses table
INSERT INTO expenses (userId, name, amount, categoryId) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'bob12'),
        'wingdom',
        30.50,
        (SELECT categoryId FROM categories WHERE categoryName = "Seattle Dine Out")
    );

INSERT INTO expenses (userId, name, amount, categoryId) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'bob12'),
        'groceries',
        35.18,
        (SELECT categoryId FROM categories WHERE categoryName = "Seattle Groceries")
    );

INSERT INTO expenses (userId, name, amount, categoryId) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'alice9'),
        'salad n go',
        12.45,
        (SELECT categoryId FROM categories WHERE categoryName = "Phx Dine Out")
    );

INSERT INTO expenses (userId, name, amount, categoryId) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'bob12'),
        'cleaning',
        5.55,
        (SELECT categoryId FROM categories WHERE categoryName = "Seattle Cleaning")
    );

INSERT INTO expenses (userId, name, amount, categoryId) 
    VALUES (
        (SELECT userId FROM users WHERE username = 'alice9'),
        'mouse food',
        15.60,
        (SELECT categoryId FROM categories WHERE categoryName = "Phx Groceries")
    );
