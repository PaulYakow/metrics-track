// Package config конфигурация модуля staticlint
package config

// ConfigData структура с конфигурацией анализаторов пакета staticcheck
//
// SA добавлены полностью. Для добавления дополнительных прописать необходимый в файле `config.json`
type ConfigData struct {
	Staticcheck []string
}
