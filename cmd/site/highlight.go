package main

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"html/template"
	"regexp"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma/v2"
	htmlfmt "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// highlightStyle tints the ink panel's fixed dark surface in both site themes.
const highlightStyleName = "github-dark"

var highlightFormatter = htmlfmt.New(
	htmlfmt.WithClasses(true),
	htmlfmt.ClassPrefix("tok-"),
	htmlfmt.PreventSurroundingPre(true),
)

// code renders src as syntax-highlighted HTML for the named language.
// Unknown languages and rendering failures fall back to escaped plain text.
func code(lang, src string) template.HTML {
	iterator, err := tokenize(lang, src)
	if err != nil {
		return template.HTML(stdhtml.EscapeString(src))
	}
	var out bytes.Buffer
	if err := highlightFormatter.Format(&out, styles.Get(highlightStyleName), iterator); err != nil {
		return template.HTML(stdhtml.EscapeString(src))
	}
	return template.HTML(out.String())
}

func tokenize(lang, src string) (chroma.Iterator, error) {
	switch lang {
	case "go":
		return goIterator(src)
	case "bash":
		return tokenIterator(shellTokens(src)), nil
	}
	lexer := lexers.Get(lang)
	if lexer == nil {
		return nil, fmt.Errorf("unknown lexer %q", lang)
	}
	return lexer.Tokenise(nil, src)
}

// goIterator promotes exported identifiers to NameClass: chroma's Go lexer
// emits types such as ui.App as NameOther, which the style leaves uncolored.
func goIterator(src string) (chroma.Iterator, error) {
	it, err := lexers.Get("go").Tokenise(nil, src)
	if err != nil {
		return nil, err
	}
	return mapTokens(it, func(t chroma.Token) chroma.Token {
		if t.Type == chroma.NameOther && isExported(t.Value) {
			t.Type = chroma.NameClass
		}
		return t
	}), nil
}

func mapTokens(it chroma.Iterator, fn func(chroma.Token) chroma.Token) chroma.Iterator {
	return func() chroma.Token {
		t := it()
		if t == chroma.EOF {
			return t
		}
		return fn(t)
	}
}

func tokenIterator(tokens []chroma.Token) chroma.Iterator {
	return func() chroma.Token {
		if len(tokens) == 0 {
			return chroma.EOF
		}
		token := tokens[0]
		tokens = tokens[1:]
		return token
	}
}

func isExported(word string) bool {
	for _, r := range word {
		return unicode.IsUpper(r)
	}
	return false
}

// shellTokens is a small line-oriented shell tokenizer. Chroma's bash lexer
// only marks builtins, leaving go get or mkdir as plain text; this one colors
// the command word, subcommands, flags, quoted strings, comments, variables,
// VAR=value prefixes, module paths and @version pins while preserving the
// source byte for byte.
func shellTokens(src string) []chroma.Token {
	var tokens []chroma.Token
	st := &shellState{command: true}
	i := 0
	for i < len(src) {
		c := src[i]
		switch {
		case c == '\n':
			tokens = append(tokens, chroma.Token{Type: chroma.TextWhitespace, Value: "\n"})
			i++
			st.nextCommand()
		case c == ' ' || c == '\t' || c == '\r':
			j := i
			for j < len(src) && (src[j] == ' ' || src[j] == '\t' || src[j] == '\r') {
				j++
			}
			tokens = append(tokens, chroma.Token{Type: chroma.TextWhitespace, Value: src[i:j]})
			i = j
		case c == '#':
			j := i
			for j < len(src) && src[j] != '\n' {
				j++
			}
			tokens = append(tokens, chroma.Token{Type: chroma.CommentSingle, Value: src[i:j]})
			i = j
		case c == '"' || c == '\'':
			token, next := shellString(src, i)
			tokens = append(tokens, token)
			i = next
		case c == '$':
			token, next := shellVariable(src, i)
			tokens = append(tokens, token)
			i = next
		case c == '&' || c == '|' || c == ';' || c == '(' || c == ')':
			j := i + 1
			for j < len(src) && src[j] == c {
				j++
			}
			tokens = append(tokens, chroma.Token{Type: chroma.Punctuation, Value: src[i:j]})
			i = j
			st.nextCommand()
		default:
			j := i
			for j < len(src) && !shellWordEnd(src[j]) {
				j++
			}
			word := src[i:j]
			i = j
			tokens = append(tokens, shellWord(word, st)...)
		}
	}
	return tokens
}

type shellState struct {
	command  bool
	context  string
	subDepth int
}

func (st *shellState) nextCommand() {
	st.command = true
	st.context = ""
	st.subDepth = 0
}

func shellWordEnd(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r', '"', '\'', '$', '&', '|', ';', '(', ')':
		return true
	}
	return false
}

