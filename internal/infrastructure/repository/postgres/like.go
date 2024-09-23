package postgres

import (
	"context"
	"database/sql"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
)

type likeRepo struct {
	db *sql.DB
}

func NewLikeRepo(db *sql.DB) repo.LikeStorageI {
	return &likeRepo{
		db: db,
	}
}

// Like method for liking or unliking tweet
func (l *likeRepo) Like(ctx context.Context, like entity.LikeAction) (bool, error) {
	countQuery := `SELECT COUNT(*) FROM likes WHERE user_id = $1 AND tweet_id = $2`

	var count int
	if err := l.db.QueryRowContext(ctx, countQuery, like.UserID, like.TweetID).Scan(&count); err != nil {
		return false, err
	}

	if count == 0 {
		insertQuery := `INSERT INTO likes (user_id, tweet_id) VALUES ($1, $2)`

		_, err := l.db.ExecContext(ctx, insertQuery, like.UserID, like.TweetID)
		if err != nil {
			return false, err
		}

		return true, nil
	} else {
		deleteQuery := `DELETE FROM likes WHERE user_id = $1 AND tweet_id = $2`

		_, err := l.db.ExecContext(ctx, deleteQuery, like.UserID, like.TweetID)
		if err != nil {
			return false, err
		}

		return false, nil
	}
}
