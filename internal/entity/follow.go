package entity

type FollowAction struct {
	UserID      string `json:"-"`
	FollowingID string `json:"following_id"`
}
