package comment_model

type CommentCreator struct {
	UserId      int
	Name        string
	ImageUrl    string
	FriendCount int
	CreatedAt   string
}

type Comment struct {
	CommentId int
	PostId    int
	UserId    int
	Comment   string
	Creator   CommentCreator
	CreatedAt string
}
