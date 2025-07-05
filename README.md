# 🐇 rabbix

This project was born out of a real need during development: to test multiple services, RabbitMQ queues and APIs with different payloads, in a quick and organized way. At first, this was a manual, repetitive and error-prone process - I would lose payloads, restart pages and spend too much time on simple tasks.

To solve this, I started with a simple script to publish messages in RabbitMQ. This script evolved into a lean CLI, which initially focused only on queues, and now I'm expanding it to also allow REST requests. The focus is to offer a tool that helps developers test manual flows in an organized, reusable and efficient way during development time.

More than a tool, this has become a serious project - made by a developer, for developers. I'm building it with attention, care and a focus on productivity. It's something that has helped me a lot on a daily basis, and I believe it can help others too.

## ⚙️ Install

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
