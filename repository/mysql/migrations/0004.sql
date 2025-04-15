
-- +migrate Up
CREATE TABLE permission (
    `id` INT PRIMARY KEY  Auto_INCREMENT,
    `title` VARCHAR(191) NOT NULL UNIQUE,
    `created_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);




-- +migrate Down
DROP TABLE permission;
