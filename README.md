# 🐇 rabbix

This project was born out of a real need during development: to test various services, RabbitMQ queues with different payloads, in a quick and organized way. At first, this was a manual, repetitive and error-prone process - I would lose payloads, restart pages and spend a lot of time on simple tasks.

To solve this, I started with a simple script to publish messages to RabbitMQ. That script evolved into a lean CLI, which initially focused only on queues, and now I'm expanding it to allow multiple executions. The focus is to offer a tool that helps developers test manual flows in an organized, reusable and efficient way during development time.

Made by a developer for developers. I'm building it with attention, care and a focus on productivity. It's something that has helped me a lot on a daily basis and I believe it can help other people too.

## ⚙️ Install

You can install directly with go:
> requests Go 1.23 ou superior instalado.

```bash
go install github.com/maxwelbm/rabbix@latest
```

[Setup Autocomplete](AUTOCOMPLETE.md)

## License

[MIT](LICENSE) License © Maxwel Mazur
