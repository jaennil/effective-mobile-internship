-- +goose Up
CREATE TABLE Projects(project_id SERIAL PRIMARY KEY NOT NULL, title TEXT NOT NULL);

-- +goose Down
DROP TABLE Projects;
