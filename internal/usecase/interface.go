package usecase

import (
	"context"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
)

type User interface {
	UniqueUsername(ctx context.Context, username string) (bool, error)
	UniqueEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user entity.CreateUserRequest) (entity.CreateUserResponse, error)
	Update(ctx context.Context, user entity.UpdateUserRequest) (entity.UpdateUserResponse, error)
	UpdatePasswd(ctx context.Context, id int, passwd string) error
	UploadImage(ctx context.Context, id int, url string) error
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, field map[string]interface{}) (entity.GetUserResponse, error)
	List(ctx context.Context, filter entity.Filter) (entity.ListUser, error)
}

type Twit interface {
	CreateTweet(ctx context.Context, tweet entity.CreateTwitRequest) (entity.CreateTweetResponse, error)
	UpdateTweet(ctx context.Context, tweet entity.UpdateTweetRequest) (entity.UpdateTweetResponse, error)
	DeleteTweet(ctx context.Context, id int) error
	GetTweet(ctx context.Context, id int) (entity.GetTweetResponse, error)
	ListTweets(ctx context.Context, filter entity.Filter) (entity.ListTweetsResponse, error)
	UserTweets(ctx context.Context, usrID int) (entity.ListTweetsResponse, error)
}

type Search interface {
	Search(ctx context.Context, data string) (entity.SearchResponse, error)
}

type Like interface {
	Like(ctx context.Context, like entity.LikeAction) (bool, error)
}

type Follow interface {
	Follow(ctx context.Context, follow entity.FollowAction) (bool, error)
	GetFollowings(ctx context.Context, id int) (entity.ListUser, error)
	GetFollowers(ctx context.Context, id int) (entity.ListUser, error)
}
