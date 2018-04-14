package controllers

import (
	"net/http"
	"strconv"

	"github.com/gedorinku/koneko-online-judge/server/models"
	"github.com/labstack/echo"
)

type submissionRequest struct {
	LanguageID uint   `json:"languageID"`
	SourceCode string `json:"sourceCode"`
}

func Submit(c echo.Context) error {
	s := getSession(c)
	problem := getProblemFromContext(c)
	if problem == nil || s == nil || !problem.CanView(s) {
		return echo.ErrNotFound
	}

	request := &submissionRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
	}

	lang := models.GetLanguage(request.LanguageID)
	if lang == nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"使用できない言語です"})
	}

	submission := &models.Submission{
		UserID:     s.UserID,
		ProblemID:  problem.ID,
		LanguageID: lang.ID,
		SourceCode: request.SourceCode,
	}

	if err := models.Submit(submission); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
	}

	fetchSubmission(submission, s)

	return c.JSON(http.StatusCreated, submission)
}

func GetSubmissions(c echo.Context) error {
	s := getSession(c)
	problem := getProblemFromContext(c)
	if problem == nil || !problem.CanView(s) {
		return echo.ErrNotFound
	}

	problem.FetchSubmissions()
	for i := range problem.Submissions {
		fetchSubmission(&problem.Submissions[i], s)
	}
	return c.JSON(http.StatusOK, problem.Submissions)
}

func Rejudge(c echo.Context) error {
	_, err := getAdminSession(c)
	if err != nil {
		return err
	}

	submission := getSubmissionFromContext(c)
	if submission == nil {
		return echo.ErrNotFound
	}

	if err := submission.Rejudge(); err != nil {
		return ErrInternalServer
	}

	return c.NoContent(http.StatusNoContent)
}

func fetchSubmission(out *models.Submission, s *models.UserSession) {
	out.FetchUser()
	out.User.Email = ""
	fetchProblem(&out.Problem, s)
	out.Problem.ContestID = new(uint)
	out.FetchJudgeSetResultsDeeply(true)
	out.FetchLanguage()
}

func getSubmissionFromContext(c echo.Context) *models.Submission {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil
	}

	return models.GetSubmission(uint(id))
}
