package validator

import (
	"testing"
	"unicode/utf8"
)

func TestUTF8CursorNormalization(t *testing.T) {
	line := "α β γ ★"
	bytecol := 9
	sub := line[:bytecol]
	charcol := utf8.RuneCountInString(sub)
	want := 6
	if charcol != want {
		t.Errorf("bytecol=%d -> charcol=%d, want %d", bytecol, charcol, want)
	}
}
