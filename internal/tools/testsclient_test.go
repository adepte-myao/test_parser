package tools_test

import (
	"testing"

	"github.com/adepte-myao/test_parser/internal/tools"
)

func TestClientHasJar(t *testing.T) {
	client, err := tools.NewTestsClient()
	if err != nil {
		t.Fatal(err)
	}
	if client.Jar == nil {
		t.Fatal("Client does not have a Jar")
	}
}
