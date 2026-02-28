-- +goose Up
-- +goose StatementBegin
CREATE TABLE department (
                            id SERIAL PRIMARY KEY ,
                            name VARCHAR(200) NOT NULL,
                            parent_id INTEGER,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,


    -- Внешний ключ на сам таблицу Department (самоссылка)
                            CONSTRAINT fk_department_parent
                                FOREIGN KEY (parent_id)
                                    REFERENCES department(id)
                                    ON DELETE SET NULL
);

CREATE TABLE employee (
                          id SERIAL PRIMARY KEY,
                          department_id INTEGER NOT NULL,
                          full_name VARCHAR(255) NOT NULL,
                          position VARCHAR(255) NOT NULL,
                          hired_at DATE,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Ограничения на непустые строковые поля
                          CONSTRAINT employee_full_name_not_empty CHECK (LENGTH(TRIM(full_name)) > 0),
                          CONSTRAINT employee_position_not_empty CHECK (LENGTH(TRIM(position)) > 0),

    -- Внешний ключ на таблицу Department
                          CONSTRAINT fk_employee_department
                              FOREIGN KEY (department_id)
                                  REFERENCES Department(id)
                                  ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table Department,Employee;
-- +goose StatementEnd
