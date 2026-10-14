package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"go/parser"
	"go/token"
	"html"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var linkPattern = regexp.MustCompile("(?:href|src)=\"([^\"]+)\"")
var idPattern = regexp.MustCompile("id=\"([^\"]+)\"")
var jsonPattern = regexp.MustCompile("(?s)<script type=\"application/ld\\+json\">(.*?)</script>")

func sourceRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func outputFile(out, path string) string {
	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}
	return filepath.Join(out, strings.TrimPrefix(path, "/"))
}

func TestBuild(t *testing.T) {
	root, out := sourceRoot(t), t.TempDir()
	if err := build(root, out); err != nil {
		t.Fatal(err)
	}
	pages, err := load(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 6 {
		t.Fatalf("got %d pages, want 5 public pages plus 404", len(pages))
	}
	seen := make(map[string]bool)
	titles := make(map[string]bool)
	for _, p := range pages {
		t.Run(p.Path, func(t *testing.T) {
			if seen[p.Path] {
				t.Fatalf("duplicate route %s", p.Path)
			}
			seen[p.Path] = true
			data, err := os.ReadFile(outputFile(out, p.Path))
			if err != nil {
				t.Fatal(err)
			}
			body := string(data)
			if strings.Contains(body, "ZgotmplZ") {
				t.Fatal("unsafe or invalid template URL")
			}
			if titles[p.Title] {
				t.Fatalf("duplicate page title %s", p.Title)
			}
			titles[p.Title] = true
			for _, want := range []string{"<html lang=\"" + p.Lang + "\">", "<title>" + html.EscapeString(p.Title) + "</title>", "<meta name=\"description\"", "id=\"main\""} {
				if !strings.Contains(body, want) {
					t.Errorf("missing %s", want)
				}
			}
			if n := strings.Count(body, "<h1>"); n != 1 {
				t.Errorf("got %d h1 elements, want one", n)
			}
			if p.Kind == "404" {
				if !strings.Contains(body, "name=\"robots\" content=\"noindex\"") || strings.Contains(body, "rel=\"canonical\"") {
					t.Fatal("404 must be noindex and have no successful-page canonical")
				}
			} else {
				if !strings.Contains(body, "rel=\"canonical\" href=\""+p.Canonical+"\"") {
					t.Fatal("missing self canonical")
				}
				m := jsonPattern.FindStringSubmatch(body)
				if len(m) != 2 || !json.Valid([]byte(m[1])) {
					t.Fatal("missing or invalid JSON-LD")
				}
			}
			if p.EnglishURL != "" {
				for lang, target := range map[string]string{"en": p.EnglishURL, "zh-CN": p.ChineseURL, "x-default": p.EnglishURL} {
					if !strings.Contains(body, "hreflang=\""+lang+"\" href=\""+target+"\"") {
						t.Errorf("missing %s alternate", lang)
					}
				}
			}
			if p.Kind == "home" || p.Kind == "start" {
				if !strings.Contains(body, "github.com/golang-gui/goui@v0.1.0") || strings.Contains(body, "goui@latest") {
					t.Fatal("tutorial must pin the released version")
				}
			}
			checkLinks(t, out, p.Path, body)
		})
	}
	robots, err := os.ReadFile(filepath.Join(out, "robots.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(robots) != "User-agent: *\nAllow: /\n\nSitemap: https://golang-gui.github.io/sitemap.xml\n" {
		t.Fatalf("unexpected robots.txt: %s", robots)
	}
	checkSitemap(t, out)
}

func checkLinks(t *testing.T, out, current, body string) {
	t.Helper()
	base, err := url.Parse(baseURL + current)
	if err != nil {
		t.Fatal(err)
	}
	for _, match := range linkPattern.FindAllStringSubmatch(body, -1) {
		u, err := url.Parse(html.UnescapeString(match[1]))
		if err != nil {
			t.Errorf("invalid URL %s: %v", match[1], err)
			continue
		}
		u = base.ResolveReference(u)
		if u.Host != base.Host {
			if u.Scheme != "https" {
				t.Errorf("external link is not HTTPS: %s", u)
			}
			continue
		}
		file := outputFile(out, u.Path)
		data, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("broken internal link %s: %v", u, err)
			continue
		}
		if u.Fragment != "" {
			found := false
			for _, id := range idPattern.FindAllStringSubmatch(string(data), -1) {
				found = found || id[1] == u.Fragment
			}
			if !found {
				t.Errorf("missing anchor %s in %s", u.Fragment, u.Path)
			}
		}
	}
}

func checkSitemap(t *testing.T, out string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(out, "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}
	d := xml.NewDecoder(bytes.NewReader(data))
	locations := make(map[string]bool)
	for {
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		if start.Name.Local == "urlset" && start.Name.Space != "http://www.sitemaps.org/schemas/sitemap/0.9" {
			t.Fatal("wrong sitemap namespace")
		}
		if start.Name.Local == "loc" {
			var location string
			if err := d.DecodeElement(&location, &start); err != nil {
				t.Fatal(err)
			}
			if locations[location] || !strings.HasPrefix(location, baseURL+"/") || strings.Contains(location, "404") {
				t.Fatalf("invalid sitemap location %s", location)
			}
			locations[location] = true
		}
	}
	for _, path := range []string{"/", "/goui/", "/goui/zh-cn/", "/goui/getting-started/", "/goui/zh-cn/getting-started/"} {
		if !locations[baseURL+path] {
			t.Errorf("missing sitemap URL %s", path)
		}
	}
	if len(locations) != 5 {
		t.Errorf("got %d sitemap URLs, want five", len(locations))
	}
}

func TestReleaseAndExample(t *testing.T) {
	pages, err := load(sourceRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	p := pages[0]
	for _, asset := range p.Downloads {
		if asset.Bytes <= 0 || asset.URL != projectURL+"/releases/download/"+p.Version+"/"+asset.File {
			t.Errorf("invalid download %s", asset.File)
		}
		digest, err := hex.DecodeString(asset.SHA256)
		if err != nil || len(digest) != 32 {
			t.Errorf("invalid SHA256 for %s", asset.File)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "main.go", p.Example, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	for _, screenshot := range p.Screenshots {
		if screenshot.Width <= 0 || screenshot.Height <= 0 {
			t.Fatal("missing native image dimensions")
		}
	}
	for _, p := range pages[:4] {
		var other *Page
		for i := range pages {
			if pages[i].Path == p.OtherPath {
				other = &pages[i]
			}
		}
		if other == nil || other.OtherPath != p.Path || other.Kind != p.Kind {
			t.Errorf("language alternatives do not reciprocate for %s", p.Path)
		}
	}
}
