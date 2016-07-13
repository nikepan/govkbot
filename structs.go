package govkbot

// Message - VK message struct
type Message struct {
	ID        int
	Date      int
	Out       int
	UserID    int `json:"user_id"`
	ChatID    int `json:"chat_id"`
	ReadState int `json:"read_state"`
	Title     string
	Body      string
	Action    string
	ActionMID int `json:"action_mid"`
}

// Messages - VK Messages
type Messages struct {
	Count int
	Items []Message
}

// MessagesResponse - VK messages response
type MessagesResponse struct {
	Response Messages
}

// User - simple VK user struct
type User struct {
	ID        int
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo     string
	InvitedBy int `json:"invited_by"`
}

// UsersResponse - VK user response
type UsersResponse struct {
	Response []*User
}

// FriendRequests - VK friend requests
type FriendRequests struct {
	Count int
	Items []int
}

// FriendRequestsResponse - VK friend requests response
type FriendRequestsResponse struct {
	Response FriendRequests
}

// SimpleResponse - simple int response
type SimpleResponse struct {
	Response int
}
