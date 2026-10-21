package generate_test

import (
	"strings"
	"testing"

	"github.com/alexeyco/helmfile/internal"
	"github.com/alexeyco/helmfile/internal/generate"
	"github.com/spf13/afero"
)

func testVersions() internal.Versions {
	return internal.Versions{
		Alpine:   "3.24",
		Kubectl:  "1.36.4",
		Helmfile: "1.7.4",
		Node:     "26.8.1",
	}
}

func TestRender(t *testing.T) {
	g := generate.New(afero.NewMemMapFs())

	got, err := g.Render(testVersions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{
		"FROM alpine:3.24 AS builder",
		"KUBECTL_VERSION=1.36.4",
		"ARG HELMFILE_VERSION=1.7.4",
		"ghcr.io/helmfile/helmfile:v${HELMFILE_VERSION}",
		"ARG NODE_VERSION=26.8.1",
		"node:${NODE_VERSION}-alpine",
		"COPY --from=node",
		"apk add --no-cache libatomic libstdc++",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestGenerate(t *testing.T) {
	fs := afero.NewMemMapFs()
	g := generate.New(fs)

	if err := g.Generate(testVersions(), "Dockerfile"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := afero.ReadFile(fs, "Dockerfile")
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	if !strings.Contains(string(content), "FROM alpine:3.24 AS builder") {
		t.Errorf("generated file missing alpine base")
	}
}
