CREATE TABLE hosts (
    "id"              SERIAL       PRIMARY KEY,
    "host_name"       varchar(255) NOT NULL,
    "canonical_name"  varchar(255) NOT NULL,
    "url"             varchar(255),
    "ip"              varchar(255),
    "ipv6"            varchar(255),
    "location"        varchar(255),
    "os"              varchar(128),
    "active"          integer      NOT NULL DEFAULT '0',
    "created_at"      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at"      timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
  
CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON hosts
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();