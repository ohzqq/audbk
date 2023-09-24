package audbk

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cast"
	"gopkg.in/ini.v1"
)

type FFMeta struct {
	Name     string
	Fields   map[string]string
	Chapters []*Chapter
}

const FFmetaHeader = ";FFMETADATA1\n"

var InvalidFFmetadata = errors.New("ffmetadata file is not valid")

func NewFFMeta(name string) *FFMeta {
	return &FFMeta{
		Name:   name,
		Fields: make(map[string]string),
	}
}

func LoadFFMeta(input string) (*FFMeta, error) {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	meta := make(map[string]any)
	ffmeta := NewFFMeta(input)

	if !IsValidFFMetadata(input) {
		return ffmeta, InvalidFFmetadata
	}

	f, err := ini.LoadSources(opts, input)
	if err != nil {
		return ffmeta, err
	}

	for f, v := range f.Section("").KeysHash() {
		meta[f] = v
		ffmeta.Fields[f] = v
	}

	if f.HasSection("chapter") {
		sections, err := f.SectionsByName("chapter")
		if err != nil {
			return ffmeta, err
		}

		for _, sec := range sections {
			c := NewChapter()
			sec.MapTo(&c)
			for _, k := range sec.KeyStrings() {
				switch k {
				case "start", "end", "title", "timebase":
				default:
					c.Fields[k] = sec.Key(k).Value()
				}
			}
			ffmeta.Chapters = append(ffmeta.Chapters, c)
		}
	}

	return ffmeta, nil
}

func DumpFFMeta(meta *FFMeta) ([]byte, error) {
	if len(meta.Fields) < 1 {
		return []byte{}, errors.New("no metadata")
	}
	ini.PrettyFormat = false

	opts := ini.LoadOptions{
		IgnoreInlineComment:    true,
		AllowNonUniqueSections: true,
	}

	ffmeta := ini.Empty(opts)

	for k, v := range meta.Fields {
		_, err := ffmeta.Section("").NewKey(k, v)
		if err != nil {
			return []byte{}, err
		}
	}

	for _, ch := range meta.Chapters {
		fmt.Printf("%#v\n", ch)
		sec, err := ffmeta.NewSection("CHAPTER")
		if err != nil {
			return []byte{}, err
		}
		if ch.Timebase != "" {
			sec.NewKey("TIMEBASE", ch.Timebase)
		}
		if ch.Start != 0 {
			sec.NewKey("START", cast.ToString(ch.Start))
		}
		if ch.End != 0 {
			sec.NewKey("END", cast.ToString(ch.End))
		}
		if ch.Title != "" {
			sec.NewKey("title", ch.Title)
		}
		for k, v := range ch.Fields {
			sec.NewKey(k, v)
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
