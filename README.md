# 🐇 rabbix

**rabbix** é uma CLI para facilitar os testes de micro-serviços que utilizam RabbitMQ, tornando o processo de criação, execução e documentação de mensagens mais simples e padronizado.

## 🎯 Objetivo

- Permitir o envio de mensagens RabbitMQ diretamente pela linha de comando.
- Salvar e reutilizar casos de teste com JSON dinâmico.
- Substituir o uso do Postman com uma interface mais simples e rápida.
- Facilitar a documentação e reexecução de mensagens usadas em ambientes de integração.

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
