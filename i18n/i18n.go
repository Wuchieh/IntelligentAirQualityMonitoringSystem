package i18n

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var (
	languages = make(map[string]map[string]string)

	////go:embed all:languages/*
	//languageConfigs embed.FS
)

func init() {
	//files, err := languageConfigs.ReadDir("languages")
	//if err != nil {
	//	log.Println(err)
	//}
	//for _, f := range files {
	//	file, err := languageConfigs.ReadFile("languages/" + f.Name())
	//	var lmap map[string]string
	//	err = json.Unmarshal(file, &lmap)
	//	if err != nil {
	//		log.Println("i18n json.Unmarshal Error", err)
	//	} else {
	//		languages[f.Name()[:len(f.Name())-len(".json")]] = lmap
	//	}
	//
	//}

	join := filepath.Join("i18n", "languages")

	dir, err := os.ReadDir(join)
	if err != nil {
		panic(err)
	}

	for _, entry := range dir {

		languageJson := filepath.Join(join, entry.Name())

		if file, err := os.ReadFile(languageJson); err != nil {
			log.Println("i18n os.ReadFile Error", err)
		} else {
			var lmap map[string]string
			err = json.Unmarshal(file, &lmap)
			if err != nil {
				log.Println("i18n json.Unmarshal Error", err)
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
