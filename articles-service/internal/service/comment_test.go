package service

import (
	"articles-service/internal/lib/types"
	"articles-service/internal/models"
	"articles-service/internal/repository"
	"articles-service/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	testArticle := &models.Article{}
	testComment := &models.Comment{
		UserId:    3,
		ArticleId: 7,
		Text:      "ploequihrgpoequirhg",
	}
	createdCommentId := 2

	testCases := []struct {
		name               string
		articleRepoMockExp func(articleRepo *mocks.ArticleRepoInterface)
		commentRepoMockExp func(commentRepo *mocks.CommentRepoInterface)
		expId              int
		expErr             error
	}{
		{
			name: "Успешный кейс, никаких ошибок. Ищем статью, она находится, все ок.",
			articleRepoMockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetOne", testComment.ArticleId).
					Return(testArticle, nil).Times(1)
			},
			commentRepoMockExp: func(commentRepo *mocks.CommentRepoInterface) {
				commentRepo.On("Create", testComment).
					Return(createdCommentId, nil).Times(1)
			},
			expId:  createdCommentId,
			expErr: nil,
		},
		{
			name: "Статья не существует",
			articleRepoMockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetOne", testComment.ArticleId).
					Return(nil, pgx.ErrNoRows).Times(1)
			},
			commentRepoMockExp: func(commentRepo *mocks.CommentRepoInterface) {},
			expId:              -1,
			expErr:             types.ErrArticleNotFound,
		},
		{
			name: "Репозиторий статей нам вернул любую другую ошибку",
			articleRepoMockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetOne", testComment.ArticleId).
					Return(nil, errors.New("some error7")).Times(1)
			},
			commentRepoMockExp: func(commentRepo *mocks.CommentRepoInterface) {},
			expId:              -1,
			expErr:             errors.New("some error7"),
		},
		{
			name: "Любая ошибка из хранилища комментариев",
			articleRepoMockExp: func(articleRepo *mocks.ArticleRepoInterface) {
				articleRepo.On("GetOne", testComment.ArticleId).
					Return(testArticle, nil).Times(1)
			},
			commentRepoMockExp: func(commentRepo *mocks.CommentRepoInterface) {
				commentRepo.On("Create", testComment).
					Return(0, errors.New("some error")).Times(1)
			},
			expId:  -1,
			expErr: errors.New("some error"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockCommentRepoInterface := &mocks.CommentRepoInterface{}
			mockArticleRepoInterface := &mocks.ArticleRepoInterface{}
			repository := &repository.Repository{
				Comment: mockCommentRepoInterface,
				Article: mockArticleRepoInterface,
			}
			service := NewService(repository)

			testCase.articleRepoMockExp(mockArticleRepoInterface)
			testCase.commentRepoMockExp(mockCommentRepoInterface)
			id, err := service.Comment.CreateComment(testComment)

			assert.Equal(t, testCase.expErr, err)
			assert.Equal(t, testCase.expId, id)
			mockArticleRepoInterface.AssertExpectations(t)
			mockCommentRepoInterface.AssertExpectations(t)
		})
	}
}

func TestGetCommentsOnArticle(t *testing.T) {
	var testComments []*models.Comment = []*models.Comment{
		{
			UserId:    3,
			ArticleId: 7,
			Text:      "ploequihrgpoequirhg",
		},
		{
			UserId:    5,
			ArticleId: 7,
			Text:      "hello, Go!",
		},
	}
	testArticleId := 7

	testCases := []struct {
		name        string
		mockExp     func(commentRepo *mocks.CommentRepoInterface)
		expComments []*models.Comment
		expErr      error
	}{
		{
			name: "Успешный кейс, все по плану",
			mockExp: func(commentRepo *mocks.CommentRepoInterface) {
				commentRepo.On("GetCommentsOnArticle",
					testArticleId,
				).Return(testComments, nil).Times(1)
			},
			expComments: testComments,
			expErr:      nil,
		},
		{
			name: "Получаем любую ошибку из хранилища",
			mockExp: func(commentRepo *mocks.CommentRepoInterface) {
				commentRepo.On("GetCommentsOnArticle", testArticleId).
					Return(nil, errors.New("some error777")).Times(1)
			},
			expComments: nil,
			expErr:      errors.New("some error777"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockCommentRepoInterface := &mocks.CommentRepoInterface{}
			repository := &repository.Repository{Comment: mockCommentRepoInterface}
			service := NewService(repository)

			testCase.mockExp(mockCommentRepoInterface)
			comments, err := service.Comment.GetCommentsOnArticle(testArticleId)

			assert.Equal(t, testCase.expComments, comments)
			assert.Equal(t, testCase.expErr, err)
			mockCommentRepoInterface.AssertExpectations(t)
		})
	}
}
