package validator

import "testing"

func TestExistsCursor(t *testing.T) {
	if !Exists("cursor") {
		t.Error("Exists(\"cursor\") should be true")
	}
}

func TestExistsUnknown(t *testing.T) {
	if Exists("banana") {
		t.Error("Exists(\"banana\") should be false")
	}
}

func TestExistsBuffer(t *testing.T) {
	if !Exists("buffer") {
		t.Error("Exists(\"buffer\") should be true")
	}
}

func TestExistsComposite(t *testing.T) {
	if !Exists("composite") {
		t.Error("Exists(\"composite\") should be true")
	}
}
