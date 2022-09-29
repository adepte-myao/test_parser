package tools

import (
	"regexp"
	"strconv"
)

type queryParams struct {
	Iter int
	Bil  int
	Test int
}

func ExtractQueryParams(link string) queryParams {
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

	return queryParams{
		Iter: iterNumb,
		Bil:  bilNumb,
		Test: testNumb,
	}
}
