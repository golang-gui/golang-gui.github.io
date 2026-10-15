// Progressive enhancement only: all content and links work without JavaScript.
"use strict";

{
  const toggle = document.querySelector("[data-theme-toggle]");
  const media = window.matchMedia("(prefers-color-scheme: dark)");
  const colors = { light: "#f6f6f8", dark: "#101014" };
  const apply = (theme) => {
    document.documentElement.dataset.theme = theme;
    for (const meta of document.querySelectorAll("[data-theme-color]")) {
      meta.content = colors[theme];
    }
  };
  if (toggle) {
    toggle.addEventListener("click", () => {
      const next = document.documentElement.dataset.theme === "dark" ? "light" : "dark";
      try {
        localStorage.setItem("theme", next);
      } catch {}
      apply(next);
    });
  }
  media.addEventListener("change", (event) => {
    let stored;
    try {
      stored = localStorage.getItem("theme");
    } catch {}
    if (stored !== "light" && stored !== "dark") {
      apply(event.matches ? "dark" : "light");
    }
  });
}

{
  const sw = document.querySelector(".lang-switch");
  if (sw) {
    const btn = sw.querySelector(".lang-switch-btn");
    let onDocClick = null;
    const close = () => {
      sw.classList.remove("open");
      if (onDocClick) {
        document.removeEventListener("click", onDocClick);
        onDocClick = null;
      }
    };
    if (btn) {
      btn.addEventListener("click", (event) => {
        event.stopPropagation();
        if (!sw.classList.toggle("open")) {
          close();
          return;
        }
        onDocClick = (ev) => {
          if (!sw.contains(ev.target)) close();
        };
        document.addEventListener("click", onDocClick);
      });
    }
    window.addEventListener("keydown", (event) => {
      if (event.key === "Escape") close();
    });
    sw.addEventListener("focusout", (event) => {
      if (!sw.contains(event.relatedTarget)) close();
    });
  }
}

