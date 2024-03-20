package repository

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	friend_model "openidea-idea-social-media-app/models/friend"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FriendsRepository interface {
	Create(ctx context.Context, userFriends friend_model.Friend) error
	Delete(ctx context.Context, userFriends friend_model.Friend) error
	Get()
}

type FriendsRepositoryImpl struct {
	DBPool *pgxpool.Pool
}

func NewFriendRepository(
	dbPool *pgxpool.Pool,
) FriendsRepository {
	return &FriendsRepositoryImpl{
		DBPool: dbPool,
	}
}

func (repository *FriendsRepositoryImpl) Create(ctx context.Context, userFriends friend_model.Friend) error {
	SQL_ADD_FRIENDS := "INSERT INTO friends(user_id_requester, user_id_accepter) VALUES ($1, $2), ($2, $1) " +
		"ON CONFLICT (user_id_requester, user_id_accepter) DO NOTHING"
	conn, err := repository.DBPool.Acquire(ctx)
	if err != nil {
		return customErr.ErrorInternalServer
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, SQL_ADD_FRIENDS, userFriends.UserIdRequester, userFriends.UserIdAccepter)
	if err != nil {
		if err == pgx.ErrNoRows {
			return customErr.ErrorBadRequest
		} else {
			return customErr.ErrorInternalServer
		}
	}

	if res.RowsAffected() == 0 {
		return customErr.ErrorBadRequest
	}

	return nil
}

func (repository *FriendsRepositoryImpl) Delete(ctx context.Context, userFriends friend_model.Friend) error {
	return nil
}

func (repository *FriendsRepositoryImpl) Get() {

}
