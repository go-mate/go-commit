package tests_test

import (
	"testing"

	"github.com/go-mate/go-commit/internal/tests"
	"github.com/pkg/errors"
)

func TestExpectPanic(t *testing.T) {
	tests.ExpectPanic(t, func() {
		panic(errors.New("expect-panic"))
	})
}
