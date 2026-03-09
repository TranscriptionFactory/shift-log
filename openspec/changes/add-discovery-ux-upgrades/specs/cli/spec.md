## MODIFIED Requirements
### Requirement: List available conversations
Users MUST be able to discover which commits have conversations.

#### Scenario: Filter and serialize listed conversations
- **WHEN** the user runs `shiftlog list` with filter flags or `--json`
- **THEN** the command filters by stored metadata where requested
- **AND** the command can emit structured JSON for downstream tooling

### Requirement: Terminal Conversation Viewer
The CLI SHALL provide a `show` command that displays conversation history for a given commit reference in the terminal.

#### Scenario: Search hands off to show
- **WHEN** the user runs `shiftlog search --show <query>`
- **THEN** the top matching conversation is rendered with the same terminal output as `shiftlog show`

### Requirement: Conversation Output Format
The terminal output SHALL be formatted for readability with clear visual separation between messages.

#### Scenario: Search highlights matching text
- **WHEN** the user runs `shiftlog search <query>`
- **THEN** matching snippets visually highlight the matching text in terminal output

#### Scenario: Search results can be serialized
- **WHEN** the user runs `shiftlog search --json <query>`
- **THEN** the command returns structured JSON search results
