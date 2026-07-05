# ready

Ferramenta de configuração de ambiente que verifica e instala aplicações e gerencia aliases do shell.

## Aplicações verificadas

| Aplicação    | Binário verificado | Pacote apt      | Instalação customizada                |
|-------------|-------------------|----------------|---------------------------------------|
| vscode      | `code`            | `code`         | —                                     |
| opencode    | `opencode`        | —              | `npm install -g opencode-ai`          |
| docker      | `docker`          | `docker.io`    | —                                     |
| podman      | `podman`          | `podman`       | —                                     |
| lazygit     | `lazygit`         | `lazygit`      | —                                     |
| fzf         | `fzf`             | `fzf`          | —                                     |
| ripgrep     | `rg`              | `ripgrep`      | —                                     |
| bat         | `batcat`          | `bat`          | —                                     |
| jq          | `jq`              | `jq`           | —                                     |
| btop        | `btop`            | `btop`         | —                                     |
| zoxide      | `zoxide`          | `zoxide`       | —                                     |
| lazydocker  | `lazydocker`      | —              | `curl ... \| bash`                    |
| gh          | `gh`              | `gh`           | —                                     |
| tldr        | `tldr`            | `tldr`         | —                                     |
| speedtest   | `speedtest-cli`   | `speedtest-cli` | `pipx install speedtest-cli`         |
| frogmouth   | `frogmouth`       | —              | `snap install frogmouth`              |
| ffmpeg      | `ffmpeg`          | `ffmpeg`       | —                                     |
| exa         | `exa`             | `exa`          | —                                     |
| whichllm    | `whichllm`        | —              | `pipx install which-llm`              |

## Aliases gerenciados

| Alias         | Comando                                                              |
|---------------|----------------------------------------------------------------------|
| `azw32pdf`    | `for i in *.azw3; do ebook-convert "$i" "${i%.*}.pdf"; done`        |
| `epub2pdf`    | `for i in *.epub; do ebook-convert "$i" "${i%.*}.pdf"; done`        |
| `k`           | `kubectl`                                                            |
| `l`           | `ls -lah`                                                            |
| `la`          | `ls -lAh`                                                            |
| `list_services` | `systemctl list-units --type=service --state=running`              |
| `ll`          | `exa --long --header --git`                                          |
| `ls`          | `ls --color=tty`                                                     |
| `lsa`         | `ls -lah`                                                            |
| `mobi2pdf`    | `for i in *.mobi; do ebook-convert "$i" "${i%.*}.pdf"; done`        |
| `ocd`         | `opencode`                                                           |
| `pdf2epub`    | `for i in *.pdf; do ebook-convert "$i" "${i%.*}.epub"; done`        |
| `postgres`    | `docker run -p 5432:5432 -e POSTGRES_PASSWORD=dbadmin -d postgres`  |
| `vi`          | `vim`                                                                |

## Uso

```bash
go build -o ready && ./ready
```

Por padrão lê `apps.txt` e `aliases.txt`. Use `--apps` e `--aliases` para caminhos customizados.
