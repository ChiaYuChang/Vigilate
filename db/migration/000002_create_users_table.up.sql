CREATE TABLE users (
    "id"              SERIAL          PRIMARY KEY,
    "first_name"      varchar(255)    NOT NULL,
    "last_name"       varchar(255)    NOT NULL,
    "user_active"     integer         NOT NULL DEFAULT '0',
    "access_level"    integer         NOT NULL DEFAULT '3',
    "email"           varchar(255)    NOT NULL,
    "password"        varchar(60)     NOT NULL,
    "created_at"      timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"      timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"      timestamp
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
   FOR EACH ROW
   EXECUTE PROCEDURE trigger_set_timestamp();