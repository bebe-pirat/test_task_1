package entity

type Subscription struct {
	Id          int    `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
}
