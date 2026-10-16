// site optionally refreshes checked-in static pages.
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	"image"
	_ "image/png"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const baseURL = "https://golang-gui.github.io"
const projectURL = "https://github.com/golang-gui/goui"

type Item struct {
	Title string
	Text  string
}

type Locale struct {
	Lang, OGLocale string
	Text           map[string]string
	Features       []Item
	FAQs           []Item
}

type Download struct {
	Platform, Target, File, URL, SHA256 string
	Bytes                               int64
}

func (d Download) Size() string { return fmt.Sprintf("%.2f MiB", float64(d.Bytes)/(1024*1024)) }

type Release struct {
	Version, Published, URL, Checksums  string
	DemoVersion, DemoPublished, DemoURL string
	Downloads                           []Download
}

type Screenshot struct {
	Platform, File string
	Width, Height  int
}

type Page struct {
	Locale
	Release
	Kind, Path, Title, Description, HomePath, StartPath, OtherPath string
	Canonical, EnglishURL, ChineseURL, OtherLabel                  string
	StructuredData                                                 template.JS
	Screenshots                                                    []Screenshot
	Example, Install, Build, WindowsBuild                          string
}

func (p Page) T(key string) (string, error) {
	value, ok := p.Text[key]
	if !ok {
		return "", fmt.Errorf("missing %s text key %q", p.Lang, key)
	}
	return value, nil
}

func main() {
	root := flag.String("root", ".", "source directory")
	out := flag.String("out", ".", "output directory")
	highlightCSSFlag := flag.Bool("highlight-css", false, "write the code highlighting stylesheet to stdout")
	flag.Parse()
	if *highlightCSSFlag {
		css, err := highlightCSS()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(css)
		return
	}
	if err := build(*root, *out); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated 4 public pages and a 404 page in %s\n", *out)
}

func readJSON(file string, value any) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}
	return nil
}

