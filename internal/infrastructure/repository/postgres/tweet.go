package postgres

import (
	"context"
	"database/sql"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/postgres"
	"github.com/lib/pq"
)

type tweetRepo struct {
	db *postgres.PostgresDB
}

func NewTweetRepo(db *postgres.PostgresDB) repo.TweetStorageI {
	return &tweetRepo{
		db: db,
	}
}

func (t *tweetRepo) CreateTweet(ctx context.Context, tweet entity.CreateTweetRequest) (entity.CreateTweetResponse, error) {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return entity.CreateTweetResponse{}, err
	}

	insertTweetQuery := `
	INSERT INTO tweets (
	    user_id,
	    parent_tweet_id,
	    content
	) VALUES ($1, $2, $3)
	RETURNING
		id,
		user_id,
		parent_tweet_id,
		content
	`

	var response entity.CreateTweetResponse

	err = tx.QueryRow(
		ctx,
		insertTweetQuery,
		tweet.UserID,
		tweet.ParentTweetID,
		tweet.Content,
	).Scan(
		&response.ID,
		&response.UserID,
		&response.ParentTweetID,
		&response.Content,
	)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return entity.CreateTweetResponse{}, err
		}
		return entity.CreateTweetResponse{}, err
	}

	for _, url := range tweet.URLs {
		insertFileURLQuery := `INSERT INTO files (tweet_id, file_url) VALUES ($1, $2) RETURNING file_url`

		var savedURL string
		if err := tx.QueryRow(ctx, insertFileURLQuery, response.ID, url).Scan(&savedURL); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return entity.CreateTweetResponse{}, err
			}
			return entity.CreateTweetResponse{}, err
		}

		response.URLs = append(response.URLs, savedURL)
	}

	if err := tx.Commit(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return entity.CreateTweetResponse{}, err
		}
		return entity.CreateTweetResponse{}, err
	}

	return response, nil
}

func (t *tweetRepo) UpdateTweet(ctx context.Context, tweet entity.UpdateTweetRequest) (entity.UpdateTweetResponse, error) {
	query := `
	UPDATE
		tweets AS t
	SET
		content = $1
	WHERE
	    t.id = $2 AND t.deleted_at IS NULL AND t.parent_tweet_id IS NULL
	RETURNING
		t.id,
	    t.user_id,
		t.parent_tweet_id,
		t.content,
		COALESCE((SELECT array_agg(file_url) FROM files WHERE tweet_id = t.id AND t.deleted_at IS NULL), '{}')
	`

	var (
		urls     []string
		response entity.UpdateTweetResponse
	)
	err := t.db.QueryRow(ctx, query, tweet.Content, tweet.ID).Scan(
		&response.ID,
		&response.UserID,
		&response.ParentTweetID,
		&response.Content,
		pq.Array(&urls),
	)
	if err != nil {
		return entity.UpdateTweetResponse{}, err
	}

	response.URLs = urls

	return response, nil
}

func (t *tweetRepo) DeleteTweet(ctx context.Context, id string) error {
	query := `UPDATE tweets SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := t.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (t *tweetRepo) GetTweet(ctx context.Context, id string) (entity.GetTweetResponse, error) {
	query := `
	SELECT
		id,
		user_id,
		parent_tweet_id,
		content,
		COALESCE((SELECT array_agg(file_url) FROM files WHERE tweet_id = $1 AND deleted_at IS NULL), '{}')
	FROM
	    tweets
	WHERE
	    id = $1 AND deleted_at IS NULL
	`

	var (
		urls     []string
		response entity.GetTweetResponse
	)
	err := t.db.QueryRow(ctx, query, id).Scan(
		&response.ID,
		&response.UserID,
		&response.ParentTweetID,
		&response.Content,
		pq.Array(&urls),
	)
	if err != nil {
		return entity.GetTweetResponse{}, err
	}

	response.URLs = urls

	return response, nil
}

func (t *tweetRepo) ListTweets(ctx context.Context, filter entity.Filter) (entity.ListTweetsResponse, error) {
	query := `
	SELECT
		t.id,
		t.user_id,
		t.parent_tweet_id,
		t.content,
		COALESCE((SELECT array_agg(file_url) FROM files WHERE tweet_id = t.id AND deleted_at IS NULL), '{}')
	FROM
	    tweets AS t
	WHERE
	    t.deleted_at IS NULL
	LIMIT $1 OFFSET $2
	`

	var response entity.ListTweetsResponse
	offset := filter.Limit * (filter.Page - 1)

	rows, err := t.db.Query(ctx, query, filter.Limit, offset)
	if err != nil {
		return entity.ListTweetsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			urls  []string
			tweet entity.GetTweetResponse
		)
		err := rows.Scan(
			&tweet.ID,
			&tweet.UserID,
			&tweet.ParentTweetID,
			&tweet.Content,
			pq.Array(&urls),
		)
		if err != nil {
			return entity.ListTweetsResponse{}, err
		}

		tweet.URLs = urls

		response.Tweets = append(response.Tweets, tweet)
	}

	countQuery := `SELECT COUNT(*) FROM tweets WHERE deleted_at IS NULL`
	if err := t.db.QueryRow(ctx, countQuery).Scan(&response.Count); err != nil {
		return entity.ListTweetsResponse{}, err
	}

	return response, nil
}

func (t *tweetRepo) UserTweets(ctx context.Context, usrID string) (entity.ListTweetsResponse, error) {
	query := `
	SELECT
		t.id,
		t.user_id,
		t.parent_tweet_id,
		t.content,
		COALESCE((SELECT array_agg(file_url) FROM files WHERE tweet_id = t.id AND deleted_at IS NULL), '{}')
	FROM
	    tweets AS t
	WHERE
	    t.deleted_at IS NULL AND t.user_id = $1
	`

	var response entity.ListTweetsResponse

	rows, err := t.db.Query(ctx, query, usrID)
	if err != nil {
		return entity.ListTweetsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			tweet entity.GetTweetResponse
			urls  []string
		)
		err := rows.Scan(
			&tweet.ID,
			&tweet.UserID,
			&tweet.ParentTweetID,
			&tweet.Content,
			pq.Array(&urls),
		)
		if err != nil {
			return entity.ListTweetsResponse{}, err
		}

		tweet.URLs = urls

		response.Tweets = append(response.Tweets, tweet)
	}

	countQuery := `SELECT COUNT(*) FROM tweets WHERE deleted_at IS NULL and user_id = $1`
	if err := t.db.QueryRow(ctx, countQuery, usrID).Scan(&response.Count); err != nil {
		return entity.ListTweetsResponse{}, err
	}

	return response, nil
}
