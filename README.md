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

## Usage

The tool will only work within a git repository that has Gerrit commits.

## Available Tools

- **get_change_id** - Get the Change-Id from the current git repository
- **get_change** - Get detailed information about a Gerrit change
- **get_comments** - Get all comments for a change
- **get_unresolved_comments** - Get only unresolved comments for a change
- **draft_comment** - Create a draft comment or reply on a change
- **publish_review** - Publish all draft comments and submit a review

### Automatic Change ID Inference

Most tools support automatic change ID detection. You can omit the `changeId` parameter and the tool will automatically extract it from your current commit. This makes it easier to work with your current change:

```
# Instead of:
"Get unresolved comments for change I1234567890abcdef"

# You can simply say:
"Get unresolved comments for my current change"
```

## Usage Examples

After setting up, you can ask Claude Code:

> "What are the unresolved comments on my current CR?"

> "Draft a reply to the comment on line 42 of main.go saying 'Fixed in latest patch set'"

> "Publish my review with message 'Addressed all feedback'"

> "Get the change information for I1234567890abcdef"
