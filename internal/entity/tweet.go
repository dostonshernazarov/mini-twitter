package entity

import "time"

type TweetRequest struct {
	ParentTweetID *string  `json:"parent_tweet_id"`
	Content       *string  `json:"content"`
	URLs          []string `json:"files"`
}

type CreateTweetRequest struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	ParentTweetID *string   `json:"parent_tweet_id"`
	Content       *string   `json:"content"`
	URLs          []string  `json:"files"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateTweetResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	ParentTweetID *string   `json:"parent_tweet_id"`
	Content       *string   `json:"content"`
	URLs          []string  `json:"urls"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdateTweetRequest struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type UpdateTweetResponse struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	ParentTweetID *string  `json:"parent_tweet_id"`
	Content       *string  `json:"content"`
	URLs          []string `json:"urls"`
}

type GetTweetResponse struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	ParentTweetID *string  `json:"parent_tweet_id"`
	Content       *string  `json:"content"`
	URLs          []string `json:"urls"`
}

type ListTweetsResponse struct {
	Tweets []GetTweetResponse `json:"tweets"`
	Count  int                `json:"count"`
}
