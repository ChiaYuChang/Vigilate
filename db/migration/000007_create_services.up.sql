CREATE TABLE services (
    id           SERIAL         PRIMARY KEY,
    service_name varchar(255)   NOT NULL,
    active       integer        NOT NULL DEFAULT '1',
    icon         varchar(255)   NOT NULL,
    created_at   timestamp      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   timestamp      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

 CREATE TRIGGER set_timestamp
 BEFORE UPDATE ON services
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();