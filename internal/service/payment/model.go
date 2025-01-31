package payment

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type CreateOrderResponse struct {
	ID     string `json:"id"`
	Links  []Link `json:"links"`
	Status string `json:"status"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}
type customData struct {
	ReservationID string `json:"reservation_id"`
	RoomID        string `json:"room_id"`
	Price         string `json:"price"`
}
