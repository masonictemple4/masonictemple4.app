package dtos

import "github.com/masonictemple4/masonictemple4.app/internal/customdate"

type CommentReturn struct {
	ID        uint                   `json:"id"`
	CreatedAt customdate.DefaultDate `json:"createdat"`
	UpdatedAt customdate.DefaultDate `json:"updatedat"`
	User      UserReturn             `json:"user"`
	Text      string                 `json:"text"`
}

type CommentInput struct {
	PostID uint   `json:"postid" validate:"required" example:"1"`
	Text   string `json:"text" validate:"required" example:"This is a comment."`
}

// TODO: Not sure we'll need an update comment just yet.
