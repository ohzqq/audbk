package audbk

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/spf13/cast"
	"gopkg.in/ini.v1"
)

type FFMetadata struct {
	Name     string
	Fields   map[string]string
	Chapters map[string]any
}

const FFmetaHeader = ";FFMETADATA1\n"

var InvalidFFmetadata = errors.New("ffmetadata file is not valid")

func LoadFFMeta(input string) (map[string]any, error) {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	ffmeta := make(map[string]any)

	if !IsValidFFMetadata(input) {
		return ffmeta, InvalidFFmetadata
	}

	f, err := ini.LoadSources(opts, input)
	if err != nil {
		return ffmeta, err
	}

	for f, v := range f.Section("").KeysHash() {
		ffmeta[f] = v
	}

	if f.HasSection("chapter") {
		var chaps []map[string]any
		sections, err := f.SectionsByName("chapter")
		if err != nil {
			return ffmeta, err
		}

		for _, sec := range sections {
			chap := make(map[string]any)
			for _, k := range sec.KeyStrings() {
				switch k {
				case "timebase":
					chap[k] = sec.Key(k).Value()
				case "start", "end":
					d, err := sec.Key(k).Int()
					if err != nil {
						d = 0
					}
					chap[k] = d
				default:
					chap[k] = sec.Key(k).Value()
				}
			}
			chaps = append(chaps, chap)
		}

		ffmeta["chapters"] = chaps
	}

	return ffmeta, nil
}

func DumpFFMeta(meta map[string]any) ([]byte, error) {
	if len(meta) < 1 {
		return []byte{}, errors.New("no metadata")
	}
	ini.PrettyFormat = false

	opts := ini.LoadOptions{
		IgnoreInlineComment:    true,
		AllowNonUniqueSections: true,
	}

	ffmeta := ini.Empty(opts)

	for k, v := range meta {
		if k != "chapters" {
			_, err := ffmeta.Section("").NewKey(k, cast.ToString(v))
			if err != nil {
				return []byte{}, err
			}
		}
		if k == "chapters" {
			for _, chapter := range v.([]map[string]any) {
				sec, err := ffmeta.NewSection("CHAPTER")
				if err != nil {
					return []byte{}, err
				}
				for ck, cv := range chapter {
					if ck == "start" || ck == "end" || ck == "timebase" {
						ck = strings.ToUpper(ck)
					}
					sec.NewKey(ck, cast.ToString(cv))
				}
			}
		}
	}

	var buf bytes.Buffer

	_, err := buf.WriteString(FFmetaHeader)
	if err != nil {
		return []byte{}, err
	}

	_, err = ffmeta.WriteTo(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func IsValidFFMetadata(f string) bool {
	contents, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer contents.Close()

	scanner := bufio.NewScanner(contents)
	line := 0
	for scanner.Scan() {
		if line == 0 && scanner.Text() == ";FFMETADATA1" {
			return true
			break
		}
	}
	return false
}
