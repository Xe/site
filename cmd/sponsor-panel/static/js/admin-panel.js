(() => {
  const STORAGE_KEY = "sponsor-panel-tab";
  const tabs = Array.from(document.querySelectorAll("nav.tabs [role=tab]"));
  const panels = Array.from(document.querySelectorAll(".tabpanel"));
  const validTabs = tabs.map((t) => t.dataset.tab);

  function activate(name, persist = true) {
    if (!validTabs.includes(name)) return;
    tabs.forEach((t) => t.setAttribute("aria-selected", t.dataset.tab === name ? "true" : "false"));
    panels.forEach((p) => p.classList.toggle("active", p.id === `panel-${name}`));
    if (persist) {
      try {
        localStorage.setItem(STORAGE_KEY, name);
      } catch (e) {}
      const url = new URL(window.location);
      url.hash = name;
      history.replaceState(null, "", url);
    }
  }

  tabs.forEach((tab) => {
    tab.addEventListener("click", () => activate(tab.dataset.tab));
  });

  const tablist = document.querySelector("nav.tabs");
  if (tablist) {
    tablist.addEventListener("keydown", (e) => {
      const idx = tabs.findIndex((t) => t === document.activeElement);
      if (idx < 0) return;
      if (e.key === "ArrowDown" || e.key === "ArrowRight") {
        e.preventDefault();
        const next = tabs[(idx + 1) % tabs.length];
        next.focus();
        activate(next.dataset.tab);
      } else if (e.key === "ArrowUp" || e.key === "ArrowLeft") {
        e.preventDefault();
        const prev = tabs[(idx - 1 + tabs.length) % tabs.length];
        prev.focus();
        activate(prev.dataset.tab);
      }
    });
  }

  const hashTab = window.location.hash.replace("#", "");
  let initial = validTabs[0] || "sponsorship";
  if (validTabs.includes(hashTab)) {
    initial = hashTab;
  } else {
    try {
      const saved = localStorage.getItem(STORAGE_KEY);
      if (validTabs.includes(saved)) initial = saved;
    } catch (e) {}
  }
  activate(initial, false);

  // ---------- Toast ----------
  const toast = document.getElementById("admin-toast");
  const toastMsg = document.getElementById("admin-toast-msg");
  let toastTimer;
  window.adminToast = function (msg) {
    if (!toast) return;
    toastMsg.textContent = msg;
    toast.classList.add("show");
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => toast.classList.remove("show"), 2200);
  };

  // ---------- Logo drop zone ----------
  const drop = document.getElementById("logo-drop");
  const fileInput = document.getElementById("logo-file");
  const dropName = document.getElementById("logo-drop-name");
  const logoReset = document.getElementById("logo-reset");

  if (drop && fileInput && dropName) {
    const defaultText = dropName.textContent;

    ["dragenter", "dragover"].forEach((ev) =>
      drop.addEventListener(ev, (e) => {
        e.preventDefault();
        drop.classList.add("drag");
      }),
    );
    ["dragleave", "drop"].forEach((ev) =>
      drop.addEventListener(ev, (e) => {
        e.preventDefault();
        drop.classList.remove("drag");
      }),
    );
    drop.addEventListener("drop", (e) => {
      if (e.dataTransfer && e.dataTransfer.files.length) {
        fileInput.files = e.dataTransfer.files;
        fileInput.dispatchEvent(new Event("change"));
      }
    });
    fileInput.addEventListener("change", () => {
      if (fileInput.files && fileInput.files[0]) {
        dropName.textContent = fileInput.files[0].name;
      }
    });
    if (logoReset) {
      logoReset.addEventListener("click", () => {
        dropName.textContent = defaultText;
      });
    }
  }

  // Surface HTMX success responses as toasts
  document.body.addEventListener("htmx:afterSwap", (e) => {
    const target = e.detail && e.detail.target;
    if (!target) return;
    if (target.id === "invite-result" && target.querySelector(".alert-success")) {
      window.adminToast("Invitation sent");
    } else if (target.id === "logo-result" && target.querySelector(".alert-success")) {
      window.adminToast("Logo submitted");
    } else if (target.id === "thoth-result" && target.querySelector(".alert-success")) {
      window.adminToast("Token generated");
    }
  });
})();
