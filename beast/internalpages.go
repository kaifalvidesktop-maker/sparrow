package main

// ---------------------------------------------------
// SPARROW INTERNAL PAGES (sparrow://settings, sparrow://bookmarks, etc.)
// ---------------------------------------------------

const settingsPageHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>SPARROW Settings</title>
<style>
  * { box-sizing: border-box; }

  body {
    margin: 0;
    background: #f5f7fa;
    color: #202124;
    font-family: 'Segoe UI', sans-serif;
    display: flex;
    min-height: 100vh;
  }

  .sidebar {
    width: 220px;
    background: #ffffff;
    border-right: 1px solid #e1e4e8;
    padding: 24px 0;
    flex-shrink: 0;
  }

  .sidebar .brand {
    font-size: 20px;
    font-weight: 700;
    letter-spacing: 2px;
    padding: 0 24px 24px 24px;
    color: #202124;
  }

  .sidebar .item {
    padding: 12px 24px;
    color: #5f6368;
    cursor: pointer;
    font-size: 13px;
    border-left: 3px solid transparent;
    transition: 0.15s;
  }

  .sidebar .item:hover {
    background: #f1f3f4;
    color: #202124;
  }

  .sidebar .item.active {
    background: #e8f0fe;
    color: #1967d2;
    border-left: 3px solid #4d90fe;
  }

  .content {
    flex: 1;
    padding: 40px 60px;
    max-width: 720px;
  }

  h1 {
    font-size: 26px;
    margin-bottom: 30px;
    font-weight: 700;
    color: #202124;
  }

  .card {
    background: #ffffff;
    border: 1px solid #e1e4e8;
    border-radius: 14px;
    padding: 22px 26px;
    margin-bottom: 20px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.06);
    transition: transform 0.15s, box-shadow 0.15s;
  }

  .card:hover {
    box-shadow: 0 5px 16px rgba(0,0,0,0.09);
    transform: translateY(-1px);
  }

  .card-title {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 4px;
    color: #202124;
  }

  .card-desc {
    font-size: 12px;
    color: #6b7280;
    margin-bottom: 16px;
  }

  .row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 0;
    border-bottom: 1px solid #edf0f2;
  }

  .row:last-child { border-bottom: none; }

  .row-label {
    font-size: 13px;
    color: #3c4043;
  }

  .switch {
    position: relative;
    width: 42px;
    height: 24px;
    background: #dadce0;
    border-radius: 20px;
    cursor: pointer;
    transition: 0.2s;
  }

  .switch.on {
    background: #4d90fe;
    box-shadow: 0 0 10px rgba(77,144,254,0.25);
  }

  .switch .knob {
    position: absolute;
    top: 2px;
    left: 2px;
    width: 20px;
    height: 20px;
    background: #ffffff;
    border-radius: 50%;
    transition: 0.2s;
    box-shadow: 0 2px 4px rgba(0,0,0,0.18);
  }

  .switch.on .knob {
    left: 20px;
  }

  select {
    background: #ffffff;
    color: #202124;
    border: 1px solid #dadce0;
    border-radius: 8px;
    padding: 8px 12px;
    font-size: 13px;
    outline: none;
    cursor: pointer;
  }

  input[type=text] {
    background: #ffffff;
    color: #202124;
    border: 1px solid #dadce0;
    border-radius: 8px;
    padding: 8px 12px;
    font-size: 13px;
    outline: none;
    width: 260px;
  }

  input[type=text]:focus,
  select:focus {
    border-color: #4d90fe;
  }

  .btn {
    background: #ffffff;
    color: #3c4043;
    border: 1px solid #dadce0;
    border-radius: 8px;
    padding: 9px 18px;
    font-size: 12px;
    cursor: pointer;
    transition: 0.15s;
  }

  .btn:hover {
    background: #f1f3f4;
    color: #202124;
  }

  .btn.danger {
    border-color: #f1b5b5;
    color: #d93025;
  }

  .btn.danger:hover {
    background: #fce8e6;
  }

  .stat {
    font-size: 12px;
    color: #188038;
    margin-top: 6px;
  }
