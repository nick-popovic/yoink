package validator

import "testing"

func TestParseResolvesBranchNamesWithSlashes(t *testing.T) {
	originalListRemoteRefs := listRemoteRefs
	listRemoteRefs = func(repoURL string) (map[string]struct{}, error) {
		t.Helper()
		if repoURL != "https://github.com/example/project" {
			t.Fatalf("unexpected repo URL: %s", repoURL)
		}

		return map[string]struct{}{
			"feature/foo": {},
			"main":        {},
		}, nil
	}
	t.Cleanup(func() {
		listRemoteRefs = originalListRemoteRefs
	})

	host, user, repo, branch, subpath, err := Parse("https://github.com/example/project/tree/feature/foo/internal/validator")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if host != "github.com" || user != "example" || repo != "project" {
		t.Fatalf("unexpected repo parts: host=%q user=%q repo=%q", host, user, repo)
	}

	if branch != "feature/foo" {
		t.Fatalf("branch = %q, want %q", branch, "feature/foo")
	}

	if subpath != "internal/validator" {
		t.Fatalf("subpath = %q, want %q", subpath, "internal/validator")
	}
}
