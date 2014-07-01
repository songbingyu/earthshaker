package earthshaker

import (
	"testing"
)


func TestLogSys(t *testing.T) {
	b := IniLog(DEBUG, true, "test")
	if !b {
		t.Errorf("IniLog Error")
	}

	b = LOG(DEBUG, "aa", 4)
	if !b {
		t.Errorf("LOG Error")
	}

	b = LOG(INFO, "aa", 4)
	if b {
		t.Errorf("LOG Error")
	}

	CloseLog()	

	b = LOG(DEBUG, "aa", 4)
	if b {
		t.Errorf("LOG Error")
	}

	OpenLog()
	b = LOG(ERROR, "aa", 4)
	if !b {
		t.Errorf("LOG Error")
	}
}

