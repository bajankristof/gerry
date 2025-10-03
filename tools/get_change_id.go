package tools

import (
	"context"
	"fmt"

	"github.com/bajankristof/gerry/git"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetChangeIDTool is the tool definition for gerrit_get_change_id
var GetChangeIDTool = mcp.NewTool("gerrit_get_change_id",
	mcp.WithDescription("Get the Gerrit Change-Id from the current git commit in the specified directory. Returns the Change-Id (e.g., I1234567890abcdef...) that can be used with other Gerrit tools."),
	mcp.WithString("directory",
		mcp.Description("The directory containing the git repository"),
	),
)

// HandleGetChangeID handles the gerrit_get_change_id tool call
func HandleGetChangeID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	directory := request.GetString("directory", "")

	changeID, err := git.GetChangeIDFromCommit(directory)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
	}

	if changeID == "" {
		return mcp.NewToolResultText("No Change-Id found in the current commit. Make sure you're in a git repository with a Gerrit commit."), nil
	}

	return mcp.NewToolResultText(changeID), nil
}
