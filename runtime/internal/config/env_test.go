package config

import "testing"

func TestLoadEnvironment_MergesDefaultsAndProcessEnvironment(t *testing.T) {
	t.Setenv("AGENT_NAME", "process-agent")
	t.Setenv("HOME", "/ignored-process-home")
	t.Setenv("USER_HOME", "/process-home")
	t.Setenv("AGENT_HOME", "/process-agent-home")
	t.Setenv("RUNTIME_HOME", "/process-runtime")
	t.Setenv("WORKSPACE", "/process-workspace")
	t.Setenv("LOCAL_UID", "2000")
	t.Setenv("LOCAL_GID", "2001")
	t.Setenv("EXTRA_PROCESS", "process")

	env := LoadEnvironment(map[string]string{
		"AGENT_NAME": "default-agent",
		"AGENT_HOME": "/default-agent-home",
		"USER_HOME":  "/default-home",
		"LOCAL_UID":  "1000",
		"DEFAULT":    "kept",
	})
	if env.AgentName != "process-agent" {
		t.Fatalf("AgentName = %q, want process-agent", env.AgentName)
	}
	if env.UserHome != "/process-home" {
		t.Fatalf("UserHome = %q, want /process-home", env.UserHome)
	}
	if env.AgentHome != "/process-agent-home" {
		t.Fatalf("AgentHome = %q, want /process-agent-home", env.AgentHome)
	}
	if env.RuntimeHome != "/process-runtime" {
		t.Fatalf("RuntimeHome = %q, want /process-runtime", env.RuntimeHome)
	}
	if env.Workspace != "/process-workspace" {
		t.Fatalf("Workspace = %q, want /process-workspace", env.Workspace)
	}
	if env.LocalUID != "2000" || env.LocalGID != "2001" {
		t.Fatalf("uid/gid = %q/%q, want 2000/2001", env.LocalUID, env.LocalGID)
	}
	if got := env.Values["DEFAULT"]; got != "kept" {
		t.Fatalf("DEFAULT = %q, want kept", got)
	}
	if got := env.Values["EXTRA_PROCESS"]; got != "process" {
		t.Fatalf("EXTRA_PROCESS = %q, want process", got)
	}
}

func TestLoadEnvironment_DefaultsOwnership(t *testing.T) {
	t.Setenv("LOCAL_UID", "")
	t.Setenv("LOCAL_GID", "")

	env := LoadEnvironment(nil)
	if env.LocalUID != "1000" || env.LocalGID != "1000" {
		t.Fatalf("uid/gid = %q/%q, want defaults 1000/1000", env.LocalUID, env.LocalGID)
	}
	if env.Values["LOCAL_UID"] != "1000" || env.Values["LOCAL_GID"] != "1000" {
		t.Fatalf("environment values do not contain ownership defaults: %#v", env.Values)
	}
}
