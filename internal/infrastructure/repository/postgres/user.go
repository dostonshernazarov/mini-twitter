package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
	"github.com/spf13/cast"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

// UniqueUsername check have or no username in users data
func (u *userRepo) UniqueUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = $1 AND deleted_at IS NULL`

	var count int
	if err := u.db.QueryRowContext(ctx, query, username).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

// UniqueEmail check have or no email in users data
func (u *userRepo) UniqueEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1 AND deleted_at IS NULL`

	var count int
	if err := u.db.QueryRowContext(ctx, query, email).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

// Create method for creating a new user
func (u *userRepo) Create(ctx context.Context, user entity.CreateUserRequest) (entity.CreateUserResponse, error) {
	query := `
	INSERT INTO users (
	    name,
	    username,
	    email,
	    password,
	    role
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING
		id,
		name,
		username,
		email,
		role
	`

	var response entity.CreateUserResponse
	err := u.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
	).Scan(
		&response.ID,
		&response.Name,
		&response.Username,
		&response.Email,
		&response.Role,
	)

	if err != nil {
		return entity.CreateUserResponse{}, err
	}

	return response, nil
}

// Update method for updating user data
func (u *userRepo) Update(ctx context.Context, user entity.UpdateUserRequest) (entity.UpdateUserResponse, error) {
	query := `
	UPDATE
		users
	SET
	    name = $1,
	    bio = $2,
	    username = $3,
	    updated_at = NOW()
	WHERE
	    id = $4
		AND deleted_at IS NULL
	RETURNING
		id,
		name,
		username,
	    email,
		bio,
	    role,
		profile_picture
	`

	var response entity.UpdateUserResponse
	err := u.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Bio,
		user.Username,
		user.ID,
	).Scan(
		&response.ID,
		&response.Name,
		&response.Username,
		&response.Email,
		&response.Bio,
		&response.Role,
		&response.ProfilePicture,
	)

	if err != nil {
		return entity.UpdateUserResponse{}, err
	}

	return response, nil
}

// UpdatePasswd method for updating user password with id
func (u *userRepo) UpdatePasswd(ctx context.Context, id int, passwd string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`

	result, err := u.db.ExecContext(ctx, query, passwd, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (u *userRepo) UploadImage(ctx context.Context, id int, url string) error {
	query := `UPDATE users SET profile_picture = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`

	result, err := u.db.ExecContext(ctx, query, url, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete method for deleting user with id
func (u *userRepo) Delete(ctx context.Context, id int) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Get method for getting user with id or email
func (u *userRepo) Get(ctx context.Context, field map[string]interface{}) (entity.GetUserResponse, error) {
	query := `
	SELECT
		u.id,
		u.name,
		u.username,
		u.email,
		u.bio,
		u.role,
		u.password,
		u.profile_picture,
		(SELECT COUNT(*) FROM follows WHERE user_id = u.id AND deleted_at IS NULL),
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id AND deleted_at is null)
	FROM
	    users AS u
	WHERE
	    u.deleted_at IS NULL
	`

	if field["id"] != nil {
		query += fmt.Sprintf(" AND u.id = '%d'", cast.ToInt(field["id"]))
	} else if field["username"] != nil {
		query += fmt.Sprintf(" AND u.username = '%s'", cast.ToString(field["username"]))
	} else if field["email"] != nil {
		query += fmt.Sprintf(" AND u.email = '%s'", cast.ToString(field["email"]))
	}
	fmt.Println(query)

	var response entity.GetUserResponse
	err := u.db.QueryRowContext(ctx, query).Scan(
		&response.ID,
		&response.Name,
		&response.Username,
		&response.Email,
		&response.Bio,
		&response.Role,
		&response.Password,
		&response.ProfilePicture,
		&response.FollowingCount,
		&response.FollowersCount,
	)

	if err != nil {
		return entity.GetUserResponse{}, err
	}

	return response, nil
}

// List method for getting list users and count with filter fields
func (u *userRepo) List(ctx context.Context, filter entity.Filter) (entity.ListUser, error) {
	query := `
	SELECT
		u.id,
		u.name,
		u.username,
		u.email,
		u.bio,
		u.role,
		u.profile_picture,
		(SELECT COUNT(*) FROM follows WHERE user_id = u.id AND deleted_at IS NULL),
		(SELECT COUNT(*) FROM follows WHERE following_id = u.id AND deleted_at is null)
	FROM
	    users AS u
	WHERE
	    u.deleted_at IS NULL
		AND u.role = 'user'
	LIMIT $1
	OFFSET $2
	`

	var response entity.ListUser

	offset := filter.Limit * (filter.Page - 1)
	rows, err := u.db.QueryContext(ctx, query, filter.Limit, offset)
	if err != nil {
		return entity.ListUser{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.GetUserResponse

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.Bio,
			&user.Role,
			&user.ProfilePicture,
			&user.FollowingCount,
			&user.FollowersCount,
		)

		if err != nil {
			return entity.ListUser{}, err
		}

		response.Users = append(response.Users, user)
	}

	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND role = 'user'`
	if err := u.db.QueryRowContext(ctx, countQuery).Scan(&response.Count); err != nil {
		return entity.ListUser{}, err
	}

	return response, nil
}
