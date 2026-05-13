// LeetDrill Firefox/Zen options page.

const $ = (id) => document.getElementById(id);

function send(type, payload) {
  return ldx.runtime
    .sendMessage({ type, payload })
    .then((res) => res || { ok: false })
    .catch((err) => ({ ok: false, error: err.message || String(err) }));
}

async function load() {
  const res = await send("LEETDRILL_GET_CONFIG");
  if (res.ok) {
    $("backend").value = res.data.backendUrl || "https://abhiy.xyz/leetdrill";
  }
}

function setStatus(msg, cls) {
  const el = $("status");
  el.textContent = msg;
  el.className = "status " + (cls || "");
}

$("save").addEventListener("click", async () => {
  const res = await send("LEETDRILL_SAVE_CONFIG", { backendUrl: $("backend").value.trim() });
  setStatus(res.ok ? "saved" : "save failed", res.ok ? "ok" : "bad");
});

$("codePage").addEventListener("click", async () => {
  await send("LEETDRILL_SAVE_CONFIG", { backendUrl: $("backend").value.trim() });
  const res = await send("LEETDRILL_OPEN_CODE_PAGE");
  setStatus(res.ok ? "opened code page" : `open failed: ${res.error || "unknown error"}`, res.ok ? "ok" : "bad");
});

$("openApp").addEventListener("click", async () => {
  await send("LEETDRILL_SAVE_CONFIG", { backendUrl: $("backend").value.trim() });
  const res = await send("LEETDRILL_OPEN_APP");
  setStatus(res.ok ? "opened LeetDrill" : `open failed: ${res.error || "unknown error"}`, res.ok ? "ok" : "bad");
});

$("testConnection").addEventListener("click", async () => {
  await send("LEETDRILL_SAVE_CONFIG", { backendUrl: $("backend").value.trim() });
  const res = await send("LEETDRILL_TEST_CONNECTION");
  if (!res.ok) {
    setStatus(`test failed: ${res.error || "unknown error"}`, "bad");
    return;
  }
  const data = res.data || {};
  if (data.connected) {
    setStatus("connection works: Zen can reach abhiy.xyz with the saved code", "ok");
  } else if (data.permission === "blocked") {
    setStatus(`Zen blocked abhiy.xyz access: ${data.message || "fetch failed"}`, "bad");
  } else {
    setStatus(`connection failed: ${data.message || "code missing or rejected"}`, "bad");
  }
});

$("saveToken").addEventListener("click", async () => {
  const res = await send("LEETDRILL_SAVE_TOKEN", { token: $("manualToken").value.trim() });
  if (!res.ok) {
    setStatus(`manual connect failed: ${res.error || "unknown error"}`, "bad");
    return;
  }
  $("manualToken").value = "";
  const test = await send("LEETDRILL_TEST_CONNECTION");
  if (test.ok && test.data && test.data.connected) {
    setStatus("connected - manual code saved and verified", "ok");
  } else {
    const msg = test.data && test.data.message ? test.data.message : "run test connection";
    setStatus(`manual code saved; test failed: ${msg}`, "bad");
  }
});

load();
