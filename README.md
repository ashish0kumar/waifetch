<h1 align="center">waifetch</h1>

<p align="center">A fun CLI that displays random anime girls holding programming books in your terminal</p>

<p align="center">
    <img width="70%" alt="showcase" src="https://github.com/user-attachments/assets/ecf67454-fe68-4eba-bfd1-4332165755c4" />
</p>

---

## Features

- **Dynamic terminal display** with images rendered directly in your terminal
- **Smart sizing** that automatically adapts to your terminal dimensions
- **Language-specific fetching** with support for all major programming languages
- **High-quality rendering** using chafa with 256 colors and ordered dithering

## Installation

### Required Dependencies

For image rendering, install [`chafa`](https://github.com/hpjansson/chafa):

> [!IMPORTANT]
> `chafa` is required for displaying images in the terminal. The CLI will not work without it.

```bash
# Ubuntu/Debian
sudo apt install chafa

# macOS
brew install chafa

# Arch Linux
sudo pacman -S chafa

# Fedora/RHEL
sudo dnf install chafa

# Nix
nix-env -iA nixpkgs.chafa
```

### Via `go install`

```bash
go install github.com/ashish0kumar/waifetch@latest
```

### Build from Source

```bash
git clone https://github.com/ashish0kumar/waifetch.git
cd waifetch/
go build
sudo mv waifetch /usr/local/bin/
waifetch --help
```

## Usage

```bash
# Random image from any language
waifetch

# List all available languages
waifetch --list

# Specific languages
waifetch --lang Go
waifetch --lang "C++"
waifetch -l Rust

# Get help
waifetch --help
```

## Configuration

### GitHub API Rate Limits

For higher rate limits, set up a GitHub Personal Access Token:

1. Go to [GitHub Settings → Developer settings → Personal access tokens](https://github.com/settings/tokens)
2. Click **"Generate new token (classic)"**
3. Give it a name like "waifetch-cli"
4. **No scopes needed** for public repositories
5. Copy the token and set it as an environment variable:
```bash
export GITHUB_TOKEN="your_token_here"
waifetch
```

### Rate Limits

| Authentication Method | Rate Limit |
| :-- | :-- |
| **Unauthenticated** | 60 requests/hour |
| **Personal Access Token** | 5,000 requests/hour |

## Contributing

Contributions are always welcome!
Please feel free to open an issue or submit a pull request.

## Dependencies

- [**chafa**](https://hpjansson.org/chafa/) - High-quality terminal image rendering
- [**Cobra**](https://github.com/spf13/cobra) - CLI framework and command structure
- [**golang.org/x/term**](https://pkg.go.dev/golang.org/x/term) - Terminal size detection

## Data Source

Images are fetched from the excellent [Anime-Girls-Holding-Programming-Books](https://github.com/cat-milk/Anime-Girls-Holding-Programming-Books) repository by [cat-milk](https://github.com/cat-milk), which contains hundreds of high-quality images of anime girls featuring various programming books.

## Acknowledgments

- Thanks to [cat-milk](https://github.com/cat-milk) for the amazing Anime-Girls-Holding-Programming-Books repository
- Built with love for the intersection of anime culture and programming

<br>
<p align="center">
<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/footers/gray0_ctp_on_line.svg?sanitize=true" />
</p>
<p align="center">
    <i><code>&copy 2025-present <a href="https://github.com/ashish0kumar">Ashish Kumar</a></code></i>
</p>
<div align="center">
<a href="https://github.com/ashish0kumar/waifetch/blob/main/LICENSE"><img src="https://img.shields.io/github/license/ashish0kumar/waifetch?style=for-the-badge&color=CBA6F7&logoColor=cdd6f4&labelColor=302D41" alt="LICENSE"></a>&nbsp;&nbsp;
</div>
