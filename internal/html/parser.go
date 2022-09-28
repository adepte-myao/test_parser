package html

import (
	"regexp"
	"strings"

	"github.com/adepte-myao/test_parser/internal/models"
)

type Parser struct {
	taskBlockReg *regexp.Regexp

	questionBlockReg         *regexp.Regexp
	excessQuestionSymbolsReg *regexp.Regexp

	optionsBlockReg         *regexp.Regexp
	excessOptionsSymbolsReg *regexp.Regexp

	answerReg              *regexp.Regexp
	excessAnswerSymbolsReg *regexp.Regexp

	nameTagReg  *regexp.Regexp
	valueTagReg *regexp.Regexp

	extraSpacesReg *regexp.Regexp
	oneSpaceReg    *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		taskBlockReg: regexp.MustCompile(`<div class="card flex-shrink-1 shadow">[[:print:][:cntrl:]А-Яа-я№«»]*?</div>`),

		questionBlockReg:         regexp.MustCompile(`[0-9]+\)</b>[[:print:][:cntrl:]А-Яа-я№«»]*?<`),
		excessQuestionSymbolsReg: regexp.MustCompile(`[0-9]+\)</b>`),

		optionsBlockReg:         regexp.MustCompile(`<span[[:print:][:cntrl:]А-Яа-я№«»]*?</span>`),
		excessOptionsSymbolsReg: regexp.MustCompile(`<b>|</b>|<span[[:print:][:cntrl:]]*?>|</span>|[0-9]+\) `),

		answerReg:              regexp.MustCompile(`color:#[2-5][c-f][a-d][6-9][1-4][c-f][[:print:][:cntrl:]А-Яа-я№«»]*?</span>`),
		excessAnswerSymbolsReg: regexp.MustCompile(`[0-9]+\)|<b>|</b>`),

		nameTagReg:  regexp.MustCompile(`name=".*?"`),
		valueTagReg: regexp.MustCompile(`value=".*?"`),

		extraSpacesReg: regexp.MustCompile(`\s+`),
		oneSpaceReg:    regexp.MustCompile(`\s*`),
	}
}

type FormDataValues struct {
	QuestionName string
	RatioName    string
	RatioValue   string
}

func (parser *Parser) ParseTasks(html string) []FormDataValues {
	foundedTasks := parser.taskBlockReg.FindAllString(html, -1)

	formDataValues := make([]FormDataValues, 0)
	for _, task := range foundedTasks {
		task = parser.oneSpaceReg.ReplaceAllString(task, "")

		names := parser.nameTagReg.FindAllString(task, -1)
		values := parser.valueTagReg.FindAllString(task, -1)

		qName := names[0][6:]
		qName = qName[:len(qName)-1]

		ratioName := names[1][6:]
		ratioName = ratioName[:len(ratioName)-1]

		ratioValue := values[1][7:]
		ratioValue = ratioValue[:len(ratioValue)-1]

		formDataValue := FormDataValues{
			QuestionName: qName,
			RatioName:    ratioName,
			RatioValue:   ratioValue,
		}

		formDataValues = append(formDataValues, formDataValue)
	}

	return formDataValues
}

func (parser *Parser) ParseSolution(html string) ([]models.Task, error) {
	foundedTasks := parser.taskBlockReg.FindAllString(html, -1)

	tasks := make([]models.Task, 0)

	for _, taskBlock := range foundedTasks {
		question := parser.parseQuestion(taskBlock)
		options := parser.parseOptions(taskBlock)
		answer := parser.parseAnswer(taskBlock)

		task := models.NewTask()
		task.Question = question
		task.Options = options
		task.Answer = answer
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (parser *Parser) parseQuestion(taskBlock string) string {
	question := parser.questionBlockReg.FindAllString(taskBlock, 1)[0]
	question = parser.excessQuestionSymbolsReg.ReplaceAllString(question, "")
	question = parser.extraSpacesReg.ReplaceAllString(question, " ")

	question = strings.TrimLeft(question, " ")
	return question[:len(question)-1]
}

func (parser *Parser) parseOptions(taskBlock string) []string {
	options := parser.optionsBlockReg.FindAllString(taskBlock, -1)[1:]
	for i := 0; i < len(options); i++ {
		options[i] = parser.extraSpacesReg.ReplaceAllString(options[i], " ")
		options[i] = parser.excessOptionsSymbolsReg.ReplaceAllString(options[i], "")
		options[i] = parser.extraSpacesReg.ReplaceAllString(options[i], " ")
	}
	return options
}

func (parser *Parser) parseAnswer(taskBlock string) string {
	answers := parser.answerReg.FindAllString(taskBlock, -1)

	var answer string
	if len(answers) > 1 {
		answer = answers[1]
	} else {
		answer = answers[0]
	}

	answer = parser.extraSpacesReg.ReplaceAllString(answer, " ")
	answer = parser.excessAnswerSymbolsReg.ReplaceAllString(answer, "")
	answer = parser.extraSpacesReg.ReplaceAllString(answer, " ")

	firstClosingBraceInd := strings.Index(answer, ">")
	lastOpeningBraceInd := strings.LastIndex(answer, "<")
	answer = answer[firstClosingBraceInd+1 : lastOpeningBraceInd]

	return answer
}
