package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/re-cinq/shift-log/internal/cli"
	"github.com/re-cinq/shift-log/internal/git"
	"github.com/re-cinq/shift-log/internal/storage"
	"github.com/spf13/cobra"
)

var (
	listAgent  string
	listBranch string
	listModel  string
	listSince  string
	listLimit  int
	listJSON   bool
)

type listResult struct {
	SHA          string          `json:"sha"`
	Date         string          `json:"date"`
	Message      string          `json:"message"`
	MessageCount int             `json:"message_count"`
	Agent        string          `json:"agent"`
	Branch       string          `json:"branch,omitempty"`
	Model        string          `json:"model,omitempty"`
	Effort       *storage.Effort `json:"effort,omitempty"`
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List commits with stored conversations",
	GroupID: "human",
	Long: `Lists all commits in the repository that have associated
conversations stored as Git Notes.

Shows:
  - Commit SHA (short)
  - Commit date
  - Commit message (truncated)
  - Number of messages in conversation

Example output:
  abc1234 2024-01-15 feat: add user auth (42 messages)
  def5678 2024-01-14 fix: login bug (15 messages)`,
	RunE: runList,
}

func init() {
	listCmd.Flags().StringVar(&listAgent, "agent", "", "filter by agent name")
	listCmd.Flags().StringVar(&listBranch, "branch", "", "filter by git branch")
	listCmd.Flags().StringVar(&listModel, "model", "", "filter by model (substring match)")
	listCmd.Flags().StringVar(&listSince, "since", "", "only include commits on or after this date (YYYY-MM-DD or RFC3339)")
	listCmd.Flags().IntVar(&listLimit, "limit", 0, "max number of results (0 = no limit)")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output results as JSON")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// Verify we're in a git repository
	if err := git.RequireGitRepo(); err != nil {
		return err
	}

	// Get list of commits with notes
	commits, err := git.ListCommitsWithNotes()
	if err != nil {
		return fmt.Errorf("could not list conversations: %w", err)
	}

	var sinceFilter time.Time
	if listSince != "" {
		sinceTime, err := parseSearchDate(listSince)
		if err != nil {
			return fmt.Errorf("invalid --since date: %w", err)
		}
		sinceFilter = sinceTime
	}

	var results []listResult
	for _, commitSHA := range commits {
		// Get commit info
		message, date, err := git.GetCommitInfo(commitSHA)
		if err != nil {
			// Skip commits we can't get info for
			continue
		}

		// Get conversation metadata
		stored, err := storage.GetStoredConversation(commitSHA)
		if err != nil || stored == nil {
			continue
		}

		if !matchListFilters(stored, date, sinceFilter) {
			continue
		}

		agentName := stored.Agent
		if agentName == "" {
			agentName = "claude"
		}

		results = append(results, listResult{
			SHA:          commitSHA,
			Date:         date,
			Message:      message,
			MessageCount: stored.MessageCount,
			Agent:        agentName,
			Branch:       stored.GitBranch,
			Model:        stored.Model,
			Effort:       stored.Effort,
		})

		if listLimit > 0 && len(results) >= listLimit {
			break
		}
	}

	if listJSON {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	if len(results) == 0 {
		fmt.Println("no conversations found")
		return nil
	}

	for _, result := range results {
		// Format date (take just the date part)
		shortDate := result.Date
		if len(shortDate) >= 10 {
			shortDate = shortDate[:10]
		}

		message := result.Message
		if len(message) > 50 {
			message = message[:47] + "..."
		}

		var details []string
		details = append(details, result.Agent)
		if result.Branch != "" {
			details = append(details, result.Branch)
		}
		if result.Model != "" {
			details = append(details, result.Model)
		}
		details = append(details, fmt.Sprintf("%d messages", result.MessageCount))
		if result.Effort != nil && result.Effort.Turns > 0 {
			details = append(details, fmt.Sprintf("%d turns", result.Effort.Turns))
		}

		fmt.Printf("%s %s %s (%s)\n",
			result.SHA[:7],
			shortDate,
			message,
			strings.Join(details, ", "),
		)
	}

	return nil
}

func matchListFilters(stored *storage.StoredConversation, commitDate string, since time.Time) bool {
	if listAgent != "" {
		agentName := stored.Agent
		if agentName == "" {
			agentName = "claude"
		}
		if !strings.EqualFold(agentName, listAgent) {
			return false
		}
	}

	if listBranch != "" && !strings.EqualFold(stored.GitBranch, listBranch) {
		return false
	}

	if listModel != "" && !strings.Contains(strings.ToLower(stored.Model), strings.ToLower(listModel)) {
		return false
	}

	if !since.IsZero() {
		commitTime, err := parseListCommitDate(commitDate)
		if err != nil {
			cli.LogDebug("list: skipping commit with unparsable date %q: %v", commitDate, err)
			return false
		}
		if commitTime.Before(since) {
			return false
		}
	}

	return true
}

func parseListCommitDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05 -0700",
		"2006-01-02",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format %q", s)
}
