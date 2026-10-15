# golang-gui website

The public home for [GOUI](https://github.com/golang-gui/goui).

The website consists of checked-in HTML, CSS, images and small vanilla JavaScript
enhancements. GitHub Pages publishes these files directly. All content is in the HTML;
there is no frontend framework or build required to preview or publish the website.

- Website: <https://golang-gui.github.io/>
- 中文: <https://golang-gui.github.io/zh-cn/>
- [Website design and acceptance plan](docs/DesignWebsite.md)
- [Search engine setup](docs/SearchSetup.md)

## Pages and preview

| File | Page |
| --- | --- |
| `index.html` | English GOUI home |
| `zh-cn/index.html` | Chinese GOUI home |
| `getting-started/index.html` | English tutorial |
| `zh-cn/getting-started/index.html` | Chinese tutorial |
| `404.html` | Not-found page |

Serve the repository root with any static HTTP server. For example:

```sh
python3 -m http.server 8000
```

Open `http://localhost:8000/` for English or `http://localhost:8000/zh-cn/` for Chinese.

JavaScript is loaded only on the GOUI landing pages and tutorials to enhance the
screenshot carousel, code copying and the light/dark theme toggle. The 404 page loads no
JavaScript. Disabling JavaScript keeps all content and ordinary navigation available,
including all three platform screenshots; without JavaScript the theme follows the system
`prefers-color-scheme` setting. The carousel advances every
5.5 seconds, supports platform selection, previous/next controls, keyboard navigation
and touch swipes, and includes a pause button. Hover, focus, hidden tabs and scrolling
out of view suspend playback; reduced-motion starts with autoplay off.

The theme toggle in the header switches between light and dark modes and remembers the
choice in `localStorage`. On first visit the site follows the system preference. An inline
head script applies the stored or system theme before first paint to avoid a flash.

The hero contains a heading, a readable description with No cgo and MIT, and two
actions. Screenshots retain their original window shape without surrounding cards;
platform and playback controls sit below the image. Maintained rendering dependencies
are mentioned in the footer. Release information appears with the demo downloads.

## Maintain HTML and translations

Edit the HTML files directly to change published pages. Shared styles and native
JavaScript are in `assets/`. Validate the actual static files with:

```sh
python3 scripts/check_site.py
```

The optional Go tool helps keep English and Chinese pages consistent when updating
translations or shared layout. Go 1.24 or newer is needed only to use this tool:

```sh
go test ./...
go run ./cmd/site
```

It writes the HTML files to the repository root by default. Review and commit those
HTML files alongside the source changes. Use `-out /tmp/goui-preview` to review a
refresh separately. Synchronize direct HTML changes into the corresponding template
or translation before regenerating, because refreshing replaces the HTML files.

- Edit `content/en.json` and `content/zh-cn.json` together; translation keys must match.
- Update `content/release.json` from actual GitHub release assets, including byte sizes and SHA-256 digests.
- Keep `content/home.go.txt` (homepage) and `content/hello.go.txt` (tutorial) compilable against the version shown on the site.
- Optional regeneration templates are in `templates/`; adding a new language also
  requires adding its routes and language links in `cmd/site/main.go` and updating
  the page list and language checks in `scripts/check_site.py`.
- Copy approved screenshots and original-color logos from GOUI. Do not overwrite framework files.
- Do not add fake search verification codes, ratings, testimonials or unsupported performance claims.

Pushes to `main` validate and publish the checked-in static files through GitHub
Actions. Pull requests validate and package files only. Deployment does not run the
Go generator. `_site/` is only a temporary copy of the files prepared for upload.

## License

Website code and prose are MIT licensed. GOUI screenshots and logos come from the project.
The GOUI logo adapts the Go gopher by [Renee French](https://go.dev/blog/gopher), under
[CC BY 4.0](https://creativecommons.org/licenses/by/4.0/); it retains its own attribution.
