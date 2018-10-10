package ext

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Time wrap time.Time
type Time struct {
	time.Time
	isSet bool
}

var timeFormats = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	"2006-01-02",
	"2006/01/02",
	"2006:01:02",
	"15:04:05",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"2006:01:02 15:04:05",
	"15:04:05 2006-01-02",
	"15:04:05 2006/01/02",
	"15:04:05 2006:01:02",
}

func (t *Time) Decode(s string) error {
	now := time.Now()
	if s == "" {
		t.Time = now
		return nil
	}
	for _, format := range timeFormats {
		v, err := time.Parse(format, s)
		if err == nil {
			newYear, newMonth, newDay := v.Year(), v.Month(), v.Day()
			reset := false
			if newYear == 0 {
				reset = true
				newYear = now.Year()
				newMonth = now.Month()
				newDay = now.Day()
			}
			if reset {
				v = time.Date(newYear, newMonth, newDay, v.Hour(), v.Minute(), v.Second(), v.Nanosecond(), v.Location())
			}
			t.Time = v
			t.isSet = true
			return nil
		}
	}
	return fmt.Errorf("unsupported time format")
}

func (t Time) Encode() string {
	return t.Format(time.RFC3339Nano)
}

func (t Time) IsSet() bool {
	return t.isSet
}

// Duration wrap time.Duration
type Duration struct {
	time.Duration
}

func (d *Duration) Decode(s string) error {
	if i, err := strconv.ParseUint(s, 10, 64); err == nil {
		d.Duration = time.Duration(i) * time.Second
		return nil
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = v
	return nil
}

func (d Duration) Encode() string {
	return d.Duration.String()
}

// File reads data from file or stdin(if filename is empty)
type File struct {
	filename string
	data     []byte
}

func (f File) Data() []byte {
	return f.data
}

func (f File) String() string {
	if f.data == nil {
		return ""
	}
	return string(f.data)
}

func (f *File) Decode(s string) error {
	var (
		data []byte
		err  error
	)
	if s == "" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(s)
	}
	if err != nil {
		return err
	}
	f.data = data
	f.filename = s
	return nil
}

func (f File) Encode() string {
	return f.filename
}

// Reader
type Reader struct {
	reader   io.Reader
	filename string
}

func (r *Reader) Decode(s string) error {
	if s == "" {
		r.reader = os.Stdin
		r.filename = os.Stdin.Name()
	} else if regexp.MustCompile(`(?i)^http`).MatchString(s) {
		r.filename = s
		response, err := http.Get(s)
		if err != nil {
			return err
		}
		r.reader = response.Body
	} else {
		r.filename = s
		file, err := os.Open(s)
		if err != nil {
			return err
		}
		r.reader = file
	}
	return nil
}

// SetReader replaces the native reader
func (r *Reader) SetReader(reader io.Reader) {
	r.Close()
	r.reader = reader
	if file, ok := reader.(*os.File); ok {
		r.filename = file.Name()
	} else {
		r.filename = ""
	}
}

// Read implementes io.Reader
func (r Reader) Read(data []byte) (n int, err error) {
	if r.reader == nil {
		return os.Stdin.Read(data)
	}
	return r.reader.Read(data)
}

func (r Reader) Close() error {
	if r.reader != nil && !r.IsStdin() {
		if c, ok := r.reader.(io.Closer); ok {
			return c.Close()
		}
	}
	return nil
}

func (r Reader) Name() string {
	if r.reader == nil {
		return os.Stdin.Name()
	}
	return r.filename
}

func (r Reader) IsStdin() bool {
	if r.reader == nil {
		return true
	}
	if stdin, ok := r.reader.(*os.File); ok {
		return uintptr(unsafe.Pointer(stdin)) == uintptr(unsafe.Pointer(os.Stdin))
	}
	return false
}

// Writer
type Writer struct {
	writer   io.Writer
	filename string
}

func (w *Writer) Decode(s string) error {
	if w.writer != nil {
		return nil
	}
	if s == "" {
		w.writer = os.Stdout
		w.filename = os.Stdout.Name()
		return nil
	}
	w.filename = s
	return nil
}

// SetWriter replaces the native writer
func (w *Writer) SetWriter(writer io.Writer) {
	w.Close()
	w.writer = writer
	if file, ok := w.writer.(*os.File); ok {
		w.filename = file.Name()
	} else {
		w.filename = ""
	}
}

// Write implementes io.Writer interface
func (w *Writer) Write(data []byte) (n int, err error) {
	if w.writer == nil {
		if w.filename == "" {
			w.writer = os.Stdout
			w.filename = os.Stdout.Name()
		} else {
			file, err := os.Create(w.filename)
			if err != nil {
				return 0, err
			}
			w.writer = file
		}
	}
	return w.writer.Write(data)
}

func (w *Writer) Close() error {
	if w.writer != nil && !w.IsStdout() {
		if c, ok := w.writer.(io.Closer); ok {
			return c.Close()
		}
	}
	return nil
}

func (w *Writer) Name() string {
	if w.writer == nil {
		return os.Stdout.Name()
	}
	return w.filename
}

func (w Writer) IsStdout() bool {
	if w.writer == nil {
		return true
	}
	if stdout, ok := w.writer.(*os.File); ok {
		return uintptr(unsafe.Pointer(stdout)) == uintptr(unsafe.Pointer(os.Stdout))
	}
	return false
}

// CSV reads one csv record
type CSVRecord struct {
	raw []string
}

func (d *CSVRecord) Decode(s string) error {
	reader := csv.NewReader(strings.NewReader(s))
	record, err := reader.Read()
	if err != nil {
		return err
	}
	d.raw = record
	return nil
}

func (d CSVRecord) Strings() []string {
	return d.raw
}

func (d CSVRecord) Bools() ([]bool, error) {
	ret := make([]bool, len(d.raw))
	for _, s := range d.raw {
		s = strings.ToLower(s)
		if s == "y" || s == "yes" || s == "true" {
			ret = append(ret, true)
		} else if s == "n" || s == "no" || s == "false" {
			ret = append(ret, false)
		} else {
			v, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("parse %s to boolean fail", s)
			}
			ret = append(ret, v != 0)
		}
	}
	return ret, nil
}

func (d CSVRecord) Ints() ([]int64, error) {
	ret := make([]int64, len(d.raw))
	for _, s := range d.raw {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func (d CSVRecord) Uints() ([]uint64, error) {
	ret := make([]uint64, len(d.raw))
	for _, s := range d.raw {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func (d CSVRecord) Floats() ([]float64, error) {
	ret := make([]float64, len(d.raw))
	for _, s := range d.raw {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}
