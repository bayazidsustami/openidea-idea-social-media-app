package repository

import (
	"context"
	"openidea-idea-social-media-app/customErr"
	"openidea-idea-social-media-app/models"
	friend_model "openidea-idea-social-media-app/models/friend"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FriendsRepository interface {
	Create(ctx context.Context, userFriends friend_model.Friend) error
	Delete(ctx context.Context, userFriends friend_model.Friend) error
	GetAll(ctx context.Context, userId int, filterFriend friend_model.FilterFriends) (friend_model.FriendDataPaging, error)
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
	SQL_DELETE_FRIENDS := "DELETE FROM friends WHERE user_id_requester = $1 AND user_id_accepter = $2 " +
		"OR user_id_accepter = $1 AND user_id_requester = $2"

	conn, err := repository.DBPool.Acquire(ctx)
	if err != nil {
		return customErr.ErrorInternalServer
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, SQL_DELETE_FRIENDS, userFriends.UserIdRequester, userFriends.UserIdAccepter)
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

func (repository *FriendsRepositoryImpl) GetAll(ctx context.Context, userId int, filterFriend friend_model.FilterFriends) (friend_model.FriendDataPaging, error) {
	query := filterFriend.BuildQuery(userId)

	conn, err := repository.DBPool.Acquire(ctx)
	if err != nil {
		return friend_model.FriendDataPaging{}, customErr.ErrorInternalServer
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return friend_model.FriendDataPaging{}, customErr.ErrorInternalServer
	}

	var friendsData []friend_model.FriendData
	var totalItem int
	for rows.Next() {
		friendData := friend_model.FriendData{}
		err := rows.Scan(
			&friendData.UserId,
			&friendData.Name,
			&friendData.ImageUrl,
			&friendData.CreatedAt,
			&totalItem,
		)
		if err != nil {
			return friend_model.FriendDataPaging{}, customErr.ErrorInternalServer
		}
		friendsData = append(friendsData, friendData)
	}

	return friend_model.FriendDataPaging{
		Data: friendsData,
		Meta: models.MetaPage{
			Limit:  filterFriend.Limit,
			Offset: filterFriend.Offset,
			Total:  totalItem,
		},
	}, nil
}
