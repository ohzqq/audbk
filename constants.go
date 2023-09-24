package audbk

//go:generate stringer -type CreatorProperty -linecomment
type CreatorProperty int

const (
	Artist      CreatorProperty = iota // artist
	Author                             // author
	Colorist                           // colorist
	Contributor                        // contributor
	Editor                             // editor
	Illustrator                        // illustrator
	Imprint                            // imprint
	Inker                              // inker
	Letterer                           // letterer
	Narrator                           // narrator
	Penciler                           // penciler
	Publisher                          // publisher
	ReadBy                             // readBy
	Translator                         // translator
)

func (cp CreatorProperty) Singular() string {
	return cp.String()
}

func (cp CreatorProperty) Plural() string {
	if cp == ReadBy {
		return cp.Singular()
	}
	return cp.Singular() + "s"
}

func (cp CreatorProperty) Matches(v string) bool {
	if cp.Singular() == v {
		return true
	}
	if cp.Plural() == v {
		return true
	}
	return false
}
