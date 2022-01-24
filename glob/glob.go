package glob

import (
	"regexp"
	"strings"
)

type Glob struct {
	regex *regexp.Regexp
	parts []string
}

func regexFromPattern(pattern string) *regexp.Regexp {
	regexParts := []string{}
	slashSplit := strings.Split(pattern, "/")
	for _, s := range slashSplit {
		if s == "*" {
			regexParts = append(regexParts, "([^/]+)")
		} else if s == "**" {
			regexParts = append(regexParts, "(.*?)")
		} else if strings.HasPrefix(s, "!") {
			// golang doesnt have negative lookahead, so match this as any, and we will post process it later
			regexParts = append(regexParts, "([^/]+)")
		} else {
			regexParts = append(regexParts, "(" + s + ")")
		}
	}
	return regexp.MustCompile("^" + strings.Join(regexParts, "/") + "$")
}

func NewGlob(pattern string) (Glob, error) {
	return Glob{
		parts: strings.Split(pattern, "/"),
		regex: regexFromPattern(pattern),
	}, nil
}

func matchPart(part, submatch string) bool {
	// Special case for negation since golang regex dont support negative lookahead. We create the regex as a "any"
	// match, and validate it here instead
	if strings.HasPrefix(part, "!") {
		compare := part[1:]
		return compare != submatch
	}
	//Every other case has already been validated by regex
	return true
}

func (g Glob) Match(path string) bool {
	if !g.regex.MatchString(path) {
		return false
	}

	submatches := g.regex.FindStringSubmatch(path)
	//skip first submatch since it's the whole expression
	for i, submatch := range submatches[1:] {
		part := g.parts[i]
		if !matchPart(part, submatch) {
			return false
		}
	}

	return true
}