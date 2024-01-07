package util

import "testing"

func TestSeparateTimeHandler(t *testing.T) {
	// shimokitazawa mosaic
	res, err := ParseTime("2023-06-13", "OPEN 18:00")
	if err != nil {
		t.Errorf("SeparateTimeHandler(\"2023-06-13\", \"OPEN 18:00\", \"/\") got error %s", err)
	}
	if res.Unix() != 1686646800 {
		t.Errorf("SeparateTimeHandler(\"2023-06-13\", \"OPEN 18:00\", \"/\") got res %d, want 1686646800", res.Unix())
	}

}
