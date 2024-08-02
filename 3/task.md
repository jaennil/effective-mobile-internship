### Задание 1. Создание базы данных
- Используя SQL, создайте базу данных с именем "Company" и тремя таблицами: "Employees", "Departments", и "Projects".
    ```sql
    CREATE DATABASE company;
    ```
    ```sql
    CREATE TABLE Departments(department_id SERIAL PRIMARY KEY NOT NULL, title TEXT NOT NULL);
    ```
    ```sql
    CREATE TABLE Projects(project_id SERIAL PRIMARY KEY NOT NULL, title TEXT NOT NULL);
    ```
    ```sql
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
    ```

### Задание 2. Вставка данных
- Вставьте несколько записей в каждую из созданных таблиц.
    ```sql
    INSERT INTO projects(title) VALUES('Time Tracker'), ('Android App');
    ```
    ```sql
    INSERT INTO departments(title) VALUES('IT'), ('Test');
    ```
    ```sql
    INSERT INTO employees(name, last_name, department_id, project_id)
    VALUES ('Nikita', 'Dubrovskih', 1, 1), ('Andrew', 'Glushkov', 2, 2);
    ```

### Задание 3. Выборка данных
- Напишите SQL-запрос для выбора всех сотрудников в отделе "IT".
    ```sql
    SELECT employees.*
    FROM employees
    NATURAL JOIN departments
    WHERE title = 'IT';
    ```

### Задание 4. Обновление данных
- Измените имя сотрудника с идентификатором 1 на "Robert".
    ```sql
    UPDATE Employees SET name = 'Robert' WHERE employee_id = 1;
    ```

### Задание 5. Удаление данных
- Удалите проект с идентификатором 2.
    ```sql
    DELETE FROM projects WHERE project_id = 2;
    ```

### Задание 6. Создание индексов
- Создайте индексы для ускорения поиска по полю "LastName" в таблице "Employees".
    ```sql
    CREATE INDEX ON employees(last_name);
    ```

### Задание 7. Использование агрегатных функций
- Напишите SQL-запрос для подсчета общего количества сотрудников в каждом отделе.
    ```sql
    SELECT department_id, title, COUNT(*)
    FROM employees
    NATURAL JOIN departments
    GROUP by department_id, title;
    ```

### Задание 8. Соединение таблиц
- Напишите SQL-запрос для получения списка сотрудников с именами, фамилиями и названиями их отделов.
    ```sql
    SELECT name, last_name, title
    FROM employees
    NATURAL JOIN departments;
    ```

### Задание 9. Транзакции
- Используйте транзакции для вставки нового отдела и проекта. Удостоверьтесь, что обе операции успешно выполнены, или обе отменены в случае ошибки.
    ```sql
    BEGIN;
    INSERT INTO projects(title) VALUES('test_transaction_projec_title');
    INSERT INTO departments(title) VALUES('test_transaction_department_title');
    COMMIT;
    ```

### Задание 10. Резервное копирование и восстановление
- Создайте резервную копию базы данных "Company".
    ```
    pg_dump company > company_dump
    ```
- Удалите случайные записи из таблиц "Employees", "Departments", и "Projects".
    ```sql
    DELETE FROM employees WHERE employee_id = 1;
    ```
    ```sql
    DELETE FROM departments WHERE department_id = 1;
    ```
    ```sql
    DELETE FROM projects WHERE project_id = 1;
    ```
- Восстановите базу данных из резервной копии.
    ```
    dropdb company
    createdb -T template0 company
    psql company < company_dump
    ```
