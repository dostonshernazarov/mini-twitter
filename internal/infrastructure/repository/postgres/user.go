package postgres

import (
	"context"
	"database/sql"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/postgres"
	"github.com/jackc/pgx"
	"github.com/spf13/cast"
)

type userRepo struct {
	db        *postgres.PostgresDB
	tableName string
}

func NewUserRepo(db *postgres.PostgresDB) repo.UserStorageI {
	return &userRepo{
		db:        db,
		tableName: "users",
	}
}

func (u *userRepo) UniqueUsername(ctx context.Context, username string) (bool, error) {
	qrBuilder := u.db.Sq.Builder.Select("COUNT(*)")
	qrBuilder = qrBuilder.From(u.tableName)
	qrBuilder = qrBuilder.Where(u.db.Sq.Equal("username", username))
	qrBuilder = qrBuilder.Where("deleted_at IS NULL")

	var count int
	query, args, err := qrBuilder.ToSql()
	if err != nil {
		return false, err
	}

	err = u.db.QueryRow(ctx, query, args...).Scan(
		&count,
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (u *userRepo) UniqueEmail(ctx context.Context, email string) (bool, error) {
	qrBuilder := u.db.Sq.Builder.Select("COUNT(*)")
	qrBuilder = qrBuilder.From(u.tableName)
	qrBuilder = qrBuilder.Where(u.db.Sq.Equal("email", email))
	qrBuilder = qrBuilder.Where("deleted_at IS NULL")

	var count int
	query, args, err := qrBuilder.ToSql()
	if err != nil {
		return false, err
	}

	err = u.db.QueryRow(ctx, query, args...).Scan(
		&count,
	)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (u *userRepo) Create(ctx context.Context, user entity.CreateUserRequest) (entity.CreateUserResponse, error) {
	clauses := map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"username":   user.Username,
		"email":      user.Email,
		"password":   user.Password,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	queryBuilder := u.db.Sq.Builder.Insert(u.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.CreateUserResponse{}, err
	}

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return entity.CreateUserResponse{}, err
	}

	if result.RowsAffected() == 0 {
		return entity.CreateUserResponse{}, pgx.ErrNoRows
	}

	return entity.CreateUserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		Bio:            new(string),
		Role:           user.Role,
		ProfilePicture: new(string),
	}, nil

}

func (u *userRepo) Update(ctx context.Context, user entity.UpdateUserRequest) error {

	clauses := map[string]interface{}{
		"name":     user.Name,
		"bio":      user.Bio,
		"username": user.Username,
	}

	queryBuilder := u.db.Sq.Builder.Update(u.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(u.db.Sq.Equal("id", user.ID))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (u *userRepo) UpdatePasswd(ctx context.Context, id int, passwd string) error {
	clauses := map[string]interface{}{
		"password": passwd,
	}

	queryBuilder := u.db.Sq.Builder.Update(u.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(u.db.Sq.Equal("id", id))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (u *userRepo) UploadImage(ctx context.Context, id int, url string) error {
	clauses := map[string]interface{}{
		"profile_picture": url,
	}

	queryBuilder := u.db.Sq.Builder.Update(u.tableName)
	queryBuilder = queryBuilder.SetMap(clauses)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(u.db.Sq.Equal("id", id))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (u *userRepo) Delete(ctx context.Context, id int) error {
	queryBuilder := u.db.Sq.Builder.Update(u.tableName)
	queryBuilder = queryBuilder.Set("deleted_at", "NOW()")
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(u.db.Sq.Equal("id", id))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (u *userRepo) Get(ctx context.Context, field map[string]interface{}) (entity.GetUserResponse, error) {
	var (
		result   entity.GetUserResponse
		NullBio  sql.NullString
		NulPhoto sql.NullString
	)
	queryBuilder := u.db.Sq.Builder.Select(
		"id, " +
			"name, " +
			"username, " +
			"email, " +
			"bio, " +
			"role, " +
			"password, " +
			"profile_picture, " +
			"(SELECT COUNT(*) FROM follows WHERE user_id = id AND deleted_at IS NULL), " +
			"(SELECT COUNT(*) FROM follows WHERE following_id = id AND deleted_at is null)")

	queryBuilder = queryBuilder.From(u.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	if field["id"] != nil {
		queryBuilder = queryBuilder.Where("id", cast.ToString(field["id"]))
	} else if field["username"] != nil {
		queryBuilder = queryBuilder.Where("username", cast.ToString(field["username"]))
	} else if field["email"] != nil {
		queryBuilder = queryBuilder.Where("email", cast.ToString(field["email"]))
	}

	selectQuery, selectArgs, err := queryBuilder.ToSql()
	if err != nil {
		return entity.GetUserResponse{}, err
	}

	err = u.db.QueryRow(ctx, selectQuery, selectArgs...).Scan(
		&result.ID,
		&result.Name,
		&result.Username,
		&result.Email,
		&NullBio,
		&result.Role,
		&result.Password,
		&NulPhoto,
		&result.FollowingCount,
		&result.FollowersCount,
	)

	if err != nil {
		return entity.GetUserResponse{}, err
	}

	if NulPhoto.Valid {
		result.ProfilePicture = &NulPhoto.String
	}

	if NullBio.Valid {
		result.Bio = &NullBio.String
	}

	return result, nil
}

func (u *userRepo) List(ctx context.Context, filter entity.Filter) (entity.ListUser, error) {
	queryBuilder := u.db.Sq.Builder.Select(
		"id, " +
			"name, " +
			"username, " +
			"email, " +
			"bio, " +
			"role, " +
			"password, " +
			"profile_picture, " +
			"(SELECT COUNT(*) FROM follows WHERE user_id = id AND deleted_at IS NULL), " +
			"(SELECT COUNT(*) FROM follows WHERE following_id = id AND deleted_at is null)")

	queryBuilder = queryBuilder.From(u.tableName)
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where("role", "user")
	queryBuilder = queryBuilder.Limit(uint64(filter.Limit))
	queryBuilder = queryBuilder.Offset(uint64(filter.Limit) * (uint64(filter.Page) - 1))

	var response entity.ListUser

	selectQuery, selectArgs, err := queryBuilder.ToSql()
	if err != nil {
		return entity.ListUser{}, err
	}

	rows, err := u.db.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return entity.ListUser{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			user     entity.GetUserResponse
			NullBio  sql.NullString
			NulPhoto sql.NullString
		)

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&NullBio,
			&user.Role,
			&NulPhoto,
			&user.FollowingCount,
			&user.FollowersCount,
		)

		if err != nil {
			return entity.ListUser{}, err
		}

		if NulPhoto.Valid {
			user.ProfilePicture = &NulPhoto.String
		}

		if NullBio.Valid {
			user.Bio = &NullBio.String
		}

		response.Users = append(response.Users, user)

		countBuilder := u.db.Sq.Builder.Select("COUNT(*)")
		countBuilder = countBuilder.From(u.tableName)
		countBuilder = countBuilder.Where("deleted_at IS NULL")
		countBuilder = countBuilder.Where("role", "user")
	}

	return response, nil
}
