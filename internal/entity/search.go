package entity

type SearchResponse struct {
	Users  []GetUserResponse  `json:"users"`
	Tweets []GetTweetResponse `json:"tweets"`
}
