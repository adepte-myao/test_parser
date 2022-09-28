package handlers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/adepte-myao/test_parser/internal/html"
	"github.com/adepte-myao/test_parser/internal/models"
	"github.com/adepte-myao/test_parser/internal/tools"
	"github.com/sirupsen/logrus"
)

type SolutionHandler struct {
	logger     *logrus.Logger
	htmlParser *html.Parser
	baseLink   string
}

func NewSolutionHandler(logger *logrus.Logger, baseLink string) *SolutionHandler {
	return &SolutionHandler{
		logger:     logger,
		htmlParser: html.NewParser(),
		baseLink:   baseLink,
	}
}

func (handler *SolutionHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Solution request received")

	f, err := os.Open("allTestReferences.txt")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("can't open the file"))
	}
	defer f.Close()

	rw.WriteHeader(http.StatusOK)
	reader := bufio.NewReader(f)
	i := -1
	for {
		i++
		if i == 4 {
			break
		}

		testLink, err := reader.ReadString('\n')
		if err != nil { // EOF?
			handler.logger.Debug("Error reading line")
			return
		}

		resp, err := handler.respFromResultPage(testLink[:len(testLink)-1])
		if err != nil {
			handler.logger.Debug("Error when getting response: ", err.Error())
			continue
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
		tasks, err := handler.htmlParser.ParseSolution(dataStr)
		if err != nil {
			handler.logger.Error("Parsing failed")

			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
		handler.writeTasksToResponse(rw, tasks)
	}

	handler.logger.Info("Solution request: processing finished")
}

func (handler *SolutionHandler) respFromResultPage(link string) (*http.Response, error) {
	resp, err := doRequest(http.MethodGet, link)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	stringifyBody := string(body)

	formData := handler.htmlParser.ParseTasks(stringifyBody)
	params := parseLinkValues(link)

	form := tools.NewForm()
	for _, formValues := range formData {
		form.Add(formValues.QuestionName, formValues.RatioName)
		form.Add(formValues.RatioName, formValues.RatioValue)
	}
	form.Add("Width", "")
	form.Add("iter", params[0])
	form.Add("bil", params[1])
	form.Add("test", params[2])

	req, err := http.NewRequest(http.MethodPost, handler.baseLink, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "")
	req.Header.Add("Cookie", "tester=%D0%98%D0%BD%D0%BA%D0%BE%D0%B3%D0%BD%D0%B8%D1%82%D0%BE")
	req.Header.Add("referer", link)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	handler.logger.Info("Making request: ", link)
	return http.DefaultClient.Do(req)
}

func parseLinkValues(link string) [3]string {
	iterReg := regexp.MustCompile(`iter=[0-9]+`)
	bilReg := regexp.MustCompile(`bil=[0-9]+`)
	testReg := regexp.MustCompile(`test=[0-9]+`)

	iter := iterReg.FindAllString(link, -1)[0]
	bil := bilReg.FindAllString(link, -1)[0]
	test := testReg.FindAllString(link, -1)[0]

	iterNumb, err := strconv.Atoi(iter[5:])
	if err != nil {
		panic(err)
	}
	bilNumb, err := strconv.Atoi(bil[4:])
	if err != nil {
		panic(err)
	}
	testNumb, err := strconv.Atoi(test[5:])
	if err != nil {
		panic(err)
	}

	return [3]string{fmt.Sprint(iterNumb + 1), fmt.Sprint(bilNumb + 1), fmt.Sprint(testNumb)}
}

func (handler *SolutionHandler) writeTasksToResponse(rw http.ResponseWriter, tasks []models.Task) {
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