func load(root string) ([]Page, error) {
	var en, zh Locale
	var release Release
	for file, value := range map[string]any{
		"en.json": &en, "zh-cn.json": &zh, "release.json": &release,
	} {
		if err := readJSON(filepath.Join(root, "content", file), value); err != nil {
			return nil, err
		}
	}
	for _, locale := range []Locale{en, zh} {
		if len(locale.Text) != len(en.Text) {
			return nil, fmt.Errorf("%s translations have different key counts", locale.Lang)
		}
		for key := range en.Text {
			if strings.TrimSpace(locale.Text[key]) == "" {
				return nil, fmt.Errorf("%s missing translation %s", locale.Lang, key)
			}
		}
	}
	if release.Version == "" || len(release.Downloads) != 3 {
		return nil, fmt.Errorf("release requires a version and three verified assets")
	}
	homeExample, err := os.ReadFile(filepath.Join(root, "content", "home.go.txt"))
	if err != nil {
		return nil, err
	}
	tutorialExample, err := os.ReadFile(filepath.Join(root, "content", "hello.go.txt"))
	if err != nil {
		return nil, err
	}
	var screenshots []Screenshot
	for _, platform := range []string{"Windows", "macOS", "Linux"} {
		file := "demo-" + strings.ToLower(platform) + ".png"
		f, err := os.Open(filepath.Join(root, "assets", file))
		if err != nil {
			return nil, err
		}
		config, _, decodeErr := image.DecodeConfig(f)
		closeErr := f.Close()
		if decodeErr != nil {
			return nil, decodeErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		screenshots = append(screenshots, Screenshot{platform, file, config.Width, config.Height})
	}
	var pages []Page
	for _, locale := range []Locale{en, zh} {
		home, start, otherHome, otherStart, otherLabel := "/", "/getting-started/", "/zh-cn/", "/zh-cn/getting-started/", "中文"
		if locale.Lang == "zh-CN" {
			home, start, otherHome, otherStart, otherLabel = otherHome, otherStart, home, start, "English"
		}
		for _, kind := range []string{"home", "start"} {
			example := homeExample
			if kind == "start" {
				example = tutorialExample
			}
			pagePath, otherPath := home, otherHome
			enURL, zhURL := baseURL+"/", baseURL+"/zh-cn/"
			if kind == "start" {
				pagePath, otherPath = start, otherStart
				enURL, zhURL = baseURL+"/getting-started/", baseURL+"/zh-cn/getting-started/"
			}
			p := Page{
				Locale: locale, Release: release, Kind: kind, Path: pagePath,
				Title: locale.Text[kind+"Title"], Description: locale.Text[kind+"Description"],
				HomePath: home, StartPath: start, OtherPath: otherPath, OtherLabel: otherLabel,
				Canonical: baseURL + pagePath, EnglishURL: enURL, ChineseURL: zhURL,
				Screenshots: screenshots, Example: string(example),
				Install:      "mkdir hello-goui\ncd hello-goui\ngo mod init example.com/hello-goui\ngo get github.com/golang-gui/goui@" + release.Version,
				Build:        "CGO_ENABLED=0 go build -trimpath -ldflags=\"-s -w\" -o hello-goui .\n./hello-goui",
				WindowsBuild: "$env:CGO_ENABLED='0'\ngo build -trimpath -ldflags=\"-s -w -H=windowsgui\" -o hello-goui.exe .\n.\\hello-goui.exe",
			}
			data, err := json.Marshal(map[string]any{
				"@context": "https://schema.org", "@type": "SoftwareSourceCode",
				"@id": baseURL + "/#project", "name": "GOUI", "url": baseURL + "/",
				"description": locale.Text["homeDescription"], "codeRepository": projectURL,
				"programmingLanguage": "Go", "runtimePlatform": "Windows, macOS, Linux / X11",
				"version": release.Version, "license": "https://opensource.org/license/mit",
				"author": map[string]any{"@type": "Organization", "name": "golang-gui", "url": baseURL + "/"},
			})
			if err != nil {
				return nil, err
			}
			// encoding/json escapes HTML delimiters; this is generated data, never user HTML.
			p.StructuredData = template.JS(data)
			pages = append(pages, p)
		}
	}
	notFound := pages[0]
	notFound.Kind, notFound.Path, notFound.Canonical = "404", "/404.html", ""
	notFound.Title, notFound.Description, notFound.StructuredData = "Page not found — GOUI", "Return to GOUI, the native Go GUI framework.", ""
	notFound.EnglishURL, notFound.ChineseURL, notFound.OtherPath, notFound.OtherLabel = "", "", "", ""
	return append(pages, notFound), nil
}

func build(root, out string) error {
	pages, err := load(root)
	if err != nil {
		return err
	}
	templates, err := template.New("site").Funcs(template.FuncMap{"code": code}).ParseGlob(filepath.Join(root, "templates", "*.html"))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(out, 0755); err != nil {
		return err
	}
	for _, p := range pages {
		path := strings.TrimPrefix(p.Path, "/")
		if strings.HasSuffix(p.Path, "/") {
			path = filepath.Join(path, "index.html")
		}
		file := filepath.Join(out, path)
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			return err
		}
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		renderErr := templates.ExecuteTemplate(f, "layout", p)
		closeErr := f.Close()
		if renderErr != nil {
			return fmt.Errorf("render %s: %w", p.Path, renderErr)
		}
		if closeErr != nil {
			return closeErr
		}
	}
	if err := copyAssets(filepath.Join(root, "assets"), filepath.Join(out, "assets")); err != nil {
		return err
	}
	if err := writeSitemap(filepath.Join(out, "sitemap.xml"), pages); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(out, "robots.txt"), []byte("User-agent: *\nAllow: /\n\nSitemap: "+baseURL+"/sitemap.xml\n"), 0644)
}

func copyAssets(src, dst string) error {
	source, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	destination, err := filepath.Abs(dst)
	if err != nil {
		return err
	}
	if source == destination {
		return nil
	}
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			return fmt.Errorf("asset symlink is not supported: %s", path)
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}

func writeSitemap(file string, pages []Page) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	writeErr := sitemap(f, pages)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

func sitemap(w io.Writer, pages []Page) error {
	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	root := xml.StartElement{Name: xml.Name{Local: "urlset"}, Attr: []xml.Attr{{Name: xml.Name{Local: "xmlns"}, Value: "http://www.sitemaps.org/schemas/sitemap/0.9"}}}
	if err := e.EncodeToken(root); err != nil {
		return err
	}
	seen := make(map[string]bool)
	for _, p := range pages {
		if p.Kind == "404" || seen[p.Canonical] {
			continue
		}
		seen[p.Canonical] = true
		start := xml.StartElement{Name: xml.Name{Local: "url"}}
		if err := e.EncodeToken(start); err != nil {
			return err
		}
		if err := e.EncodeElement(p.Canonical, xml.StartElement{Name: xml.Name{Local: "loc"}}); err != nil {
			return err
		}
		if err := e.EncodeToken(start.End()); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(root.End()); err != nil {
		return err
	}
	return e.Flush()
}
