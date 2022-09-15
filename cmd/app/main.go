package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/adepte-myao/test_parser/internal/models"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handleSimpleHmtl)
	router.HandleFunc("/file", handleFile)

	http.ListenAndServe(":9095", router)
}

func handleSimpleHmtl(rw http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://tests24.ru/?iter=3&test=726")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte("Response from given source is not OK"))
		return
	}

	_, err = io.Copy(rw, resp.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func handleFile(rw http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("src.html")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	dataStr := string(data)

	tasks := parseHtml(dataStr)
	rw.WriteHeader(http.StatusOK)
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

func parseHtml(html string) []models.Task {
	taskReg := regexp.MustCompile(`<div class="card flex-shrink-1 shadow">[[:print:][:cntrl:]А-Яа-я№«»]*?</div>`)
	foundedQuestionsWithAnswers := taskReg.FindAllString(html, -1)

	questionReg := regexp.MustCompile(`[0-9]+\)</b>[[:print:][:cntrl:]А-Яа-я№«»]*?<`)
	excessQuestionSymbolsReg := regexp.MustCompile(`[0-9]+\)</b>`)

	optionsReg := regexp.MustCompile(`<span[[:print:][:cntrl:]А-Яа-я№«»]*?</span>`)
	excessOptionsSymbolsReg := regexp.MustCompile(`<b>|</b>|<span[[:print:][:cntrl:]]*?>|</span>|[0-9]+\) `)

	answerReg := regexp.MustCompile(`color:#3ea82e[[:print:][:cntrl:]А-Яа-я№«»]*?</span>`)
	excessAnswerSymbolsReg := regexp.MustCompile(`[0-9]+\)|<b>|</b>`)

	extraSpacesReg := regexp.MustCompile(`\s+`)

	tasks := make([]models.Task, 0)

	for _, element := range foundedQuestionsWithAnswers {
		question := questionReg.FindAllString(element, 1)[0]
		question = excessQuestionSymbolsReg.ReplaceAllString(question, "")
		question = extraSpacesReg.ReplaceAllString(question, " ")
		question = question[:len(question)-1]

		options := optionsReg.FindAllString(element, -1)[1:]
		for i := 0; i < len(options); i++ {
			options[i] = excessOptionsSymbolsReg.ReplaceAllString(options[i], "")
			options[i] = extraSpacesReg.ReplaceAllString(options[i], " ")
		}

		answers := answerReg.FindAllString(element, -1)
		var answer string
		if len(answers) > 1 {
			answer = answers[1]
		} else {
			answer = answers[0]
		}
		answer = excessAnswerSymbolsReg.ReplaceAllString(answer, "")
		answer = extraSpacesReg.ReplaceAllString(answer, " ")
		firstClosingBraceInd := strings.Index(answer, ">")
		lastOpeningBraceInd := strings.LastIndex(answer, "<")
		answer = answer[firstClosingBraceInd+1 : lastOpeningBraceInd]

		task := models.NewTask()
		task.Question = question
		task.Options = options
		task.Answer = answer
		tasks = append(tasks, *task)
	}

	return tasks
}