for (const gallery of document.querySelectorAll("[data-gallery]")) {
  const switcher = gallery.querySelector(".platform-switch");
  const buttons = Array.from(gallery.querySelectorAll("[data-platform]"));
  const figures = Array.from(gallery.querySelectorAll("[data-screen]"));
  const controls = gallery.querySelector(".gallery-controls");
  const previous = gallery.querySelector("[data-carousel-prev]");
  const next = gallery.querySelector("[data-carousel-next]");
  const toggle = gallery.querySelector("[data-carousel-toggle]");
  const status = gallery.querySelector("[data-gallery-status]");
  const stage = gallery.querySelector(".screenshots");
  if (
    !switcher || !buttons.length || !figures.length || !controls ||
    !previous || !next || !toggle || !status || !stage
  ) continue;

  const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
  let current = 0;
  let playing = !reducedMotion.matches;
  let hovered = false;
  let visible = false;
  let timer;
  let touchStart;
  let lastSwipeAt = 0;
  const schedule = () => {
    window.clearTimeout(timer);
    if (
      !playing || hovered || touchStart || !visible || document.hidden ||
      gallery.contains(document.activeElement)
    ) return;
    timer = window.setTimeout(() => select(current + 1), 5500);
  };
  const select = (index, announce = false) => {
    current = (index + figures.length) % figures.length;
    const platform = figures[current].dataset.screen;
    for (const button of buttons) {
      button.setAttribute(
        "aria-pressed",
        String(button.dataset.platform === platform),
      );
    }
    for (const [i, figure] of figures.entries()) {
      figure.hidden = i !== current;
    }
    const count = `${String(current + 1).padStart(2, "0")} / ${String(figures.length).padStart(2, "0")}`;
    if (announce) status.textContent = `${platform} · ${count}`;
    schedule();
  };
  const updatePlayback = () => {
    const label = playing
      ? toggle.dataset.pauseLabel
      : toggle.dataset.playLabel;
    toggle.setAttribute("aria-label", label);
    toggle.title = label;
    toggle.querySelector("[data-pause-icon]").hidden = !playing;
    toggle.querySelector("[data-play-icon]").hidden = playing;
    schedule();
  };

  gallery.setAttribute("role", "region");
  gallery.setAttribute("aria-roledescription", gallery.dataset.carouselRole);
  for (const [i, figure] of figures.entries()) {
    figure.setAttribute("role", "group");
    figure.setAttribute("aria-roledescription", gallery.dataset.slideRole);
    figure.setAttribute(
      "aria-label",
      `${figure.dataset.screen} · ${i + 1} / ${figures.length}`,
    );
  }
  // Load every image before its first automatic appearance.
  for (const image of stage.querySelectorAll("img")) image.loading = "eager";
  switcher.hidden = false;
  controls.hidden = false;
  gallery.classList.add("gallery-ready");
  select(0);
  updatePlayback();

  for (const button of buttons) {
    button.addEventListener("click", () => {
      select(
        figures.findIndex((figure) => figure.dataset.screen === button.dataset.platform),
        true,
      );
    });
  }
  previous.addEventListener("click", () => select(current - 1, true));
  next.addEventListener("click", () => select(current + 1, true));
  toggle.addEventListener("click", () => {
    playing = !playing;
    updatePlayback();
  });
  gallery.addEventListener("keydown", (event) => {
    if (
      !event.target.closest("button") || event.altKey || event.ctrlKey ||
      event.metaKey || event.shiftKey
    ) return;
    const destinations = {
      ArrowLeft: current - 1,
      ArrowRight: current + 1,
      Home: 0,
      End: figures.length - 1,
    };
    if (!Object.hasOwn(destinations, event.key)) return;
    event.preventDefault();
    select(destinations[event.key], true);
  });
  gallery.addEventListener("pointerenter", (event) => {
    if (event.pointerType !== "mouse") return;
    hovered = true;
    schedule();
  });
  gallery.addEventListener("pointerleave", (event) => {
    if (event.pointerType !== "mouse") return;
    hovered = false;
    schedule();
  });
  gallery.addEventListener("focusin", schedule);
  gallery.addEventListener("focusout", () => window.setTimeout(schedule, 0));
  document.addEventListener("visibilitychange", schedule);
  reducedMotion.addEventListener("change", () => {
    playing = !reducedMotion.matches;
    updatePlayback();
  });

  if ("IntersectionObserver" in window) {
    const observer = new IntersectionObserver(
      ([entry]) => {
        visible = entry.isIntersecting;
        schedule();
      },
      { threshold: 0.15 },
    );
    observer.observe(gallery);
  } else {
    visible = true;
    schedule();
  }

  stage.addEventListener("pointerdown", (event) => {
    if (event.pointerType === "mouse" || !event.isPrimary) return;
    touchStart = { id: event.pointerId, x: event.clientX, y: event.clientY };
    schedule();
  });
  stage.addEventListener("pointerup", (event) => {
    if (touchStart?.id !== event.pointerId) return;
    const dx = event.clientX - touchStart.x;
    const dy = event.clientY - touchStart.y;
    touchStart = null;
    if (Math.abs(dx) > 48 && Math.abs(dx) > Math.abs(dy) * 1.25) {
      lastSwipeAt = Date.now();
      select(current + (dx < 0 ? 1 : -1), true);
    } else {
      schedule();
    }
  });
  stage.addEventListener("pointercancel", () => {
    touchStart = null;
    schedule();
  });
  stage.addEventListener("click", (event) => {
    if (Date.now() - lastSwipeAt < 350) event.preventDefault();
  });
}

for (const button of document.querySelectorAll("[data-copy]")) {
  button.hidden = false;
  button.addEventListener("click", async () => {
    const code = button.closest(".code-block")?.querySelector("pre code");
    const status = document.getElementById("copy-status");
    if (!code || !status) return;
    try {
      await navigator.clipboard.writeText(code.textContent);
      button.textContent = document.body.dataset.copiedLabel;
      status.textContent = document.body.dataset.copiedLabel;
      window.setTimeout(() => {
        button.textContent = document.body.dataset.copyLabel;
      }, 1600);
    } catch {
      status.textContent = document.body.dataset.copyFailed;
    }
  });
}
