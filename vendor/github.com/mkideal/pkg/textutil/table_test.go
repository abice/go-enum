package textutil

import (
	"bytes"
	"testing"
)

func TestWriteTable(t *testing.T) {
	table := StringMatrix([][]string{
		[]string{"A", "B", "C"},
		[]string{"AA", "BB", "CC"},
		[]string{"AAA", "BBB", "CCC"},
	})
	expectedStr := `+-----+-----+-----+
| A   | B   | C   |
+-----+-----+-----+
| AA  | BB  | CC  |
+-----+-----+-----+
| AAA | BBB | CCC |
+-----+-----+-----+`
	buf := bytes.NewBufferString("")
	WriteTable(buf, table)
	if got := buf.String(); got != expectedStr {
		t.Errorf("want:\n`%s`\ngot:\n`%s`\n", expectedStr, got)
	}

}
