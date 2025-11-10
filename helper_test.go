package simpledi_test

import (
	"errors"
	"testing"
)

type testService1 interface{ doSomething1() }
type testServiceImpl1 struct{}

func (t *testServiceImpl1) doSomething1() {}

type testServiceImpl2 struct{ data string }

type testServiceImpl3 struct{ service1 *testServiceImpl1 }

func assertOrder[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	gotCount, wantCount := len(got), len(want)
	if len(got) != len(want) {
		t.Errorf("got: %d count, want: %d count", gotCount, wantCount)
		return
	}
	for i := 0; i < gotCount; i++ {
		if got[i] != want[i] {
			t.Errorf("[%d] got: %v, want: %v", i, got[i], want[i])
			t.Errorf("got: %#v, want: %#v", got, want)
			return
		}
	}
}

func assertSameValue[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func assertSamePointer[T comparable](t *testing.T, got, want *T) {
	t.Helper()
	if got != want {
		t.Errorf("got: %p, want: %p", got, want)
	}
}

func assertError(t *testing.T, fn func() error, wants ...error) {
	t.Helper()
	err := fn()
	for _, want := range wants {
		if !errors.Is(err, want) {
			t.Errorf("got: %v, want: %v", err, want)
		}
	}
}

func assertNoError(t *testing.T, fn func() error) {
	t.Helper()
	if err := fn(); err != nil {
		t.Errorf("got: %v, want: no error", err)
	}
}

func assertPanic(t *testing.T, fn func(), want error) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("got: no panic, want: panic(%v)", want)
			return
		}

		err, ok := r.(error)
		if !ok {
			t.Errorf("got: %T, want: error", r)
			return
		}

		if !errors.Is(err, want) {
			t.Errorf("got: %v, want: %v", err, want)
		}
	}()
	fn()
}

func assertNoPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("got: panic(%v), want: no panic", r)
		}
	}()
	fn()
}
