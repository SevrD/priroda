-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tgID BIGINT UNIQUE,
    login TEXT NOT NULL,
    name TEXT NOT NULL,
    createData TIMESTAMP,
    phone TEXT,
    chatID BIGINT,
    ban BOOLEAN
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE chatStatuses (
    tgID BIGINT UNIQUE,
    status BIGINT,
    annID BIGINT
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE announcements (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tgID BIGINT,
    chatID BIGINT,
    txt TEXT,
    status BIGINT,
    admMsgID BIGINT,
    fileID TEXT,
    publicID BIGINT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE chatStatuses;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE announcements;
-- +goose StatementEnd
