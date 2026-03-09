## MODIFIED Requirements
### Requirement: Resume API endpoint
The API MUST provide an endpoint to trigger session resume.

#### Scenario: Resume uses stored agent
- **WHEN** the client sends `POST /api/resume/:sha`
- **AND** the stored conversation was captured from a supported non-Claude agent
- **THEN** the server restores the session with that agent
- **AND** launches that agent's resume command instead of always launching Claude

### Requirement: Conversation viewer
The web UI MUST display conversation content in a readable format.

#### Scenario: Search stored conversations in the web UI
- **WHEN** the user submits a search query in the web UI
- **THEN** the UI calls a search API endpoint
- **AND** displays matching commits with match previews
- **AND** allows selecting a result to open its conversation

### Requirement: Git graph visualization
The web UI MUST display the commit graph visually.

#### Scenario: Load more branch history
- **WHEN** the user requests more branch history in the overview UI
- **THEN** the client fetches more commits per branch
- **AND** redraws the overview with the larger history window

#### Scenario: Branch counts include full branch history
- **WHEN** the client requests branch summaries
- **THEN** each branch's conversation count reflects all reachable commits on that branch
- **AND** is not truncated to a fixed recent-commit window
