package tools

import (
	"context"
	"fmt"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/gerrit"
	"github.com/mark3labs/mcp-go/mcp"
)

// PostCommentReplyTool is the tool definition for gerrit_post_comment_reply
var PostCommentReplyTool = mcp.NewTool("gerrit_post_comment_reply",
	mcp.WithDescription("Post a reply to a specific comment on a Gerrit change. This can be used to mark comments as resolved or provide updates."),
	mcp.WithString("changeId",
		mcp.Required(),
		mcp.Description("The Gerrit Change-Id (e.g., I1234567890abcdef...)"),
	),
	mcp.WithString("commentId",
		mcp.Required(),
		mcp.Description("The ID of the comment to reply to"),
	),
	mcp.WithString("message",
		mcp.Required(),
		mcp.Description("The reply message"),
	),
	mcp.WithString("directory",
		mcp.Description("The directory containing the git repository (used to determine Gerrit host)"),
	),
)

// HandlePostCommentReply handles the gerrit_post_comment_reply tool call
func HandlePostCommentReply(cfg *config.Config) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		changeID, err := request.RequireString("changeId")
		if err != nil {
			return mcp.NewToolResultError("changeId is required"), nil
		}

		commentID, err := request.RequireString("commentId")
		if err != nil {
			return mcp.NewToolResultError("commentId is required"), nil
		}

		message, err := request.RequireString("message")
		if err != nil {
			return mcp.NewToolResultError("message is required"), nil
		}

		directory := request.GetString("directory", "")
		client, err := gerrit.NewClientFromGit(directory, cfg.GerritUsername, cfg.GerritPassword)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		if err := client.PostCommentReply(changeID, commentID, message); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}

		return mcp.NewToolResultText("Reply posted successfully."), nil
	}
}
