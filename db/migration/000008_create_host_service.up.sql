CREATE TABLE IF NOT EXISTS host_services (
    id              SERIAL       PRIMARY KEY,
    host_id         integer      NOT NULL,
    service_id      integer      NOT NULL,
    active          integer      NOT NULL DEFAULT '1',
    schedule_number integer      NOT NULL DEFAULT '3',
    schedule_unit   varchar(8)   NOT NULL DEFAULT 'm',
    status          integer      NOT NULL DEFAULT '0',
    last_message    varchar(255) NOT NULL DEFAULT '',
    last_check      timestamp    NOT NULL DEFAULT '0001-01-01 00:00:01',
    created_at      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON host_services
   FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

ALTER TABLE host_services 
  ADD FOREIGN KEY (host_id) 
    REFERENCES hosts (id)
    ON DELETE CASCADE 
    ON UPDATE CASCADE;

ALTER TABLE host_services 
  ADD FOREIGN KEY (service_id) 
    REFERENCES services (id)
    ON DELETE CASCADE 
    ON UPDATE CASCADE;