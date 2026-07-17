/* ============================================================
   theme-loader.js — Generic Theme Manager
   ============================================================
   Standalone JS. Drop it on any page alongside themes.css.
   Usage:
     1. Include <link rel="stylesheet" href="themes.css">
     2. Include <script src="theme-loader.js"></script>
     3. Add data-theme-switch attributes (see below)
   ============================================================ */

(function () {
  "use strict";

  var STORAGE_KEY = "ld-theme";
  var THEMES = [
    "system",
    "light",
    "dark",
    "high-contrast",
    "night",
    "dracula",
    "solarized",
    "catppuccin",
    "tokyo-night",
    "gruvbox",
  ];

  function clean(t) {
    return THEMES.indexOf(t) === -1 ? "system" : t;
  }

  function wantsDark(t) {
    return (
      t === "dark" ||
      t === "night" ||
      t === "dracula" ||
      t === "solarized" ||
      t === "catppuccin" ||
      t === "tokyo-night" ||
      t === "gruvbox" ||
      (t === "system" &&
        window.matchMedia &&
        window.matchMedia("(prefers-color-scheme: dark)").matches)
    );
  }

  function effective(t) {
    if (t === "system") {
      return wantsDark(t) ? "dark" : "light";
    }
    return t;
  }

  function label(t) {
    switch (t) {
      case "high-contrast":
        return "High Contrast";
      case "catppuccin":
        return "Catppuccin";
      case "tokyo-night":
        return "Tokyo Night";
      default:
        return t.charAt(0).toUpperCase() + t.slice(1);
    }
  }

  function apply(t) {
    t = clean(t);
    var e = effective(t);
    var d = document.documentElement;
    d.classList.toggle("dark", wantsDark(t));
    d.setAttribute("data-theme", e);

    document.querySelectorAll("[data-theme-label]").forEach(function (el) {
      el.textContent = label(t);
    });
    document
      .querySelectorAll("[data-theme-toggle]")
      .forEach(function (el) {
        el.setAttribute("aria-label", "Theme: " + label(t));
      });
    document
      .querySelectorAll("[data-theme-picker]")
      .forEach(function (el) {
        el.value = t;
      });
  }

  window.leetdrillTheme = {
    themes: THEMES,
    apply: apply,
    set: function (t) {
      t = clean(t);
      try {
        localStorage.setItem(STORAGE_KEY, t);
      } catch (_) {}
      apply(t);
    },
    next: function () {
      var current = clean(
        (function () {
          try {
            return localStorage.getItem(STORAGE_KEY) || "system";
          } catch (_) {
            return "system";
          }
        })()
      );
      var idx = THEMES.indexOf(current);
      var next = THEMES[(idx + 1) % THEMES.length];
      this.set(next);
    },
  };

  // Apply saved theme on load
  try {
    apply(localStorage.getItem(STORAGE_KEY) || "system");
  } catch (_) {
    apply("system");
  }

  // Re-apply on DOM ready (in case elements render after this script)
  document.addEventListener("DOMContentLoaded", function () {
    try {
      apply(localStorage.getItem(STORAGE_KEY) || "system");
    } catch (_) {
      apply("system");
    }
    // Wire up simple toggle buttons (auth pages, mobile header)
    document.querySelectorAll("[data-theme-toggle]").forEach(function (el) {
      // Only add click handler if not already handled by Alpine (no x-data parent)
      if (!el.closest("[x-data]")) {
        el.addEventListener("click", function () {
          window.leetdrillTheme.next();
        });
      }
    });
  });

  // Listen for OS-level scheme changes
  if (window.matchMedia) {
    window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", function () {
      try {
        apply(localStorage.getItem(STORAGE_KEY) || "system");
      } catch (_) {
        apply("system");
      }
    });
  }
})();
