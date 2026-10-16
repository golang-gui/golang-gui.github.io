package main

import (
	"html"
	"regexp"
	"strings"
	"testing"
)

var tagPattern = regexp.MustCompile("<[^>]*>")

func plainText(s string) string {
	return html.UnescapeString(tagPattern.ReplaceAllString(s, ""))
}

func TestCodeHighlightsGo(t *testing.T) {
	src := "package main\n\nimport \"fmt\"\n\n// hi\nfunc run(app ui.App) ui.RootView {\n\tfmt.Println(\"hello\")\n\treturn nil\n}\n"
	out := string(code("go", src))
	for _, want := range []string{"tok-kn", "tok-kd", "tok-s", "tok-c1", "tok-nf", "tok-nc"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing token class %s in %s", want, out)
		}
	}
	if strings.Contains(out, " style=") {
		t.Error("class mode must not emit inline styles")
	}
	if got := plainText(out); got != src {
		t.Errorf("highlighted text changed:\ngot  %q\nwant %q", got, src)
	}
}

func TestCodeHighlightsShell(t *testing.T) {
	for _, tc := range []struct {
		lang, src string
		wants     []string
	}{
		{
			lang: "bash",
			src:  "CGO_ENABLED=0 go build -trimpath -ldflags=\"-s -w\" -o hello-goui .\ngo mod init example.com/hello-goui\ngo get github.com/golang-gui/goui@v0.1.1\n./hello-goui\n",
			// env var, number, commands, subcommands, flags, module path, version, string
			wants: []string{"tok-nv", "tok-mi", "tok-nb", "tok-nt", "tok-na", "tok-s2"},
		},
		{
			lang:  "powershell",
			src:   "$env:CGO_ENABLED='0'\ngo build -o app .\n",
			wants: []string{"tok-nv", "tok-s1", "tok-n"},
		},
	} {
		out := string(code(tc.lang, tc.src))
		for _, want := range tc.wants {
			if !strings.Contains(out, want) {
				t.Errorf("%s missing token class %s in %s", tc.lang, want, out)
			}
		}
		if got := plainText(out); got != tc.src {
			t.Errorf("%s highlighted text changed:\ngot  %q\nwant %q", tc.lang, got, tc.src)
		}
	}
}

func TestCodeEscapesMarkup(t *testing.T) {
	src := "x := \"<script>alert(1)</script>\"\n"
	out := string(code("go", src))
	if strings.Contains(out, "<script>") || strings.Contains(out, "</script>") {
		t.Fatalf("code output contains raw markup: %s", out)
	}
	if got := plainText(out); got != src {
		t.Errorf("escaped text changed:\ngot  %q\nwant %q", got, src)
	}
}

func TestCodeUnknownLexerFallsBack(t *testing.T) {
	src := "a < b && c > d\n"
	out := string(code("no-such-lexer", src))
	if strings.Contains(out, "tok-") {
		t.Errorf("fallback must not emit token classes: %s", out)
	}
	if out != "a &lt; b &amp;&amp; c &gt; d\n" {
		t.Errorf("unexpected fallback: %q", out)
	}
}

func TestHighlightCSSScoped(t *testing.T) {
	css, err := highlightCSS()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(css, ".code-block .tok-k, .install-line .tok-k {") {
		t.Error("token rules are not scoped to .code-block and .install-line")
	}
	if strings.Contains(css, ".tok-chroma") {
		t.Error("chroma wrapper selectors must be rewritten")
	}
	if strings.Contains(css, "background-color: #0d1117") {
		t.Error("style background must not override the ink panel surface")
	}
}
