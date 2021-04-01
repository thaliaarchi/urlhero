package wwiki

import "testing"

func TestGetDumps(t *testing.T) {
	dumps, err := GetDumps()
	if err != nil {
		t.Fatal(err)
	}
	if len(dumps) == 0 {
		t.Fatalf("no dumps")
	}
}
