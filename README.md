# 🐇 rabbix

**rabbix** é uma CLI para facilitar os testes de micro-serviços que utilizam RabbitMQ, tornando o processo de criação, execução e documentação de mensagens mais simples e padronizado.

## 🎯 Objetivo

- Permitir o envio de mensagens RabbitMQ diretamente pela linha de comando.
- Salvar e reutilizar casos de teste com JSON dinâmico.
- Substituir o uso do Postman com uma interface mais simples e rápida.
- Facilitar a documentação e reexecução de mensagens usadas em ambientes de integração.

## 📁 Estrutura do Projeto

```
rabbix/
├── main.go                 # Ponto de entrada da aplicação
├── pkg/cmd/               # Comandos da CLI
│   ├── root.go           # Comando raiz
│   ├── add.go            # Adicionar testes
│   ├── list.go           # Listar testes
│   ├── run.go            # Executar teste individual
│   ├── batch.go          # Executar testes em lote
│   ├── config.go         # Configurações
│   └── ui.go             # Interface web avançada
└── web/                  # Assets da interface web
    ├── embed.go          # Sistema de embed
    ├── templates/        # Templates HTML
    │   └── index.html    # Interface principal
    └── static/           # Arquivos estáticos
        ├── css/
        │   └── style.css # Estilos da interface
        └── js/
            └── script.js # JavaScript da interface
```

## ⚙️ Instalação

Você pode instalar diretamente com:

```bash
go install github.com/maxwelbm/rabbix@latest
```

> Requer Go 1.18 ou superior instalado.

## 📁 Configuração

Use o comando `rabbix config` para definir o host RabbitMQ e o diretório onde os testes serão salvos:

```bash
# Define o host base
rabbix config set --host http://localhost:15672

# Define o diretório onde os testes serão salvos
rabbix config set --output ./vaca
```

Você pode verificar as configurações atuais com:

```bash
rabbix config get
```

## 💡 Comandos disponíveis

```bash
rabbix add --file exemplo.json --routeKey minha.fila --name teste_simples
rabbix list
rabbix ui
```

## 🔄 [Setup Autocomplete](README_AUTOCOMPLETE.md)

## 🪪 Licença

[MIT](LICENSE) License © Maxwel Mazur
