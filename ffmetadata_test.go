package audbk

import (
	"errors"
	"path/filepath"
	"testing"
)

var testFFmetaWithChaps = `testdata/ffmeta-with-chap.ini`

func TestLoadFFmeta(t *testing.T) {
	loadFFmeta(t)
}

func TestDumpFFmeta(t *testing.T) {
	meta := loadFFmeta(t)
	for _, m := range meta {
		d, err := DumpFFMeta(m)
		if err != nil {
			t.Error(err)
		}
		println(string(d))
	}
}

func loadFFmeta(t *testing.T) []*Meta {
	files, err := filepath.Glob("testdata/ffmeta*")
	if err != nil {
		t.Error(err)
	}

	var m []*Meta
	for _, file := range files {
		meta, err := LoadFFMeta(file)
		if err != nil && !errors.Is(err, InvalidFFmetadata) {
			t.Error(err)
		}
		m = append(m, meta)
	}
	return m
}
