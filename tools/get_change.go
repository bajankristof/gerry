package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/gerrit"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetChangeTool is the tool definition for get_change
var GetChangeTool = mcp.NewTool("get_change",
	mcp.WithDescription("Get information about a Gerrit change by its Change-Id. Returns details like project, branch, subject, status, and current revision."),
	mcp.WithString("changeId",
		mcp.Description("The Gerrit Change-Id (e.g., I1234567890abcdef...). If omitted, will auto-detect from current git commit."),
	),
	mcp.WithString("directory",
		mcp.Description("The directory containing the git repository (used to determine Gerrit host)"),
	),
)

// HandleGetChange handles the get_change tool call
func HandleGetChange(cfg *config.Config) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		change, err := client.GetChange(changeID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		changeJSON, err := json.MarshalIndent(change, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText(string(changeJSON)), nil
	}
}
