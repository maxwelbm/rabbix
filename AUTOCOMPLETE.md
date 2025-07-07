## 🔄 Autocomplete

### 🐚 Bash

Para habilitar o autocomplete no Bash, você precisa ter o pacote `bash-completion` instalado:

#### 📦 Instale o bash-completion

- **Ubuntu/Debian**:
  ```bash
  sudo apt install bash-completion
  ```

- **Arch Linux**:
  ```bash
  sudo pacman -S bash-completion
  ```

#### ⚙️ Configure o autocomplete

Gere o script e salve no seu diretório pessoal:

```bash
mkdir -p ~/.rabbix
rabbix completion bash > ~/.rabbix/rabbix.bash
```

Adicione ao seu `~/.bashrc`:

```bash
echo 'source ~/.rabbix/rabbix.bash' >> ~/.bashrc
```

Reinicie o terminal ou rode:

```bash
source ~/.bashrc
```

---

### 🧞 Zsh

No Zsh, o autocomplete é mais direto e não requer dependências extras.

#### ⚙️ Configure o autocomplete

Adicione essa linha ao final do seu `~/.zshrc`:

```zsh
autoload -U compinit; compinit
source <(rabbix completion zsh); compdef _rabbix rabbix
```

Reinicie o terminal ou rode:

```bash
source ~/.zshrc
```

---

### 🧠 Cache Inteligente

O Rabbix agora possui um sistema de cache inteligente que melhora significativamente o autocomplete:

#### ✨ Funcionalidades do Cache

- **Autocomplete Dinâmico**: Os comandos `batch` e `run` agora sugerem automaticamente os testes disponíveis
- **Sincronização Automática**: O cache é atualizado automaticamente quando você:
  - Adiciona novos testes com `rabbix add`
  - Modifica configurações com `rabbix config set`
- **Performance**: Sugestões rápidas sem precisar escanear o sistema de arquivos a cada vez

#### 🔧 Gerenciamento do Cache

Comandos disponíveis para gerenciar o cache:

```bash
# Ver estatísticas do cache
rabbix config cache stats

# Sincronizar manualmente com os arquivos
rabbix config cache sync

# Limpar o cache completamente
rabbix config cache clear
```

#### 🎯 Exemplos de Uso

Após configurar o autocomplete, você pode usar:

```bash
# Autocomplete para comando batch
rabbix batch [TAB][TAB]          # Lista todos os testes disponíveis
rabbix batch teste1 [TAB][TAB]   # Lista testes restantes (excluindo já selecionados)

# Autocomplete para comando run
rabbix run [TAB][TAB]            # Lista todos os testes disponíveis
```

#### 🔄 Como Funciona

1. **Adicionar Teste**: Quando você usa `rabbix add`, o teste é automaticamente adicionado ao cache
2. **Configurar**: Quando você usa `rabbix config set`, o cache é sincronizado com o sistema de arquivos
3. **Autocomplete**: Os comandos `batch` e `run` consultam o cache para fornecer sugestões instantâneas

O cache é armazenado em `~/.rabbix/cache.json` e contém informações sobre nome, route key e timestamps dos testes.

---

Após isso, comandos como `rabbix [TAB][TAB]` devem exibir sugestões corretamente, incluindo sugestões inteligentes para nomes de testes nos comandos `batch` e `run`.
