package audbk

type Chapter struct {
	Timebase string `ini:"TIMEBASE"`
	Start    int    `ini:"START"`
	End      int    `ini:"END"`
	Title    string `ini:"title,omitempty"`
	Fields   map[string]string
}

func NewChapter() *Chapter {
	return &Chapter{
		Fields: make(map[string]string),
	}
}
