package audbk

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ohzqq/cdb"
)

var testFFmetaWithChaps = `testdata/ffmeta-with-chap.ini`

func TestLoadFFmeta(t *testing.T) {
	meta := reloadFFmeta(t)
	for _, m := range meta {
		fmt.Printf("test load %#v\n", m)
	}
}

func TestDumpFFmeta(t *testing.T) {
	println("TestDumpFFmeta")
	meta := reloadFFmeta(t)
	for _, m := range meta {
		book := cdb.Book{}
		FFMetaToBook(&book, m)
		ff := NewFFMeta()
		err := BookToFFMeta(ff, book.StringMap())
		if err != nil {
			t.Error(err)
		}

		_, err = m.WriteTo(os.Stdout)
		if err != nil {
			t.Error(err)
		}
	}
}

func reloadFFmeta(t *testing.T) []*FFMeta {
	files, err := filepath.Glob("testdata/ffmeta*")
	if err != nil {
		t.Error(err)
	}

	var m []*FFMeta
	for _, file := range files {
		ff := NewFFMeta()
		meta, err := ff.ReadFile(file)
		if err != nil && !errors.Is(err, InvalidFFmetadata) {
			t.Error(err)
		}
		m = append(m, meta)
	}
	return m
}

func TestParseGrouping(t *testing.T) {
	testStr := []string{
		"Series Title, book 2",
		"Series Title, Book 2",
		"Series Title",
	}

	for _, str := range testStr {
		s, i := parseGrouping(str)
		if s != "Series Title" {
			t.Errorf("expected 'Series Title' got %s\n", s)
		}
		if i != 2 && i != 0 {
			t.Errorf("expected 2 or 0 got %f\n", i)
		}
	}
}
