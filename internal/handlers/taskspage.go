package handlers

import (
	"io"
	"net/http"

	"github.com/adepte-myao/test_parser/internal/html"
	"github.com/adepte-myao/test_parser/internal/models"
	"github.com/adepte-myao/test_parser/internal/tools"
	"github.com/sirupsen/logrus"
)

type TaskPageHandler struct {
	logger     *logrus.Logger
	htmlParser *html.Parser
}

func NewTaskPageHandler(logger *logrus.Logger) *TaskPageHandler {
	return &TaskPageHandler{
		logger:     logger,
		htmlParser: html.NewParser(),
	}
}

func (handler *TaskPageHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Received handle simple HTML request")

	resp, err := tools.DoRequest("GET", "https://tests24.ru/?iter=3&test=726")
	if err != nil {
		handler.logger.Error("Request problem: ", err)

		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handler.logger.Error("Status code from site is not OK")

		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source is not OK"))
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		handler.logger.Error("Reading from response failed")

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataStr := string(data)
	tasks, err := handler.htmlParser.ParseHtml(dataStr)
	if err != nil {
		handler.logger.Error("Parsing failed")

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	handler.writeTasksToResponse(rw, tasks)
}

func (handler *TaskPageHandler) writeTasksToResponse(rw http.ResponseWriter, tasks []models.Task) {
	for _, task := range tasks {
		rw.Write([]byte(task.Question))

		rw.Write([]byte("\n"))
		for _, option := range task.Options {
			rw.Write([]byte("\t"))
			rw.Write([]byte(option))
			rw.Write([]byte("\n"))
		}

		rw.Write([]byte("\tRight answer: "))
		rw.Write([]byte(task.Answer))

		rw.Write([]byte("\n\n"))
	}
}
