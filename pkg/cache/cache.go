package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/maxwelbm/rabbix/pkg/sett"
)

func getCachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rabbix", "cache.json")
}

func loadCache() *Cache {
	path := getCachePath()
	cache := &Cache{
		Tests:   []CacheEntry{},
		Version: "1.0",
	}

	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, cache)
	}

	return cache
}

func saveCache(cache *Cache) error {
	path := getCachePath()
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func addToCache(name, routeKey string) {
	cache := loadCache()

	// Verifica se já existe
	for i, entry := range cache.Tests {
		if entry.Name == name {
			// Atualiza entrada existente
			cache.Tests[i].RouteKey = routeKey
			cache.Tests[i].UpdatedAt = time.Now()
			if err := saveCache(cache); err != nil {
				fmt.Printf("Erro ao salvar cache: %v\n", err)
			}
			return
		}
	}

	// Adiciona nova entrada
	entry := CacheEntry{
		Name:      name,
		RouteKey:  routeKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cache.Tests = append(cache.Tests, entry)
	if err := saveCache(cache); err != nil {
		fmt.Printf("Erro ao salvar cache: %v\n", err)
	}
}

func GetCachedTests() []string {
	cache := loadCache()
	var tests []string

	for _, entry := range cache.Tests {
		tests = append(tests, entry.Name)
	}

	return tests
}

func SyncCacheWithFileSystem() {
	settings := sett.LoadSettings()
	outputDir := settings["output_dir"]
	if outputDir == "" {
		home, _ := os.UserHomeDir()
		outputDir = filepath.Join(home, ".rabbix", "tests")
	}

	// Carrega cache atual
	cache := loadCache()

	// Mapeia testes existentes no cache
	cacheMap := make(map[string]CacheEntry)
	for _, entry := range cache.Tests {
		cacheMap[entry.Name] = entry
	}

	// Verifica arquivos no sistema
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return
	}

	var newTests []CacheEntry

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			fileName := strings.TrimSuffix(file.Name(), ".json")

			// Tenta carregar detalhes do arquivo
			testPath := filepath.Join(outputDir, file.Name())
			if data, err := os.ReadFile(testPath); err == nil {
				var testCase TestCase
				if err := json.Unmarshal(data, &testCase); err == nil {
					// Se já existe no cache, mantém as datas
					if existing, exists := cacheMap[fileName]; exists {
						existing.RouteKey = testCase.RouteKey
						existing.UpdatedAt = time.Now()
						newTests = append(newTests, existing)
					} else {
						// Novo teste encontrado - usa nome do arquivo, não o campo "name" do JSON
						entry := CacheEntry{
							Name:      fileName,
							RouteKey:  testCase.RouteKey,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}
						newTests = append(newTests, entry)
					}
				}
			}
		}
	}

	// Atualiza cache
	cache.Tests = newTests
	if err := saveCache(cache); err != nil {
		fmt.Printf("❌ Erro ao salvar cache: %v\n", err)
	}
}

func clearCache() {
	cache := &Cache{
		Tests:   []CacheEntry{},
		Version: "1.0",
	}

	if err := saveCache(cache); err != nil {
		fmt.Printf("❌ Erro ao limpar cache: %v\n", err)
	} else {
		fmt.Println("✅ Cache limpo com sucesso")
	}
}
