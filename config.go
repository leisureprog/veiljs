package veiljs

type Config struct {
	Compact                        bool
	ControlFlowFlattening          bool
	ControlFlowFlatteningThreshold float64
	NumbersToExpressions           bool
	Simplify                       bool
	StringArray                    bool
	StringArrayEncoding            []string
	StringArrayThreshold           float64
	UnicodeEscapeSequence          bool
	RenameGlobals                  bool
}

func DefaultConfig() Config {
	return Config{
		Compact:                        true,
		ControlFlowFlattening:          true,
		ControlFlowFlatteningThreshold: 0.75,
		NumbersToExpressions:           false,
		Simplify:                       true,
		StringArray:                    true,
		StringArrayEncoding:            []string{"base64"},
		StringArrayThreshold:           0.75,
		UnicodeEscapeSequence:          false,
		RenameGlobals:                  false,
	}
}
