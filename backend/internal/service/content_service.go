package service

import (
	"context"

	"cesizen/internal/domain"
)

type ContentRepo interface {
	FindAllPublished(ctx context.Context) ([]domain.Content, error)
	FindAll(ctx context.Context) ([]domain.Content, error)
	FindByID(ctx context.Context, id int) (*domain.Content, error)
	Create(ctx context.Context, c *domain.Content) error
	Update(ctx context.Context, c *domain.Content) error
	Delete(ctx context.Context, id int) error
}

type ContentService struct {
	contents ContentRepo
}

func NewContentService(contents ContentRepo) *ContentService {
	return &ContentService{contents: contents}
}

func (s *ContentService) ListPublished(ctx context.Context) ([]domain.Content, error) {
	return s.contents.FindAllPublished(ctx)
}

func (s *ContentService) ListAll(ctx context.Context) ([]domain.Content, error) {
	return s.contents.FindAll(ctx)
}

func (s *ContentService) GetByID(ctx context.Context, id int) (*domain.Content, error) {
	return s.contents.FindByID(ctx, id)
}

func (s *ContentService) Create(ctx context.Context, input domain.CreateContentInput) (*domain.Content, error) {
	c := &domain.Content{
		Title:   input.Title,
		Content: input.Content,
		Author:  input.Author,
	}
	if err := s.contents.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ContentService) Update(ctx context.Context, id int, input domain.UpdateContentInput) (*domain.Content, error) {
	c, err := s.contents.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Title = input.Title
	c.Content = input.Content
	c.Author = input.Author
	c.IsPublished = input.IsPublished
	if err := s.contents.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ContentService) Delete(ctx context.Context, id int) error {
	return s.contents.Delete(ctx, id)
}
