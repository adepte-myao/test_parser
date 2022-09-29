package handlers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/adepte-myao/test_parser/internal/html"
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

	outFile, err := os.Create("TestsFile.txt")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("can't create the file"))
	}
	defer outFile.Close()

	rw.WriteHeader(http.StatusOK)
	reader := bufio.NewReader(f)
	writer := bufio.NewWriter(outFile)
	i := -1
	for {
		i++
		// if i == 2 {
		// 	break
		// }

		testLink, errFileRead := reader.ReadString('\n')
		if errFileRead != nil && errFileRead != io.EOF {
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

		tools.WriteTasks(writer, tasks)

		if errFileRead == io.EOF {
			break
		}
	}

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
