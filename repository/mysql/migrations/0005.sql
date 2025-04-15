
-- +migrate Up
CREATE TABLE access_control(
                            `id` INT PRIMARY KEY  Auto_INCREMENT,
                            `actor_id` VARCHAR(191) NOT NULL UNIQUE,
                            `actor_type` ENUM('role','user') NOT NULL,
                            `permission_id` int NOT NULL,
                            `created_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (`permission_id`) REFERENCES `permission`(`id`)
);


-- +migrate Down
DROP TABLE access_control;
