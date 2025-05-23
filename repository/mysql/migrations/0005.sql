
-- +migrate Up
CREATE TABLE access_control(
                            `id` INT PRIMARY KEY  Auto_INCREMENT,
                            `actor_id` int NOT NULL,
                            `actor_type` ENUM('role','mysqluser') NOT NULL,
                            `permission_id` int NOT NULL,
                            `created_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            FOREIGN KEY (`permission_id`) REFERENCES `permission`(`id`)
);


-- +migrate Down
DROP TABLE access_control;
