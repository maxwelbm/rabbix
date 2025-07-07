package cache

import (
	"fmt"

	"github.com/spf13/cobra"
)

var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Gerencia o cache de autocomplete",
}

var cacheStats = &cobra.Command{
	Use:   "stats",
	Short: "Exibe estatísticas do cache",
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

var cacheSync = &cobra.Command{
	Use:   "sync",
	Short: "Sincroniza o cache com os arquivos de teste",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Sincronizando cache...")
		SyncCacheWithFileSystem()
		fmt.Println("✅ Cache sincronizado com sucesso.")
	},
}

var cacheClear = &cobra.Command{
	Use:   "clear",
	Short: "Limpa o cache de autocomplete",
	Run: func(cmd *cobra.Command, args []string) {
		cache := &Cache{
			Tests:   []CacheEntry{},
			Version: "1.0",
		}

		if err := saveCache(cache); err != nil {
			fmt.Printf("❌ Erro ao limpar cache: %v\n", err)
		} else {
			fmt.Println("✅ Cache limpo com sucesso")
		}
	},
}

func init() {
	CacheCmd.AddCommand(cacheStats)
	CacheCmd.AddCommand(cacheSync)
	CacheCmd.AddCommand(cacheClear)
}
