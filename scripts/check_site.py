#!/usr/bin/env python3
"""Validate checked-in static pages without generating or changing any files."""

import json
import re
import struct
import sys
import xml.etree.ElementTree as ET
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import unquote, urljoin, urlsplit


ROOT = Path(__file__).resolve().parents[1]
BASE = "https://golang-gui.github.io"
PAGES = {
    "/": ("index.html", "en"),
    "/zh-cn/": ("zh-cn/index.html", "zh-CN"),
    "/getting-started/": ("getting-started/index.html", "en"),
    "/zh-cn/getting-started/": ("zh-cn/getting-started/index.html", "zh-CN"),
    "/404.html": ("404.html", "en"),
}


class Page(HTMLParser):
    def __init__(self, text):
        super().__init__(convert_charrefs=True)
        self.lang = None
        self.title = ""
        self.in_title = False
        self.headings = 0
        self.ids = []
        self.links = []
        self.meta = {}
        self.canonical = None
        self.alternates = {}
        self.images = []
        self.screens = []
        self.scripts = []
        self.theme_colors = []
        self.theme_toggle = False
        self.structured = []
        self.json_text = None
        self.feed(text)
        self.close()

    def handle_starttag(self, tag, attributes):
        attrs = dict(attributes)
        if tag == "html":
            self.lang = attrs.get("lang")
        if tag == "title":
            self.in_title = True
        if tag == "h1":
            self.headings += 1
        if "id" in attrs:
            self.ids.append(attrs["id"])
        for key in ("href", "src"):
            if attrs.get(key):
                self.links.append(attrs[key])
        if tag == "meta" and "name" in attrs:
            self.meta[attrs["name"]] = attrs.get("content", "")
            if attrs["name"] == "theme-color":
                self.theme_colors.append(attrs)
        if tag == "button" and "data-theme-toggle" in attrs:
            self.theme_toggle = True
        if tag == "link":
            if attrs.get("rel") == "canonical":
                self.canonical = attrs.get("href")
            if attrs.get("rel") == "alternate" and "hreflang" in attrs:
                self.alternates[attrs["hreflang"]] = attrs.get("href")
        if tag == "img":
            self.images.append(attrs)
        if tag == "figure" and "data-screen" in attrs:
            self.screens.append(attrs)
        if tag == "script":
            if attrs.get("type") == "application/ld+json":
                self.json_text = []
            else:
                self.scripts.append(attrs.get("src"))

    handle_startendtag = handle_starttag

    def handle_data(self, text):
        if self.in_title:
            self.title += text
        if self.json_text is not None:
            self.json_text.append(text)

    def handle_endtag(self, tag):
        if tag == "title":
            self.in_title = False
        if tag == "script" and self.json_text is not None:
            self.structured.append(json.loads("".join(self.json_text)))
            self.json_text = None


