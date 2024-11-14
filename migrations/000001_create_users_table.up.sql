CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'Active',
    referrer_id INT REFERENCES users(id);
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,  
    price INT NOT NULL, 
    is_active BOOLEAN DEFAULT TRUE 
);

INSERT INTO tasks (name, description, price, is_active)
VALUES 
    ('Telegram', 'Sub', 100, TRUE),
    ('YouTube', 'Watch the video', 150, TRUE),
    ('Twitter', 'Sub', 50, TRUE),
    ('Referrer', 'Add', 200, TRUE);

CREATE TABLE user_tasks (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    task_id INT NOT NULL,
    completed BOOLEAN DEFAULT FALSE,  
    completion_date TIMESTAMP, 
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    CONSTRAINT unique_user_task UNIQUE (user_id, task_id) 
);
