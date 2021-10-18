package ue

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
)

type BlogService struct {
	*dependencies.Dependencies
}

func NewBlogService(d *dependencies.Dependencies) *BlogService {
	return &BlogService{
		Dependencies: d,
	}
}

func (b *BlogService) GetTopNBlogs() {

}

func (b *BlogService) GetBlogCount() {

}

func (b *BlogService) GetBlogsByLimitation() {

}

func (b *BlogService) GetBlogsByPageLimitation() {

}

func (b *BlogService) GetBlogDetail() {

}
