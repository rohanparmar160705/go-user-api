/*
Package repository defines the interface for data access operations.
The UserRepository interface acts as a contract for interacting with the database,
allowing for dependency injection and easier unit testing (by mocking the repository).
*/
package repository

import (
	"context"

	db "github.com/rohanparmar/go-user-api/db/sqlc/generated"
)

type UserRepository interface {
	Create(ctx context.Context, name string, dob string) (db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
	List(ctx context.Context, limit, offset int32) ([]db.User, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, id int32, name string, dob string) (db.User, error)
	Delete(ctx context.Context, id int32) error
}
