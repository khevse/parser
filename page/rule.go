package page

import (
	"fmt"
	"regexp"
	"strings"
)

type Rule struct {
	Tag   string
	Attrs map[string][]*regexp.Regexp
}

var reSymply = regexp.MustCompile(`^(\d|\w)+$`)

func NewRule(tag string, attrsRules ...string) *Rule {

	if len(attrsRules)%2 != 0 {
		panic("html.NewRule: odd argument 'attrsRules' count for tag" + tag)
	}

	attrs := make(map[string][]*regexp.Regexp)

	for i := 0; i < len(attrsRules)/2; i++ {
		name := attrsRules[i*2]
		regexStr := attrsRules[i*2+1]
		if reSymply.MatchString(regexStr) {
			regexStr = fmt.Sprintf(`(^|\s)%s($|\s)`, regexStr)
		}

		var regex *regexp.Regexp
		if len(regexStr) > 0 {
			if val, err := regexp.Compile(regexStr); err != nil {
				panic(fmt.Sprintf("html.NewRule: Wrong tag '%s' attribute rule: %s", tag, err.Error()))
			} else {
				regex = val
			}
		}

		rules, ok := attrs[name]
		if !ok {
			rules = make([]*regexp.Regexp, 0, 1)
		}
		attrs[name] = append(rules, regex)
	}

	r := &Rule{
		Tag:   tag,
		Attrs: attrs,
	}

	return r
}

func (r Rule) Equal(node *Node) bool {

	if strings.Compare(node.Tag, r.Tag) != 0 {
		return false
	}

	if len(r.Attrs) == 0 {
		return true
	}

	countChecks := 0

	for i, _ := range node.Attr {
		attr := &node.Attr[i]
		attrRules, ok := r.Attrs[attr.Key]
		if !ok {
			continue
		}

		if attrEqual(attr.Val, attrRules) {
			countChecks += 1
		} else {
			break
		}
	}

	return countChecks == len(r.Attrs)
}

func attrEqual(value string, attrRules []*regexp.Regexp) bool {

	var count int
	for _, r := range attrRules {
		if r == nil {
			if len(value) == 0 {
				count += 1
			}
		} else if r.MatchString(value) {
			count += 1
		} else {
			return false
		}
	}

	return count == len(attrRules)
}
