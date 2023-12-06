package query

import (
	"testing"
)

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name      string
		pathParts []string
		pattern   []string
		want      bool
	}{
		{
			name:      "Direct Parent-Child Relationship",
			pathParts: []string{"parent", "child", "grandchild"},
			pattern:   []string{"parent", "/", "child", "/", "grandchild"},
			want:      true,
		},
		{
			name:      "Wildcard Single Level",
			pathParts: []string{"parent", "random", "child"},
			pattern:   []string{"parent", "/", "*", "/", "child"},
			want:      true,
		},
		{
			name:      "Ancestor-Descendant Relationship",
			pathParts: []string{"ancestor", "intermediate", "descendant"},
			pattern:   []string{"ancestor", "//", "descendant"},
			want:      true,
		},
		{
			name:      "No Match Due to Incorrect Hierarchy",
			pathParts: []string{"parent", "child", "grandchild"},
			pattern:   []string{"parent", "/", "grandchild", "/", "child"},
			want:      false,
		},
		{
			name:      "Empty Pattern",
			pathParts: []string{"parent", "child", "grandchild"},
			pattern:   []string{},
			want:      false,
		},
		{
			name:      "Empty Path",
			pathParts: []string{},
			pattern:   []string{"parent", "/", "child", "/", "grandchild"},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchPattern(tt.pathParts, tt.pattern); got != tt.want {
				t.Errorf("matchPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
