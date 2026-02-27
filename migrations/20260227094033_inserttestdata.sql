-- +goose Up
-- +goose StatementBegin
INSERT INTO Department (name, parent_id) VALUES
                                             ('Головной офис', NULL),
                                             ('IT-департамент', 1),
                                             ('Отдел разработки', 2),
                                             ('Отдел тестирования', 2),
                                             ('HR-департамент', 1),
                                             ('Бухгалтерия', 1);
INSERT INTO Employee (department_id, full_name, position, hired_at) VALUES
                                             (1, 'Иванов Иван Иванович', 'Генеральный директор', '2020-01-15'),
                                             (3, 'Петров Петр Петрович', 'Ведущий разработчик', '2020-03-20'),
                                             (3, 'Сидорова Мария Ивановна', 'Разработчик', '2021-06-10'),
                                             (4, 'Козлов Андрей Сергеевич', 'Тестировщик', '2021-09-05'),
                                             (5, 'Смирнова Елена Владимировна', 'HR-менеджер', '2020-02-01'),
                                             (6, 'Васильева Ольга Николаевна', 'Главный бухгалтер', '2020-01-20');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table Department;
truncate table Employee;
-- +goose StatementEnd
