package tools

import (
	"context"
	"fmt"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/gerrit"
	"github.com/mark3labs/mcp-go/mcp"
)

// PublishReviewTool is the tool definition for publish_review
var PublishReviewTool = mcp.NewTool("publish_review",
	mcp.WithDescription("Submit and publish a review for a Gerrit change. This publishes all draft comments (making them visible to others) and optionally includes a review message."),
	mcp.WithString("changeId",
		mcp.Description("The Gerrit Change-Id (e.g., I1234567890abcdef...). If omitted, will auto-detect from current git commit."),
	),
	mcp.WithString("message",
		mcp.Description("Optional review message to include with the published comments"),
	),
	mcp.WithString("directory",
		mcp.Description("The directory containing the git repository (used to determine Gerrit host)"),
	),
)

// HandlePublishReview handles the publish_review tool call
func HandlePublishReview(cfg *config.Config) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		input := gerrit.PublishReviewInput{
			Message: request.GetString("message", ""),
		}

		if err := client.PublishReview(changeID, input); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText("Success."), nil
	}
}
