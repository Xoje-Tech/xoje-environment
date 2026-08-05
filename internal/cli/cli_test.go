package cli

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantCmd     string
		wantTool    string
		wantErr     bool
	}{
		{"No args - default to tui", []string{}, "tui", "", false},
		{"Diagnose", []string{"diagnose"}, "diagnose", "", false},
		{"Bootstrap alias", []string{"bootstrap"}, "diagnose", "", false},
		{"Install with tool", []string{"install", "gentle-ai"}, "install", "gentle-ai", false},
		{"Install without tool - error", []string{"install"}, "", "", true},
		{"Unknown command - error", []string{"frobnicate"}, "", "", true},
		{"Explicit tui", []string{"tui"}, "tui", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, tool, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if cmd != tt.wantCmd {
				t.Errorf("Parse() cmd = %v, want %v", cmd, tt.wantCmd)
			}
			if tool != tt.wantTool {
				t.Errorf("Parse() tool = %v, want %v", tool, tt.wantTool)
			}
		})
	}
}
