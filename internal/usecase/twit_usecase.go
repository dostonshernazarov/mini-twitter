package usecase

import (
	"context"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
)

type tweetService struct {
	ctxTimeout time.Duration
	repo       repo.TweetStorageI
}

func NewTweetService(timeout time.Duration, repository repo.TweetStorageI) Twit {
	return &tweetService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (t *tweetService) CreateTweet(ctx context.Context, tweet entity.CreateTwitRequest) (entity.CreateTweetResponse, error) {
	return t.repo.CreateTweet(ctx, tweet)
}

func (t *tweetService) UpdateTweet(ctx context.Context, tweet entity.UpdateTweetRequest) (entity.UpdateTweetResponse, error) {
	return t.repo.UpdateTweet(ctx, tweet)
}

func (t *tweetService) DeleteTweet(ctx context.Context, id int) error {
	return t.repo.DeleteTweet(ctx, id)
}

func (t *tweetService) GetTweet(ctx context.Context, id int) (entity.GetTweetResponse, error) {
	return t.repo.GetTweet(ctx, id)
}

func (t *tweetService) ListTweets(ctx context.Context, filter entity.Filter) (entity.ListTweetsResponse, error) {
	return t.repo.ListTweets(ctx, filter)
}

func (t *tweetService) UserTweets(ctx context.Context, usrID int) (entity.ListTweetsResponse, error) {
	return t.repo.UserTweets(ctx, usrID)
}