def validate():
    errors = []
    pages = {}
    titles = set()

    def check(condition, message):
        if not condition:
            errors.append(message)

    for route, (filename, lang) in PAGES.items():
        text = (ROOT / filename).read_text(encoding="utf-8")
        page = Page(text)
        pages[route] = page
        check(text.lower().startswith("<!doctype html>"), f"{filename}: missing doctype")
        check(not re.search(r"{{(?:\.|\$|if|range|template|define|end)", text),
              f"{filename}: contains unrendered template markup")
        check(page.lang == lang, f"{filename}: incorrect language")
        check(page.title.strip() and page.title not in titles,
              f"{filename}: missing/duplicate title")
        titles.add(page.title)
        check(page.headings == 1, f"{filename}: expected one h1")
        check(bool(page.meta.get("description")), f"{filename}: missing description")
        check("main" in page.ids, f"{filename}: missing main anchor")
        check(len(page.ids) == len(set(page.ids)), f"{filename}: duplicate IDs")
        if route == "/404.html":
            check(page.canonical is None and page.meta.get("robots") == "noindex",
                  f"{filename}: 404 must be noindex without canonical")
        else:
            check(page.canonical == BASE + route, f"{filename}: incorrect canonical")
            check(bool(page.structured), f"{filename}: missing JSON-LD")
        interactive = route != "/404.html"
        external_scripts = [s for s in page.scripts if s]
        check(external_scripts == (["/assets/site.js"] if interactive else []),
              f"{filename}: unexpected JavaScript dependency")
        if interactive:
            contents = [m.get("content", "") for m in page.theme_colors]
            check(len(page.theme_colors) == 2 and "#f6f6f8" in contents and "#101014" in contents,
                  f"{filename}: missing light/dark theme-color metas")
            check(page.theme_toggle, f"{filename}: missing theme toggle button")
        else:
            check(not page.theme_toggle, f"{filename}: 404 must not have a theme toggle")
        if route in ("/", "/zh-cn/"):
            check([screen.get("data-screen") for screen in page.screens] ==
                  ["Windows", "macOS", "Linux"], f"{filename}: missing platform screenshots")
            check(all("hidden" not in screen for screen in page.screens),
                  f"{filename}: screenshots must be readable without JavaScript")
        for image in page.images:
            check("alt" in image, f"{filename}: image missing alt attribute")
            if image.get("src", "").endswith(".png"):
                image_file = ROOT / image["src"].lstrip("/")
                if image_file.is_file():
                    width, height = struct.unpack(">II", image_file.read_bytes()[16:24])
                    check(image.get("width") == str(width) and image.get("height") == str(height),
                          f"{filename}: incorrect image dimensions for {image['src']}")

    for route, page in pages.items():
        for link in page.links:
            url = urlsplit(urljoin(BASE + route, link))
            if url.netloc != urlsplit(BASE).netloc:
                check(url.scheme == "https", f"{route}: external link is not HTTPS: {link}")
                continue
            target = (ROOT / unquote(url.path).lstrip("/")).resolve()
            if target.is_dir():
                target /= "index.html"
            check(target.is_relative_to(ROOT) and target.is_file(),
                  f"{route}: broken internal link: {link}")
            if url.fragment and url.path in pages:
                check(url.fragment in pages[url.path].ids, f"{route}: missing anchor: {link}")
        if route != "/404.html":
            suffix = "getting-started/" if "getting-started" in route else ""
            expected = {
                "en": BASE + "/" + suffix,
                "zh-CN": BASE + "/zh-cn/" + suffix,
                "x-default": BASE + "/" + suffix,
            }
            check(page.alternates == expected, f"{route}: incorrect language links")

    sitemap = ET.parse(ROOT / "sitemap.xml")
    namespace = "{http://www.sitemaps.org/schemas/sitemap/0.9}"
    locations = [item.text for item in sitemap.findall(f"{namespace}url/{namespace}loc")]
    check(not (ROOT / "goui").exists(), "Obsolete goui page directory still exists")
    check(len(locations) == 4 and set(locations) ==
          {BASE + route for route in PAGES if route != "/404.html"},
          "Incorrect sitemap URLs")
    check((ROOT / "robots.txt").read_text() ==
          f"User-agent: *\nAllow: /\n\nSitemap: {BASE}/sitemap.xml\n", "Incorrect robots.txt")
    check((ROOT / ".nojekyll").is_file(), "Missing .nojekyll")
    google_files = sorted(ROOT.glob("google*.html"))
    check(bool(google_files), "Missing Google verification file")
    for path in google_files:
        check(path.read_text(encoding="utf-8").strip() ==
              f"google-site-verification: {path.name}",
              f"{path.name}: invalid Google verification response")
    bing = ET.parse(ROOT / "BingSiteAuth.xml").getroot()
    check(bing.tag == "users" and bool(bing.findall("user")) and
          all((user.text or "").strip() for user in bing.findall("user")),
          "BingSiteAuth.xml: invalid Bing verification response")
    if errors:
        raise ValueError("\n".join(errors))
    print("Validated 5 checked-in static HTML pages, links, images, metadata, sitemap and verification files.")


if __name__ == "__main__":
    try:
        validate()
    except (OSError, ValueError, ET.ParseError) as error:
        print(error, file=sys.stderr)
        sys.exit(1)
