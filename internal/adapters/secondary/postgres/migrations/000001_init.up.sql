CREATE TABLE IF NOT EXISTS employee (
    username VARCHAR PRIMARY KEY,
    hashed_password VARCHAR(255),
    balance INT NOT NULL DEFAULT 1000,

    CONSTRAINT balance_check CHECK (balance >= 0)
);

-- ограничения на username

CREATE TABLE IF NOT EXISTS transaction (
    id SERIAL PRIMARY KEY,
    sender_username VARCHAR NOT NULL,
    receiver_username VARCHAR NOT NULL,
    amount INT NOT NULL,

    FOREIGN KEY (sender_username) REFERENCES employee (username),
    FOREIGN KEY (receiver_username) REFERENCES employee (username),
    CONSTRAINT amount_check CHECK (amount > 0),
    CONSTRAINT username_check CHECK (sender_username != receiver_username)
);

CREATE TABLE IF NOT EXISTS merch (
    name VARCHAR PRIMARY KEY,
    cost INT,

    CONSTRAINT cost_check CHECK (cost > 0)
);

CREATE TABLE IF NOT EXISTS employee_merch (
    employee_username VARCHAR NOT NULL,
    merch_name VARCHAR NOT NULL,
    count INT NOT NULL,

    PRIMARY KEY (employee_username, merch_name),
    FOREIGN KEY (employee_username) REFERENCES employee(username),
    FOREIGN KEY (merch_name) REFERENCES merch(name),

    CONSTRAINT count_check CHECK (count >= 0)
);

INSERT INTO merch (name, cost)
VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);