CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name TEXT NOT NULL,
                          price NUMERIC NOT NULL,
                          category_id INTEGER REFERENCES categories(id)
);
