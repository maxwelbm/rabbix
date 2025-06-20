package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CacheEntry struct {
	Name      string    `json:"name"`
	RouteKey  string    `json:"route_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cache struct {
	Tests   []CacheEntry `json:"tests"`
	Version string       `json:"version"`
}

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

// func removeFromCache(name string) {
// 	cache := loadCache()

// 	for i, entry := range cache.Tests {
// 		if entry.Name == name {
// 			// Remove entrada
// 			cache.Tests = append(cache.Tests[:i], cache.Tests[i+1:]...)
// 			if err := saveCache(cache); err != nil {
// 				fmt.Printf("Erro ao salvar cache: %v\n", err)
// 			}
// 			return
// 		}
// 	}
// }

func getCachedTests() []string {
	cache := loadCache()
	var tests []string

	for _, entry := range cache.Tests {
		tests = append(tests, entry.Name)
	}

	return tests
}

// func getCachedTestsWithRouteKey() []CacheEntry {
// 	cache := loadCache()
// 	return cache.Tests
// }

func syncCacheWithFileSystem() {
	settings := loadSettings()
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

func printCacheStats() {
	cache := loadCache()
	fmt.Printf("📊 Cache Statistics:\n")
	fmt.Printf("   Total tests: %d\n", len(cache.Tests))
	fmt.Printf("   Cache version: %s\n", cache.Version)

	if len(cache.Tests) > 0 {
		fmt.Printf("   Tests available for autocomplete:\n")
		for _, entry := range cache.Tests {
			fmt.Printf("     • %s (route: %s)\n", entry.Name, entry.RouteKey)
		}
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
