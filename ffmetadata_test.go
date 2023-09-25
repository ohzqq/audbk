package audbk

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/ini.v1"
)

var testFFmetaWithChaps = `testdata/ffmeta-with-chap.ini`

func TestLoadFFmeta(t *testing.T) {
	meta := reloadFFmeta(t)
	for _, m := range meta {
		//var ff FFMeta
		//mapstructure.Decode(m.Fields, &ff)
		fmt.Printf("test load %#v\n", m)
	}
}

func TestDumpFFmeta(t *testing.T) {
	meta := reloadFFmeta(t)
	for _, m := range meta {
		ini.PrettyFormat = false

		opts := ini.LoadOptions{
			IgnoreInlineComment:    true,
			AllowNonUniqueSections: true,
		}

		ffmeta := ini.Empty(opts)
		err := ini.ReflectFrom(ffmeta, &m)

		for _, chapter := range m.Chapters {
			sec, err := ffmeta.NewSection("CHAPTER")
			if err != nil {
				//return []byte{}, err
				t.Error(err)
			}
			sec.ReflectFrom(&chapter)
			//for ck, cv := range chapter {
			//  if ck == "start" || ck == "end" || ck == "timebase" {
			//    ck = strings.ToUpper(ck)
			//  }
			//  sec.NewKey(ck, cast.ToString(cv))
			//}

		}

		//d, err := DumpFFMeta(m)
		if err != nil {
			t.Error(err)
		}
		//println(string(d))

		//fmt.Printf("%#v\n", ff)
		_, err = ffmeta.WriteTo(os.Stdout)
		if err != nil {
			t.Error(err)
		}
	}
}

func reloadFFmeta(t *testing.T) []FFMeta {
	files, err := filepath.Glob("testdata/ffmeta*")
	if err != nil {
		t.Error(err)
	}

	var m []FFMeta
	for _, file := range files {
		meta, err := LoadToStruct(file)
		if err != nil && !errors.Is(err, InvalidFFmetadata) {
			t.Error(err)
		}
		m = append(m, meta)
	}
	return m
}
