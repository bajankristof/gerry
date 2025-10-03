# Gerry

MCP server for interacting with Gerrit Code Review from Claude Code.

## Prerequisites

- Go 1.25.1 or later
- Access to a Gerrit instance
- Gerrit HTTP credentials (username and password)

## Installation

```bash
make build
sudo make link
```

## Configuration

Create a configuration file at `~/.config/gerry.json`:

```json
{
  "gerritUsername": "your-username",
  "gerritPassword": "your-http-password"
}
```

To get your Gerrit HTTP password:
1. Go to your Gerrit instance
2. Navigate to Settings → HTTP Credentials
3. Generate a new password if needed

## Adding to Claude Code

Run `claude mcp add gerry gerry` to add Gerry to your Claude Code instance.

## Available Tools

- **gerrit_get_change_id** - Get the Change-Id from the current git repository
- **gerrit_get_change** - Get detailed information about a Gerrit change
- **gerrit_get_comments** - Get all comments for a change
- **gerrit_get_unresolved_comments** - Get only unresolved comments for a change
- **gerrit_post_comment_reply** - Post a reply to a comment

## Usage Example

After setting up, you can ask Claude Code:

> "What are the unresolved comments on change I1234567890abcdef?"

> "Get the change information for the current repository"

> "Reply to comment abc123 with 'Fixed in latest patch set'"
