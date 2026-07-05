package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type appInfo struct {
	Binary     string
	Pkg        string
	InstallCmd []string
}

var appMap = map[string]appInfo{
	"vscode":     {"code", "code", nil},
	"opencode":   {"opencode", "", []string{"npm", "install", "-g", "opencode-ai"}},
	"docker":     {"docker", "docker.io", nil},
	"podman":     {"podman", "podman", nil},
	"lazygit":    {"lazygit", "lazygit", nil},
	"fzf":        {"fzf", "fzf", nil},
	"ripgrep":    {"rg", "ripgrep", nil},
	"bat":        {"batcat", "bat", nil},
	"jq":         {"jq", "jq", nil},
	"btop":       {"btop", "btop", nil},
	"zoxide":     {"zoxide", "zoxide", nil},
	"lazydocker": {"lazydocker", "", []string{"bash", "-c", "curl -fsSL https://raw.githubusercontent.com/jesseduffield/lazydocker/master/scripts/install_update_linux.sh | bash"}},
	"gh":         {"gh", "gh", nil},
	"tldr":       {"tldr", "tldr", nil},
	"speedtest":  {"speedtest-cli", "speedtest-cli", []string{"pipx", "install", "speedtest-cli"}},
	"frogmouth":  {"frogmouth", "", []string{"snap", "install", "frogmouth"}},
	"ffmpeg":     {"ffmpeg", "ffmpeg", nil},
	"exa":        {"exa", "exa", nil},
	"whichllm":   {"whichllm", "", []string{"pipx", "install", "which-llm"}},
}

var aliasRegex = regexp.MustCompile(`^\s*alias\s+([\w-]+)\s*=`)

func main() {
	appsFile := flag.String("apps", "apps.txt", "caminho para lista de apps")
	aliasesFile := flag.String("aliases", "aliases.txt", "caminho para lista de aliases")
	flag.Parse()

	if err := ensureApps(*appsFile); err != nil {
		fmt.Fprintf(os.Stderr, "erro ao instalar apps: %v\n", err)
		os.Exit(1)
	}

	if err := configureAliases(*aliasesFile); err != nil {
		fmt.Fprintf(os.Stderr, "erro ao configurar aliases: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ setup concluído com sucesso")
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, s.Err()
}

func ensureApps(appsPath string) error {
	apps, err := readLines(appsPath)
	if err != nil {
		return fmt.Errorf("lendo %s: %w", appsPath, err)
	}

	for _, name := range apps {
		info, ok := appMap[name]
		if !ok {
			fmt.Printf("⚠  %s: sem mapeamento, verificando como '%s'...\n", name, name)
			info = appInfo{Binary: name, Pkg: name}
		}

		if isInstalled(info.Binary) {
			fmt.Printf("✓ %s\n", name)
			continue
		}

		fmt.Printf("→ %s: instalando...\n", name)
		if err := installApp(name, info); err != nil {
			fmt.Fprintf(os.Stderr, "✗ %s: %v\n", name, err)
		}
	}
	return nil
}

func isInstalled(binary string) bool {
	return exec.Command("sh", "-c", "command -v "+binary).Run() == nil
}

func installApp(name string, info appInfo) error {
	hasSudo := exec.Command("sudo", "-n", "true").Run() == nil

	if info.Pkg != "" && hasSudo {
		cmd := exec.Command("sudo", "apt", "install", "-y", info.Pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
		fmt.Printf("  apt falhou, tentando alternativa...\n")
	}

	if info.InstallCmd != nil {
		cmd := exec.Command(info.InstallCmd[0], info.InstallCmd[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if !hasSudo {
		return fmt.Errorf("precisa de sudo para instalar via apt, e não há método alternativo")
	}
	return fmt.Errorf("apt falhou e não há método alternativo")
}

func configureAliases(aliasesPath string) error {
	desired, err := parseAliasesFile(aliasesPath)
	if err != nil {
		return fmt.Errorf("lendo %s: %w", aliasesPath, err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	aliasesFile := filepath.Join(home, ".opencode_aliases")
	if err := writeManagedAliases(aliasesFile, desired); err != nil {
		return fmt.Errorf("escrevendo %s: %w", aliasesFile, err)
	}

	rcFile := filepath.Join(home, ".zshrc")
	if _, err := os.Stat(rcFile); os.IsNotExist(err) {
		rcFile = filepath.Join(home, ".bashrc")
	}

	rcData, err := os.ReadFile(rcFile)
	if err != nil {
		return fmt.Errorf("lendo %s: %w", rcFile, err)
	}

	lines := strings.Split(string(rcData), "\n")
	currentAliases := extractAliases(lines)

	toRemove := make(map[string]bool)
	for name := range currentAliases {
		if _, ok := desired[name]; !ok {
			toRemove[name] = true
		}
	}

	var newLines []string
	inManaged := false
	skipNext := false
	for i, line := range lines {
		if skipNext {
			skipNext = false
			continue
		}

		trimmed := strings.TrimSpace(line)

		if trimmed == "# >>> opencode aliases >>>" {
			inManaged = true
			continue
		}
		if trimmed == "# <<< opencode aliases <<<" {
			inManaged = false
			continue
		}
		if inManaged {
			continue
		}

		if strings.TrimSpace(line) == fmt.Sprintf("source %s", aliasesFile) {
			continue
		}

		if m := aliasRegex.FindStringSubmatch(trimmed); m != nil {
			name := m[1]
			if _, exists := desired[name]; exists {
				if i+1 < len(lines) && strings.HasPrefix(lines[i+1], "\\") {
					skipNext = true
				}
				continue
			}
			if toRemove[name] {
				if i+1 < len(lines) && strings.HasPrefix(lines[i+1], "\\") {
					skipNext = true
				}
				continue
			}
			if i+1 < len(lines) && strings.HasPrefix(lines[i+1], "\\") {
				skipNext = true
			}
		}

		newLines = append(newLines, line)
	}

	newLines = append(newLines,
		"",
		"# >>> opencode aliases >>>",
		fmt.Sprintf("source %s", aliasesFile),
		"# <<< opencode aliases <<<",
	)

	if err := os.WriteFile(rcFile, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("escrevendo %s: %w", rcFile, err)
	}

	removed := make([]string, 0, len(toRemove))
	for name := range toRemove {
		removed = append(removed, name)
	}
	if len(removed) > 0 {
		fmt.Printf("↻ aliases removidos do shell: %s\n", strings.Join(removed, ", "))
	}

	fmt.Printf("✓ %d aliases configurados em %s\n", len(desired), aliasesFile)
	return nil
}

func parseAliasesFile(path string) (map[string]string, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}

	aliases := make(map[string]string, len(lines))
	for _, line := range lines {
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		name := line[:eq]
		val := line[eq+1:]
		val = strings.Trim(val, `"'`)
		aliases[name] = val
	}
	return aliases, nil
}

func writeManagedAliases(path string, aliases map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var names []string
	for name := range aliases {
		names = append(names, name)
	}

	for _, name := range names {
		fmt.Fprintf(f, "alias %s='%s'\n", name, aliases[name])
	}
	return nil
}

func extractAliases(lines []string) map[string]string {
	aliases := make(map[string]string)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		m := aliasRegex.FindStringSubmatch(trimmed)
		if m == nil {
			continue
		}
		name := m[1]
		eq := strings.IndexByte(trimmed, '=')
		if eq < 0 {
			continue
		}
		val := strings.TrimSpace(trimmed[eq+1:])
		val = strings.Trim(val, `"'`)
		aliases[name] = val
	}
	return aliases
}
