CREATE TABLE remember_tokens (
  id             SERIAL       PRIMARY KEY,
  user_id        integer      NOT NULL,
  created_at     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  remember_token varchar(100) NOT NULL
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON remember_tokens
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

ALTER TABLE remember_tokens 
  ADD FOREIGN KEY (user_id) 
    REFERENCES users (id)
    ON DELETE SET NULL 
    ON UPDATE CASCADE;