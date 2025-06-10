# Web Assets - Rabbix UI

Este diretório contém todos os arquivos relacionados à interface web do Rabbix, organizados seguindo as convenções do Go.

## 📁 Estrutura de Diretórios

```
web/
├── README.md           # Este arquivo
├── embed.go           # Configuração do embed para arquivos estáticos
├── templates/         # Templates HTML
│   └── index.html     # Template principal da interface
└── static/           # Arquivos estáticos
    ├── css/          # Folhas de estilo
    │   └── style.css # CSS principal da interface
    └── js/           # Scripts JavaScript
        └── script.js # JavaScript principal da interface
```

## 🔧 Como Funciona

### Embed (Go 1.16+)
Os arquivos estáticos são **embebidos diretamente no binário** usando a diretiva `//go:embed`, eliminando a necessidade de distribuir arquivos separados.

```go
//go:embed static/css/*.css static/js/*.js templates/*.html
var Assets embed.FS
```

### Vantagens do Embed
- ✅ **Binário único**: Todos os assets incluídos no executável
- ✅ **Deploy simplificado**: Não precisa copiar arquivos separados
- ✅ **Performance**: Acesso direto aos arquivos na memória
- ✅ **Segurança**: Assets não podem ser modificados externamente

## 📝 Arquivos

### `embed.go`
Configura o sistema de embed e expõe funções para acessar os assets:
- `GetTemplate()`: Carrega templates HTML
- `GetStaticHandler()`: Handler para arquivos estáticos
- `GetStaticFile()`: Acesso direto a arquivos

### `templates/index.html`
Template principal da interface web com:
- Estrutura HTML responsiva
- Sistema de abas (Logs, Resultados, Gráficos)
- Lista de testes com checkboxes
- Configurações de execução em lote
- Placeholders para conteúdo dinâmico

### `static/css/style.css`
Estilos modernos incluindo:
- Tema escuro com gradientes
- Layout responsivo (desktop/tablet/mobile)
- Animações e transições suaves
- Componentes estilizados (botões, formulários, cards)
- Sistema de cores consistente

### `static/js/script.js`
JavaScript vanilla com funcionalidades:
- Gerenciamento de abas
- Execução de testes individuais e em lote
- Logs em tempo real via Server-Sent Events
- Atualização dinâmica de resultados
- Gráficos simples em canvas

## 🎨 Design System

### Cores Principais
- **Primária**: `#4facfe` (Azul)
- **Sucesso**: `#55efc4` (Verde)
- **Aviso**: `#fdcb6e` (Amarelo)
- **Erro**: `#ff7675` (Vermelho)
- **Fundo**: `#0f0f23` → `#1a1a2e` (Gradiente)

### Tipografia
- **Fonte**: Segoe UI, Tahoma, Geneva, Verdana, sans-serif
- **Logs**: Consolas, Monaco, monospace

### Breakpoints
- **Desktop**: 1024px+
- **Tablet**: 768px - 1024px
- **Mobile**: < 768px

## 🔄 Fluxo de Desenvolvimento

### Modificando a Interface

1. **HTML**: Edite `templates/index.html`
2. **CSS**: Modifique `static/css/style.css`
3. **JS**: Altere `static/js/script.js`
4. **Rebuild**: Execute `go build` para reembutir os assets

### Testando Mudanças

```bash
# Rebuild com novos assets
go build -o rabbix .

# Teste a interface
./rabbix ui
```

### Adicionando Novos Assets

1. Adicione arquivos nos diretórios apropriados
2. Atualize a diretiva `//go:embed` em `embed.go` se necessário
3. Rebuild o projeto

## 🚀 APIs Utilizadas

A interface consome as seguintes APIs REST:

- `GET /api/tests` - Lista testes disponíveis
- `POST /api/run/{teste}` - Executa teste individual
- `POST /api/batch` - Executa lote de testes
- `GET /api/execution/{id}` - Status da execução
- `GET /api/logs/{id}` - Logs em tempo real (SSE)

## 📱 Funcionalidades

### Execução Individual
- Botão ▶ ao lado de cada teste
- Feedback visual (loading state)
- Resultado imediato nos logs e resultados

### Execução em Lote
- Seleção múltipla com checkboxes
- Configuração de concorrência (1-20)
- Configuração de delay (0-5000ms)
- Logs em tempo real durante execução

### Monitoramento
- **Logs**: Stream em tempo real com cores por tipo
- **Resultados**: Lista detalhada com estatísticas
- **Gráficos**: Visualizações simples de status e timing

## 🔧 Extensibilidade

### Adicionando Novas Funcionalidades

1. **Novas APIs**: Adicione endpoints no `ui.go`
2. **Novos Templates**: Crie arquivos em `templates/`
3. **Componentes CSS**: Adicione estilos em `style.css`
4. **Interações JS**: Implemente em `script.js`

### Melhorias Futuras

- [ ] Gráficos avançados com Chart.js
- [ ] WebSockets para comunicação bidirecional
- [ ] Filtros e busca na lista de testes
- [ ] Temas customizáveis
- [ ] Exportação de relatórios
- [ ] Histórico de execuções

## 📚 Referências

- [Go embed](https://pkg.go.dev/embed)
- [HTML Templates](https://pkg.go.dev/html/template)
- [Server-Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
- [Responsive Design](https://developer.mozilla.org/en-US/docs/Learn/CSS/CSS_layout/Responsive_Design)

---

**Mantido pelo time Rabbix** 🐰