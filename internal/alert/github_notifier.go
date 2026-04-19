package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/eliziario/vaultwatch/internal/vault"
)

type GitHubNotifier struct {
	w io.Writer
}

func NewGitHubNotifier(w io.Writer) *GitHubNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &GitHubNotifier{w: w}
}

func (n *GitHubNotifier) Notify(teams []vault.GitHubTeam) error {
	if len(teams) == 0 {
		fmt.Fprintln(n.w, "[github] no teams mapped")
		return nil
	}
	fmt.Fprintf(n.w, "[github] %d team mapping(s) found:\n", len(teams))
	for _, t := range teams {
		fmt.Fprintf(n.w, "  team=%-30s policy=%s\n", t.Name, t.Policy)
	}
	return nil
}
