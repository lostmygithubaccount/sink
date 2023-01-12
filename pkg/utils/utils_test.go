package utils

import (
	"testing"

	"github.com/cli/go-gh"
)

func TestIssue(t *testing.T) {

	client, err := gh.RESTClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: how to do this?
	// set the org to the user
	org := GetUser(client)
	source := "sink_test_source"
	targets := []string{"sink_test_target"}

	t.Log("Creating issue")
	t.Log("Org: ", org)
	t.Log("Source: ", source)
	t.Log("Targets: ", targets)
}

func TestMockFail(t *testing.T) {
	t.Error("a test error in the call")
}
