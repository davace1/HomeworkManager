package MyLog

import "testing"

func TestClose(t *testing.T) {
	err:=close()
	if err!=nil {
		t.Fatal(err)
	}
}
