// This file is generated by "./lib/proto/generate"

package proto

/*

PWA

This domain allows interacting with the browser to control PWAs.

*/

// PWAFileHandlerAccept The following types are the replica of
// https://crsrc.org/c/chrome/browser/web_applications/proto/web_app_os_integration_state.proto;drc=9910d3be894c8f142c977ba1023f30a656bc13fc;l=67
type PWAFileHandlerAccept struct {
	// MediaType New name of the mimetype according to
	// https://www.iana.org/assignments/media-types/media-types.xhtml
	MediaType string `json:"mediaType"`

	// FileExtensions ...
	FileExtensions []string `json:"fileExtensions"`
}

// PWAFileHandler ...
type PWAFileHandler struct {
	// Action ...
	Action string `json:"action"`

	// Accepts ...
	Accepts []*PWAFileHandlerAccept `json:"accepts"`

	// DisplayName ...
	DisplayName string `json:"displayName"`
}

// PWADisplayMode If user prefers opening the app in browser or an app window.
type PWADisplayMode string

const (
	// PWADisplayModeStandalone enum const.
	PWADisplayModeStandalone PWADisplayMode = "standalone"

	// PWADisplayModeBrowser enum const.
	PWADisplayModeBrowser PWADisplayMode = "browser"
)

// PWAGetOsAppState Returns the following OS state for the given manifest id.
type PWAGetOsAppState struct {
	// ManifestID The id from the webapp's manifest file, commonly it's the url of the
	// site installing the webapp. See
	// https://web.dev/learn/pwa/web-app-manifest.
	ManifestID string `json:"manifestId"`
}

// ProtoReq name.
func (m PWAGetOsAppState) ProtoReq() string { return "PWA.getOsAppState" }

// Call the request.
func (m PWAGetOsAppState) Call(c Client) (*PWAGetOsAppStateResult, error) {
	var res PWAGetOsAppStateResult
	return &res, call(m.ProtoReq(), m, &res, c)
}

// PWAGetOsAppStateResult ...
type PWAGetOsAppStateResult struct {
	// BadgeCount ...
	BadgeCount int `json:"badgeCount"`

	// FileHandlers ...
	FileHandlers []*PWAFileHandler `json:"fileHandlers"`
}

// PWAInstall Installs the given manifest identity, optionally using the given install_url
// or IWA bundle location.
//
// TODO(crbug.com/337872319) Support IWA to meet the following specific
// requirement.
// IWA-specific install description: If the manifest_id is isolated-app://,
// install_url_or_bundle_url is required, and can be either an http(s) URL or
// file:// URL pointing to a signed web bundle (.swbn). The .swbn file's
// signing key must correspond to manifest_id. If Chrome is not in IWA dev
// mode, the installation will fail, regardless of the state of the allowlist.
type PWAInstall struct {
	// ManifestID ...
	ManifestID string `json:"manifestId"`

	// InstallURLOrBundleURL (optional) The location of the app or bundle overriding the one derived from the
	// manifestId.
	InstallURLOrBundleURL string `json:"installUrlOrBundleUrl,omitempty"`
}

// ProtoReq name.
func (m PWAInstall) ProtoReq() string { return "PWA.install" }

// Call sends the request.
func (m PWAInstall) Call(c Client) error {
	return call(m.ProtoReq(), m, nil, c)
}

// PWAUninstall Uninstalls the given manifest_id and closes any opened app windows.
type PWAUninstall struct {
	// ManifestID ...
	ManifestID string `json:"manifestId"`
}

// ProtoReq name.
func (m PWAUninstall) ProtoReq() string { return "PWA.uninstall" }

// Call sends the request.
func (m PWAUninstall) Call(c Client) error {
	return call(m.ProtoReq(), m, nil, c)
}

// PWALaunch Launches the installed web app, or an url in the same web app instead of the
// default start url if it is provided. Returns a page Target.TargetID which
// can be used to attach to via Target.attachToTarget or similar APIs.
type PWALaunch struct {
	// ManifestID ...
	ManifestID string `json:"manifestId"`

	// URL (optional) ...
	URL string `json:"url,omitempty"`
}

