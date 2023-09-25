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

const FFmetaHeader = ";FFMETADATA1\n"

const (
	MediaArtist = "artist"
	AlbumArtist = "album_artist"
	Composer    = "composer"
	Album       = "album"
	Date        = "date"
	Genre       = "genre"
	Comment     = "comment"
	Grouping    = "grouping"
	Disc        = "disc"
)

type FFMeta struct {
	Album       string          `mapstructure:"title,omitempty" ini:"album,omitempty"`
	Title       string          `mapstructure:"title" ini:"title"`
	Artist      string          `mapstructure:"authors,omitempty" ini:"artist,omitempty"`
	AlbumArtist string          `mapstructure:"authors,omitempty" ini:"album_artist,omitempty"`
	Composer    string          `mapstructure:"narrators,omitempty" ini:"composer,omitempty"`
	Grouping    string          `mapstructure:"series,omitempty" ini:"grouping,omitempty"`
	Disc        float64         `mapstructure:"series_index,omitempty" ini:"disc,omitempty"`
	Genre       string          `mapstructure:"tags,omitempty" ini:"genre,omitempty"`
	Date        string          `mapstructure:"pubdate,omitempty" ini:"date,omitempty"`
	Comment     string          `mapstructure:"comments,omitempty" ini:"comment,omitempty"`
	Other       map[string]any  `mapstructure:",remain" ini:"-"`
	Chapters    []FFMetaChapter `mapstructure:"chapters,omitempty" ini:"-"`
}

type FFMetaChapter struct {
	Timebase string         `ini:"TIMEBASE" mapstructure:"timebase,omitempty"`
	Start    int            `ini:"START" mapstructure:"start,omitempty"`
	End      int            `ini:"END" mapstructure:"end,omitempty"`
	Title    string         `ini:"title,omitempty" mapstructure:"title,omitempty"`
	Fields   map[string]any `mapstructure:",remain" ini:"-"`
}

func NewFFMeta() *FFMeta {
	return &FFMeta{}
}

func NewChapter() *FFMetaChapter {
	return &FFMetaChapter{
		Fields: make(map[string]any),
	}
}

var InvalidFFmetadata = errors.New("ffmetadata file is not valid")

func LoadToStruct(input string) (FFMeta, error) {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	var ff FFMeta

	if !IsValidFFMetadata(input) {
		return ff, InvalidFFmetadata
	}

	f, err := ini.LoadSources(opts, input)
	if err != nil {
		return ff, err
	}

	terr := f.MapTo(&ff)
	if terr != nil {
		log.Fatal(terr)
	}

	if f.HasSection("chapter") {
		sections, err := f.SectionsByName("chapter")
		if err != nil {
			return ff, err
		}

		for _, sec := range sections {
			var ch FFMetaChapter
			sec.MapTo(&ch)
			ff.Chapters = append(ff.Chapters, ch)
		}
	}
	return ff, nil
}

func audioFieldToBookField(k string) string {
	return ""
}

func DumpFFMeta(meta *Meta) ([]byte, error) {
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
