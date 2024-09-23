package entity

type FollowAction struct {
	UserID      int `json:"-"`
	FollowingID int `json:"following_id"`
}
