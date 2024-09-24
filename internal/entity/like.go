package entity

type LikeAction struct {
	UserID  string `json:"-"`
	TweetID string `json:"tweet_id"`
}