</style>
</head>
<body>

  <div class="sidebar">
    <div class="brand">SPARROW</div>
    <div class="item active">General</div>
    <div class="item" onclick="window.realNavigate('sparrow://bookmarks')">Bookmarks</div>
    <div class="item" onclick="window.realNavigate('sparrow://downloads')">Downloads</div>
    <div class="item" onclick="window.realNavigate('sparrow://history')">History</div>
  </div>

  <div class="content">
    <h1>Settings</h1>

    <div class="card">
      <div class="card-title">Appearance</div>
      <div class="card-desc">Choose how SPARROW looks</div>
      <div class="row">
        <div class="row-label">Dark Theme</div>
        <div class="switch on" id="themeSwitch" onclick="toggleTheme()"><div class="knob"></div></div>
      </div>
    </div>

    <div class="card">
      <div class="card-title">Search Engine</div>
      <div class="card-desc">Default engine used when typing a search in the address bar</div>
      <div class="row">
        <div class="row-label">Engine</div>
        <select id="engineSelect" onchange="changeEngine()">
          <option value="google">Google</option>
          <option value="bing">Bing</option>
          <option value="duckduckgo">DuckDuckGo</option>
        </select>
      </div>
    </div>

    <div class="card">
      <div class="card-title">Homepage</div>
      <div class="card-desc">Page shown when you open a new tab</div>
      <div class="row">
        <input type="text" id="homepageInput" placeholder="sparrow://home">
        <button class="btn" onclick="saveHomepage()">Save</button>
      </div>
    </div>

    <div class="card">
      <div class="card-title">Privacy & Security Shield</div>
      <div class="card-desc">Control what SPARROW blocks automatically</div>

      <div class="row">
        <div class="row-label">Block Ads</div>
        <div class="switch on" id="adSwitch" onclick="toggleShield('ad')"><div class="knob"></div></div>
      </div>

      <div class="row">
        <div class="row-label">Block Trackers</div>
        <div class="switch on" id="trackerSwitch" onclick="toggleShield('tracker')"><div class="knob"></div></div>
      </div>

      <div class="row">
        <div class="row-label">HTTPS Only Mode</div>
        <div class="switch on" id="httpsSwitch" onclick="toggleShield('https')"><div class="knob"></div></div>
      </div>

      <div class="stat" id="blockedStat">0 threats blocked this session</div>
    </div>

    <div class="card">
      <div class="card-title">Content</div>
      <div class="card-desc">Control what pages are allowed to run</div>

      <div class="row">
        <div class="row-label">JavaScript</div>
        <div class="switch on" id="jsSwitch" onclick="toggleContent('js')"><div class="knob"></div></div>
      </div>

      <div class="row">
        <div class="row-label">Images</div>
        <div class="switch on" id="imgSwitch" onclick="toggleContent('img')"><div class="knob"></div></div>
      </div>
    </div>

    <div class="card">
      <div class="card-title">Data & Privacy</div>
      <div class="card-desc">SPARROW keeps nothing on disk. History and cache live only in RAM.</div>

      <div class="row">
        <div class="row-label">Clear all browsing history</div>
        <button class="btn danger" onclick="clearAllHistory()">Clear Now</button>
      </div>

      <div class="row">
        <div class="row-label">Reset all settings to default</div>
        <button class="btn danger" onclick="resetAll()">Reset</button>
      </div>
    </div>

  </div>

