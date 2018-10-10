package ext

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/mkideal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	type argT struct {
		When Time `cli:"w" dft:"2016-01-02"`
	}
	for _, tt := range []struct {
		src string
		t   time.Time
	}{
		{"-w2016-01-02T15:04:05+00:00", time.Date(2016, 1, 2, 15, 4, 5, 0, time.UTC)},
	} {
		args := []string{}
		if tt.src != "" {
			args = append(args, tt.src)
		}
		argv := new(argT)
		assert.Nil(t, cli.Parse(args, argv))
		assert.Equal(t, tt.t.Unix(), argv.When.Unix())
	}
}

func TestDuration(t *testing.T) {
	type argT struct {
		Long Duration `cli:"d" dft:"1s"`
	}
	for _, tt := range []struct {
		args []string
		long time.Duration
	}{
		{[]string{}, time.Second},
		{[]string{"-d10s"}, time.Second * 10},
		{[]string{"-d10ms"}, time.Millisecond * 10},
	} {
		argv := new(argT)
		assert.Nil(t, cli.Parse(tt.args, argv))
		assert.Equal(t, tt.long, argv.Long.Duration)
	}
}

func TestFile(t *testing.T) {
	type argT struct {
		File File `cli:"f"`
	}
	filename := "yXLLBhNHkv9VdAarIF87"
	content := "hello,world"
	require.Nil(t, ioutil.WriteFile(filename, []byte(content), 0644))
	defer os.Remove(filename)
	argv := new(argT)
	assert.Nil(t, cli.Parse([]string{"-f", filename}, argv))
	assert.Equal(t, content, argv.File.String())
}

func TestReader(t *testing.T) {
	type argT struct {
		Reader Reader `cli:"r"`
	}

	// read from file
	filename := "yXLLBhNHkv9VdAarIF87"
	content := "hello,world"
	require.Nil(t, ioutil.WriteFile(filename, []byte(content), 0644))
	defer os.Remove(filename)
	argv := new(argT)
	assert.Nil(t, cli.Parse([]string{"-r", filename}, argv))
	data, err := ioutil.ReadAll(argv.Reader)
	require.Nil(t, err)
	assert.Equal(t, string(data), content)
	assert.Nil(t, argv.Reader.Close())

	// read from reader
	content = "dlrow,olleh"
	r := bytes.NewBufferString(content)
	argv.Reader.SetReader(r)
	data, err = ioutil.ReadAll(argv.Reader)
	require.Nil(t, err)
	assert.Equal(t, string(data), content)
	assert.Nil(t, argv.Reader.Close())
}

func TestWriter(t *testing.T) {
	type argT struct {
		Writer Writer `cli:"w"`
	}
	// write to file
	filename := "yXLLBhNHkv9VdAarIF87"
	content := "hello,world"
	argv := new(argT)
	require.Nil(t, cli.Parse([]string{"-w", filename}, argv))
	n, err := argv.Writer.Write([]byte(content))
	defer os.Remove(filename)
	assert.Nil(t, err)
	assert.Equal(t, n, len(content))
	assert.Nil(t, argv.Writer.Close())
	data, err := ioutil.ReadFile(filename)
	require.Nil(t, err)
	assert.Equal(t, string(data), content)

	// write to writer
	w := bytes.NewBufferString("")
	argv.Writer.SetWriter(w)
	n, err = argv.Writer.Write([]byte(content))
	assert.Nil(t, err)
	assert.Equal(t, n, len(content))
	assert.Nil(t, argv.Writer.Close())
	assert.Equal(t, w.String(), content)
}
