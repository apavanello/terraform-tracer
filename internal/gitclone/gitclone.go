package gitclone

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsGitSource returns true if the module source points to a git repository.
func IsGitSource(source string) bool {
	if source == "" {
		return false
	}

	// Explicit git:: prefix
	if strings.HasPrefix(source, "git::") {
		return true
	}

	// GitHub/GitLab/Bitbucket shorthand (e.g. github.com/org/repo)
	for _, host := range []string{"github.com/", "gitlab.com/", "bitbucket.org/"} {
		if strings.HasPrefix(source, host) {
			return true
		}
	}

	// Generic HTTPS ending with .git
	if strings.HasPrefix(source, "https://") && strings.Contains(source, ".git") {
		return true
	}

	return false
}

// CloneModule clones the git module into cacheDir and returns the local path
// to the module directory (accounting for subdirectories specified with //).
func CloneModule(source string, cacheDir string) (string, error) {
	repoURL, subDir, ref := parseSource(source)

	// Create a deterministic cache directory name from the repo URL + ref
	cacheKey := repoURL
	if ref != "" {
		cacheKey += "@" + ref
	}
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(cacheKey)))[:12]
	cloneDir := filepath.Join(cacheDir, hash)

	// Skip clone if already cached
	if _, err := os.Stat(cloneDir); err == nil {
		localPath := cloneDir
		if subDir != "" {
			localPath = filepath.Join(cloneDir, subDir)
		}
		return localPath, nil
	}

	// Build git clone command
	args := []string{"clone", "--depth", "1"}
	if ref != "" {
		args = append(args, "--branch", ref)
	}
	args = append(args, repoURL, cloneDir)

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	fmt.Fprintf(os.Stderr, "📦 Cloning module: %s\n", repoURL)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failed for %s: %w", repoURL, err)
	}

	localPath := cloneDir
	if subDir != "" {
		localPath = filepath.Join(cloneDir, subDir)
	}

	// Verify the path exists
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return "", fmt.Errorf("subdirectory %q not found in cloned repo %s", subDir, repoURL)
	}

	return localPath, nil
}

// parseSource extracts the repo URL, subdirectory, and ref from a Terraform module source.
//
// Examples:
//
//	"git::https://github.com/org/repo.git//modules/vpc?ref=v1.0" → ("https://github.com/org/repo.git", "modules/vpc", "v1.0")
//	"github.com/org/repo//modules/vpc"                           → ("https://github.com/org/repo.git", "modules/vpc", "")
//	"https://github.com/org/repo.git"                            → ("https://github.com/org/repo.git", "", "")
func parseSource(source string) (repoURL, subDir, ref string) {
	// Strip git:: prefix
	source = strings.TrimPrefix(source, "git::")

	// Extract ?ref=xxx query parameter
	if idx := strings.Index(source, "?"); idx != -1 {
		query := source[idx+1:]
		source = source[:idx]
		if params, err := url.ParseQuery(query); err == nil {
			ref = params.Get("ref")
		}
	}

	// Extract //subdir — but skip the protocol's :// (e.g., https://)
	searchStart := 0
	if schemeEnd := strings.Index(source, "://"); schemeEnd != -1 {
		searchStart = schemeEnd + 3
	}
	if idx := strings.Index(source[searchStart:], "//"); idx != -1 {
		realIdx := searchStart + idx
		subDir = source[realIdx+2:]
		source = source[:realIdx]
	}

	// Convert GitHub/GitLab/Bitbucket shorthand to HTTPS
	for _, host := range []string{"github.com/", "gitlab.com/", "bitbucket.org/"} {
		if strings.HasPrefix(source, host) {
			source = "https://" + source
			if !strings.HasSuffix(source, ".git") {
				source += ".git"
			}
			break
		}
	}

	repoURL = source
	return
}