<script>
  async function loadSettings() {
    const s = await window.getSettings();

    document.getElementById('themeSwitch').className =
      s.theme === 'dark' ? 'switch on' : 'switch';

    document.getElementById('engineSelect').value = s.searchEngine;
    document.getElementById('homepageInput').value = s.homepage;

    document.getElementById('jsSwitch').className =
      s.jsEnabled ? 'switch on' : 'switch';

    document.getElementById('imgSwitch').className =
      s.imagesEnabled ? 'switch on' : 'switch';

    const stats = await window.getShieldStats();

    document.getElementById('adSwitch').className =
      stats.adBlock ? 'switch on' : 'switch';

    document.getElementById('trackerSwitch').className =
      stats.trackerBlock ? 'switch on' : 'switch';

    document.getElementById('httpsSwitch').className =
      stats.httpsOnly ? 'switch on' : 'switch';

    document.getElementById('blockedStat').innerText =
      stats.blockedCount + ' threats blocked this session';
  }

  async function toggleTheme() {
    const el = document.getElementById('themeSwitch');
    const isOn = el.classList.contains('on');

    const newTheme = isOn ? 'light' : 'dark';

    await window.setTheme(newTheme);

    el.className = !isOn ? 'switch on' : 'switch';
  }

  async function changeEngine() {
    const val = document.getElementById('engineSelect').value;
    await window.setSearchEngine(val);
  }

  async function saveHomepage() {
    const val = document.getElementById('homepageInput').value;

    await window.setHomepage(val);

    alert('Homepage saved');
  }

  async function toggleShield(kind) {
    let id, fn;

    if (kind === 'ad') {
      id = 'adSwitch';
      fn = window.toggleAdBlock;
    }

    if (kind === 'tracker') {
      id = 'trackerSwitch';
      fn = window.toggleTrackerBlock;
    }

    if (kind === 'https') {
      id = 'httpsSwitch';
      fn = window.toggleHTTPSOnly;
    }

    const result = await fn();

    document.getElementById(id).className =
      result ? 'switch on' : 'switch';
  }

  async function toggleContent(kind) {
    let id, fn;

    if (kind === 'js') {
      id = 'jsSwitch';
      fn = window.toggleJS;
    }

    if (kind === 'img') {
      id = 'imgSwitch';
      fn = window.toggleImages;
    }

    const result = await fn();

    document.getElementById(id).className =
      result ? 'switch on' : 'switch';
  }

  async function clearAllHistory() {
    if (!confirm('Clear all browsing history?')) return;

    await window.clearHistory('all');

    alert('History cleared');
  }

  async function resetAll() {
    if (!confirm('Reset all settings to default?')) return;

    await window.resetSettings();

    loadSettings();
  }

  loadSettings();
</script>

