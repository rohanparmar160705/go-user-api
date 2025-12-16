/*
Package repository implements the data access layer.
The userRepository struct uses the generated SQLC code (`db.Queries`) to execute SQL queries
against the PostgreSQL database. It handles type conversions and data retrieval.
*/
package repository

import (
	"context"
	"time"

	db "github.com/rohanparmar/go-user-api/db/sqlc/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) Create(ctx context.Context, name string, dob string) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Dob:  parsePGDate(dob),
	})
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]db.User, error) {
	return r.queries.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *userRepository) Update(ctx context.Context, id int32, name string, dob string) (db.User, error) {
	return r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  parsePGDate(dob),
	})
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

// parsePGDate converts "YYYY-MM-DD" string to pgtype.Date
func parsePGDate(d string) pgtype.Date {
	t, _ := time.Parse("2006-01-02", d)
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}
