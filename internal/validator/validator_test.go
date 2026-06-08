package validator

import (
	"testing"

	"github.com/divijg19/Praxis/internal/content"
)

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

func TestAllChallengesHaveValidVerify(t *testing.T) {
	for _, c := range content.All() {
		if c.Verify == "" {
			t.Errorf("challenge %s has empty Verify", c.ID)
		}
		if !Exists(c.Verify) {
			t.Errorf("challenge %s has unknown Verify: %s", c.ID, c.Verify)
		}
	}
}

func TestExistsBuffer(t *testing.T) {
	if !Exists("buffer") {
		t.Error("Exists(\"buffer\") should be true")
	}
}

func TestBufferChallengesHaveResult(t *testing.T) {
	for _, c := range content.All() {
		if c.Verify == "buffer" {
			if len(c.Result) == 0 {
				t.Errorf("buffer challenge %s has empty Result", c.ID)
			}
		}
	}
}
