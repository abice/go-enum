package textutil

import (
	"bytes"
	"fmt"
	"io"

	"github.com/labstack/gommon/color"
)

const (
	borderCross = "+"
	borderRow   = "-"
	borderCol   = "|"
)

// Table represents a string-matrix
type Table interface {
	RowCount() int
	ColCount() int
	Get(i, j int) string
}

// TableStyler represents s render style for Table
type TableStyler interface {
	CellRender(row, col int, cell string, w *ColorWriter)
	BorderRender(border string, w *ColorWriter)
}

// ColorWriter
type ColorWriter struct {
	w io.Writer
	*color.Color
}

func (cw *ColorWriter) Write(b []byte) (int, error) {
	return cw.w.Write(b)
}

type colorable interface {
	Color() *color.Color
}

func newColorWriter(w io.Writer) *ColorWriter {
	if cw, ok := w.(colorable); ok {
		ret := &ColorWriter{w: w}
		ret.Color = cw.Color()
		return ret
	}
	ret := &ColorWriter{w: w}
	ret.Color = color.New()
	ret.Disable()
	return ret
}

// DefaultStyle ...
type DefaultStyle struct{}

func (style DefaultStyle) CellRender(row, col int, cell string, w *ColorWriter) {
	fmt.Fprint(w, cell)
}

func (style DefaultStyle) BorderRender(border string, w *ColorWriter) {
	fmt.Fprint(w, w.Grey(border))
}

var defaultStyle = DefaultStyle{}

// WriteTable formats table to writer with specified style
func WriteTable(w io.Writer, table Table, style TableStyler) {
	if style == nil {
		style = defaultStyle
	}
	rowCount, colCount := table.RowCount(), table.ColCount()
	if rowCount <= 0 || colCount <= 0 {
		return
	}
	widthArray := make([]int, colCount)
	for j := 0; j < colCount; j++ {
		maxWidth := 0
		for i := 0; i < rowCount; i++ {
			width := len(table.Get(i, j))
			if i == 0 || width > maxWidth {
				maxWidth = width
			}
		}
		widthArray[j] = maxWidth
	}
	cw := newColorWriter(w)
	rowBorder := rowBorderLine(widthArray)
	style.BorderRender(rowBorder, cw)
	for i := 0; i < rowCount; i++ {
		fmt.Fprint(cw, "\n")
		writeTableRow(cw, table, i, widthArray, style)
		fmt.Fprint(cw, "\n")
		style.BorderRender(rowBorder, cw)
	}
	fmt.Fprint(cw, "\n")
}

func rowBorderLine(widthArray []int) string {
	buf := bytes.NewBufferString(borderCross)
	for _, width := range widthArray {
		repeatWriteString(buf, borderRow, width+2)
		buf.WriteString(borderCross)
	}
	return buf.String()
}

func writeTableRow(cw *ColorWriter, table Table, rowIndex int, widthArray []int, style TableStyler) {
	style.BorderRender(borderCol, cw)
	colCount := table.ColCount()
	for j := 0; j < colCount; j++ {
		fmt.Fprint(cw, " ")
		format := fmt.Sprintf("%%-%ds", widthArray[j]+1)
		s := fmt.Sprintf(format, table.Get(rowIndex, j))
		style.CellRender(rowIndex, j, s, cw)
		style.BorderRender(borderCol, cw)
	}
}

func repeatWriteString(w io.Writer, s string, count int) {
	for i := 0; i < count; i++ {
		fmt.Fprint(w, s)
	}
}

// TableView represents a view of table, it implements Table interface, too
type TableView struct {
	table              Table
	rowIndex, colIndex int
	rowCount, colCount int
}

func (tv TableView) RowCount() int {
	return tv.rowCount
}

func (tv TableView) ColCount() int {
	return tv.colCount
}

func (tv TableView) Get(i, j int) string {
	return tv.table.Get(tv.rowIndex+i, tv.colIndex+j)
}

// ClipTable creates a view of table
func ClipTable(table Table, i, j, m, n int) Table {
	minR, minC := i, j
	maxR, maxC := i+m, j+n
	if minR < 0 || minC < 0 || minR > maxR || minC > maxC || maxR >= table.RowCount() || maxC >= table.ColCount() {
		panic("out of bound")
	}
	return &TableView{table, i, j, m, n}
}

// TableWithHeader add header for table
type TableWithHeader struct {
	table  Table
	header []string
}

func (twh TableWithHeader) RowCount() int { return twh.table.RowCount() + 1 }
func (twh TableWithHeader) ColCount() int { return twh.table.ColCount() }
func (twh TableWithHeader) Get(i, j int) string {
	if i == 0 {
		return twh.header[j]
	}
	return twh.table.Get(i-1, j)
}

func AddTableHeader(table Table, header []string) Table {
	return &TableWithHeader{table, header}
}

// 2-Array string
type StringMatrix [][]string

func (m StringMatrix) RowCount() int { return len(m) }
func (m StringMatrix) ColCount() int {
	if len(m) == 0 {
		return 0
	}
	return len(m[0])
}
func (m StringMatrix) Get(i, j int) string { return m[i][j] }