func shellWord(word string, st *shellState) []chroma.Token {
	switch {
	case word == "":
		return nil
	case strings.HasPrefix(word, "-"):
		if body, ok := strings.CutSuffix(word, "="); ok {
			return []chroma.Token{
				{Type: chroma.NameAttribute, Value: body},
				{Type: chroma.Punctuation, Value: "="},
			}
		}
		return []chroma.Token{{Type: chroma.NameAttribute, Value: word}}
	case st.command:
		if name, value, ok := strings.Cut(word, "="); ok && isShellName(name) {
			tokens := []chroma.Token{
				{Type: chroma.NameVariable, Value: name},
				{Type: chroma.Punctuation, Value: "="},
			}
			if value != "" {
				tokens = append(tokens, chroma.Token{Type: shellValueToken(value), Value: value})
			}
			return tokens // the command itself still follows
		}
		st.command = false
		st.context = word
		st.subDepth = 0
		if _, ok := shellSubcommands[word]; ok {
			st.subDepth = 1
		}
		return []chroma.Token{{Type: chroma.NameBuiltin, Value: word}}
	case st.subDepth > 0:
		if st.subDepth == 2 { // third word after go mod / git remote ...
			st.subDepth = 0
			if isShellName(word) {
				return []chroma.Token{{Type: chroma.NameTag, Value: word}}
			}
		} else if shellSubcommands[st.context][word] {
			if shellSecondLevel[word] {
				st.subDepth, st.context = 2, word
			} else {
				st.subDepth = 0
			}
			return []chroma.Token{{Type: chroma.NameTag, Value: word}}
		} else {
			st.subDepth = 0
		}
	}
	if path, version, ok := strings.Cut(word, "@"); ok && path != "" && isVersionPin(version) {
		return []chroma.Token{
			{Type: chroma.NameVariable, Value: path},
			{Type: chroma.Punctuation, Value: "@"},
			{Type: chroma.NameTag, Value: version},
		}
	}
	if strings.Contains(word, "/") {
		return []chroma.Token{{Type: chroma.NameVariable, Value: word}}
	}
	return []chroma.Token{{Type: chroma.Text, Value: word}}
}

func isVersionPin(s string) bool {
	if s == "" {
		return false
	}
	if s[0] == 'v' {
		s = s[1:]
	}
	if s == "" || s[0] < '0' || s[0] > '9' {
		return false
	}
	return true
}

// shellSubcommands colors the word after known multi-word commands, and
// shellSecondLevel keeps the cursor open for a third word (go mod init).
var shellSubcommands = map[string]map[string]bool{
	"go":        {"get": true, "mod": true, "build": true, "run": true, "test": true, "install": true, "fmt": true, "vet": true, "generate": true, "list": true, "version": true, "clean": true, "tool": true, "work": true, "env": true, "doc": true, "cover": true},
	"git":       {"clone": true, "init": true, "add": true, "commit": true, "push": true, "pull": true, "fetch": true, "checkout": true, "branch": true, "merge": true, "rebase": true, "log": true, "status": true, "diff": true, "remote": true, "tag": true, "stash": true, "reset": true, "show": true},
	"docker":    {"build": true, "run": true, "push": true, "pull": true, "images": true, "ps": true, "exec": true, "stop": true, "rm": true, "logs": true, "compose": true},
	"npm":       {"install": true, "ci": true, "run": true, "start": true, "test": true, "build": true, "publish": true, "update": true, "init": true},
	"cargo":     {"build": true, "run": true, "test": true, "check": true, "add": true, "update": true, "publish": true, "fmt": true},
	"brew":      {"install": true, "uninstall": true, "update": true, "upgrade": true, "search": true, "list": true},
	"apt":       {"install": true, "remove": true, "update": true, "upgrade": true, "search": true},
	"systemctl": {"start": true, "stop": true, "restart": true, "enable": true, "disable": true, "status": true},
	"kubectl":   {"get": true, "apply": true, "delete": true, "describe": true, "logs": true, "exec": true, "create": true, "rollout": true},
}

var shellSecondLevel = map[string]bool{
	"mod": true, "remote": true, "container": true, "image": true,
	"network": true, "volume": true, "compose": true,
}

func shellValueToken(value string) chroma.TokenType {
	for _, r := range value {
		if r < '0' || r > '9' {
			return chroma.Text
		}
	}
	return chroma.LiteralNumberInteger
}

func isShellName(word string) bool {
	if word == "" {
		return false
	}
	for _, r := range word {
		if r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func shellString(src string, start int) (chroma.Token, int) {
	quote := src[start]
	typ := chroma.LiteralStringSingle
	if quote == '"' {
		typ = chroma.LiteralStringDouble
	}
	j := start + 1
	for j < len(src) {
		if quote == '"' && src[j] == '\\' && j+1 < len(src) {
			j += 2
			continue
		}
		if src[j] == quote {
			j++
			break
		}
		if src[j] == '\n' {
			break // an unterminated string ends at the line
		}
		j++
	}
	return chroma.Token{Type: typ, Value: src[start:j]}, j
}

func shellVariable(src string, start int) (chroma.Token, int) {
	j := start + 1
	if j < len(src) && src[j] == '{' {
		for j < len(src) && src[j] != '}' {
			j++
		}
		if j < len(src) {
			j++
		}
	} else {
		for j < len(src) {
			c := src[j]
			if c != '_' && c != '?' && !unicode.IsLetter(rune(c)) && !unicode.IsDigit(rune(c)) {
				break
			}
			j++
		}
	}
	return chroma.Token{Type: chroma.NameVariable, Value: src[start:j]}, j
}

var chromaSelectorRE = regexp.MustCompile(`\.tok-chroma \.tok-([\w-]+)`)

// highlightCSS re-scopes the chroma stylesheet to the site's code blocks and
// drops wrapper, background and line-table rules the markup never emits.
func highlightCSS() (string, error) {
	var raw bytes.Buffer
	if err := highlightFormatter.WriteCSS(&raw, styles.Get(highlightStyleName)); err != nil {
		return "", err
	}
	var out strings.Builder
	for _, line := range strings.Split(raw.String(), "\n") {
		if unwantedChromaRule(line) {
			continue
		}
		out.WriteString(chromaSelectorRE.ReplaceAllString(line, ".code-block .tok-$1, .install-line .tok-$1"))
		out.WriteByte('\n')
	}
	return out.String(), nil
}

func unwantedChromaRule(line string) bool {
	for _, name := range []string{
		"Background", "PreWrapper", "LineLink", "LineTableTD", "LineTable",
		"LineHighlight", "LineNumbersTable", "LineNumbers", "Line",
	} {
		if strings.HasPrefix(strings.TrimSpace(line), "/* "+name+" */") {
			return true
		}
	}
	return false
}
