package postgres

import (
	"context"
	"fmt"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/postgres"
	"github.com/lib/pq"
)

type searchRepo struct {
	db *postgres.PostgresDB
}

func NewSearchRepo(db *postgres.PostgresDB) repo.SearchStorageI {
	return &searchRepo{
		db: db,
	}
}

// Search method for searching users or tweets with text
func (s *searchRepo) Search(ctx context.Context, data string) (entity.SearchResponse, error) {
	var response entity.SearchResponse

	searchUsers := fmt.Sprintf(`
	SELECT
		id,
		name,
		username,
		email,
		bio,
		role,
		profile_picture
	FROM
	    users
	WHERE
	    deleted_at IS NULL
		AND role = 'user'
		AND name ILIKE '%s'
	`, "%"+data+"%")

	userRows, err := s.db.QueryContext(ctx, searchUsers)
	if err != nil {
		return entity.SearchResponse{}, err
	}
	defer userRows.Close()

	for userRows.Next() {
		var user entity.GetUserResponse

		err := userRows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.Bio,
			&user.Role,
			&user.ProfilePicture,
		)
		if err != nil {
			return entity.SearchResponse{}, err
		}

		response.Users = append(response.Users, user)
	}

	searchTweets := fmt.Sprintf(`
	SELECT
		id,
		user_id,
		parent_tweet_id,
		content,
		COALESCE((SELECT array_agg(file_url) FROM files WHERE tweet_id = id AND deleted_at IS NULL), '{}')
	FROM
		tweets
	WHERE
		deleted_at IS NULL AND content ILIKE '%s'
	`, "%"+data+"%")

	tweetRows, err := s.db.QueryContext(ctx, searchTweets)
	if err != nil {
		return entity.SearchResponse{}, err
	}
	defer tweetRows.Close()

	for tweetRows.Next() {
		var (
			urls  []string
			tweet entity.GetTweetResponse
		)

		err := tweetRows.Scan(
			&tweet.ID,
			&tweet.UserID,
			&tweet.ParentTweetID,
			&tweet.Content,
			pq.Array(&urls),
		)
		if err != nil {
			return entity.SearchResponse{}, err
		}

		tweet.URLs = urls

		response.Tweets = append(response.Tweets, tweet)
	}

	return response, nil
}
