package config

import (
	"errors"
	"flag"
	"testing"
)

func TestParseAgent(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want AgentConfig
	}{
		{
			name: "defaults",
			args: nil,
			want: AgentConfig{
				Collector: CollectorConfig{PollInterval: 2},
				Sender:    SenderConfig{Addr: "localhost:8080", ReportInterval: 10},
			},
		},
		{
			name: "all flags",
			args: []string{"-a", "example.com:9090", "-p", "5", "-r", "30"},
			want: AgentConfig{
				Collector: CollectorConfig{PollInterval: 5},
				Sender:    SenderConfig{Addr: "example.com:9090", ReportInterval: 30},
			},
		},
		{
			name: "partial flags",
			args: []string{"-r=20"},
			want: AgentConfig{
				Collector: CollectorConfig{PollInterval: 2},
				Sender:    SenderConfig{Addr: "localhost:8080", ReportInterval: 20},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAgent(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseAgent_Errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "unknown flag", args: []string{"-x"}},
		{name: "not a number", args: []string{"-p", "abc"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseAgent(tt.args); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestParseServer(t *testing.T) {
	got, err := ParseServer(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Addr != "localhost:8080" {
		t.Errorf("default Addr = %q, want %q", got.Addr, "localhost:8080")
	}

	got, err = ParseServer([]string{"-a", ":9000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Addr != ":9000" {
		t.Errorf("Addr = %q, want %q", got.Addr, ":9000")
	}
}

func TestParseServer_RejectsAgentFlags(t *testing.T) {
	if _, err := ParseServer([]string{"-p", "5"}); err == nil {
		t.Error("expected error for agent-only flag, got nil")
	}
}

func TestParse_Help(t *testing.T) {
	if _, err := ParseAgent([]string{"-h"}); !errors.Is(err, flag.ErrHelp) {
		t.Errorf("ParseAgent(-h) err = %v, want flag.ErrHelp", err)
	}
	if _, err := ParseServer([]string{"-h"}); !errors.Is(err, flag.ErrHelp) {
		t.Errorf("ParseServer(-h) err = %v, want flag.ErrHelp", err)
	}
}
