// Package gitutil holds small helpers around the git CLI.
package gitutils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return string(out), nil
}

// Prefix returns dir's path relative to the top of its repository, with forward
// slashes and a trailing slash ("" when dir is the top level).
func Prefix(dir string) (string, error) {
	out, err := run(dir, "rev-parse", "--show-prefix")
	return strings.TrimSpace(out), err
}

// AddWorktree checks ref out into a fresh temporary directory (detached, so the
// caller's working tree, index and branches are never touched) and returns that
// directory's path together with a cleanup function that removes it again.
func AddWorktree(dir, ref string) (path string, cleanup func(), err error) {
	if strings.HasPrefix(ref, "-") {
		return "", nil, fmt.Errorf("invalid git ref %q", ref)
	}
	tmp, err := os.MkdirTemp("", "archlens-worktree-*")
	if err != nil {
		return "", nil, err
	}
	if _, err := run(dir, "worktree", "add", "--detach", tmp, ref); err != nil {
		_ = os.RemoveAll(tmp)
		return "", nil, err
	}
	cleanup = func() {
		_, _ = run(dir, "worktree", "remove", "--force", tmp)
		_ = os.RemoveAll(tmp)
		_, _ = run(dir, "worktree", "prune")
	}
	return tmp, cleanup, nil
}
