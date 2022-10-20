package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adepte-myao/test_parser/internal/html"
	"github.com/adepte-myao/test_parser/internal/storage"
	"github.com/adepte-myao/test_parser/internal/tools"
	"github.com/sirupsen/logrus"
)

type SolutionHandler struct {
	logger            *logrus.Logger
	sitemapRepository *storage.SitemapRepository
	taskRepository    *storage.TaskRepository
	htmlParser        *html.Parser
	baseLink          string
}

func NewSolutionHandler(logger *logrus.Logger, baseLink string, store *storage.Store) *SolutionHandler {
	sitemapRepository := storage.NewSitemapRepository(store)
	taskRepo := storage.NewTaskRepository(store)

	return &SolutionHandler{
		logger:            logger,
		sitemapRepository: sitemapRepository,
		taskRepository:    taskRepo,
		htmlParser:        html.NewParser(),
		baseLink:          baseLink,
	}
}

func (handler *SolutionHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Solution request received")

	testLinks, err := handler.sitemapRepository.GetAllLinks()
	if err != nil {
		handler.logger.Error("Error when getting links from database")
		return
	}

	handler.taskRepository.DeleteAll()

	for _, testLink := range testLinks {

		resp, err := handler.respFromResultPage(string(testLink))
		if err != nil {
			handler.logger.Debug("Error when getting response: ", err.Error())
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			handler.logger.Warn("Status code from site is ", resp.StatusCode)

			rw.Write([]byte(fmt.Sprint("Response from ", testLink, "is not OK\n")))
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			handler.logger.Error("Reading from response failed")

			rw.Write([]byte(fmt.Sprint("Error when processing ", testLink, "\n")))
			continue
		}

		dataStr := string(data)
		tasks, err := handler.htmlParser.ParseSolution(dataStr)
		if err != nil {
			handler.logger.Error("Parsing failed")

			rw.Write([]byte(fmt.Sprint("Error when parsing", testLink, "\n")))
			continue
		}

		for _, task := range tasks {
			err = handler.taskRepository.CreateTask(task)
			if err != nil {
				handler.logger.Error(err.Error())
			}
		}
	}

	rw.WriteHeader(http.StatusOK)

	handler.logger.Info("Solution request: processing finished")
}

func (handler *SolutionHandler) respFromResultPage(link string) (*http.Response, error) {
	resp, err := tools.DoProperRequest(http.MethodGet, link)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stringifyBody := string(body)

	formData := handler.htmlParser.ParseTasks(stringifyBody)
	params := tools.ExtractQueryParams(link)

	form := tools.NewForm()
	for _, formValues := range formData {
		form.Add(formValues.QuestionName, formValues.RatioName)
		form.Add(formValues.RatioName, formValues.RatioValue)
	}
	form.Add("Width", "")
	form.Add("iter", fmt.Sprint(params.Iter+1))
	form.Add("bil", fmt.Sprint(params.Bil))
	form.Add("test", fmt.Sprint(params.Test))

	req, err := http.NewRequest(http.MethodPost, handler.baseLink, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "tester=%D0%98%D0%BD%D0%BA%D0%BE%D0%B3%D0%BD%D0%B8%D1%82%D0%BE")
	req.Header.Add("referer", link)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	handler.logger.Info("Making request: ", link)
	return http.DefaultClient.Do(req)
}
