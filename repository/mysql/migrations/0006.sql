-- +migrate Up
CREATE TABLE players(
                               `id` INT PRIMARY KEY  Auto_INCREMENT,
                               `user_id` int NOT NULL,
                               `game_id` int NOT NULL,
                               `score` int NOT NULL DEFAULT 0,
                               `created_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE player_answers (
                                id INT  PRIMARY KEY Auto_INCREMENT,
                                player_id INT NOT NULL,
                                question_id INT NOT NULL,
                                choice INT NOT NULL,
                                FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);



-- +migrate Down
DROP TABLE palyers;
DROP TABLE palyer_answers;
