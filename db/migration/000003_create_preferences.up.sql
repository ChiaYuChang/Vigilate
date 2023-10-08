CREATE TABLE preferences (
    "id"         SERIAL         PRIMARY KEY,
    "name"       VARCHAR(255)   NOT NULL,
    "preference" text,
    "created_at" timestamp      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON preferences
   FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();