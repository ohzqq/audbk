package audbk

import (
	"bufio"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

const FFmetaHeader = ";FFMETADATA1\n"

func Load(input string) (map[string]any, error) {
	opts := ini.LoadOptions{}
	opts.Insensitive = true
	opts.InsensitiveSections = true
	opts.IgnoreInlineComment = true
	opts.AllowNonUniqueSections = true

	//if !IsFFMeta(ffmeta.File) {
	//  return ffmeta, fmt.Errorf("not an ffmetadata file")
	//}

	f, err := ini.LoadSources(opts, input)
	if err != nil {
		return ffmeta, err
	}

	meta := make(map[string]any)

	for f, v := range f.Section("").KeysHash() {
		meta[f] = v
	}

	if f.HasSection("chapter") {
		var chaps []map[string]string
		sections, err := f.SectionsByName("chapter")
		if err != nil {
			return ffmeta, err
		}

		for _, sec := range sections {
			chaps = append(chaps, sec.KeysHash())
		}

		meta["chapters"] = chaps
	}

	return meta, nil
}

//func Dump(meta avtools.Meta) []byte {
//  ini.PrettyFormat = false

//  opts := ini.LoadOptions{
//    IgnoreInlineComment:    true,
//    AllowNonUniqueSections: true,
//  }

//  ffmeta := ini.Empty(opts)

//  for k, v := range meta.Tags() {
//    _, err := ffmeta.Section("").NewKey(k, v)
//    if err != nil {
//      log.Fatal(err)
//    }
//  }

//  for _, chapter := range meta.Chapters() {
//    sec, err := ffmeta.NewSection("CHAPTER")
//    if err != nil {
//      log.Fatal(err)
//    }
//    sec.NewKey("TIMEBASE", "1/1000")
//    ss := strconv.Itoa(int(chapter.StartStamp.Dur.Milliseconds()))
//    sec.NewKey("START", ss)

//    to := strconv.Itoa(int(chapter.EndStamp.Dur.Milliseconds()))
//    sec.NewKey("END", to)
//    sec.NewKey("title", chapter.ChapTitle)
//    for k, v := range chapter.Tags {
//      sec.NewKey(k, v)
//    }
//  }

//  var buf bytes.Buffer
//  _, err := buf.WriteString(FFmetaComment)
//  _, err = ffmeta.WriteTo(&buf)
//  if err != nil {
//    log.Fatal(err)
//  }

//  return buf.Bytes()
//}

func IsFFMeta(f string) bool {
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
