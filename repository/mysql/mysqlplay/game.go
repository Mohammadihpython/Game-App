package mysqlplay

import (
	"GameApp/entity"
	"context"
	"time"
)

func (d *DB) createGame(ctx context.Context, category entity.Category, questions []entity.Question) (entity.Game, error) {
	tx, err := d.conn.Conn().BeginTx(ctx, nil)
	if err != nil {
		return entity.Game{}, err
	}
	res, err := tx.Exec(
		`INSERT INTO games(category,start_time) VALUES (?,?)`, category, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return entity.Game{}, err
	}
	var GameID, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return entity.Game{}, err
	}
	stmt, err := tx.Prepare(`INSERT INTO games_questions(game_id,question_id) VALUES (?,?)`)
	if err != nil {
		tx.Rollback()
		return entity.Game{}, err
	}
	defer stmt.Close()

	for _, q := range questions {
		_, err = stmt.Exec(GameID, q.ID)
		if err != nil {
			tx.Rollback()
			return entity.Game{}, err
		}

	}
	err = tx.Commit()
	if err != nil {
		return entity.Game{}, err
	}
	return entity.Game{
		ID:          uint(GameID),
		Category:    string(category),
		QuestionIDs: questions,
		StartTime:   time.Now(),
		WinnerID:    0,
	}, nil

}

func (d *DB) createPlayer(ctx context.Context, userID entity.User, gameID uint) (int, error) {
	tx, err := d.conn.Conn().BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	res, err := tx.Exec(
		"INSERT INTO players(user_id,game_id,score) VALUES(?,?,?)", userID, gameID,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	playerID, _ := res.LastInsertId()
	return int(playerID), nil
}
