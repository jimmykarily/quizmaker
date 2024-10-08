package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/models"
	"gorm.io/gorm/clause"
)

type (
	QuestionController struct{}
)

func (c *QuestionController) Answer(gctx *gin.Context) {
	session, err := currentSession(gctx)
	if handleError(gctx.Writer, err, http.StatusUnauthorized) {
		return
	}

	err = gctx.Request.ParseForm()
	if handleError(gctx.Writer, err, http.StatusBadRequest) {
		return
	}

	qid := gctx.Param("id")
	selectedAnswer := gctx.Request.FormValue("answer")

	var question models.Question
	err = Settings.DB.First(&question, "ID = ?", qid).Error
	if handleError(gctx.Writer, err, http.StatusNotFound) {
		return
	}

	// If the question doesn't belong to the current session
	if question.SessionEmail != session.Email {
		handleError(gctx.Writer, errors.New("question doesn't belong to session"), http.StatusUnauthorized)
	}

	// Don't allow answering expired or already answered questions
	if question.Expired() || question.UserAnswer != 0 {
		// TODO: Flash error
	} else {
		question.UserAnswer, err = strconv.Atoi(selectedAnswer)
		if handleError(gctx.Writer, err, http.StatusBadRequest) {
			return
		}

		err = Settings.DB.Save(&question).Error
		if handleError(gctx.Writer, err, http.StatusInternalServerError) {
			return
		}

		// reload session
		err = Settings.DB.Preload(clause.Associations).Find(&session).Error
		if handleError(gctx.Writer, err, http.StatusInternalServerError) {
			return
		}
		s := &session
		s.UpdateCacheColumns()
		err = Settings.DB.Save(s).Error
		if handleError(gctx.Writer, err, http.StatusInternalServerError) {
			return
		}
	}

	redirectURL, err := GetFullURL(gctx.Request, "QuizShow", nil)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	gctx.Redirect(http.StatusFound, redirectURL)
}
