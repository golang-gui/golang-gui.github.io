// Progressive enhancement only: all content and links work without JavaScript.
"use strict";

for (const gallery of document.querySelectorAll("[data-gallery]")) {
  const switcher = gallery.querySelector(".platform-switch");
  const buttons = Array.from(gallery.querySelectorAll("[data-platform]"));
  const figures = Array.from(gallery.querySelectorAll("[data-screen]"));
  if (!switcher || !buttons.length || !figures.length) continue;
  const select = (platform) => {
    for (const button of buttons)
      button.setAttribute(
        "aria-pressed",
        String(button.dataset.platform === platform),
      );
    for (const figure of figures)
      figure.hidden = figure.dataset.screen !== platform;
  };
  switcher.hidden = false;
  select(buttons[0].dataset.platform);
  for (const button of buttons)
    button.addEventListener("click", () => select(button.dataset.platform));
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
