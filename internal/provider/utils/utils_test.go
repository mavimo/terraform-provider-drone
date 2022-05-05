package utils

import (
	"testing"
)

func TestParseOrgId(t *testing.T) {
	for _, test := range []struct {
		name, str, organization, id string
		is_error                    bool
	}{
		{"Test valid organization", "octocat/hello-world", "octocat", "hello-world", false},
		{"Test another valid organization", "drone/drone", "drone", "drone", false},
		{"Test an organization with many slashes", "foo/bar/baz", "foo", "bar/baz", false},
		{"Test an in valid organization without ID", "drone", "", "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			organization, id, err := ParseOrgId(test.str, "")

			if (test.is_error == true) && (err == nil) {
				t.Errorf("expected error")
			}

			if (test.is_error == false) && (err != nil) {
				t.Errorf("unexpected error")
			}

			if test.organization != organization {
				t.Errorf("unexpected organization")
			}

			if test.id != id {
				t.Errorf("unexpected repo")
			}
		})
	}
}

func TestParseRepo(t *testing.T) {
	for _, test := range []struct {
		name, str, user, repo string
		is_error              bool
	}{
		{"Test valid repository", "octocat/hello-world", "octocat", "hello-world", false},
		{"Test another valid repository", "drone/drone", "drone", "drone", false},
		{"Test invalid repository without slash", "foobar", "", "", true},
		{"Test invalid repository with too many slashes", "foo/bar/baz", "", "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			user, repo, err := ParseRepo(test.str)

			if (test.is_error == true) && (err == nil) {
				t.Errorf("expected error")
			}

			if (test.is_error == false) && (err != nil) {
				t.Errorf("unexpected error")
			}

			if test.user != user {
				t.Errorf("unexpected user")
			}

			if test.repo != repo {
				t.Errorf("unexpected repo")
			}
		})
	}
}

func TestParseId(t *testing.T) {
	for _, test := range []struct {
		name, str, user, repo, id string
		is_error                  bool
	}{
		{"Test valid ID", "octocat/hello-world/resource-1", "octocat", "hello-world", "resource-1", false},
		{"Test another valid ID", "drone/drone/resource-1/others-2", "drone", "drone", "resource-1/others-2", false},
		{"Test invalid ID without slash", "foobar/demo", "", "", "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			user, repo, id, err := ParseId(test.str, "")

			if (test.is_error == true) && (err == nil) {
				t.Errorf("expected error")
			}

			if (test.is_error == false) && (err != nil) {
				t.Errorf("unexpected error")
			}

			if test.user != user {
				t.Errorf("unexpected user")
			}

			if test.repo != repo {
				t.Errorf("unexpected repo")
			}

			if test.id != id {
				t.Errorf("unexpected id")
			}
		})
	}
}
