-- +goose Up
CREATE TABLE Departments(department_id SERIAL PRIMARY KEY NOT NULL, title TEXT NOT NULL);

-- +goose Down
DROP TABLE Departments;
