CREATE TABLE IF NOT EXISTS role
(
    id          INTEGER     NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(20) NOT NULL DEFAULT '',
    status_code TINYINT     NOT NULL DEFAULT 0,
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE user
    ADD CONSTRAINT role_fk
        FOREIGN KEY (id_role) REFERENCES role (id);
