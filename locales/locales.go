package locales

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Xuanwo/go-locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed config/*
var translationFS embed.FS

var (
	bundle          *i18n.Bundle
	localizer       *i18n.Localizer
	defaultLanguage = language.English
	currentLanguage = defaultLanguage
)

// 初始化 i18n bundle
func init() {
	tag, _ := locale.Detect()
	if tag != language.Und {
		currentLanguage = tag
	}

	bundle = i18n.NewBundle(defaultLanguage)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	err := fs.WalkDir(translationFS, "config", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".toml") {
			_, err = bundle.LoadMessageFileFS(translationFS, path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	localizer = i18n.NewLocalizer(bundle, currentLanguage.String(), defaultLanguage.String())
}

// / MustLocalize 翻译字符串
func MustLocalize(lc *i18n.LocalizeConfig) string {
	return localizer.MustLocalize(lc)
}

// / MustLocalizeMessage 翻译字符串
func MustLocalizeMessage(msg *i18n.Message) string {
	return localizer.MustLocalizeMessage(msg)
}

// Tr 翻译没有参数的字符串
func Tr(messageID string) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID
	}
	return msg
}

// TrParams 翻译包含参数模板的字符串
func TrParams(messageID string, params map[string]any) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: params,
	})
	if err != nil {
		return messageID
	}
	return msg
}

// TrPlural 翻译复数形式的字符串
func TrPlural(messageID string, count any) string {
	return TrPluralParams(messageID, count, nil)
}

// TrPluralParams 翻译复数形式且含有参数模板的字符串
func TrPluralParams(messageID string, count any, params map[string]any) string {
	if params == nil {
		params = make(map[string]any)
	}
	params["Count"] = count
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		PluralCount:  count,
		TemplateData: params,
	})
	if err != nil {
		return messageID
	}
	return msg
}
