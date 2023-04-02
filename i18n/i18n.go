package i18n

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var (
	languages = make(map[string]map[string]string)
)

func init() {
	join := filepath.Join("i18n")

	dir, err := os.ReadDir(join)
	if err != nil {
		panic(err)
	}

	for _, entry := range dir {

		if entry.IsDir() {
			continue
		}

		languageJson := filepath.Join("i18n", entry.Name())
		if file, err := os.ReadFile(languageJson); err != nil {
			log.Println(err)
		} else {
			var lmap map[string]string
			err = json.Unmarshal(file, &lmap)
			if err != nil {
				log.Println(err)
			} else {
				languages[entry.Name()[:len(entry.Name())-len(".json")]] = lmap
			}
		}
	}
}

func Get(language, arg string) string {
	if _, ok := languages[language]; ok {
		return languages[language][arg]
	} else {
		return languages["en"][arg]
	}
}
