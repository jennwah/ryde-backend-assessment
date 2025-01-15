package repository

type userFriendPair struct {
	UserID   string `db:"user_id"`
	FriendID string `db:"friend_id"`
}
