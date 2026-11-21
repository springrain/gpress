(function () {
  var root = document.documentElement;
  var header = document.querySelector("[data-header]");
  var menuButton = document.querySelector("[data-menu-toggle]");
  var mobilePanel = document.querySelector("[data-mobile-panel]");
  var themeButton = document.querySelector("[data-theme-toggle]");
  var searchForm = document.querySelector("[data-search-form]");
  var searchToggle = document.querySelector("[data-search-toggle]");
  var backtop = document.querySelector("[data-backtop]");

  function applyTheme(theme) {
    root.classList.toggle("dark", theme === "dark");
  }

  var storageKey = "fixit-theme";
  var storedTheme = localStorage.getItem(storageKey);
  var prefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
  applyTheme(storedTheme || (prefersDark ? "dark" : "light"));

  if (themeButton) {
    themeButton.addEventListener("click", function () {
      var nextTheme = root.classList.contains("dark") ? "light" : "dark";
      localStorage.setItem(storageKey, nextTheme);
      applyTheme(nextTheme);
    });
  }

  if (menuButton && mobilePanel) {
    menuButton.addEventListener("click", function () {
      var isHidden = mobilePanel.hasAttribute("hidden");
      if (isHidden) {
        mobilePanel.removeAttribute("hidden");
      } else {
        mobilePanel.setAttribute("hidden", "");
      }
      menuButton.setAttribute("aria-expanded", String(isHidden));
    });
  }

  if (searchForm && searchToggle) {
    var searchInput = searchForm.querySelector(".tm-search__input");

    function openSearch() {
      searchForm.classList.add("is-open");
      if (searchInput) {
        window.setTimeout(function () {
          searchInput.focus();
        }, 40);
      }
    }

    function closeSearch() {
      if (searchInput && searchInput.value.trim()) {
        return;
      }
      searchForm.classList.remove("is-open");
    }

    if (searchInput) {
      searchInput.addEventListener("focus", openSearch);
    }

    searchForm.addEventListener("submit", function (event) {
      var hasQuery = searchInput && searchInput.value.trim();
      if (!searchForm.classList.contains("is-open") && !hasQuery) {
        event.preventDefault();
        openSearch();
        return;
      }
    });

    document.addEventListener("click", function (event) {
      if (!searchForm.contains(event.target)) {
        closeSearch();
      }
    });

    document.addEventListener("keydown", function (event) {
      if (event.key === "Escape") {
        closeSearch();
      }
    });
  }

  function handleScroll() {
    var scrolled = window.scrollY > 8;
    if (header) {
      header.classList.toggle("is-scrolled", scrolled);
    }
    if (backtop) {
      backtop.classList.toggle("is-visible", window.scrollY > 480);
    }
  }

  window.addEventListener("scroll", handleScroll, { passive: true });
  handleScroll();

  if (backtop) {
    backtop.addEventListener("click", function () {
      window.scrollTo({ top: 0, behavior: "smooth" });
    });
  }

  var toc = document.querySelector("[data-toc]");
  if (toc) {
    var tocContent = toc.querySelector("[data-toc-content]");
    var tocToggle = toc.querySelector("[data-toc-toggle]");
    var tocLinks = Array.prototype.slice.call(toc.querySelectorAll("a[href^='#']"));
    var tocItems = tocLinks
      .map(function (link) {
        var rawId = link.getAttribute("href").slice(1);
        var id = rawId;
        try {
          id = decodeURIComponent(rawId);
        } catch (error) {
          id = rawId;
        }
        return { link: link, heading: document.getElementById(id) };
      })
      .filter(function (item) {
        return Boolean(item.heading);
      });
    var activeTocLink = null;
    var tocFrame = 0;

    function revealTocBranch(link) {
      toc.querySelectorAll("li.has-active").forEach(function (item) {
        item.classList.remove("has-active");
      });

      var parent = link.closest("li");
      while (parent && toc.contains(parent)) {
        parent.classList.add("has-active");
        parent = parent.parentElement ? parent.parentElement.closest("li") : null;
      }
    }

    function scrollTocLinkIntoView(link) {
      if (!tocContent || tocContent.hidden) {
        return;
      }

      var containerRect = tocContent.getBoundingClientRect();
      var linkRect = link.getBoundingClientRect();
      if (linkRect.top < containerRect.top) {
        tocContent.scrollTop += linkRect.top - containerRect.top;
      } else if (linkRect.bottom > containerRect.bottom) {
        tocContent.scrollTop += linkRect.bottom - containerRect.bottom;
      }
    }

    function updateToc() {
      tocFrame = 0;
      if (!tocItems.length) {
        return;
      }

      var headerOffset = header ? header.offsetHeight : 0;
      var threshold = window.scrollY + headerOffset + 24;
      var activeItem = tocItems[0];

      tocItems.forEach(function (item) {
        var headingTop = window.scrollY + item.heading.getBoundingClientRect().top;
        if (headingTop <= threshold) {
          activeItem = item;
        }
      });

      if (activeTocLink === activeItem.link) {
        return;
      }

      tocLinks.forEach(function (link) {
        link.classList.toggle("is-active", link === activeItem.link);
        if (link === activeItem.link) {
          link.setAttribute("aria-current", "location");
        } else {
          link.removeAttribute("aria-current");
        }
      });
      activeTocLink = activeItem.link;
      revealTocBranch(activeItem.link);
      scrollTocLinkIntoView(activeItem.link);
    }

    function scheduleTocUpdate() {
      if (!tocFrame) {
        tocFrame = window.requestAnimationFrame(updateToc);
      }
    }

    if (tocToggle && tocContent) {
      tocToggle.addEventListener("click", function () {
        var expanded = tocToggle.getAttribute("aria-expanded") === "true";
        tocToggle.setAttribute("aria-expanded", String(!expanded));
        toc.classList.toggle("is-collapsed", expanded);
        tocContent.hidden = expanded;
      });
    }

    window.addEventListener("scroll", scheduleTocUpdate, { passive: true });
    window.addEventListener("resize", scheduleTocUpdate);
    window.addEventListener("load", scheduleTocUpdate);
    window.addEventListener("hashchange", scheduleTocUpdate);
    updateToc();
  }

  Array.prototype.slice.call(document.querySelectorAll("[data-pagination]")).forEach(function (pagination) {
    pagination.addEventListener("keydown", function (event) {
      if (["ArrowLeft", "ArrowRight", "Home", "End"].indexOf(event.key) === -1) {
        return;
      }

      var items = Array.prototype.slice.call(
        pagination.querySelectorAll("a[href], [aria-current='page']")
      );
      var currentIndex = items.indexOf(document.activeElement);
      if (currentIndex === -1 || !items.length) {
        return;
      }

      var nextIndex = currentIndex;
      if (event.key === "ArrowLeft") {
        nextIndex = Math.max(0, currentIndex - 1);
      } else if (event.key === "ArrowRight") {
        nextIndex = Math.min(items.length - 1, currentIndex + 1);
      } else if (event.key === "Home") {
        nextIndex = 0;
      } else if (event.key === "End") {
        nextIndex = items.length - 1;
      }

      if (nextIndex !== currentIndex) {
        event.preventDefault();
        items[nextIndex].focus();
      }
    });
  });

  var languageNames = {
    bash: "Bash",
    c: "C",
    cpp: "C++",
    csharp: "C#",
    css: "CSS",
    diff: "Diff",
    go: "Go",
    html: "HTML",
    java: "Java",
    javascript: "JavaScript",
    js: "JavaScript",
    json: "JSON",
    jsx: "JSX",
    markdown: "Markdown",
    md: "Markdown",
    php: "PHP",
    powershell: "PowerShell",
    ps1: "PowerShell",
    python: "Python",
    py: "Python",
    ruby: "Ruby",
    rust: "Rust",
    rs: "Rust",
    scss: "SCSS",
    shell: "Shell",
    sh: "Shell",
    sql: "SQL",
    swift: "Swift",
    toml: "TOML",
    ts: "TypeScript",
    tsx: "TSX",
    typescript: "TypeScript",
    xml: "XML",
    yaml: "YAML",
    yml: "YAML"
  };

  function icon(name) {
    var attrs = 'viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"';
    var paths = {
      chevron: '<path d="m6 9 6 6 6-6"></path>',
      copy: '<rect x="9" y="9" width="11" height="11" rx="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>',
      check: '<path d="M20 6 9 17l-5-5"></path>',
      download: '<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><path d="M7 10l5 5 5-5"></path><path d="M12 15V3"></path>',
      expand: '<path d="M15 3h6v6"></path><path d="M9 21H3v-6"></path><path d="m21 3-7 7"></path><path d="m3 21 7-7"></path>',
      hash: '<path d="M4 9h16"></path><path d="M4 15h16"></path><path d="M10 3 8 21"></path><path d="m16 3-2 18"></path>',
      wrap: '<path d="M3 7h13a4 4 0 0 1 0 8H7"></path><path d="m10 12-3 3 3 3"></path>'
    };
    return '<svg ' + attrs + ">" + paths[name] + "</svg>";
  }

  function detectLanguage(pre, code) {
    var source = [
      code && code.getAttribute("class"),
      pre && pre.getAttribute("class"),
      code && code.getAttribute("data-lang"),
      pre && pre.getAttribute("data-lang")
    ].filter(Boolean).join(" ");
    var match = source.match(/(?:language|lang)-([a-z0-9_+#.-]+)/i);
    var language = match ? match[1].toLowerCase() : "";
    return language.replace(/^text$/, "");
  }

  function countCodeLines(text) {
    if (!text) {
      return 1;
    }
    return text.replace(/\n$/, "").split("\n").length;
  }

  function makeLineNumbers(count) {
    var lines = [];
    for (var i = 1; i <= count; i += 1) {
      lines.push(String(i));
    }
    return lines.join("\n");
  }

  function copyToClipboard(text) {
    if (navigator.clipboard && window.isSecureContext) {
      return navigator.clipboard.writeText(text);
    }
    return new Promise(function (resolve, reject) {
      var textarea = document.createElement("textarea");
      textarea.value = text;
      textarea.setAttribute("readonly", "");
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      document.body.appendChild(textarea);
      textarea.select();
      try {
        document.execCommand("copy") ? resolve() : reject(new Error("Copy failed"));
      } catch (error) {
        reject(error);
      } finally {
        document.body.removeChild(textarea);
      }
    });
  }

  function downloadCode(text, language) {
    var extension = language || "txt";
    if (extension === "javascript") extension = "js";
    if (extension === "typescript") extension = "ts";
    if (extension === "markdown") extension = "md";
    var blob = new Blob([text], { type: "text/plain;charset=utf-8" });
    var url = URL.createObjectURL(blob);
    var link = document.createElement("a");
    link.href = url;
    link.download = "code." + extension;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }

  function setPressed(button, pressed) {
    button.setAttribute("aria-pressed", String(pressed));
  }

  function saveMermaidSource() {
    document.querySelectorAll('pre.mermaid').forEach(function (pre) {
      if (pre.hasAttribute('data-source')) return;
      var text = pre.textContent;
      if (text.trim()) pre.setAttribute('data-source', text);
    });
  }
  function initMermaidCopy() {
    document.querySelectorAll(".tm-post__body pre.mermaid").forEach(function (pre) {
      if (pre.closest(".tm-mermaid") || !pre.hasAttribute("data-source")) return;
      var wrapper = document.createElement("div");
      wrapper.className = "tm-mermaid";
      var button = document.createElement("button");
      button.className = "tm-mermaid__copy";
      button.type = "button";
      button.setAttribute("aria-label", "Copy diagram source");
      button.title = "Copy diagram source";
      button.innerHTML = icon("copy");
      var source = pre.getAttribute("data-source") || "";
      button.addEventListener("click", function () {
        copyToClipboard(source.replace(/\s+$/, "")).then(function () {
          button.dataset.copied = "true";
          button.innerHTML = icon("check");
          window.setTimeout(function () {
            button.dataset.copied = "false";
            button.innerHTML = icon("copy");
          }, 1400);
        });
      });
      pre.parentNode.insertBefore(wrapper, pre);
      wrapper.appendChild(pre);
      wrapper.appendChild(button);
    });
  }
  function initCodeBlocks() {
    var postBodies = Array.prototype.slice.call(document.querySelectorAll(".tm-post__body"));
    if (!postBodies.length) {
      return;
    }

    var preBlocks = [];
    postBodies.forEach(function (postBody) {
      preBlocks = preBlocks.concat(Array.prototype.slice.call(postBody.querySelectorAll("pre")));
    });

    preBlocks.forEach(function (pre) {
      if (pre.closest(".tm-code-block") || pre.classList.contains("tm-code-lines") || pre.classList.contains("mermaid")) {
        return;
      }

      var code = pre.querySelector("code");
      if (!code) {
        code = document.createElement("code");
        code.textContent = pre.textContent || "";
        pre.textContent = "";
        pre.appendChild(code);
      }

      var language = detectLanguage(pre, code);
      var languageLabel = languageNames[language] || (language ? language.toUpperCase() : "Code");
      var text = code.textContent || pre.textContent || "";
      var block = document.createElement("div");
      var header = document.createElement("div");
      var title = document.createElement("button");
      var titleText = document.createElement("span");
      var actions = document.createElement("div");
      var body = document.createElement("div");
      var lines = document.createElement("pre");
      var scroller = document.createElement("div");

      block.className = "tm-code-block highlight";
      block.dataset.language = language || "code";

      header.className = "tm-code-header";
      title.className = "tm-code-title";
      title.type = "button";
      title.setAttribute("aria-label", "Toggle code block");
      title.innerHTML = icon("chevron");
      titleText.textContent = languageLabel;
      title.appendChild(titleText);

      actions.className = "tm-code-actions";

      var lineButton = document.createElement("button");
      lineButton.className = "tm-code-action";
      lineButton.type = "button";
      lineButton.title = "Toggle line numbers";
      lineButton.setAttribute("aria-label", "Toggle line numbers");
      lineButton.innerHTML = icon("hash");
      setPressed(lineButton, true);

      var wrapButton = document.createElement("button");
      wrapButton.className = "tm-code-action";
      wrapButton.type = "button";
      wrapButton.title = "Toggle line wrap";
      wrapButton.setAttribute("aria-label", "Toggle line wrap");
      wrapButton.innerHTML = icon("wrap");
      setPressed(wrapButton, false);

      var copyButton = document.createElement("button");
      copyButton.className = "tm-code-action";
      copyButton.type = "button";
      copyButton.title = "Copy code";
      copyButton.setAttribute("aria-label", "Copy code");
      copyButton.innerHTML = icon("copy");

      var downloadButton = document.createElement("button");
      downloadButton.className = "tm-code-action";
      downloadButton.type = "button";
      downloadButton.title = "Download code";
      downloadButton.setAttribute("aria-label", "Download code");
      downloadButton.innerHTML = icon("download");

      var fullscreenButton = document.createElement("button");
      fullscreenButton.className = "tm-code-action";
      fullscreenButton.type = "button";
      fullscreenButton.title = "Fullscreen";
      fullscreenButton.setAttribute("aria-label", "Fullscreen");
      fullscreenButton.innerHTML = icon("expand");

      actions.appendChild(lineButton);
      actions.appendChild(wrapButton);
      actions.appendChild(copyButton);
      actions.appendChild(downloadButton);
      actions.appendChild(fullscreenButton);

      header.appendChild(title);
      header.appendChild(actions);

      body.className = "tm-code-body";
      lines.className = "tm-code-lines";
      lines.setAttribute("aria-hidden", "true");
      lines.textContent = makeLineNumbers(countCodeLines(text));
      scroller.className = "tm-code-scroller";

      pre.parentNode.insertBefore(block, pre);
      scroller.appendChild(pre);
      body.appendChild(lines);
      body.appendChild(scroller);
      block.appendChild(header);
      block.appendChild(body);

      title.addEventListener("click", function () {
        block.classList.toggle("is-collapsed");
      });

      lineButton.addEventListener("click", function () {
        var hidden = block.classList.toggle("no-lines");
        setPressed(lineButton, !hidden);
      });

      wrapButton.addEventListener("click", function () {
        var wrapped = block.classList.toggle("is-wrapped");
        setPressed(wrapButton, wrapped);
      });

      copyButton.addEventListener("click", function () {
        copyToClipboard(text.replace(/\s+$/, "")).then(function () {
          copyButton.dataset.copied = "true";
          copyButton.innerHTML = icon("check");
          window.setTimeout(function () {
            copyButton.dataset.copied = "false";
            copyButton.innerHTML = icon("copy");
          }, 1400);
        });
      });

      downloadButton.addEventListener("click", function () {
        downloadCode(text, language);
      });

      fullscreenButton.addEventListener("click", function () {
        var fullscreen = block.classList.toggle("is-fullscreen");
        document.body.style.overflow = fullscreen ? "hidden" : "";
        fullscreenButton.classList.toggle("is-active", fullscreen);
      });
    });
  }

  function initImagePreview() {
    if (!window.jQuery || !window.jQuery.fancybox) {
      return;
    }

    var $ = window.jQuery;
    $(".tm-post__body").each(function () {
      $(this).find("img").each(function () {
        var image = this;
        var link = image.parentElement && image.parentElement.matches("a")
          ? image.parentElement
          : null;

        if (!link) {
          link = document.createElement("a");
          link.href = image.currentSrc || image.src;
          image.parentNode.insertBefore(link, image);
          link.appendChild(image);
        } else if (!link.getAttribute("href")) {
          link.href = image.currentSrc || image.src;
        }

        link.classList.add("fancybox");
        link.setAttribute("data-fancybox", "post-gallery");
        link.setAttribute(
          "data-caption",
          image.getAttribute("title") || image.getAttribute("alt") || ""
        );
      });
    });

    $(".fancybox").fancybox({
      selector: ".fancybox",
      protect: true
    });
  }

  document.addEventListener("keydown", function (event) {
    if (event.key !== "Escape") {
      return;
    }
    var fullscreenBlock = document.querySelector(".tm-code-block.is-fullscreen");
    if (fullscreenBlock) {
      fullscreenBlock.classList.remove("is-fullscreen");
      var button = fullscreenBlock.querySelector(".tm-code-action.is-active");
      if (button) {
        button.classList.remove("is-active");
      }
      document.body.style.overflow = "";
    }
  });

  function initTypeIt() {
    var target = document.querySelector("[data-fi-typeit]");
    if (!target) {
      return;
    }

    var text = target.getAttribute("data-fi-typeit") || "";
    if (!text) {
      return;
    }

    var chars = Array.from(text);
    var index = 0;
    var deleting = false;
    var wait = 0;
    target.textContent = "";

    function tick() {
      if (wait > 0) {
        wait -= 1;
      } else if (deleting) {
        index -= 1;
        if (index <= 0) {
          index = 0;
          deleting = false;
          wait = 8;
        }
      } else {
        index += 1;
        if (index >= chars.length) {
          index = chars.length;
          deleting = true;
          wait = 26;
        }
      }

      target.textContent = chars.slice(0, index).join("");
      window.setTimeout(tick, deleting ? 54 : 96);
    }

    window.setTimeout(tick, 260);
  }

  initTypeIt();
  initImagePreview();
  saveMermaidSource();
  initCodeBlocks();
  initMermaidCopy();
})();
