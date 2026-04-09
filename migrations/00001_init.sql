-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     user_role TEXT NOT NULL,
                                     first_name TEXT NOT NULL,
                                     last_name TEXT NOT NULL,
                                     position TEXT NOT NULL,
                                     email TEXT UNIQUE NOT NULL,
                                     hashed_password TEXT NOT NULL,
                                     is_active BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS forms (
                                     id UUID PRIMARY KEY,
                                     user_id UUID NOT NULL,
                                     title TEXT NOT NULL,
                                     description TEXT,
                                     start_date TIMESTAMPTZ,
                                     end_date TIMESTAMPTZ,
                                     created_at TIMESTAMPTZ NOT NULL,
                                     reviewed_at TIMESTAMPTZ,
                                     status TEXT NOT NULL,
                                     comment TEXT,
                                     executor_id UUID,
                                     CONSTRAINT fk_forms_user FOREIGN KEY (user_id) REFERENCES users(id),
                                     CONSTRAINT fk_forms_executor FOREIGN KEY (executor_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_forms_executor_id ON forms(executor_id);
CREATE INDEX IF NOT EXISTS idx_forms_status ON forms(status);
CREATE INDEX IF NOT EXISTS idx_forms_user_id ON forms(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS forms;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
