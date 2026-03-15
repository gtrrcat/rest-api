-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tasks (title, description, completed) VALUES
('Изучить Go', 'Разобраться с миграциями', true),
('Написать API', 'Сделать эндпоинт для задач', false);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
