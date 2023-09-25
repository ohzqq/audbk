//go:build exclude

package audbk

const (
	ConformsTo           = "conformsTo"
	Context              = "@context"
	ReadingOrder         = "readingOrder"
	Name                 = "name"
	Abridged             = "abridged"
	AccessibilityFeature = "accessibilityFeature"
	AccessibilityHazard  = "accessibilityHazard"
	AccessibilitySummary = "accessibilitySummary"
	AccessMode           = "accessMode"
	AccessModeSufficient = "accessModeSufficient"
	Cover                = "cover"
	Duration             = "duration"
	DateModified         = "dateModified"
	DatePublished        = "datePublished"
	ID                   = "id"
	InLanguage           = "inLanguage"
	ReadingProgression   = "readingProgression"
	Resources            = "resources"
	URL                  = "url"
)

type Manifest struct {
	ConformsTo           string           `json:"conformsTo" mapstructure:"conformsTo"`
	Context              string           `json:"@context" mapstructure:"context"`
	ReadingOrder         []LinkedResource `json:"readingOrder" mapstructure:"readingOrder"`
	Name                 string           `json:"name" mapstructure:"name"`
	Abridged             bool             `json:"abridged,omitempty" mapstructure:"abridged,omitempty"`
	AccessibilityFeature []string         `json:"accessibilityFeature,omitempty" mapstructure:"accessibilityFeature,omitempty"`
	AccessibilityHazard  []string         `json:"accessibilityHazard,omitempty" mapstructure:"accessibilityHazard,omitempty"`
	AccessibilitySummary string           `json:"accessibilitySummary,omitempty" mapstructure:"accessibilitySummary,omitempty"`
	AccessMode           []string         `json:"accessMode,omitempty" mapstructure:"accessMode,omitempty"`
	AccessModeSufficient []any            `json:"accessModeSufficient,omitempty" mapstructure:"accessModeSufficient,omitempty"`
	Author               []string         `json:"author,omitempty" mapstructure:"author,omitempty"`
	Cover                LinkedResource   `json:"cover,omitempty" mapstructure:"cover,omitempty"`
	Duration             string           `json:"duration,omitempty" mapstructure:"duration,omitempty"`
	DateModified         string           `json:"dateModified,omitempty" mapstructure:"dateModified,omitempty"`
	DatePublished        string           `json:"datePublished,omitempty" mapstructure:"datePublished,omitempty"`
	ID                   string           `json:"id,omitempty" mapstructure:"id,omitempty"`
	InLanguage           []string         `json:"inLanguage,omitempty" mapstructure:"inLanguage,omitempty"`
	ReadBy               []string         `json:"readBy,omitempty" mapstructure:"readBy,omitempty"`
	ReadingProgression   string           `json:"readingProgression,omitempty" mapstructure:"readingProgression,omitempty"`
	Resources            []LinkedResource `json:"resources,omitempty" mapstructure:"resources,omitempty"`
	URL                  string           `json:"url,omitempty" mapstructure:"url,omitempty"`
	AdditionalProperties map[string]string
}

type LinkedResource struct {
	Type           []string `json:"type,omitempty"`
	URL            string   `json:"url"`
	EncodingFormat string   `json:"encodingFormat,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	Rel            []string `json:"rel,omitempty"`
	Integrity      string   `json:"integrity,omitempty"`
	Duration       string   `json:"duration,omitempty"`
}

type Creator struct {
	Name   string
	Entity *Entity
}

type Entity struct {
	Type       []string
	Name       []LocalizableString
	ID         string
	URL        string
	Identifier []string
}

type LocalizableString struct {
	Monolingual  string
	Multilingual MultilingualString
}

type MultilingualString struct {
	Value     string
	Language  string
	Direction string
}

//func (m Manifest) MarshalJSON() ([]byte, error) {
//}
