-- +migrate Up
CREATE TABLE game(
                        `id` INT PRIMARY KEY  Auto_INCREMENT,
                        `category` int NOT NULL,
                        `created_at`  TIMESTAMP DEFAULT CURRENT_TIMESTAMP



);


CREATE TABLE possible_answers (
                                  `id` INT  PRIMARY KEY Auto_INCREMENT,
                                  `question_id` INT NOT NULL,
                                  `player_id` INT NOT NULL ,
                                  `text` TEXT not null ,
                                  `choice` TINYINT NOT NULL  CHECK ( choice BETWEEN 1 AND 4),

                                  FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
                                  FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);


CREATE TABLE questions (
                                `id` INT  PRIMARY KEY Auto_INCREMENT,
                                `text` TEXT NOT NULL,
                                correct_answer_id INT DEFAULT NULL, -- FK to possible_answers.id
                                difficulty TINYINT NOT NULL CHECK ( difficulty BETWEEN 1 AND 3),
                                category INT not null

);





-- Add foreign key constraint for correct_answer_id (after both tables exist)
ALTER TABLE questions
ADD CONSTRAINT fk_correct_answer
FOREIGN KEY (correct_answer_id) REFERENCES possible_answers(id)
        ON DELETE SET NULL;



-- +migrate Down
DROP TABLE game;
DROP TABLE questions;
DROP TABLE possible_answers;
