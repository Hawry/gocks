package settings

import (
	"regexp"
	"strings"
)

// Mock describes what path that should return what status & body
type Mock struct {
	Path         string            `json:"path"`
	Method       string            `json:"method"`
	ResponseCode int               `json:"response_code"`
	ResponseBody string            `json:"response_body"`
	Headers      map[string]string `json:"headers"`
	DisableAuth  bool              `json:"disable_auth"`
}

//Wildcard returns all possible wildcards in a path
func (m *Mock) Wildcard() (bool, []string) {
	r := regexp.MustCompile("{[[:alnum:]]*}")
	matches := r.FindAllString(m.Path, -1)
	mv := make([]string, len(matches))
	for i, v := range matches {
		v = strings.TrimPrefix(v, "{")
		v = strings.TrimSuffix(v, "}")
		mv[i] = v
	}
	return (len(mv) > 0), mv
}