// ProtoReq name.
func (m PWALaunch) ProtoReq() string { return "PWA.launch" }

// Call the request.
func (m PWALaunch) Call(c Client) (*PWALaunchResult, error) {
	var res PWALaunchResult
	return &res, call(m.ProtoReq(), m, &res, c)
}

// PWALaunchResult ...
type PWALaunchResult struct {
	// TargetID ID of the tab target created as a result.
	TargetID TargetTargetID `json:"targetId"`
}

// PWALaunchFilesInApp Opens one or more local files from an installed web app identified by its
// manifestId. The web app needs to have file handlers registered to process
// the files. The API returns one or more page Target.TargetIDs which can be
// used to attach to via Target.attachToTarget or similar APIs.
// If some files in the parameters cannot be handled by the web app, they will
// be ignored. If none of the files can be handled, this API returns an error.
// If no files are provided as the parameter, this API also returns an error.
//
// According to the definition of the file handlers in the manifest file, one
// Target.TargetID may represent a page handling one or more files. The order
// of the returned Target.TargetIDs is not guaranteed.
//
// TODO(crbug.com/339454034): Check the existences of the input files.
type PWALaunchFilesInApp struct {
	// ManifestID ...
	ManifestID string `json:"manifestId"`

	// Files ...
	Files []string `json:"files"`
}

// ProtoReq name.
func (m PWALaunchFilesInApp) ProtoReq() string { return "PWA.launchFilesInApp" }

// Call the request.
func (m PWALaunchFilesInApp) Call(c Client) (*PWALaunchFilesInAppResult, error) {
	var res PWALaunchFilesInAppResult
	return &res, call(m.ProtoReq(), m, &res, c)
}

// PWALaunchFilesInAppResult ...
type PWALaunchFilesInAppResult struct {
	// TargetIDs IDs of the tab targets created as the result.
	TargetIDs []TargetTargetID `json:"targetIds"`
}

// PWAOpenCurrentPageInApp Opens the current page in its web app identified by the manifest id, needs
// to be called on a page target. This function returns immediately without
// waiting for the app to finish loading.
type PWAOpenCurrentPageInApp struct {
	// ManifestID ...
	ManifestID string `json:"manifestId"`
}

// ProtoReq name.
func (m PWAOpenCurrentPageInApp) ProtoReq() string { return "PWA.openCurrentPageInApp" }

// Call sends the request.
func (m PWAOpenCurrentPageInApp) Call(c Client) error {
	return call(m.ProtoReq(), m, nil, c)
}

// PWAChangeAppUserSettings Changes user settings of the web app identified by its manifestId. If the
// app was not installed, this command returns an error. Unset parameters will
// be ignored; unrecognized values will cause an error.
//
// Unlike the ones defined in the manifest files of the web apps, these
// settings are provided by the browser and controlled by the users, they
// impact the way the browser handling the web apps.
//
// See the comment of each parameter.
type PWAChangeAppUserSettings struct {
	// ManifestID ...
	ManifestID string `json:"manifestId"`

	// LinkCapturing (optional) If user allows the links clicked on by the user in the app's scope, or
	// extended scope if the manifest has scope extensions and the flags
	// `DesktopPWAsLinkCapturingWithScopeExtensions` and
	// `WebAppEnableScopeExtensions` are enabled.
	//
	// Note, the API does not support resetting the linkCapturing to the
	// initial value, uninstalling and installing the web app again will reset
	// it.
	//
	// TODO(crbug.com/339453269): Setting this value on ChromeOS is not
	// supported yet.
	LinkCapturing bool `json:"linkCapturing,omitempty"`

	// DisplayMode (optional) ...
	DisplayMode PWADisplayMode `json:"displayMode,omitempty"`
}

// ProtoReq name.
func (m PWAChangeAppUserSettings) ProtoReq() string { return "PWA.changeAppUserSettings" }

// Call sends the request.
func (m PWAChangeAppUserSettings) Call(c Client) error {
	return call(m.ProtoReq(), m, nil, c)
}