</body>
</html>
`

const bookmarksPageHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>SPARROW Bookmarks</title>
<style>
  * { box-sizing: border-box; }

  body {
    margin: 0;
    background: #f5f7fa;
    color: #202124;
    font-family: 'Segoe UI', sans-serif;
    padding: 40px 60px;
  }

  h1 {
    font-size: 26px;
    margin-bottom: 6px;
    color: #202124;
  }

  .sub {
    color: #6b7280;
    font-size: 13px;
    margin-bottom: 26px;
  }

  .search {
    width: 320px;
    padding: 10px 16px;
    border-radius: 20px;
    border: 1px solid #dadce0;
    background: #ffffff;
    color: #202124;
    outline: none;
    margin-bottom: 24px;
    font-size: 13px;
  }

  .search:focus {
    border-color: #4d90fe;
  }

  .bm-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #ffffff;
    border: 1px solid #e1e4e8;
    border-radius: 12px;
    padding: 14px 18px;
    margin-bottom: 10px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.06);
    transition: 0.15s;
  }

  .bm-item:hover {
    transform: translateY(-1px);
    box-shadow: 0 5px 16px rgba(0,0,0,0.09);
    border-color: #c7cdd4;
  }

  .bm-info {
    cursor: pointer;
  }

  .bm-title {
    font-size: 14px;
    color: #202124;
    margin-bottom: 3px;
  }

  .bm-url {
    font-size: 11px;
    color: #6b7280;
  }

  .bm-remove {
    color: #5f6368;
    cursor: pointer;
    font-size: 12px;
    padding: 6px 10px;
    border-radius: 6px;
  }

  .bm-remove:hover {
    color: #d93025;
    background: #fce8e6;
  }

  .empty {
    color: #6b7280;
    font-size: 13px;
    margin-top: 40px;
    text-align: center;
  }

  .action-btn {
    background: #ffffff;
    color: #3c4043;
    border: 1px solid #dadce0;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 12px;
  }

  .action-btn:hover {
    background: #f1f3f4;
  }
</style>
</head>
<body>

<h1>Bookmarks</h1>

<div class="sub" id="countText">Loading...</div>

<input
  class="search"
  placeholder="Search bookmarks..."
  oninput="filterBookmarks(this.value)"
>

<div style="display:flex; gap:10px; margin-bottom:20px;">
  <button class="action-btn" onclick="exportBm()">Export Bookmarks</button>

  <button
    class="action-btn"
    onclick="document.getElementById('importFile').click()"
  >
    Import Bookmarks
  </button>

  <input
    type="file"
    id="importFile"
    accept=".json"
    style="display:none;"
    onchange="importBm(event)"
  >
</div>

<div id="bmList"></div>

<script>
  let allBookmarks = [];

  async function exportBm() {
    const json = await window.exportBookmarks();

    const blob = new Blob(
      [json],
      { type: 'application/json' }
    );

    const url = URL.createObjectURL(blob);

    const a = document.createElement('a');
    a.href = url;
    a.download = 'sparrow-bookmarks.json';
    a.click();

    URL.revokeObjectURL(url);
  }

  function importBm(event) {
    const file = event.target.files[0];

    if (!file) return;

    const reader = new FileReader();

    reader.onload = async function(e) {
      const count =
        await window.importBookmarks(e.target.result);

      alert(count + ' bookmark(s) imported');

      loadBookmarks();
    };

    reader.readAsText(file);
  }

  async function loadBookmarks() {
    allBookmarks =
      await window.getBookmarks() || [];

    document.getElementById('countText').innerText =
      allBookmarks.length + ' saved pages';

    renderList(allBookmarks);
  }

  function renderList(list) {
    const container =
      document.getElementById('bmList');

    if (list.length === 0) {
      container.innerHTML =
        '<div class="empty">No bookmarks yet. Click the star icon on any page to save it here.</div>';

      return;
    }

    container.innerHTML = list.map(function(b) {
      return '<div class="bm-item">' +
        '<div class="bm-info" onclick="window.realNavigate(\'' + b.URL + '\')">' +
          '<div class="bm-title">' + b.Title + '</div>' +
          '<div class="bm-url">' + b.URL + '</div>' +
        '</div>' +
        '<div class="bm-remove" onclick="removeBm(' + b.ID + ')">Remove</div>' +
      '</div>';
    }).join('');
  }

  async function filterBookmarks(keyword) {
    if (keyword.trim() === '') {
      renderList(allBookmarks);
      return;
    }

    const results =
      await window.searchBookmarks(keyword);

    renderList(results || []);
  }

  async function removeBm(id) {
    await window.removeBookmark(id);
    loadBookmarks();
  }

  loadBookmarks();
</script>

</body>
</html>
`

const downloadsPageHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>SPARROW Downloads</title>
<style>
  * { box-sizing: border-box; }

  body {
    margin: 0;
    background: #f5f7fa;
    color: #202124;
    font-family: 'Segoe UI', sans-serif;
    padding: 40px 60px;
  }

  h1 {
    font-size: 26px;
    margin-bottom: 24px;
    color: #202124;
  }

  .dl-item {
    background: #ffffff;
    border: 1px solid #e1e4e8;
    border-radius: 12px;
    padding: 16px 20px;
    margin-bottom: 10px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  }

  .dl-name {
    font-size: 14px;
    color: #202124;
    margin-bottom: 6px;
  }

  .dl-meta {
    font-size: 11px;
    color: #6b7280;
    margin-bottom: 10px;
  }

  .bar-bg {
    width: 100%;
    height: 6px;
    background: #e8eaed;
    border-radius: 4px;
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    background: linear-gradient(90deg, #4d90fe, #6bb0ff);
    border-radius: 4px;
    transition: width 0.3s;
  }

  .status-completed {
    color: #188038;
  }

  .status-failed {
    color: #d93025;
  }

  .status-cancelled {
    color: #6b7280;
  }

  .empty {
    color: #6b7280;
    font-size: 13px;
    margin-top: 40px;
    text-align: center;
  }
</style>
</head>
<body>

<h1>Downloads</h1>

<div id="dlList"></div>

<script>
  async function loadDownloads() {
    const items =
      await window.getDownloads() || [];

    const container =
      document.getElementById('dlList');

    if (items.length === 0) {
      container.innerHTML =
        '<div class="empty">No downloads yet.</div>';

      return;
    }

    container.innerHTML = items.map(function(d) {
      let statusClass =
        'status-' + d.Status;

      return '<div class="dl-item">' +
        '<div class="dl-name">' +
          d.FileName +
        '</div>' +

        '<div class="dl-meta ' +
          statusClass +
        '">' +
          d.Status +
          ' — ' +
          d.Progress +
          '%' +
        '</div>' +

        '<div class="bar-bg">' +
          '<div class="bar-fill" style="width:' +
            d.Progress +
            '%"></div>' +
        '</div>' +

      '</div>';
    }).join('');
  }

  loadDownloads();

  setInterval(
    loadDownloads,
    1500
  );
</script>

</body>
</html>
`

const historyPageHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>SPARROW History</title>
<style>
  * { box-sizing: border-box; }

  body {
    margin: 0;
    background: #f5f7fa;
    color: #202124;
    font-family: 'Segoe UI', sans-serif;
    padding: 40px 60px;
  }

  h1 {
    font-size: 26px;
    margin-bottom: 6px;
    color: #202124;
  }

  .sub {
    color: #6b7280;
    font-size: 13px;
    margin-bottom: 24px;
  }

  .btnbar {
    display: flex;
    gap: 10px;
    margin-bottom: 24px;
  }

  .btn {
    background: #ffffff;
    color: #3c4043;
    border: 1px solid #dadce0;
    border-radius: 8px;
    padding: 8px 16px;
    font-size: 12px;
    cursor: pointer;
  }

  .btn:hover {
    background: #f1f3f4;
    color: #202124;
  }

  .hist-item {
    display: flex;
    justify-content: space-between;
    padding: 12px 16px;
    background: #ffffff;
    border: 1px solid #e1e4e8;
    border-radius: 10px;
    margin-bottom: 8px;
    cursor: pointer;
    font-size: 13px;
    transition: 0.15s;
  }

  .hist-item:hover {
    background: #f8f9fa;
    box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  }

  .hist-url {
    color: #3c4043;
  }

  .hist-time {
    color: #6b7280;
    font-size: 11px;
  }

  .empty {
    color: #6b7280;
    font-size: 13px;
    margin-top: 40px;
    text-align: center;
  }
</style>
</head>
<body>

<h1>History</h1>

<div class="sub">
  Stored only in RAM — cleared automatically when SPARROW closes
</div>

<div class="btnbar">
  <button class="btn" onclick="clearScope('1hour')">
    Clear Last Hour
  </button>

  <button class="btn" onclick="clearScope('today')">
    Clear Today
  </button>

  <button class="btn" onclick="clearScope('all')">
    Clear All
  </button>
</div>

<div id="histList"></div>

<script>
  async function loadHistory() {
    const items =
      await window.getRecentHistory(200) || [];

    const container =
      document.getElementById('histList');

    if (items.length === 0) {
      container.innerHTML =
        '<div class="empty">No history yet.</div>';

      return;
    }

    container.innerHTML =
      items.slice().reverse().map(function(h) {

        return '<div class="hist-item" onclick="window.realNavigate(\'' +
          h.URL +
          '\')">' +

          '<span class="hist-url">' +
            h.URL +
          '</span>' +

          '<span class="hist-time">' +
            h.Time +
          '</span>' +

        '</div>';

      }).join('');
  }

  async function clearScope(scope) {
    await window.clearHistory(scope);
    loadHistory();
  }

  loadHistory();
</script>

</body>
</html>
`