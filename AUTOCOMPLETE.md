## 🔄 Autocomplete

### 🐚 Bash

To enable autocomplete in Bash, you need to have the `bash-completion` package installed:

#### 📦 Install bash-completion

- **Ubuntu/Debian**:
  ```bash
  sudo apt install bash-completion
  ```

- **Arch Linux**:
  ```bash
  sudo pacman -S bash-completion
  ```

#### Configure the autocomplete

Get the script and save it in your personal directory:

```bash
mkdir -p ~/.rabbix
rabbix completion bash > ~/.rabbix/rabbix.bash
```

Add to your `~/.bashrc`:

```bash
echo 'source ~/.rabbix/rabbix.bash' >> ~/.bashrc
```

Restart the terminal or run:

```bash
source ~/.bashrc
```

---

### 🧞 Zsh

In Zsh, autocomplete is more straightforward and doesn't require extra dependencies.

#### Configure the autocomplete

Add this line to the end of your `~/.zshrc`:

```zsh
autoload -U compinit; compinit
source <(rabbix completion zsh); compdef _rabbix rabbix
```

Restart the terminal or run:

```bash
source ~/.zshrc
```