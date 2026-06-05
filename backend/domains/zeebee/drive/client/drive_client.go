package client

import (
	"context"
	"io"
	"strings"
	"sync"

	"golang.org/x/oauth2"
	driveapi "google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"zeebee-backend/common/errorx"
)

const folderMimeType = "application/vnd.google-apps.folder"

// managedAppPropertyKey marks files that Zeebee uploaded, so reconciliation can
// safely remove our own stale files without touching files a user added manually.
const managedAppPropertyKey = "zeebeeManaged"
const managedAppPropertyValue = "1"

// DriveInfo is a lightweight description of a Google Shared Drive (or My Drive).
type DriveInfo struct {
	Id   string
	Name string
}

// DriveFileInfo is a lightweight description of a Drive file.
type DriveFileInfo struct {
	Id   string
	Name string
}

// DriveFolderInfo is a lightweight description of a Drive folder.
type DriveFolderInfo struct {
	Id          string
	Name        string
	WebViewLink string
}

// DriveClient wraps the Google Drive v3 API for a single authorized account.
//
// Drive scope convention: every file/folder operation carries a driveId.
//   - driveId == ""  -> the authorized user's My Drive (default user corpus)
//   - driveId != ""  -> a Shared Drive identified by driveId
type DriveClient struct {
	mu     sync.Mutex
	config *oauth2.Config
	token  *oauth2.Token
}

// NewDriveClient exchanges an OAuth code for a token and returns a client.
func NewDriveClient(config *oauth2.Config, code string) (client *DriveClient, err error) {
	token, err := config.Exchange(context.Background(), code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, errorx.Wrap(err, "Exchange token failed")
	}
	return &DriveClient{config: config, token: token}, nil
}

// NewDriveClientFromToken builds a client from a previously stored token.
func NewDriveClientFromToken(config *oauth2.Config, token *oauth2.Token) (client *DriveClient, err error) {
	if token == nil {
		return nil, errorx.New("Token is nil")
	}
	return &DriveClient{config: config, token: token}, nil
}

// Token returns the current (possibly refreshed) token.
func (cli *DriveClient) Token() *oauth2.Token {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	return cli.token
}

// RefreshToken refreshes the token when expired and returns the current token.
func (cli *DriveClient) RefreshToken(ctx context.Context) (token *oauth2.Token, err error) {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	if cli.token == nil {
		return nil, errorx.New("Token is nil")
	}
	if cli.token.RefreshToken == "" {
		return nil, errorx.New("Refresh token is nil")
	}
	if cli.token.Valid() {
		return cli.token, nil
	}
	ts := cli.config.TokenSource(ctx, cli.token)
	token, err = ts.Token()
	if err != nil {
		return nil, errorx.Wrap(err, "Refresh token failed")
	}
	if token == nil {
		return nil, errorx.New("New token is nil")
	}
	cli.token = token
	return token, nil
}

// GetEmail returns the email address of the authorized Google account.
func (cli *DriveClient) GetEmail(ctx context.Context) (email string, err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return "", err
	}
	about, err := srv.About.Get().Fields("user(emailAddress)").Context(ctx).Do()
	if err != nil {
		return "", errorx.WithStack(err)
	}
	if about.User == nil {
		return "", nil
	}
	return about.User.EmailAddress, nil
}

// ListDrives returns the Shared Drives accessible to the account. My Drive is NOT
// included here; callers prepend it manually when needed.
func (cli *DriveClient) ListDrives(ctx context.Context) (drives []DriveInfo, err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return nil, err
	}
	pageToken := ""
	for {
		call := srv.Drives.List().PageSize(100).Fields("nextPageToken, drives(id, name)")
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		res, err := call.Context(ctx).Do()
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		for _, d := range res.Drives {
			drives = append(drives, DriveInfo{Id: d.Id, Name: d.Name})
		}
		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}
	return drives, nil
}

