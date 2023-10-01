package audbk

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/cdb"
	"github.com/spf13/cast"
	"gopkg.in/ini.v1"
)

const FFMetaHeader = ";FFMETADATA1"

var InvalidFFmetadata = errors.New("ffmetadata file is not valid")

type FFMetaDecoder struct {
	ini    *ini.File
	reader io.Reader
}

type FFMeta struct {
	ini         *ini.File
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
	ini.PrettyFormat = false

	opts := ini.LoadOptions{
		IgnoreInlineComment:    true,
		AllowNonUniqueSections: true,
	}

	return &FFMeta{
		ini: ini.Empty(opts),
	}
}

func DecodeINI(s *cdb.Serializer) cdb.BookDecoder {
	s.Format = ".ini"
	return func(r io.Reader) cdb.Decoder {
		ff := NewFFMeta()
		//if err != nil {
		//  return err
		//}
		return ff
	}
}

func NewChapter() *FFMetaChapter {
	return &FFMetaChapter{
		Fields: make(map[string]any),
	}
}

func (ff *FFMeta) ReadFile(input string) (*FFMeta, error) {
	data, err := OpenFFMeta(input)
	if err != nil {
		return ff, err
	}

	file := io.NopCloser(bytes.NewReader(data))

	err = LoadFFMeta(ff, file)
	if err != nil {
		return ff, err
	}

	return ff, nil
}

func (ff *FFMeta) Decode(v any) error {
	//meta := make(map[string]any)
	//switch m := v.(type) {
	//case map[string]any:
	//meta = m
	//case cdb.Book:
	//meta = m.StringMap()
	//}
	//mapstructure.Decode(v, &meta)
	//BookToFFMeta(ff, meta)
	if b, ok := v.(*cdb.Book); ok {
		println("poot")
		FFMetaToBook(b, ff)
	}
	return nil
}

func (ff *FFMeta) WriteTo(w io.Writer) (int64, error) {
	if ff.Title == "" {
		return 0, nil
	}
	meta, err := FFMetaToIniFile(ff)
	if err != nil {
		return 0, err
	}
	ff.ini = meta

	head := []byte(FFMetaHeader)
	head = append(head, '\n')
	i, err := w.Write(head)
	if err != nil {
		return int64(i), err
	}

	return meta.WriteTo(w)
}

func FFMetaToIniFile(ff *FFMeta) (*ini.File, error) {
	ini.PrettyFormat = false

	opts := ini.LoadOptions{
		IgnoreInlineComment:    true,
		AllowNonUniqueSections: true,
	}

	ffmeta := ini.Empty(opts)

	if ff.Title == "" {
		return ffmeta, nil
	}

	ffmeta.ReflectFrom(ff)

	for _, chapter := range ff.Chapters {
		sec, err := ffmeta.NewSection("CHAPTER")
		if err != nil {
			return ffmeta, err
		}
		if chapter.Timebase != "" {
			sec.NewKey("TIMEBASE", chapter.Timebase)
		}
		sec.NewKey("START", cast.ToString(chapter.Start))
		if chapter.End > 0 {
			sec.NewKey("END", cast.ToString(chapter.End))
		}
		if chapter.Title != "" {
			sec.NewKey("title", chapter.Title)
		}
	}

	return ffmeta, nil
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

func BookToFFMeta(ff *FFMeta, meta map[string]any) error {
	var pub time.Time
	if d, ok := meta[cdb.Pubdate]; ok {
		pub = d.(time.Time)
	}
	delete(meta, cdb.Pubdate)

	err := mapstructure.Decode(meta, ff)
	if err != nil {
		return err
	}
	ff.Date = pub.Format(time.DateOnly)

	if ff.Grouping != "" {
		idx := cast.ToString(ff.Disc)
		ff.Grouping = ff.Grouping + ", Book " + idx
	}

	return nil
}

func LoadFFMeta(ff *FFMeta, rc io.Reader) error {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	f, err := ini.LoadSources(opts, rc)
	if err != nil {
		return err
	}

	terr := f.MapTo(ff)
	if terr != nil {
		return terr
	}

	if f.HasSection("chapter") {
		sections, err := f.SectionsByName("chapter")
		if err != nil {
			return err
		}

		for _, sec := range sections {
			var ch FFMetaChapter
			sec.MapTo(&ch)
			ff.Chapters = append(ff.Chapters, ch)
		}
	}

	ff.ini = f
	return nil
}

func IsValidFFMeta(b []byte) bool {
	head := []byte(FFMetaHeader)
	return bytes.Equal(head, b[:len(head)])
}

func OpenFFMeta(f string) ([]byte, error) {
	file, err := os.Open(f)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	return ReadFFMeta(file)
}

func ReadFFMeta(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return b, err
	}

	head := []byte(FFMetaHeader)
	if !bytes.Equal(head, b[:len(head)]) {
		return b, InvalidFFmetadata
	}

	return b, nil
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
