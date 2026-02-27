-- +goose Up
-- +goose StatementBegin
create table Department(
                           id int primary key,
                           name varchar(200) not null,
                           parent_id int references Department (id),
                           created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
create table Employee(
                         id int primary key,
                         parent_id int references Department (id),
                         full_name varchar(200) not null,
                         position varchar(200) not null,
                         hired_at DATE,
                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table Department,Employee;
-- +goose StatementEnd
