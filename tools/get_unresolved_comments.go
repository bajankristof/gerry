package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/gerrit"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetUnresolvedCommentsTool is the tool definition for gerrit_get_unresolved_comments
var GetUnresolvedCommentsTool = mcp.NewTool("gerrit_get_unresolved_comments",
	mcp.WithDescription("Get all unresolved comments for a Gerrit change. Returns a list of comments with their file path, line number, message, and author. These are the comments that need to be addressed."),
	mcp.WithString("changeId",
		mcp.Required(),
		mcp.Description("The Gerrit Change-Id (e.g., I1234567890abcdef...)"),
	),
	mcp.WithString("directory",
		mcp.Description("The directory containing the git repository (used to determine Gerrit host)"),
	),
)

// HandleGetUnresolvedComments handles the gerrit_get_unresolved_comments tool call
func HandleGetUnresolvedComments(cfg *config.Config) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		changeID, err := request.RequireString("changeId")
		if err != nil {
			return mcp.NewToolResultError("changeId is required"), nil
		}

		directory := request.GetString("directory", "")
		client, err := gerrit.NewClientFromGit(directory, cfg.GerritUsername, cfg.GerritPassword)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		comments, err := client.GetUnresolvedComments(changeID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		if len(comments) == 0 {
			return mcp.NewToolResultText("No unresolved comments found."), nil
		}

		commentsJSON, err := json.MarshalIndent(comments, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText(string(commentsJSON)), nil
	}
}
