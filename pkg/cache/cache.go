package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getCachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".rabbix", "cache.json")
}

func loadCache() *CacheStr {
	path := getCachePath()
	cache := &CacheStr{
		Tests:   []CacheEntry{},
		Version: "1.0",
	}

	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, cache)
	}

	return cache
}

func saveCache(cache *CacheStr) error {
	path := getCachePath()
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c *Cache) GetCachedTests() []string {
	cache := loadCache()
	var tests []string

	for _, entry := range cache.Tests {
		tests = append(tests, entry.Name)
	}

	return tests
}

func (c *Cache) SyncCacheWithFileSystem() {
	settings := c.settings.LoadSettings()
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
