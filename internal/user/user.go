package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/db"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles operations related to users.
type UserService struct {
	db db.Querier
}

// New creates a new instance of UserService.
func New(db db.Querier) UserService {
	return UserService{
		db: db,
	}
}

// CreateUser adds a new user to the database.
func (u *UserService) CreateUser(ctx context.Context, user *db.CreateUserParams) (*db.User, error) {
	user.ID = xid.New().String()

	var err error
	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user.Email = strings.ToLower(user.Email)

	createdUser, err := u.db.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

// GetUserByID retrieves a user by their ID.
func (u *UserService) GetUserByID(ctx context.Context, id string) (*db.User, error) {
	return u.db.GetUserByID(ctx, id)
}

// GetUserByEmail retrieves a user by their email.
func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	return u.db.GetUserByEmail(ctx, strings.ToLower(email))
}

// CheckUserPassword checks if a users password is valid.
func (u *UserService) CheckUserPassword(ctx context.Context, email, password string) (bool, *db.User, error) {
	user, err := u.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return false, nil, err
	}

	return CheckPasswordHash(password, user.Password), user, nil
}

// UpdateUserPassword updates a user's password.
func (u *UserService) UpdateUserPassword(ctx context.Context, id, newPassword string) error {
	return u.db.UpdateUserPassword(ctx, &db.UpdateUserPasswordParams{
		ID:       id,
		Password: newPassword,
	})
}

// IsUserValidated checks if a user is validated.
func (u *UserService) IsUserValidated(ctx context.Context, id string) (bool, error) {
	result, err := u.db.IsUserValidated(ctx, id)
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// DeleteUser logically deletes a user by setting the deleted_at timestamp.
func (u *UserService) DeleteUser(ctx context.Context, id string) error {
	return u.db.DeleteUser(ctx, id)
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plaintext password with a bcrypt hashed password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
