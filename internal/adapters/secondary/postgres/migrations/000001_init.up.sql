CREATE TABLE IF NOT EXISTS employee (
    id SERIAL PRIMARY KEY,
    balance INT NOT NULL,

    CONSTRAINT CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS transaction (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    amount INT NOT NULL,

    FOREIGN KEY (sender_id) REFERENCES employee (id),
    FOREIGN KEY (receiver_id) REFERENCES employee (id),
    CONSTRAINT CHECK (amount > 0),
    CONSTRAINT CHECK (sender_id != receiver_id)
);

CREATE TABLE IF NOT EXISTS merch (
    name VARCHAR PRIMARY KEY,
    cost INT,

    CONSTRAINT CHECK (cost > 0)
);

CREATE TABLE IF NOT EXISTS employee_merch (
    employee_id INT NOT NULL,
    merch_name VARCHAR NOT NULL,
    count INT NOT NULL,

    PRIMARY KEY (employee_id, merch_name),
    FOREIGN KEY (employee_id) REFERENCES employee(id),
    FOREIGN KEY (merch_name) REFERENCES merch(name),

    CONSTRAINT CHECK (count >= 0)
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