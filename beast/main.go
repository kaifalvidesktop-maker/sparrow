package main

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	webview "github.com/webview/webview_go"
)

// ---------------------------------------------------
// HISTORY SYSTEM (RAM ONLY - NEVER TOUCHES DISK)
// ---------------------------------------------------

type HistoryEntry struct {
	URL   string
	Title string
	Time  time.Time
}

type RAMHistory struct {
	mu      sync.Mutex
	Entries []HistoryEntry
}

var history = &RAMHistory{}

type NavResult struct {
	Blocked bool   `json:"blocked"`
	URL     string `json:"url"`
}

type NavState struct {
	URL        string `json:"url"`
	CanBack    bool   `json:"canBack"`
	CanForward bool   `json:"canForward"`
	TabID      int    `json:"tabId"`
}

func dataURI(html string) string {
	return "data:text/html;charset=utf-8;base64," + base64.StdEncoding.EncodeToString([]byte(html))
}

func internalPageHTML(sparrowURL string) (string, bool) {
	switch sparrowURL {
	case "sparrow://home":
		return homePageHTML, true
	case "sparrow://settings":
		return settingsPageHTML, true
	case "sparrow://bookmarks":
		return bookmarksPageHTML, true
	case "sparrow://downloads":
		return downloadsPageHTML, true
	case "sparrow://history":
		return historyPageHTML, true
	case "sparrow://shortcuts":
		return shortcutsPageHTML, true
	case "sparrow://about":
		return aboutPageHTML, true
	case "sparrow://site-settings":
		return siteSettingsPageHTML, true
	case "sparrow://cookies":
		return cookiesPageHTML, true
	case "sparrow://welcome":
		return welcomePageHTML, true
	case "sparrow://autofill":
		return autofillPageHTML, true
	case "sparrow://feedback":
		return feedbackPageHTML, true
	case "sparrow://passwords":
		return passwordsPageHTML, true
	case "sparrow://updates":
		return updatesPageHTML, true
	case "sparrow://backup":
		return backupPageHTML, true
	}
	return "", false
}

func navigateTo(w webview.WebView, target string) {
	if html, ok := internalPageHTML(target); ok {
		w.Dispatch(func() { w.Navigate(dataURI(html)) })
		return
	}
	w.Dispatch(func() { w.Navigate(target) })
}

func (h *RAMHistory) Add(url string, title string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Entries = append(h.Entries, HistoryEntry{
		URL:   url,
		Title: title,
		Time:  time.Now(),
	})
}

func (h *RAMHistory) DeleteSince(cutoff time.Time) {
	h.mu.Lock()
	defer h.mu.Unlock()

	newEntries := []HistoryEntry{}
	for _, entry := range h.Entries {
		if entry.Time.Before(cutoff) {
			newEntries = append(newEntries, entry)
		}
	}
	h.Entries = newEntries
}

func (h *RAMHistory) DeleteAll() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Entries = []HistoryEntry{}
}

func (h *RAMHistory) GetAll() []HistoryEntry {
	h.mu.Lock()
	defer h.mu.Unlock()

	entries := make([]HistoryEntry, len(h.Entries))
	copy(entries, h.Entries)
	return entries
}

func (h *RAMHistory) GetRecent(limit int) []HistoryEntry {
	h.mu.Lock()
	defer h.mu.Unlock()

	if limit <= 0 {
		return []HistoryEntry{}
	}

	start := 0
	if len(h.Entries) > limit {
		start = len(h.Entries) - limit
	}

	entries := make([]HistoryEntry, len(h.Entries)-start)
	copy(entries, h.Entries[start:])
	return entries
}

// ---------------------------------------------------
// SPARROW HOME PAGE
// ---------------------------------------------------

const homePageHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>SPARROW</title>
<style>
  * { box-sizing: border-box; }
  body {
    margin: 0;
    height: 100vh;
    background: linear-gradient(180deg, #0b0b0d 0%, #121214 100%);
    color: #e8e8e8;
    font-family: 'Segoe UI', sans-serif;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .logo {
    font-size: 46px;
    font-weight: 700;
    letter-spacing: 3px;
    margin-bottom: 6px;
    background: linear-gradient(90deg, #ffffff, #9b9b9b);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .tagline {
    color: #777;
    margin-bottom: 34px;
    font-size: 13px;
  }

  .search-bar {
    width: 560px;
    max-width: 82%;
    padding: 15px 22px;
    border-radius: 30px;
    border: 1px solid #2a2a2a;
    background: #17181a;
    color: white;
    font-size: 15px;
    outline: none;
    box-shadow: 0 4px 20px rgba(0,0,0,0.4);
  }

  .search-bar:focus { border-color: #4d90fe; }

  .top-sites {
    display: flex;
    gap: 22px;
    margin-top: 42px;
  }

  .site {
    width: 76px;
    height: 76px;
    background: #17181a;
    border: 1px solid #232323;
    border-radius: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 11px;
    color: #999;
    cursor: pointer;
    transition: 0.15s;
  }

  .site:hover {
    background: #202124;
    color: white;
    transform: translateY(-2px);
  }

  .history-controls {
    position: fixed;
    bottom: 20px;
    display: flex;
    gap: 10px;
  }

  .history-controls button {
    background: #17181a;
    color: #ccc;
    border: 1px solid #2a2a2a;
    padding: 8px 14px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 12px;
  }

  .history-controls button:hover { background: #202124; color: white; }
</style>
</head>
<body>

  <div class="logo">SPARROW</div>
  <div class="tagline">Fast. Private. Yours.</div>

  <input class="search-bar" id="searchBar" placeholder="Search Google or type a URL"
         onkeydown="if(event.key==='Enter'){ handleSearch() }">

  <div class="top-sites" id="topSitesContainer"></div>

  <div class="history-controls">
    <button onclick="go('sparrow://history')">History</button>
    <button onclick="go('sparrow://bookmarks')">Bookmarks</button>
    <button onclick="go('sparrow://settings')">Settings</button>
  </div>

<script>
  function handleSearch() {
    const input = document.getElementById('searchBar').value;
    if (input.trim() === '') return;
    go(input);
  }

  async function loadTopSites() {
    if (window.getTopSites) {
      const sites = await window.getTopSites(4);
      const container = document.getElementById('topSitesContainer');
      container.innerHTML = sites.map(function(s) {
        return '<div class="site" onclick="go(\'' + s.URL + '\')">' + s.Domain + '</div>';
      }).join('');
    }
  }
  loadTopSites();
</script>

</body>
</html>
`

// ---------------------------------------------------
// MAIN FUNCTION - CREATES BROWSER INSTANCE
// ---------------------------------------------------

func main() {
	// Initialize disk persistence engine
	InitPersistence()
	defer StopPersistence()

	var autofillEnabled = true
	var autofillEnabledMu sync.Mutex

	w := webview.New(true)
	defer w.Destroy()

	windowHandle := uintptr(w.Window())
	configureNativeWindow(windowHandle)

	w.SetTitle("SPARROW")
	w.SetSize(1200, 800, 0)

	// -------------------------------
	// GLOBAL SCRIPTS & INJECTIONS
	// -------------------------------

	w.Init("window.__SHORTCUTS__ = " + shortcutsJSON() + ";")
	w.Init(chromeInjectionJS)
	w.Init(autocompleteHTML())

	w.Bind("closeWindow", func() {
		closeNativeWindow(windowHandle)
	})

	w.Bind("minimizeWindow", func() {
		minimizeNativeWindow(windowHandle)
	})

	w.Bind("toggleMaximizeWindow", func() {
		toggleMaximizeNativeWindow(windowHandle)
	})

	w.Bind("toggleFullscreenWindow", func() {
		toggleFullscreenNativeWindow(windowHandle)
	})

	w.Bind("beginWindowDrag", func() {
		beginNativeWindowDrag(windowHandle)
	})

	// -------------------------------
	// HISTORY CONTROLS
	// -------------------------------

	w.Bind("clearHistory", func(scope string) {
		switch scope {
		case "1hour":
			history.DeleteSince(time.Now().Add(-1 * time.Hour))
		case "today":
			now := time.Now()
			startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			history.DeleteSince(startOfDay)
		case "all":
			history.DeleteAll()
		}
		fmt.Println("History cleared:", scope)
	})

	w.Bind("getRecentHistory", func(limit int) []HistoryEntry {
		return history.GetRecent(limit)
	})

	// -------------------------------
	// REAL NAVIGATION WITH SECURITY SHIELD
	// -------------------------------

	w.Bind("realNavigate", func(input string) NavResult {
		trimmed := strings.TrimSpace(input)
		if trimmed == "" {
			return NavResult{}
		}

		if _, ok := internalPageHTML(trimmed); ok {
			tab := tabManager.GetActiveTab()
			if tab == nil {
				tab = tabManager.NewTab(trimmed)
			}
			tabManager.RecordNavigation(tab.ID, trimmed)
			navigateTo(w, trimmed)
			return NavResult{URL: trimmed}
		}

		target := resolveInput(input)

		if shield.ShouldBlock(target) {
			domain := extractDomain(target)
			tab := tabManager.GetActiveTab()
			if tab == nil {
				tab = tabManager.NewTab(target)
			}
			blockedURL := "sparrow://blocked?" + domain
			tabManager.RecordNavigation(tab.ID, blockedURL)
			w.Dispatch(func() { w.Navigate(dataURI(buildBlockedPageHTML(domain))) })
			return NavResult{Blocked: true, URL: target}
		}

		target = shield.EnforceHTTPS(target)

		tab := tabManager.GetActiveTab()
		if tab == nil {
			tab = tabManager.NewTab(target)
		}
		tabManager.RecordNavigation(tab.ID, target)

		if !incognito.IsActive() {
			history.Add(target, target)
		}

		navigateTo(w, target)
		return NavResult{URL: target}
	})

	// -------------------------------
	// TAB CONTROLS
	// -------------------------------

	w.Bind("openNewTab", func() int {
		tab := tabManager.NewTab("")
		tabManager.RecordNavigation(tab.ID, "sparrow://home")
		navigateTo(w, "sparrow://home")
		return tab.ID
	})

	w.Bind("closeTab", func(id int) {
		var closedTitle, closedURL string
		for _, t := range tabManager.GetAllTabs() {
			if t.ID == id {
				closedTitle = t.Title
				closedURL = t.URL
				break
			}
		}
		tabManager.CloseTab(id)
		sessionManager.RecordClosed(closedTitle, closedURL)
		pinManager.Clear(id)
		readerMode.Clear(id)

		active := tabManager.GetActiveTab()
		if active == nil {
			tab := tabManager.NewTab("")
			tabManager.RecordNavigation(tab.ID, "sparrow://home")
			navigateTo(w, "sparrow://home")
			return
		}
		target := active.URL
		if target == "" {
			target = "sparrow://home"
		}
		navigateTo(w, target)
	})

	w.Bind("switchTab", func(id int) {
		tab := tabManager.SwitchTab(id)
		if tab == nil {
			return
		}
		target := tab.URL
		if target == "" {
			target = "sparrow://home"
		}
		navigateTo(w, target)
	})

	w.Bind("goBackNav", func() {
		tab := tabManager.GetActiveTab()
		if tab == nil {
			return
		}
		if url, ok := tabManager.GoBack(tab.ID); ok {
			navigateTo(w, url)
		}
	})

	w.Bind("goForwardNav", func() {
		tab := tabManager.GetActiveTab()
		if tab == nil {
			return
		}
		if url, ok := tabManager.GoForward(tab.ID); ok {
			navigateTo(w, url)
		}
	})

	w.Bind("getNavState", func() NavState {
		tab := tabManager.GetActiveTab()
		if tab == nil {
			return NavState{}
		}
		return NavState{
			URL:        tab.URL,
			CanBack:    tabManager.CanGoBack(tab.ID),
			CanForward: tabManager.CanGoForward(tab.ID),
			TabID:      tab.ID,
		}
	})

	w.Bind("getAllTabsUI", func() []*Tab {
		return tabManager.GetAllTabs()
	})

	w.Bind("updateTabTitle", func(id int, title string, url string) {
		tabManager.SetTitle(id, title)
	})

	// -------------------------------
	// SECURITY SHIELD CONTROLS
	// -------------------------------

	w.Bind("toggleAdBlock", func() bool {
		return shield.ToggleAdBlock()
	})

	w.Bind("toggleTrackerBlock", func() bool {
		return shield.ToggleTrackerBlock()
	})

	w.Bind("toggleHTTPSOnly", func() bool {
		return shield.ToggleHTTPSOnly()
	})

	w.Bind("getShieldStats", func() map[string]any {
		return shield.GetStats()
	})

	// -------------------------------
	// SETTINGS CONTROLS
	// -------------------------------

	w.Bind("getSettings", func() map[string]any {
		return settings.GetAll()
	})

	w.Bind("setTheme", func(theme string) string {
		return settings.SetTheme(theme)
	})

	w.Bind("setSearchEngine", func(engine string) string {
		return settings.SetSearchEngine(engine)
	})

	w.Bind("setHomepage", func(url string) string {
		return settings.SetHomepage(url)
	})

	w.Bind("toggleJS", func() bool {
		return settings.ToggleJS()
	})

	w.Bind("toggleImages", func() bool {
		return settings.ToggleImages()
	})

	w.Bind("resetSettings", func() map[string]any {
		return settings.ResetToDefault()
	})

	// -------------------------------
	// BOOKMARK CONTROLS
	// -------------------------------

	w.Bind("toggleBookmark", func(title string, url string) bool {
		return bookmarkManager.Toggle(title, url)
	})

	w.Bind("isBookmarked", func(url string) bool {
		return bookmarkManager.IsBookmarked(url)
	})

	w.Bind("getBookmarks", func() []*Bookmark {
		return bookmarkManager.GetAll()
	})

	w.Bind("removeBookmark", func(id int) bool {
		return bookmarkManager.Remove(id)
	})

	w.Bind("searchBookmarks", func(keyword string) []*Bookmark {
		return bookmarkManager.Search(keyword)
	})

	// -------------------------------
	// DOWNLOAD CONTROLS
	// -------------------------------

	w.Bind("getDownloads", func() []*DownloadItem {
		return downloadManager.GetAll()
	})

	w.Bind("clearFinishedDownloads", func() {
		downloadManager.ClearFinished()
	})

	w.Bind("cancelDownload", func(id int) {
		downloadManager.Cancel(id)
	})

	// -------------------------------
	// INCOGNITO MODE
	// -------------------------------

	w.Bind("enableIncognito", func() {
		incognito.Enable()
	})

	w.Bind("disableIncognito", func() {
		incognito.Disable()
	})

	w.Bind("isIncognitoActive", func() bool {
		return incognito.IsActive()
	})

	w.Bind("getIncognitoDuration", func() string {
		return incognito.Duration()
	})

	// -------------------------------
	// PERMISSIONS
	// -------------------------------

	w.Bind("getSitePermissions", func(domain string) map[string]string {
		perms := permissionManager.GetFor(domain)
		return map[string]string{
			"camera":        string(perms.Camera),
			"microphone":    string(perms.Microphone),
			"location":      string(perms.Location),
			"notifications": string(perms.Notifications),
		}
	})

	w.Bind("setSitePermission", func(domain string, permType string, state string) map[string]string {
		return permissionManager.Set(domain, permType, state)
	})

	w.Bind("resetSitePermissions", func(domain string) {
		permissionManager.ResetDomain(domain)
	})

	w.Bind("resetAllPermissions", func() {
		permissionManager.ResetAll()
	})

	// -------------------------------
	// READER MODE
	// -------------------------------

	w.Bind("toggleReaderMode", func(tabID int) bool {
		isOn := readerMode.Toggle(tabID)
		if isOn {
			w.Eval(readerModeJS)
		}
		return isOn
	})

	// -------------------------------
	// SHORTCUTS
	// -------------------------------

	w.Bind("getShortcuts", func() []Shortcut {
		return getAllShortcuts()
	})

	// -------------------------------
	// SESSION (RECENTLY CLOSED TABS)
	// -------------------------------

	w.Bind("getRecentlyClosed", func() []ClosedTab {
		return sessionManager.GetAll()
	})

	w.Bind("reopenLastClosed", func() *ClosedTab {
		return sessionManager.PopLastClosed()
	})

	w.Bind("clearRecentlyClosed", func() {
		sessionManager.Clear()
	})

	// -------------------------------
	// ZOOM CONTROLS
	// -------------------------------

	w.Bind("zoomIn", func(tabID int) float64 {
		return zoomManager.ZoomIn(tabID)
	})

	w.Bind("zoomOut", func(tabID int) float64 {
		return zoomManager.ZoomOut(tabID)
	})

	w.Bind("zoomReset", func(tabID int) float64 {
		return zoomManager.Reset(tabID)
	})

	w.Bind("getZoom", func(tabID int) float64 {
		return zoomManager.Get(tabID)
	})

	// -------------------------------
	// REAL DOWNLOADS
	// -------------------------------

	w.Bind("downloadFile", func(url string) *DownloadItem {
		return StartDownload(url, settings.DownloadPath)
	})

	w.Bind("cancelActiveDownload", func(id int) {
		CancelActiveDownload(id)
		downloadManager.Cancel(id)
	})

	w.Bind("isDownloadableLink", func(url string) bool {
		return IsDownloadableURL(url)
	})

	// -------------------------------
	// CONTEXT MENU HELPERS
	// -------------------------------

	w.Bind("openLinkInNewTab", func(url string) int {
		tab := tabManager.NewTab(url)
		history.Add(url, url)
		return tab.ID
	})

	w.Bind("openDevTools", func() {
		fmt.Println("DevTools requested")
	})

	// -------------------------------
	// TAB PINNING
	// -------------------------------

	w.Bind("togglePinTab", func(tabID int) bool {
		return pinManager.Toggle(tabID)
	})

	w.Bind("isPinnedTab", func(tabID int) bool {
		return pinManager.IsPinned(tabID)
	})

	w.Bind("getOrderedTabs", func() []*Tab {
		return GetOrderedTabs()
	})

	// -------------------------------
	// NOTIFICATION CENTER
	// -------------------------------

	w.Bind("pushNotification", func(title string, message string, kind string) *NotificationItem {
		return notifCenter.Push(title, message, kind)
	})

	w.Bind("getNotifications", func() []*NotificationItem {
		return notifCenter.GetAll()
	})

	w.Bind("markNotifRead", func(id int) {
		notifCenter.MarkRead(id)
	})

	w.Bind("markAllNotifsRead", func() {
		notifCenter.MarkAllRead()
	})

	w.Bind("getUnreadNotifCount", func() int {
		return notifCenter.UnreadCount()
	})

	w.Bind("clearAllNotifs", func() {
		notifCenter.ClearAll()
	})

	// -------------------------------
	// AUTOCOMPLETE / SUGGESTIONS
	// -------------------------------

	w.Bind("getAutocomplete", func(query string) []AutocompleteResult {
		return getAutocompleteResults(query)
	})

	w.Bind("getSuggestions", func(partial string) []AutocompleteResult {
		results := getAutocompleteResults(partial)
		if len(results) > 6 {
			return results[:6]
		}
		return results
	})

	// -------------------------------
	// SITE INFO POPUP
	// -------------------------------

	w.Bind("getSiteInfo", func(url string) SiteInfo {
		return BuildSiteInfo(url)
	})

	// -------------------------------
	// FIND IN PAGE & CONTEXT MENU
	// -------------------------------

	w.Bind("injectFindInPage", func() {
		w.Eval(findInPageJS)
	})

	w.Bind("injectContextMenu", func() {
		w.Eval(contextMenuJS)
	})

	// -------------------------------
	// BOOKMARK EXPORT / IMPORT
	// -------------------------------

	w.Bind("exportBookmarks", func() string {
		return ExportBookmarksJSON()
	})

	w.Bind("importBookmarks", func(jsonStr string) int {
		return ImportBookmarksJSON(jsonStr)
	})

	// -------------------------------
	// SITE SETTINGS PAGE
	// -------------------------------

	w.Bind("getAllSitePermissions", func() map[string]map[string]string {
		return permissionManager.GetAllAsMap()
	})

	// -------------------------------
	// SMART TOP SITES & SHORTCUTS
	// -------------------------------

	w.Bind("getTopSites", func(limit int) []TopSite {
		return topSitesEngine.GetTopSites(limit)
	})

	w.Bind("addShortcut", func(title string, url string) []PinnedShortcut {
		return shortcutManagerNewTab.AddShortcut(title, url)
	})

	w.Bind("removeShortcut", func(index int) []PinnedShortcut {
		return shortcutManagerNewTab.RemoveShortcut(index)
	})

	w.Bind("getShortcutsNewTab", func() []PinnedShortcut {
		return shortcutManagerNewTab.GetShortcuts()
	})

	// -------------------------------
	// COOKIE MANAGER
	// -------------------------------

	w.Bind("reportCookie", func(domain string, name string) bool {
		return cookieManager.Record(domain, name)
	})

	w.Bind("getCookieSummary", func() map[string]int {
		return cookieManager.GetSummary()
	})

	w.Bind("clearDomainCookies", func(domain string) {
		cookieManager.ClearDomain(domain)
	})

	w.Bind("clearAllCookies", func() {
		cookieManager.ClearAll()
	})

	w.Bind("setCookieBlocked", func(domain string, blocked bool) {
		cookieManager.SetBlocked(domain, blocked)
	})

	// -------------------------------
	// FORCE DARK MODE
	// -------------------------------

	w.Bind("toggleForceDarkGlobal", func() bool {
		return forceDark.ToggleGlobal()
	})

	w.Bind("toggleForceDarkTab", func(tabID int) bool {
		on := forceDark.ToggleForTab(tabID)
		w.Eval(forceDarkJS)
		return on
	})

	// -------------------------------
	// PERFORMANCE / BATTERY SAVER
	// -------------------------------

	w.Bind("toggleSaverMode", func() bool {
		return perfManager.ToggleSaver()
	})

	w.Bind("isSaverModeOn", func() bool {
		return perfManager.IsSaverOn()
	})

	w.Bind("suspendTab", func(tabID int) {
		perfManager.SuspendTab(tabID)
	})

	w.Bind("wakeTab", func(tabID int) {
		perfManager.WakeTab(tabID)
	})

	w.Bind("getSuspendedCount", func() int {
		return perfManager.SuspendedCount()
	})

	// -------------------------------
	// PRINT
	// -------------------------------

	w.Bind("printPage", func() {
		w.Eval(printPageJS)
	})

	// -------------------------------
	// AUTOFILL
	// -------------------------------

	w.Bind("saveAutofillProfile", func(data AutofillData) AutofillData {
		data.Name = strings.TrimSpace(data.Name)
		data.Email = strings.TrimSpace(data.Email)
		data.Phone = strings.TrimSpace(data.Phone)
		data.Address = strings.TrimSpace(data.Address)
		saveAutofillProfile(data)
		return data
	})

	w.Bind("getAutofillProfile", func() AutofillData {
		return getAutofillProfile()
	})

	w.Bind("clearAutofillProfile", func() {
		clearAutofillProfile()
	})

	w.Bind("toggleAutofillEnabled", func() bool {
		autofillEnabledMu.Lock()
		defer autofillEnabledMu.Unlock()
		autofillEnabled = !autofillEnabled
		return autofillEnabled
	})

	// -------------------------------
	// WALLPAPER
	// -------------------------------

	w.Bind("setWallpaper", func(id string) string {
		return wallpaperManager.SetWallpaper(id)
	})

	w.Bind("getWallpaperCSS", func() string {
		return wallpaperManager.GetSelectedCSS()
	})

	w.Bind("getWallpaperOptions", func() []WallpaperOption {
		return wallpaperManager.GetAllOptions()
	})

	// -------------------------------
	// ACCESSIBILITY
	// -------------------------------

	w.Bind("setFontScale", func(scale float64) float64 {
		return accessibility.SetFontScale(scale)
	})

	w.Bind("toggleHighContrast", func() bool {
		return accessibility.ToggleHighContrast()
	})

	w.Bind("toggleReduceMotion", func() bool {
		return accessibility.ToggleReduceMotion()
	})

	w.Bind("toggleUnderlineLinks", func() bool {
		return accessibility.ToggleUnderlineLinks()
	})

	w.Bind("getAccessibilitySettings", func() map[string]any {
		return accessibility.GetAll()
	})

	w.Bind("applyAccessibility", func() {
		s := accessibility
		w.Eval(buildAccessibilityJS(s.FontScale, s.HighContrast, s.ReduceMotion, s.UnderlineLinks))
	})

	// -------------------------------
	// FEEDBACK
	// -------------------------------

	w.Bind("submitFeedback", func(message string, category string) *FeedbackEntry {
		return feedbackManager.Submit(message, category)
	})

	w.Bind("getAllFeedback", func() []*FeedbackEntry {
		return feedbackManager.GetAll()
	})

	// -------------------------------
	// PASSWORD MANAGER
	// -------------------------------

	w.Bind("unlockVault", func(masterPassword string) bool {
		return passwordVault.Unlock(masterPassword)
	})

	w.Bind("lockVault", func() {
		passwordVault.Lock()
	})

	w.Bind("isVaultUnlocked", func() bool {
		return passwordVault.IsUnlocked()
	})

	w.Bind("savePasswordEntry", func(domain string, username string, password string) (*PasswordEntry, error) {
		return passwordVault.Save(domain, username, password)
	})

	w.Bind("revealPasswordEntry", func(id int) (string, error) {
		return passwordVault.Reveal(id)
	})

	w.Bind("deletePasswordEntry", func(id int) bool {
		return passwordVault.Remove(id)
	})

	w.Bind("getPasswordList", func() []PasswordMetadata {
		return passwordVault.GetMetadataOnly()
	})

	w.Bind("findPasswordsForDomain", func(domain string) []PasswordMetadata {
		return passwordVault.FindForDomain(domain)
	})

	w.Bind("wipeAllPasswords", func() {
		passwordVault.WipeAll()
	})

	w.Bind("generatePassword", func(length int) string {
		return GenerateStrongPassword(length)
	})

	// -------------------------------
	// UPDATE CHECKER
	// -------------------------------

	w.Bind("getChangelog", func() []ChangelogEntry {
		return sparrowChangelog
	})

	// -------------------------------
	// PAGE EXPORT
	// -------------------------------

	w.Bind("savePageText", func(title string, text string) (string, error) {
		return SavePageText(title, text)
	})

	w.Bind("extractPageText", func() {
		w.Eval(extractPageTextJS)
	})

	// -------------------------------
	// BACKUP / RESTORE
	// -------------------------------

	w.Bind("exportBackup", func() string {
		return ExportBackupJSON()
	})

	w.Bind("saveBackupToFile", func() (string, error) {
		return SaveBackupToFile()
	})

	w.Bind("restoreBackup", func(jsonStr string) (int, error) {
		return RestoreBackupJSON(jsonStr)
	})

	// -------------------------------
	// PICTURE-IN-PICTURE
	// -------------------------------

	w.Bind("requestPiP", func(tabID int) {
		pipManager.SetActive(tabID, true)
		w.Eval(pipRequestJS)
	})

	w.Bind("isPiPActive", func(tabID int) bool {
		return pipManager.IsActive(tabID)
	})

	// -------------------------------
	// CUSTOM THEME
	// -------------------------------

	w.Bind("setAccentColor", func(hex string) string {
		return customTheme.SetAccentColor(hex)
	})

	w.Bind("setTabRadius", func(radius int) int {
		return customTheme.SetTabRadius(radius)
	})

	w.Bind("toggleCompactToolbar", func() bool {
		return customTheme.ToggleCompact()
	})

	w.Bind("getCustomTheme", func() map[string]any {
		return customTheme.GetAll()
	})

	w.Bind("getThemeCSS", func() string {
		return customTheme.BuildThemeCSS()
	})

	// -------------------------------
	// CRASH RECOVERY
	// -------------------------------

	w.Bind("updateRecoverySnapshot", func(urls []string) {
		recoverySnapshot.UpdateSnapshot(urls)
	})

	w.Bind("getRecoverySnapshot", func() []string {
		return recoverySnapshot.GetSnapshot()
	})

	w.Bind("hasRecoverableSession", func() bool {
		return recoverySnapshot.HasRecoverableSession()
	})

	w.Bind("clearRecoverySnapshot", func() {
		recoverySnapshot.ClearSnapshot()
	})

	w.Bind("toggleCrashRecovery", func() bool {
		return recoverySnapshot.ToggleEnabled()
	})

	// -------------------------------
	// SCREENSHOT
	// -------------------------------

	w.Bind("captureScreenshot", func() {
		w.Eval(screenshotCaptureJS)
	})

	w.Bind("reportScreenshot", func(base64Data string, title string) (string, error) {
		return SaveScreenshot(base64Data, title)
	})

	// -------------------------------
	// TAB SEARCH
	// -------------------------------

	w.Bind("searchTabs", func(query string) []TabSearchResult {
		return tabSearchEngine.Search(query)
	})

	// -------------------------------
	// CUSTOM SEARCH ENGINES
	// -------------------------------

	w.Bind("addCustomSearchEngine", func(name string, template string) (*CustomSearchEngine, bool) {
		return customSearchManager.AddEngine(name, template)
	})

	w.Bind("removeCustomSearchEngine", func(id int) bool {
		return customSearchManager.RemoveEngine(id)
	})

	w.Bind("setActiveCustomSearchEngine", func(id int) bool {
		return customSearchManager.SetActive(id)
	})

	w.Bind("getCustomSearchEngines", func() []*CustomSearchEngine {
		return customSearchManager.GetAll()
	})

	// -------------------------------
	// SPEED DIAL
	// -------------------------------

	w.Bind("addSpeedDialTile", func(title string, url string) *SpeedDialTile {
		return speedDialManager.AddTile(title, url)
	})

	w.Bind("removeSpeedDialTile", func(id int) bool {
		return speedDialManager.RemoveTile(id)
	})

	w.Bind("reorderSpeedDial", func(fromIndex int, toIndex int) bool {
		return speedDialManager.Reorder(fromIndex, toIndex)
	})

	w.Bind("getSpeedDialTiles", func() []*SpeedDialTile {
		return speedDialManager.GetAll()
	})

	// -------------------------------
	// OPEN INITIAL TAB
	// -------------------------------

	firstTab := tabManager.NewTab("sparrow://home")
	tabManager.RecordNavigation(firstTab.ID, "sparrow://home")
	w.Navigate(dataURI(homePageHTML))

	w.Run()
}
