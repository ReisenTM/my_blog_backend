package message_type_enum

type Type int8

const (
	CommentType          Type = 1
	ReplyType            Type = 2
	FavorArticleType     Type = 3
	UnFavorArticleType   Type = 4
	FavorCommentType     Type = 5
	UnFavorCommentType   Type = 6
	CollectArticleType   Type = 7
	UnCollectArticleType Type = 8
	SystemType           Type = 9
)
