package mysqlQuestion

import (
	"GameApp/entity"
	"log"
)

func (db *DB) getQuestion(category string, questionLimit uint) ([]entity.Question, error) {

	rows, err := db.conn.Conn().Query(`(SELECT id, text, correct_answer_id, difficulty, category
         FROM questions WHERE difficulty = 1 AND category = ?
         ORDER BY RAND() LIMIT ?)`, category, questionLimit)
	if err != nil {
		return nil, err
	}
	var questions []entity.Question
	for rows.Next() {
		var q entity.Question
		err := rows.Scan(&q.ID, &q.Text, &q.CorrectAnswerID, &q.Difficulty, &q.Category)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		questions = append(questions, q)
		return questions, nil
	}
	return []entity.Question{}, nil

}
