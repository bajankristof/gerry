package tools

import (
	"github.com/bajankristof/gerry/config"
	"github.com/mark3labs/mcp-go/server"
)

// Inject registers all Gerrit MCP tools with the server
func Inject(s *server.MCPServer, cfg *config.Config) {
	s.AddTool(GetChangeIDTool, HandleGetChangeID)
	s.AddTool(GetChangeTool, HandleGetChange(cfg))
	s.AddTool(GetCommentsTool, HandleGetComments(cfg))
	s.AddTool(GetUnresolvedCommentsTool, HandleGetUnresolvedComments(cfg))
	s.AddTool(PostCommentReplyTool, HandlePostCommentReply(cfg))
}
