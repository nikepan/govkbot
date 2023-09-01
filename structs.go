package govkbot

import (
	"strings"
)

// Mention - user mention in message
type Mention struct {
	ID   int
	Name string
}

// Button for keyboard, which sends to user
type Button struct {
	Action struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
		Label   string `json:"label"`
	} `json:"action"`
	Color string `json:"color"`
}

// Keyboard to send for user
type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
}

//Reply for message
type Reply struct {
	Msg      string
	Keyboard *Keyboard
}

// Message - VK message struct
type Message struct {
	ID          int
	Date        int
	Out         int
	UserID      int   `json:"user_id"`
	ChatID      int   `json:"chat_id"`
	PeerID      int64 `json:"peer_id"`
	ReadState   int   `json:"read_state"`
	Title       string
	Body        string
	Action      string
	ActionMID   int `json:"action_mid"`
	Flags       int
	Timestamp   int64
	Payload     string
	FwdMessages []Message `json:"fwd_messages"`
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
	ID              int    `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ScreenName      string `json:"screen_name"`
	Photo           string `json:"photo"`
	InvitedBy       int    `json:"invited_by"`
	City            Geo    `json:"city"`
	Country         Geo    `json:"country"`
	Sex             int    `json:"sex"`
	BDate           string `json:"bdate"`
	Photo50         string `json:"photo_50"`
	Photo100        string `json:"photo_100"`
	Status          string `json:"status"`
	About           string `json:"about"`
	Relation        int    `json:"relation"`
	Hidden          int    `json:"hidden"`
	Closed          bool    `json:"is_closed"`
	CanAccessClosed bool    `json:"can_access_closed"`
	Deactivated     string `json:"deactivated"`
	IsAdmin         bool   `json:"is_admin"`
	IsOwner         bool   `json:"is_owner"`
}

// FullName - returns full name of user
func (u *User) FullName() string {
	if u != nil {
		return strings.Trim(u.FirstName+" "+u.LastName, " ")
	}
	return ""
}

// VKUsers - Users list. Can be sort by full name
type VKUsers []*User

// MemberItem - conversation item
type MemberItem struct {
	MemberID  int  `json:"member_id"`
	JoinDate  int  `json:"join_date"`
	IsOwner   bool `json:"is_owner"`
	IsAdmin   bool `json:"is_admin"`
	InvitedBy int  `json:"invited_by"`
}

// UserProfile - conversation user profile
type UserProfile struct {
	ID              int
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	IsClosed        bool   `json:"is_closed"`
	CanAccessClosed bool   `json:"can_access_closed"`
	Sex             int
	ScreenName      string `json:"screen_name"`
	BDate           string `json:"bdate"`
	Photo           string
	Online          int
	City            Geo
	Country         Geo
}

// GroupProfile - conversation group profile
type GroupProfile struct {
	ID       int
	Name     string
	IsClosed int `json:"is_closed"`
	Type     string
	Photo50  string
	Photo100 string
	Photo200 string
}

// VKMembers - conversation members info
type VKMembers struct {
	Items    []MemberItem
	Profiles []UserProfile
	Groups   []GroupProfile
}

// UsersResponse - VK user response
type UsersResponse struct {
	Response VKUsers
	Error    *VKError
}

// MembersResponse - VK user response
type MembersResponse struct {
	Response VKMembers
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

// FriendDeleteResponse - VK friend delete response
type FriendDeleteResponse struct {
	Response map[string]int
	Error    *VKError
}

// SimpleResponse - simple int response
type SimpleResponse struct {
	Response int
	Error    *VKError
}

// ErrorResponse - need to parse VK error
type ErrorResponse struct {
	Error *VKError
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

// VKError - error with response content
type ResponseError struct {
	err     error
	content string
}

func (err ResponseError) Error() string {
	return err.err.Error()
}

//Content - error content
func (err ResponseError) Content() string {
	return err.content
}

//ConversationInfo - conversation info
type ConversationInfo struct {
	Peer struct {
		ID      int
		Type    string
		LocalID int `json:"local_id"`
	}
	InRead        int `json:"in_read"`
	OutRead       int `json:"out_read"`
	LastMessageID int `json:"last_message_id"`
	CanWrite      struct {
		Allowed bool
	} `json:"can_write"`
	ChatSettings struct {
		Title        string
		MembersCount int `json:"members_count"`
		State        string
		ActiveIDs    []int `json:"active_ids"`
		ACL          struct {
			CanInvite           bool `json:"can_invite"`
			CanChangeInfo       bool `json:"can_change_info"`
			CanChangePin        bool `json:"can_change_pin"`
			CanPromoteUsers     bool `json:"can_promote_users"`
			CanSeeInviteLink    bool `json:"can_see_invite_link"`
			CanChangeInviteLink bool `json:"can_change_invite_link"`
		}
		IsGroupChannel bool `json:"is_group_channel"`
		OwnerID        int  `json:"owner_id"`
	} `json:"chat_settings"`
}

//ConversationsResponse - resonse of confersations info
type ConversationsResponse struct {
	Response struct {
		Items    []ConversationInfo
		Profiles []UserProfile
	}
	Error *VKError
}

// ChatInfoResponse - chat info vk struct
type ChatInfoResponse struct {
	Response ChatInfo
	Error    *VKError
}

func (err VKError) Error() string {
	return "vk: " + err.ErrorMsg
}

func (a VKUsers) Len() int           { return len(a) }
func (a VKUsers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VKUsers) Less(i, j int) bool { return a[i].FullName() < a[j].FullName() }
