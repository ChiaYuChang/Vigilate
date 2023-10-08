CREATE TABLE IF NOT EXISTS events (
    id              SERIAL          PRIMARY KEY,
    type            VARCHAR(128),
    host_service_id integer         NOT NULL,
    host_id         integer         NOT NULL,
    service_id      integer         NOT NULL,
    message         VARCHAR(512),
    created_at      timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp       NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON events
   FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();