package main

import (
	"errors"
	"os"
	"testing"
)

// fakeProc is a test helper implementing processSignaler.
type fakeProc struct {
	shouldErr  bool
	lastSignal os.Signal
}

func (f *fakeProc) Signal(sig os.Signal) error {
	f.lastSignal = sig
	if f.shouldErr {
		return errors.New("signal failed")
	}
	return nil
}

func TestSendSealSignalWithFinder_Success(t *testing.T) {
	fake := &fakeProc{shouldErr: false}
	find := func(pid int) (processSignaler, error) { return fake, nil }

	if err := SendSealSignalWithFinder(find, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.lastSignal == nil {
		t.Fatalf("expected signal to be set")
	}
}

func TestSendSealSignalWithFinder_FinderError(t *testing.T) {
	find := func(pid int) (processSignaler, error) { return nil, errors.New("not found") }
	if err := SendSealSignalWithFinder(find, 1); err == nil {
		t.Fatalf("expected error when finder fails")
	}
}

func TestSendSealSignalWithFinder_SignalError(t *testing.T) {
	fake := &fakeProc{shouldErr: true}
	find := func(pid int) (processSignaler, error) { return fake, nil }
	if err := SendSealSignalWithFinder(find, 1); err == nil {
		t.Fatalf("expected error when Signal fails")
	}
}