// EnsureFolder returns an existing folder by name under parent, creating it if absent.
func (cli *DriveClient) EnsureFolder(ctx context.Context, name, parentId, driveId string) (id, webViewLink string, err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return "", "", err
	}
	parent := resolveParent(parentId, driveId)

	q := "name = '" + escape(name) + "' and mimeType = '" + folderMimeType +
		"' and '" + escape(parent) + "' in parents and trashed = false"
	list := applyDriveScope(srv.Files.List().Q(q).Fields("files(id, webViewLink)").PageSize(1), driveId)
	res, err := list.Context(ctx).Do()
	if err != nil {
		return "", "", errorx.WithStack(err)
	}
	if len(res.Files) > 0 {
		return res.Files[0].Id, res.Files[0].WebViewLink, nil
	}

	created, err := srv.Files.Create(&driveapi.File{
		Name:     name,
		MimeType: folderMimeType,
		Parents:  []string{parent},
	}).SupportsAllDrives(true).Fields("id, webViewLink").Context(ctx).Do()
	if err != nil {
		return "", "", errorx.WithStack(err)
	}
	return created.Id, created.WebViewLink, nil
}

// ListFolders returns the non-trashed sub-folders directly under parent.
func (cli *DriveClient) ListFolders(ctx context.Context, parentId, driveId string) (folders []DriveFolderInfo, err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return nil, err
	}
	parent := resolveParent(parentId, driveId)
	q := "mimeType = '" + folderMimeType + "' and '" + escape(parent) + "' in parents and trashed = false"
	pageToken := ""
	for {
		call := applyDriveScope(srv.Files.List().Q(q).Fields("nextPageToken, files(id, name, webViewLink)").PageSize(100), driveId)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		res, err := call.Context(ctx).Do()
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		for _, f := range res.Files {
			folders = append(folders, DriveFolderInfo{Id: f.Id, Name: f.Name, WebViewLink: f.WebViewLink})
		}
		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}
	return folders, nil
}

// RenameFile renames a file or folder by id. The id (and thus webViewLink) is unchanged.
func (cli *DriveClient) RenameFile(ctx context.Context, fileId, newName string) (err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return err
	}
	_, err = srv.Files.Update(fileId, &driveapi.File{Name: newName}).
		SupportsAllDrives(true).Context(ctx).Do()
	return errorx.WithStack(err)
}

// FindFile returns the first non-trashed file with name under parent, or empty id if absent.
func (cli *DriveClient) FindFile(ctx context.Context, name, parentId, driveId string) (id string, err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return "", err
	}
	parent := resolveParent(parentId, driveId)
	q := "name = '" + escape(name) + "' and '" + escape(parent) +
		"' in parents and trashed = false"
	list := applyDriveScope(srv.Files.List().Q(q).Fields("files(id)").PageSize(1), driveId)
	res, err := list.Context(ctx).Do()
	if err != nil {
		return "", errorx.WithStack(err)
	}
	if len(res.Files) > 0 {
		return res.Files[0].Id, nil
	}
	return "", nil
}

// ListManagedFiles returns the non-folder files under parent that Zeebee uploaded
// (carrying the managed app-property marker). User-added files are excluded, so
// callers may safely reconcile/remove these without affecting manual uploads.
func (cli *DriveClient) ListManagedFiles(ctx context.Context, parentId, driveId string) (files []DriveFileInfo, err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return nil, err
	}
	parent := resolveParent(parentId, driveId)
	q := "'" + escape(parent) + "' in parents and trashed = false and mimeType != '" +
		folderMimeType + "' and appProperties has { key = '" + managedAppPropertyKey +
		"' and value = '" + managedAppPropertyValue + "' }"
	pageToken := ""
	for {
		call := applyDriveScope(srv.Files.List().Q(q).Fields("nextPageToken, files(id, name)").PageSize(100), driveId)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		res, err := call.Context(ctx).Do()
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		for _, f := range res.Files {
			files = append(files, DriveFileInfo{Id: f.Id, Name: f.Name})
		}
		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}
	return files, nil
}

