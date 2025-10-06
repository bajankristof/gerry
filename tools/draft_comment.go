package tools

import (
	"context"
	"fmt"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/gerrit"
	"github.com/mark3labs/mcp-go/mcp"
)

// DraftCommentTool is the tool definition for draft_comment
var DraftCommentTool = mcp.NewTool("draft_comment",
	mcp.WithDescription("Create a draft comment or reply on a Gerrit change. Drafts are not visible until published with publish_review."),
	mcp.WithString("changeId",
		mcp.Description("The Gerrit Change-Id (e.g., I1234567890abcdef...). If omitted, will auto-detect from current git commit."),
	),
	mcp.WithString("message",
		mcp.Required(),
		mcp.Description("The comment or reply message"),
	),
	mcp.WithString("path",
		mcp.Required(),
		mcp.Description("The file path for the comment"),
	),
	mcp.WithNumber("line",
		mcp.Description("The line number for the comment (omit for file-level comments)"),
	),
	mcp.WithString("inReplyTo",
		mcp.Description("The comment ID to reply to (omit for new comments)"),
	),
	mcp.WithBoolean("unresolved",
		mcp.Description("Whether the comment should be marked as unresolved (default: false)"),
	),
	mcp.WithString("directory",
		mcp.Description("The directory containing the git repository (used to determine Gerrit host)"),
	),
)

// HandleDraftComment handles the draft_comment tool call
func HandleDraftComment(cfg *config.Config) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		changeID, err := inferChangeID(request)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		message, err := request.RequireString("message")
		if err != nil {
			return mcp.NewToolResultError("message is required"), nil
		}

		path, err := request.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError("path is required"), nil
		}

		directory := request.GetString("directory", "")
		client, err := gerrit.NewClientFromGit(directory, cfg.GerritUsername, cfg.GerritPassword)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		input := gerrit.DraftCommentInput{
			Message:    message,
			Path:       path,
			Line:       request.GetInt("line", 0),
			InReplyTo:  request.GetString("inReplyTo", ""),
			Unresolved: request.GetBool("unresolved", false),
		}

		if err := client.DraftComment(changeID, input); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText("Success."), nil
	}
}
