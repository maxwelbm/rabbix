package cache

import (
	"fmt"

	"github.com/maxwelbm/rabbix/pkg/sett"
	"github.com/spf13/cobra"
)

type CacheItf interface {
	GetCachedTests() []string
	SyncCacheWithFileSystem()
	CmdCache() *cobra.Command
}

// CacheItf implementation
var _ CacheItf = (*Cache)(nil)

type Cache struct {
	settings sett.SettItf
}

func New(settings sett.SettItf) CacheItf {
	return &Cache{
		settings: settings,
	}
}

func (c *Cache) CmdCache() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Gerencia o cache de autocomplete",
	}

	cmd.AddCommand(c.cmdStats())
	cmd.AddCommand(c.cmdClear())
	cmd.AddCommand(c.cmdSync())

	return cmd
}

func (c *Cache) cmdStats() *cobra.Command {
	return &cobra.Command{
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
}

func (c *Cache) cmdSync() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sincroniza o cache com os arquivos de teste",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔄 Sincronizando cache...")
			c.SyncCacheWithFileSystem()
			fmt.Println("✅ Cache sincronizado com sucesso.")
		},
	}
}

func (c *Cache) cmdClear() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Limpa o cache de autocomplete",
		Run: func(cmd *cobra.Command, args []string) {
			cache := &CacheStr{
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
}
