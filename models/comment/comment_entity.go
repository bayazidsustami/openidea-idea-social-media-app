package comment_model

type CommentCreator struct {
	UserId      string
	Name        string
	ImageUrl    string
	FriendCount int
	CreatedAt   string
}

type Comment struct {
	CommentId int
	PostId    int
	UserId    string
	Comment   string
	Creator   CommentCreator
	CreatedAt string
}
