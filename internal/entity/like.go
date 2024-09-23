package entity

type LikeAction struct {
	UserID  int `json:"-"`
	TweetID int `json:"tweet_id"`
}
