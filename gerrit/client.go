// Package gerrit provides a client for interacting with Gerrit Code Review REST API
package gerrit

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/bajankristof/gerry/git"
	"resty.dev/v3"
)

var (
	// ErrNoGerritHost is returned when the Gerrit host cannot be determined
	ErrNoGerritHost = errors.New("could not determine Gerrit host. Please provide a directory with a git remote configured")
)

// Author represents a comment author
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Range represents a comment range
type Range struct {
	StartLine      int `json:"start_line"`
	StartCharacter int `json:"start_character"`
	EndLine        int `json:"end_line"`
	EndCharacter   int `json:"end_character"`
}

// Comment represents a comment on a change
type Comment struct {
	ID         string `json:"id"`
	Author     Author `json:"author"`
	Message    string `json:"message"`
	Path       string `json:"path,omitempty"`
	Line       int    `json:"line,omitempty"`
	Range      *Range `json:"range,omitempty"`
	Unresolved bool   `json:"unresolved"`
	PatchSet   int    `json:"patch_set"`
	InReplyTo  string `json:"in_reply_to,omitempty"`
}

// Change represents a Gerrit change
type Change struct {
	ID              string `json:"id"`
	ChangeID        string `json:"change_id"`
	Project         string `json:"project"`
	Branch          string `json:"branch"`
	Subject         string `json:"subject"`
	Status          string `json:"status"`
	CurrentRevision string `json:"current_revision"`
}

// Client provides methods to interact with Gerrit
type Client struct {
	host     string
	username string
	password string
	client   *resty.Client
}

// NewClient creates a new Gerrit client
func NewClient(host, username, password string) *Client {
	client := resty.New()
	client.SetBasicAuth(username, password)
	client.SetBaseURL(fmt.Sprintf("https://%s/a", host))
	client.SetHeader("Content-Type", "application/json")
	client.SetDoNotParseResponse(true)
	client.AddResponseMiddleware(autoErrorMiddleware)
	client.AddResponseMiddleware(autoParseMiddleware)

	return &Client{
		host:     host,
		username: username,
		password: password,
		client:   client,
	}
}

// NewClientFromGit creates a Gerrit client by extracting the host from a git repository
func NewClientFromGit(directory, username, password string) (*Client, error) {
	var host string
	var err error
	if directory == "" {
		directory = "."
	}

	host, err = git.GetHostFromRemote(directory)
	if err != nil {
		return nil, err
	}

	if host == "" {
		return nil, ErrNoGerritHost
	}

	return NewClient(host, username, password), nil
}

// Host returns the Gerrit host
func (c *Client) Host() string {
	return c.host
}

// GetChange gets change information by Change-Id
func (c *Client) GetChange(changeID string) (Change, error) {
	path := fmt.Sprintf("/changes/%s", url.PathEscape(changeID))

	resp, err := c.client.R().SetResult(Change{}).Get(path)
	if err != nil {
		return Change{}, err
	}

	return *resp.Result().(*Change), nil
}

// GetComments gets all comments for a change
func (c *Client) GetComments(changeID string) ([]Comment, error) {
	path := fmt.Sprintf("/changes/%s/comments", url.PathEscape(changeID))

	resp, err := c.client.R().SetResult(map[string][]Comment{}).Get(path)
	if err != nil {
		return nil, err
	}

	rawComments := *resp.Result().(*map[string][]Comment)

	var result []Comment
	for path, comments := range rawComments {
		for _, comment := range comments {
			comment.Path = path
			result = append(result, comment)
		}
	}

	return result, nil
}

// GetUnresolvedComments gets all unresolved comments for a change
func (c *Client) GetUnresolvedComments(changeID string) ([]Comment, error) {
	comments, err := c.GetComments(changeID)
	if err != nil {
		return nil, err
	}

	var result []Comment
	for _, comment := range comments {
		if comment.Unresolved {
			result = append(result, comment)
		}
	}

	return result, nil
}

// DraftCommentInput represents a draft comment or reply
type DraftCommentInput struct {
	Message    string `json:"message"`
	Path       string `json:"path"`
	Line       int    `json:"line,omitempty"`
	InReplyTo  string `json:"in_reply_to,omitempty"`
	Unresolved bool   `json:"unresolved,omitempty"`
}

// DraftComment creates a draft comment or reply
func (c *Client) DraftComment(changeID string, input DraftCommentInput) error {
	path := fmt.Sprintf("/changes/%s/revisions/current/drafts", url.PathEscape(changeID))

	_, err := c.client.R().
		SetBody(input).
		Put(path)

	return err
}

// PublishReviewInput represents a review to be published
type PublishReviewInput struct {
	Message string `json:"message,omitempty"`
}

// PublishReview publishes all draft comments for a change
func (c *Client) PublishReview(changeID string, input PublishReviewInput) error {
	path := fmt.Sprintf("/changes/%s/revisions/current/review", url.PathEscape(changeID))

	body := map[string]any{"drafts":  "PUBLISH_ALL_REVISIONS"}
	if input.Message != "" {
		body["comments"] = map[string][]map[string]any{
			"/PATCHSET_LEVEL": {{"message": input.Message}},
		}
	}

	_, err := c.client.R().
		SetBody(body).
		Post(path)

	return err
}
