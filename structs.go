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
	Items []*Message
}

// MessagesResponse - VK messages response
type MessagesResponse struct {
	Response Messages
	Error    *VKError
}

// Geo - City and Country info
type Geo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// User - simple VK user struct
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo     string `json:"photo"`
	InvitedBy int    `json:"invited_by"`
	City      Geo    `json:"city"`
	Country   Geo    `json:"country"`
}

// FullName - returns full name of user
func (u *User) FullName() string {
	if u != nil {
		return u.FirstName + " " + u.LastName
	}
	return ""
}

// VKUsers - Users list. Can be sort by full name
type VKUsers []*User

// UsersResponse - VK user response
type UsersResponse struct {
	Response VKUsers
	Error    *VKError
}

// FriendRequests - VK friend requests
type FriendRequests struct {
	Count int
	Items []int
}

// FriendRequestsResponse - VK friend requests response
type FriendRequestsResponse struct {
	Response FriendRequests
	Error    *VKError
}

// SimpleResponse - simple int response
type SimpleResponse struct {
	Response int
	Error    *VKError
}

// ChatInfo - chat info
type ChatInfo struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Kicked  int    `json:"kicked"`
	AdminID int    `json:"admin_id"`
	Users   VKUsers
}

// VKError - error info
type VKError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	//	RequestParams
}

// ChatInfoResponse - chat info vk struct
type ChatInfoResponse struct {
	Response ChatInfo
	Error    *VKError
}

func (err *VKError) Error() string {
	return err.ErrorMsg
}

func (a VKUsers) Len() int           { return len(a) }
func (a VKUsers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VKUsers) Less(i, j int) bool { return a[i].FullName() < a[j].FullName() }
