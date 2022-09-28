package tools

import "strings"

type Record [2]string

type Form struct {
	records []Record
}

func NewForm() *Form {
	return &Form{
		records: make([]Record, 0),
	}
}

func (f *Form) Add(key string, value string) {
	f.records = append(f.records, [2]string{key, value})
}

func (f *Form) Encode() string {
	var buf strings.Builder
	for _, v := range f.records {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(v[0])
		buf.WriteByte('=')
		buf.WriteString(v[1])
	}

	return buf.String()
}
