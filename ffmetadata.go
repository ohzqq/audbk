package audbk

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cast"
	"gopkg.in/ini.v1"
)

const FFmetaHeader = ";FFMETADATA1"

const (
	MediaArtist = "artist"
	AlbumArtist = "album_artist"
	Composer    = "composer"
	Album       = "album"
	MediaTitle  = "title"
	Date        = "date"
	Genre       = "genre"
	Comment     = "comment"
	Grouping    = "grouping"
	Disc        = "disc"
)

var InvalidFFmetadata = errors.New("ffmetadata file is not valid")

type FFMeta struct {
	Title       string          `mapstructure:"title" ini:"title"`
	Album       string          `mapstructure:"title,omitempty" ini:"album,omitempty"`
	Artist      []string        `mapstructure:"authors,omitempty" ini:"artist,omitempty"`
	AlbumArtist []string        `mapstructure:"authors,omitempty" ini:"album_artist,omitempty"`
	Composer    []string        `mapstructure:"narrators,omitempty" ini:"composer,omitempty"`
	Grouping    string          `mapstructure:"series,omitempty" ini:"grouping,omitempty"`
	Disc        float64         `mapstructure:"series_index,omitempty" ini:"disc,omitempty"`
	Genre       []string        `mapstructure:"tags,omitempty" ini:"genre,omitempty"`
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

func (ff *FFMeta) LoadFile(input string) (*FFMeta, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return ff, err
	}

	if !IsValidFFMetadata(data) {
		return ff, InvalidFFmetadata
	}

	file := io.NopCloser(bytes.NewReader(data))

	return LoadFFMeta(ff, file)
}

func FFMetaToBook(book *cdb.Book, ff *FFMeta) error {
	if ff.Date != "" {
		t, err := time.Parse(time.DateOnly, ff.Date)
		if err != nil {
			t = time.Now()
		}
		book.Pubdate = t
	}

	book.Authors = ff.Artist
	book.Authors = ff.AlbumArtist
	book.Narrators = ff.Composer
	book.Tags = ff.Genre
	s, i := parseGrouping(ff.Grouping)
	book.Series = s
	book.SeriesIndex = ff.Disc
	if ff.Disc != i {
		book.SeriesIndex = i
	}
	book.Comments = ff.Comment
	book.Title = ff.Title
	book.Title = ff.Album

	return nil
}

func BookToFFMeta(ff *FFMeta, book *cdb.Book) error {
	meta := book.StringMap()
	delete(meta, cdb.Pubdate)

	err := mapstructure.Decode(meta, ff)
	if err != nil {
		return err
	}
	ff.Date = book.Pubdate.Format(time.DateOnly)

	if ff.Grouping != "" {
		ff.Grouping = ff.Grouping
	}

	return nil
}

func LoadFFMeta(ff *FFMeta, rc io.ReadCloser) (*FFMeta, error) {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	f, err := ini.LoadSources(opts, rc)
	if err != nil {
		return ff, err
	}

	terr := f.MapTo(ff)
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

	_, err := buf.WriteString(FFmetaHeader + "\n")
	if err != nil {
		return []byte{}, err
	}

	_, err = ffmeta.WriteTo(&buf)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func IsValidFFMetadata(b []byte) bool {
	head := []byte(FFmetaHeader)
	return bytes.Equal(head, b[:len(head)])
}

var groupRegexp = regexp.MustCompile(`(?P<series>.*), [b|B]ook (?P<index>.*)`)

func parseGrouping(g string) (string, float64) {
	matches := groupRegexp.FindStringSubmatch(g)
	s := groupRegexp.SubexpIndex("series")
	si := groupRegexp.SubexpIndex("index")

	if len(matches) < 1 {
		return g, 0
	}

	series := matches[s]

	seriesIdx := matches[si]
	if seriesIdx != "" {
		f, err := cast.ToFloat64E(seriesIdx)
		if err != nil {
			return series, 0
		}
		return series, f
	}

	return series, 0
}
