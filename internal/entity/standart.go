package entity

type Error struct {
	Message string `json:"message"`
}

type ResponseWithMessage struct {
	Message string `json:"message"`
}

type ResponseWithStatus struct {
	Status bool `json:"status"`
}
