package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	translations = make(map[string]map[string]string)
	mu           sync.RWMutex
	localesDir   = getDefaultLocalesDir() // or read from ENV
)

func getDefaultLocalesDir() string {
	if dir := os.Getenv("LOCALES_DIR"); dir != "" {
		return dir
	}
	return "locales"
}

// SetLocalesDir optionally set directory path (call in main if needed)
func SetLocalesDir(dir string) {
	localesDir = dir
}

// LoadLanguage loads and caches a single language (noop if already loaded)
func LoadLanguage(lang string) error {
	mu.RLock()
	if _, ok := translations[lang]; ok {
		mu.RUnlock()
		return nil
	}
	mu.RUnlock()

	// read file
	path := filepath.Join(localesDir, fmt.Sprintf("%s.json", lang))
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	mu.Lock()
	translations[lang] = m
	mu.Unlock()
	return nil
}

// LoadAllLanguages load list of languages (use at startup)
func LoadAllLanguages(langs []string) error {
	for _, l := range langs {
		if err := LoadLanguage(l); err != nil {
			return err
		}
	}
	return nil
}

// T returns translation for given lang and key, fallback to "en" then key itself
func T(lang, key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if mp, ok := translations[lang]; ok {
		if v, ok2 := mp[key]; ok2 {
			return v
		}
	}
	// fallback to english if available
	if mp, ok := translations["en"]; ok {
		if v, ok2 := mp[key]; ok2 {
			return v
		}
	}
	return key
}