// TrashFile moves a file to the trash (reversible) rather than deleting it
// permanently, so an accidental reconciliation can be recovered from Drive.
func (cli *DriveClient) TrashFile(ctx context.Context, fileId string) (err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return err
	}
	_, err = srv.Files.Update(fileId, &driveapi.File{
		Trashed:         true,
		ForceSendFields: []string{"Trashed"},
	}).SupportsAllDrives(true).Context(ctx).Do()
	return errorx.WithStack(err)
}

// UpsertFile creates the file or replaces its content if a same-named file already
// exists under parent. replaced reports whether an existing file was overwritten.
func (cli *DriveClient) UpsertFile(ctx context.Context, name, contentType string, body io.Reader, parentId, driveId string) (id, webViewLink string, replaced bool, err error) {
	existingId, err := cli.FindFile(ctx, name, parentId, driveId)
	if err != nil {
		return "", "", false, err
	}
	if existingId != "" {
		id, webViewLink, err = cli.updateFileContent(ctx, existingId, contentType, body)
		return id, webViewLink, true, err
	}
	id, webViewLink, err = cli.uploadFile(ctx, name, contentType, body, parentId, driveId)
	return id, webViewLink, false, err
}

func (cli *DriveClient) uploadFile(ctx context.Context, name, contentType string, body io.Reader, parentId, driveId string) (id, webViewLink string, err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return "", "", err
	}
	parent := resolveParent(parentId, driveId)
	created, err := srv.Files.Create(&driveapi.File{
		Name:          name,
		Parents:       []string{parent},
		AppProperties: map[string]string{managedAppPropertyKey: managedAppPropertyValue},
	}).Media(body, googleapi.ContentType(contentType)).
		SupportsAllDrives(true).Fields("id, webViewLink").Context(ctx).Do()
	if err != nil {
		return "", "", errorx.WithStack(err)
	}
	return created.Id, created.WebViewLink, nil
}

func (cli *DriveClient) updateFileContent(ctx context.Context, fileId, contentType string, body io.Reader) (id, webViewLink string, err error) {
	srv, err := cli.getSrv(ctx)
	if err != nil {
		return "", "", err
	}
	updated, err := srv.Files.Update(fileId, &driveapi.File{
		AppProperties: map[string]string{managedAppPropertyKey: managedAppPropertyValue},
	}).
		Media(body, googleapi.ContentType(contentType)).
		SupportsAllDrives(true).Fields("id, webViewLink").Context(ctx).Do()
	if err != nil {
		return "", "", errorx.WithStack(err)
	}
	return updated.Id, updated.WebViewLink, nil
}

func (cli *DriveClient) getSrv(ctx context.Context) (srv *driveapi.Service, err error) {
	cli.mu.Lock()
	token := cli.token
	cli.mu.Unlock()
	httpClient := cli.config.Client(ctx, token)
	srv, err = driveapi.NewService(ctx, option.WithHTTPClient(httpClient))
	return srv, errorx.WithStack(err)
}

// resolveParent picks the effective parent folder id: an explicit parent wins;
// otherwise a Shared Drive root (driveId) or My Drive root ("root").
func resolveParent(parentId, driveId string) string {
	if parentId != "" {
		return parentId
	}
	if driveId != "" {
		return driveId
	}
	return "root"
}

// applyDriveScope narrows a files list query to a Shared Drive when driveId is set.
func applyDriveScope(call *driveapi.FilesListCall, driveId string) *driveapi.FilesListCall {
	call = call.SupportsAllDrives(true)
	if driveId != "" {
		call = call.Corpora("drive").DriveId(driveId).IncludeItemsFromAllDrives(true)
	}
	return call
}

// escape escapes a literal for use inside a Drive query string.
func escape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return s
}
