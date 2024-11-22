package user

import (
	"context"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/pkg/errors"
)

type UserService struct {
	DB *db.DBPool
}

func GetUser(ctx context.Context) (*db.User, error) {
	// TODO: Entender se pode dar problema usar esses context
	user, ok := ctx.Value("user").(*db.User)
	if !ok {
		return nil, errors.New("could not get user instance from context")
	}

	return user, nil
}
