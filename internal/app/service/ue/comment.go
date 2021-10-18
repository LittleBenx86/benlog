package ue

import "github.com/LittleBenx86/Benlog/internal/global/dependencies"

type CommentService struct {
	*dependencies.Dependencies
}

func NewCommentService(d *dependencies.Dependencies) *CommentService {
	return &CommentService{
		Dependencies: d,
	}
}

func (c *CommentService) GetCommentCount() {

}

func (c *CommentService) GetCommentsByBlogId() {

}

func (c *CommentService) GetCommentsWithPageLimitation() {

}

func (c *CommentService) GetCommentsByBlogIdWithPageLimitation() {

}
