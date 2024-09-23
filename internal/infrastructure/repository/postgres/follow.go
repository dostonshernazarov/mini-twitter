package postgres

import (
	"context"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/postgres"
)

type followRepo struct {
	db *postgres.PostgresDB
}

func NewFollowRepo(db *postgres.PostgresDB) repo.FollowStorageI {
	return &followRepo{
		db: db,
	}
}

func (f *followRepo) Follow(ctx context.Context, follow entity.FollowAction) (bool, error) {
	countQuery := `SELECT COUNT(*) FROM follows WHERE user_id = $1 AND following_id = $2`

	var count int
	if err := f.db.QueryRowContext(ctx, countQuery, follow.UserID, follow.FollowingID).Scan(&count); err != nil {
		return false, err
	}

	if count == 0 {
		insertQuery := `INSERT INTO follows (user_id, following_id) VALUES ($1, $2)`

		_, err := f.db.ExecContext(ctx, insertQuery, follow.UserID, follow.FollowingID)
		if err != nil {
			return false, err
		}

		return true, nil
	} else {
		deleteQuery := `DELETE FROM follows WHERE user_id = $1 AND following_id = $2`

		_, err := f.db.ExecContext(ctx, deleteQuery, follow.UserID, follow.FollowingID)
		if err != nil {
			return false, err
		}

		return false, nil
	}
}

func (f *followRepo) GetFollowings(ctx context.Context, id int) (entity.ListUser, error) {
	query := `
	SELECT
		u.id,
		u.name,
		u.username,
		u.email,
		u.role,
		u.bio,
		u.profile_picture
	FROM
	    users AS u
	INNER JOIN
	    follows AS f ON u.id = f.following_id
	WHERE
	    u.deleted_at IS NULL AND f.user_id = $1
	`

	rows, err := f.db.QueryContext(ctx, query, id)
	if err != nil {
		return entity.ListUser{}, err
	}
	defer rows.Close()

	var response entity.ListUser
	for rows.Next() {
		var user entity.GetUserResponse
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.Bio,
			&user.ProfilePicture,
		)

		if err != nil {
			return entity.ListUser{}, err
		}

		response.Users = append(response.Users, user)
	}

	countQuery := `SELECT COUNT(*) FROM follows WHERE user_id = $1`
	if err := f.db.QueryRowContext(ctx, countQuery, id).Scan(&response.Count); err != nil {
		return entity.ListUser{}, err
	}

	return response, nil
}

func (f *followRepo) GetFollowers(ctx context.Context, id int) (entity.ListUser, error) {
	query := `
	SELECT
		u.id,
		u.name,
		u.username,
		u.email,
		u.role,
		u.bio,
		u.profile_picture
	FROM
	    users AS u
	INNER JOIN
	    follows AS f ON u.id = f.user_id
	WHERE
	    u.deleted_at IS NULL AND f.following_id = $1
	`

	rows, err := f.db.QueryContext(ctx, query, id)
	if err != nil {
		return entity.ListUser{}, err
	}
	defer rows.Close()

	var response entity.ListUser
	for rows.Next() {
		var user entity.GetUserResponse
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.Bio,
			&user.ProfilePicture,
		)

		if err != nil {
			return entity.ListUser{}, err
		}

		response.Users = append(response.Users, user)
	}

	countQuery := `SELECT COUNT(*) FROM follows WHERE following_id = $1`
	if err := f.db.QueryRowContext(ctx, countQuery, id).Scan(&response.Count); err != nil {
		return entity.ListUser{}, err
	}

	return response, nil
}
