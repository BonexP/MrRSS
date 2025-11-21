<div align="center">
  <img src="frontend/logo.svg" alt="MrRSS Logo" width="120" height="120">
  <h1>MrRSS</h1>
  <p><strong>A Modern, Cross-Platform Desktop RSS Reader</strong></p>
  
  [![Version](https://img.shields.io/badge/version-1.1.0-blue.svg)](https://github.com/WCY-dt/MrRSS/releases)
  [![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
  [![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
  [![Wails](https://img.shields.io/badge/Wails-v2-red)](https://wails.io/)
  [![Vue.js](https://img.shields.io/badge/Vue.js-3.x-4FC08D?logo=vue.js)](https://vuejs.org/)
  
  <p>Built with Wails, Go, and Vue.js</p>
</div>

---

## ✨ Features

- 🖥️ **Cross-Platform**: Native desktop app for Windows, macOS, and Linux
- 📰 **RSS/Atom Support**: Subscribe to unlimited feeds
- 🌐 **Auto-Translation**: Translate article titles using Google Translate or DeepL API
- 📁 **Smart Organization**: Organize feeds into categories that expand by default
- ⭐ **Favorites & Reading Tracking**: Save articles and track read/unread status
- 🎨 **Modern UI**: Clean, responsive interface with dark mode support
- 🌍 **Multi-Language**: Full support for English and Chinese (中文)
- 📦 **OPML Import/Export**: Easy migration from other RSS readers
- 🔄 **Auto-Refresh**: Configurable automatic feed updates
- 🧹 **Smart Cleanup**: Optional auto-cleanup of old articles
- 🎯 **Progressive Loading**: Articles update progressively as feeds are fetched
- 🚀 **Viewport Translation**: Auto-translate articles as they enter the viewport

## 📸 Screenshots

<div align="center">
  <p><em>Coming soon - screenshots will be added in future releases</em></p>
</div>

## 🚀 Quick Start

### Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://go.dev/) (1.21 or higher)
- [Node.js](https://nodejs.org/) (LTS version with npm)
- [Wails](https://wails.io/docs/gettingstarted/installation) CLI

### Installation

1. **Clone the repository**
   \`\`\`bash
   git clone https://github.com/WCY-dt/MrRSS.git
   cd MrRSS
   \`\`\`

2. **Install frontend dependencies**
   \`\`\`bash
   cd frontend
   npm install
   cd ..
   \`\`\`

3. **Build the application**
   \`\`\`bash
   wails build
   \`\`\`
   
   The executable will be created in the \`build/bin\` directory.

4. **Run the application**
   - Windows: \`build/bin/MrRSS.exe\`
   - macOS: \`build/bin/MrRSS.app\`
   - Linux: \`build/bin/MrRSS\`

## 🛠️ Development

### Running in Development Mode

Start the application with hot reloading:

\`\`\`bash
wails dev
\`\`\`

This will:
- Start the Go backend with live reload
- Start the Vue.js frontend dev server
- Open the application in a development window

### Building for Production

\`\`\`bash
# Build for current platform
wails build

# Build with specific flags
wails build -clean -ldflags "-s -w"
\`\`\`

For Windows-specific builds, see [BUILD_WINDOWS.md](BUILD_WINDOWS.md).

### Running Tests

\`\`\`bash
# Backend tests
go test ./...

# Frontend tests
cd frontend
npm test
\`\`\`

## 📁 Project Structure

\`\`\`
MrRSS/
├── main.go                    # Application entry point
├── wails.json                 # Wails configuration
├── internal/                  # Backend Go code
│   ├── database/             # SQLite database layer
│   ├── feed/                 # RSS/Atom feed fetching & parsing
│   ├── handlers/             # Application logic handlers
│   ├── models/               # Data models
│   ├── opml/                 # OPML import/export
│   └── translation/          # Translation services
├── frontend/                  # Vue.js frontend
│   ├── src/
│   │   ├── components/       # Vue components
│   │   ├── store.js          # Global state management
│   │   ├── i18n.js           # Internationalization
│   │   └── App.vue           # Root component
│   └── wailsjs/              # Auto-generated Go bindings
└── build/                     # Build output directory
\`\`\`

## 🎯 Usage

### Adding Feeds

1. Click the **+** button in the sidebar
2. Enter the feed URL
3. Optionally set a custom title and category
4. Click **Add Feed**

### Managing Feeds

- **Right-click** on a feed to edit or unsubscribe
- **Right-click** on a category to rename it
- Use the **Settings → Feeds** tab to bulk manage subscriptions

### Keyboard Shortcuts

- \`Ctrl/Cmd + R\` - Refresh all feeds
- \`Esc\` - Close modals
- Arrow keys - Navigate articles

## ⚙️ Configuration

Access settings via the gear icon in the sidebar:

### General Settings

- **Dark Mode**: Toggle between light and dark themes
- **Language**: Choose between English and Chinese
- **Auto-update Interval**: Set how often feeds refresh (in minutes)
- **Auto Cleanup**: Automatically remove old articles

### Translation Settings

- **Enable Translation**: Toggle automatic translation
- **Provider**: Choose between Google Translate (free) or DeepL API
- **Target Language**: Select translation target language
- **DeepL API Key**: Required if using DeepL

### Data Management

- **Import OPML**: Import subscriptions from other RSS readers
- **Export OPML**: Export your subscriptions
- **Clean Database**: Remove all articles except read and favorited ones

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

Before contributing:
1. Read the [Code of Conduct](CODE_OF_CONDUCT.md)
2. Check existing issues or create a new one
3. Fork the repository and create a feature branch
4. Make your changes and add tests
5. Submit a pull request

## 🔒 Security

If you discover a security vulnerability, please follow our [Security Policy](SECURITY.md).

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Wails](https://wails.io/) - Go + Web framework
- UI powered by [Vue.js 3](https://vuejs.org/) and [Tailwind CSS](https://tailwindcss.com/)
- Icons from [Phosphor Icons](https://phosphoricons.com/)
- RSS parsing with [gofeed](https://github.com/mmcdole/gofeed)

## 📮 Contact & Support

- **Issues**: [GitHub Issues](https://github.com/WCY-dt/MrRSS/issues)
- **Discussions**: [GitHub Discussions](https://github.com/WCY-dt/MrRSS/discussions)
- **Repository**: [github.com/WCY-dt/MrRSS](https://github.com/WCY-dt/MrRSS)

---

<div align="center">
  <p>Made with ❤️ by the MrRSS Team</p>
  <p>⭐ Star us on GitHub if you find this project useful!</p>
</div>
