package tools

import (
	"fmt"

	"github.com/bajankristof/gerry/config"
	"github.com/bajankristof/gerry/git"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Inject registers all Gerrit MCP tools with the server
func Inject(s *server.MCPServer, cfg *config.Config) {
	s.AddTool(GetChangeIDTool, HandleGetChangeID)
	s.AddTool(GetChangeTool, HandleGetChange(cfg))
	s.AddTool(GetCommentsTool, HandleGetComments(cfg))
	s.AddTool(GetUnresolvedCommentsTool, HandleGetUnresolvedComments(cfg))
	s.AddTool(DraftCommentTool, HandleDraftComment(cfg))
	s.AddTool(PublishReviewTool, HandlePublishReview(cfg))
}

// inferChangeID extracts changeId from the request or auto-detects it from git
func inferChangeID(request mcp.CallToolRequest) (string, error) {
	changeID := request.GetString("changeId", "")

	if changeID != "" {
		return changeID, nil
	}

	// Auto-detect from git
	directory := request.GetString("directory", "")
	changeID, err := git.GetChangeIDFromCommit(directory)
	if err != nil {
		return "", fmt.Errorf("could not auto-detect changeId from git: %w", err)
	}

	if changeID == "" {
		return "", fmt.Errorf("no Change-Id found in current commit")
	}

	return changeID, nil
}
