package lang

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"os"
)

// translations is a map of language tags to a map of translation keys to translation values.
var translations = map[language.Tag]map[string]string{}

func init() {
	register(language.English)
	register(language.Indonesian)
	register(language.Chinese)
	register(language.Japanese)
	register(language.Vietnamese)
	register(language.Hindi)
}

// Register registers a language file.
func register(lang language.Tag) {
	translationData := map[string]string{}

	rawBytes, err := os.ReadFile(fmt.Sprintf("assets/lang/%s.yml", lang.String()))
	if err != nil {
		panic(fmt.Errorf("failed to read lang file: %w", err))
	}

	if err := yaml.Unmarshal(rawBytes, &translationData); err != nil {
		panic(fmt.Errorf("failed to unmarshal lang: %w", err))
	}

	translations[lang] = translationData
}

func ToLangTag(locale string) language.Tag {
	switch locale {
	case "id", "id-ID":
		return language.Indonesian
	case "zh", "zh-CN", "zh-TW", "zh-HK":
		return language.Chinese
	case "vi", "vi-VN":
		return language.Vietnamese
	case "ja", "ja-JP":
		return language.Japanese
	case "hi", "hi-IN":
		return language.Hindi
	}
	return language.English
}

func Translatef(lang language.Tag, key string, args ...interface{}) string {
	return text.Colourf(Translate(lang, key), args...)
}

func Translate(lang language.Tag, key string) string {
	t, ok := translations[lang]
	if !ok {
		if lang == language.English {
			return key
		}
		return Translate(language.English, key)
	}
	return t[key]
}
