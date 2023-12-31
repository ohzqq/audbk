// Code generated by "stringer -type CreatorProperty -linecomment"; DO NOT EDIT.

package audbk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Artist-0]
	_ = x[Author-1]
	_ = x[Colorist-2]
	_ = x[Contributor-3]
	_ = x[Editor-4]
	_ = x[Illustrator-5]
	_ = x[Imprint-6]
	_ = x[Inker-7]
	_ = x[Letterer-8]
	_ = x[Narrator-9]
	_ = x[Penciler-10]
	_ = x[Publisher-11]
	_ = x[ReadBy-12]
	_ = x[Translator-13]
}

const _CreatorProperty_name = "artistauthorcoloristcontributoreditorillustratorimprintinkerletterernarratorpencilerpublisherreadBytranslator"

var _CreatorProperty_index = [...]uint8{0, 6, 12, 20, 31, 37, 48, 55, 60, 68, 76, 84, 93, 99, 109}

func (i CreatorProperty) String() string {
	if i < 0 || i >= CreatorProperty(len(_CreatorProperty_index)-1) {
		return "CreatorProperty(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CreatorProperty_name[_CreatorProperty_index[i]:_CreatorProperty_index[i+1]]
}
