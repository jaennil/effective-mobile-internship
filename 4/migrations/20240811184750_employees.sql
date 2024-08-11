-- +goose Up
CREATE TABLE Employees(
    employee_id SERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    last_name TEXT,
    department_id INT,
    project_id INT,
    CONSTRAINT fk_department
        FOREIGN KEY (department_id)
        REFERENCES departments(department_id) ON DELETE SET NULL,
    CONSTRAINT fk_project
        FOREIGN KEY (project_id)
        REFERENCES projects(project_id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE Employees;
