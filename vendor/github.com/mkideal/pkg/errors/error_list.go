package errors

import (
	"bytes"
	"sort"
)

type ErrorList struct {
	errs []error
}

func (list *ErrorList) Add(err error) {
	if err != nil {
		list.errs = append(list.errs, err)
	}
}

func (list *ErrorList) Len() int { return len(list.errs) }

func (list *ErrorList) Err() error {
	if len(list.errs) > 0 {
		return list
	}
	return nil
}

func (list *ErrorList) Error() string {
	var buf bytes.Buffer
	prefix := ""
	sep := "\n"
	for i, err := range list.errs {
		if i > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(prefix)
		buf.WriteString(err.Error())
	}
	return buf.String()
}

func (list *ErrorList) swap(i, j int) { list.errs[i], list.errs[j] = list.errs[j], list.errs[i] }

func (list *ErrorList) Sort(less func(error, error) bool) {
	sort.Sort(sortErrors{list, less})
}

type sortErrors struct {
	errors *ErrorList
	less   func(error, error) bool
}

func (se sortErrors) Len() int           { return se.errors.Len() }
func (se sortErrors) Less(i, j int) bool { return se.less(se.errors.errs[i], se.errors.errs[j]) }
func (se sortErrors) Swap(i, j int)      { se.errors.swap(i, j) }
