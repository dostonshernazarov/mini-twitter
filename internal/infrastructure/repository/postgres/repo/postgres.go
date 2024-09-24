package repo

import (
	"context"
	"github.com/dostonshernazarov/mini-twitter/internal/entity"
)

type UserStorageI interface {
	UniqueUsername(ctx context.Context, username string) (bool, error)
	UniqueEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user entity.CreateUserRequest) (entity.CreateUserResponse, error)
	Update(ctx context.Context, user entity.UpdateUserRequest) error
	UpdatePasswd(ctx context.Context, id string, passwd string) error
	UploadImage(ctx context.Context, id string, url string) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, field map[string]interface{}) (entity.GetUserResponse, error)
	List(ctx context.Context, filter entity.Filter) (entity.ListUser, error)
}

type TweetStorageI interface {
	CreateTweet(ctx context.Context, tweet entity.CreateTweetRequest) (entity.CreateTweetResponse, error)
	UpdateTweet(ctx context.Context, tweet entity.UpdateTweetRequest) (entity.UpdateTweetResponse, error)
	DeleteTweet(ctx context.Context, id string) error
	GetTweet(ctx context.Context, id string) (entity.GetTweetResponse, error)
	ListTweets(ctx context.Context, filter entity.Filter) (entity.ListTweetsResponse, error)
	UserTweets(ctx context.Context, userID string) (entity.ListTweetsResponse, error)
}

type SearchStorageI interface {
	Search(ctx context.Context, data string) (entity.SearchResponse, error)
}

type LikeStorageI interface {
	Like(ctx context.Context, like entity.LikeAction) (bool, error)
}

type FollowStorageI interface {
	Follow(ctx context.Context, follow entity.FollowAction) (bool, error)
	GetFollowings(ctx context.Context, id string) (entity.ListUser, error)
	GetFollowers(ctx context.Context, id string) (entity.ListUser, error)
}
