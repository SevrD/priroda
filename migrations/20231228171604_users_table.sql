-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    tgID BIGINT UNIQUE,
    login TEXT NOT NULL,
    name TEXT NOT NULL,
    createData TIMESTAMP,
    phone TEXT,
    chatID BIGINT,
    ban BOOLEAN
);
CREATE TABLE chatStatuses (
    tgID BIGINT UNIQUE,
    status BIGINT,
    annID BIGINT
);
CREATE TABLE announcements (
    id INTEGER PRIMARY KEY,
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
DROP TABLE chatStatuses;
DROP TABLE announcements;
-- +goose StatementEnd
