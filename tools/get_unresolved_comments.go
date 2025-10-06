package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/gerrit"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetUnresolvedCommentsTool is the tool definition for get_unresolved_comments
var GetUnresolvedCommentsTool = mcp.NewTool("get_unresolved_comments",
	mcp.WithDescription("Get all unresolved comments for a Gerrit change. Returns a list of comments with their file path, line number, message, and author. These are the comments that need to be addressed."),
	mcp.WithString("changeId",
		mcp.Description("The Gerrit Change-Id (e.g., I1234567890abcdef...). If omitted, will auto-detect from current git commit."),
	),
	mcp.WithString("directory",
		mcp.Description("The directory containing the git repository (used to determine Gerrit host)"),
	),
)

// HandleGetUnresolvedComments handles the get_unresolved_comments tool call
func HandleGetUnresolvedComments(cfg *config.Config) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		changeID, err := inferChangeID(request)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
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
