package middleware

import (
	"log"
	"p2_hcm_go_dat_thanh/api-gateway/internal/i18n"
	"strings"

	"github.com/gin-gonic/gin"
)

var supportedLangs = []string{"en", "vi"}

func isSupported(lang string) bool {
	for _, s := range supportedLangs {
		if s == lang {
			return true
		}
	}
	return false
}

func parseAcceptLanguageHeader(h string) string {
	if h == "" {
		return ""
	}
	parts := strings.Split(h, ",")
	if len(parts) == 0 {
		return ""
	}
	langTag := strings.TrimSpace(parts[0]) // ex: "vi-VN" or "vi"
	if idx := strings.Index(langTag, "-"); idx != -1 {
		langTag = langTag[:idx]
	}
	return strings.ToLower(langTag)
}

func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := strings.ToLower(c.Query("lang"))

		if lang == "" {
			lang = parseAcceptLanguageHeader(c.GetHeader("Accept-Language"))
		}

		if lang == "" || !isSupported(lang) {
			lang = "en"
		}

		if err := i18n.LoadLanguage(lang); err != nil {
			log.Printf("i18n: failed load language %s: %v", lang, err)
			lang = "en"
			_ = i18n.LoadLanguage("en")
		}

		c.Set("lang", lang)
		c.Set("T", func(key string) string {
			return i18n.T(lang, key)
		})

		c.Next()
	}
}
