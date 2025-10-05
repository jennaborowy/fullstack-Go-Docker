CREATE TABLE IF NOT EXISTS lists (
    id SERIAL PRIMARY KEY,            -- matches List.ID in Go
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,            -- matches Item.ID in Go
    title VARCHAR(255) NOT NULL,
    content TEXT,
    item_date DATE,                   -- matches Item.Date in Go
    list_id INT REFERENCES lists(id) ON DELETE CASCADE,  -- matches Item.ListID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sample data
INSERT INTO lists (title) VALUES
    ('Daily Tasks'),
    ('Goals');

INSERT INTO items (title, content, item_date, list_id) VALUES
    ('First Item', 'This is a test item', '2025-10-03', 1),
    ('Second Item', 'Another test item', '2025-10-07', 1),
    ('Third Item', 'One more item', '2025-12-25', 2);
