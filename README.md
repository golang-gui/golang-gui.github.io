# golang-gui website

The public home for [GOUI](https://github.com/golang-gui/goui) and related projects.

- Website: <https://golang-gui.github.io/goui/>
- 中文: <https://golang-gui.github.io/goui/zh-cn/>
- [Website design and acceptance plan](docs/DesignWebsite.md)
- [Search engine setup](docs/SearchSetup.md)

## Develop

Requires Go 1.24 or newer. There are no third-party Go modules or npm build dependencies.

```sh
go test ./...
go run ./cmd/site
```

The generator creates `_site/`, which is not committed. Serve that directory with any
static HTTP server to preview. Complete HTML is generated before publication; JavaScript
only enhances screenshot switching and code copying.

## Maintain

- Edit `content/en.json` and `content/zh-cn.json` together; translation keys must match.
- Update `content/release.json` from actual GitHub release assets, including byte sizes and SHA-256 digests.
- Keep `content/hello.go.txt` compilable against the version shown in the tutorial.
- Templates are in `templates/`; shared styles and JavaScript are in `assets/`.
- Copy approved screenshots and original-color logos from GOUI. Do not overwrite framework files.
- Do not add fake search verification codes, ratings, testimonials or unsupported performance claims.

Pushes to `main` build and deploy through GitHub Actions. Pull requests build and test only.
Pages uses the Actions publishing source, not a generated source branch.

## License

Website code and prose are MIT licensed. GOUI screenshots and logos come from the project.
The GOUI logo adapts the Go gopher by [Renee French](https://go.dev/blog/gopher), under
[CC BY 4.0](https://creativecommons.org/licenses/by/4.0/); it retains its own attribution.
