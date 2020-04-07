package language

var Command map[string]string

func init() {
	supportedLanguages := [...]map[string]string{LangEn, LangEs}
	Command = make(map[string]string)

	for _, mp := range supportedLanguages {
		for k, v := range mp {
			Command[k] = v
		}
	}
}
