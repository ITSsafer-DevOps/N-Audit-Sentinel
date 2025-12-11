package main

import (
	"errors"
	"testing"
)

// fakeCmdRunner implements cmdRunner for tests.
type fakeCmdRunner struct {
	shouldErr  bool
	cmdsCalled int
}

func (f *fakeCmdRunner) Run(name string, args ...string) error {
	f.cmdsCalled++
	if f.shouldErr {
		return errors.New("command failed")
	}
	return nil
}

func TestBuildBinariesWithRunner_Success(t *testing.T) {
	runner := &fakeCmdRunner{shouldErr: false}
	sentinelBin, cliBin, err := BuildBinariesWithRunner(runner, "/tmp/build")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sentinelBin == "" || cliBin == "" {
		t.Fatalf("expected non-empty binary paths")
	}
	if runner.cmdsCalled != 2 {
		t.Fatalf("expected 2 commands to be run, got %d", runner.cmdsCalled)
	}
}

func TestBuildBinariesWithRunner_BuildFails(t *testing.T) {
	runner := &fakeCmdRunner{shouldErr: true}
	_, _, err := BuildBinariesWithRunner(runner, "/tmp/build")

	if err == nil {
		t.Fatalf("expected error when build fails")
	}
	if runner.cmdsCalled != 1 {
		t.Fatalf("expected 1 command attempt before failure, got %d", runner.cmdsCalled)
	}
}
