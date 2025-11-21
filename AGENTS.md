# AI Agent Guidelines for MrRSS

This document provides comprehensive guidance for AI agents (like GitHub Copilot, ChatGPT, Claude, etc.) working on the MrRSS project.

## Project Overview

**MrRSS** is a modern, cross-platform desktop RSS reader built with:
- **Backend**: Go 1.21+ with Wails v2 framework
- **Frontend**: Vue.js 3 (Composition API) with Tailwind CSS
- **Database**: SQLite for local storage
- **Build Tool**: Wails CLI

### Core Functionality
- RSS/Atom feed subscription and parsing
- Article management (read/unread, favorites)
- Category-based organization
- Auto-translation (Google Translate/DeepL)
- OPML import/export
- Multi-language support (English/Chinese)
- Auto-refresh with configurable intervals

## Project Structure

\`\`\`
MrRSS/
├── main.go                      # Application entry point
├── wails.json                   # Wails configuration, version info
├── internal/                    # Backend Go code (not exposed directly)
│   ├── database/
│   │   └── sqlite.go           # SQLite operations, schema management
│   ├── feed/
│   │   └── fetcher.go          # RSS/Atom parsing with gofeed
│   ├── handlers/
│   │   └── handlers.go         # App logic, exposed to frontend
│   ├── models/
│   │   └── models.go           # Data structures
│   ├── opml/
│   │   └── opml.go             # OPML import/export
│   └── translation/
│       └── translator.go        # Translation services
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── modals/         # Modal dialogs
│   │   │   ├── ArticleList.vue
│   │   │   ├── ArticleDetail.vue
│   │   │   ├── Sidebar.vue
│   │   │   ├── ContextMenu.vue
│   │   │   └── Toast.vue
│   │   ├── store.js            # Global state management
│   │   ├── i18n.js             # Internationalization
│   │   ├── App.vue             # Root component
│   │   └── style.css           # Global styles
│   └── wailsjs/                # Auto-generated Go→JS bindings
└── test/                        # Test files
\`\`\`

## Key Technologies & Patterns

### Backend (Go)

**Wails Framework**:
- Use \`context.Context\` for all exported methods
- Methods are automatically exposed to frontend
- Use struct methods for organization

**Database**:
- SQLite with \`modernc.org/sqlite\` driver
- Use prepared statements
- Migrations removed (development phase)
- Default settings in \`Init()\`

**RSS Parsing**:
- Use \`github.com/mmcdole/gofeed\` for parsing
- Support both RSS and Atom formats
- Handle malformed feeds gracefully

**Translation**:
- Google Translate: Free, no API key
- DeepL: Requires API key
- Store translations in database

### Frontend (Vue.js)

**Vue 3 Composition API**:
\`\`\`vue
<script setup>
import { ref, computed, onMounted } from 'vue';
import { store } from '../store.js';

const items = ref([]);
const filteredItems = computed(() => /* ... */);

onMounted(async () => {
  // Initialize
});
</script>
\`\`\`

**State Management**:
- Use \`store.js\` for global state
- Reactive with Vue's \`reactive()\`
- Methods for state mutations

**Styling**:
- Tailwind CSS utility classes
- Theme variables in CSS
- Dark mode support via \`data-theme\`

**Internationalization**:
- \`i18n.js\` provides translation function
- Use \`store.i18n.t('key')\` in templates
- Support English (\`en\`) and Chinese (\`zh\`)

## Development Guidelines

### Code Style

**Go**:
- Follow \`gofmt\` formatting
- Use meaningful variable names
- Handle errors explicitly
- Add comments for exported functions
- Keep functions focused and small

**Vue.js**:
- Use \`<script setup>\` syntax
- Props validation with \`defineProps\`
- Emit declarations with \`defineEmits\`
- Composition API over Options API
- Keep components under 300 lines

**Commit Messages**:
Follow Conventional Commits:
\`\`\`
<type>(<scope>): <description>

[optional body]
\`\`\`

Types: \`feat\`, \`fix\`, \`docs\`, \`style\`, \`refactor\`, \`test\`, \`chore\`

### Testing

**Backend**:
\`\`\`bash
go test ./...
go test -cover ./internal/database
\`\`\`

**Frontend**:
\`\`\`bash
cd frontend
npm test
npm run lint
\`\`\`

**Manual Testing**:
\`\`\`bash
wails dev  # Development mode with hot reload
wails build  # Production build
\`\`\`

### Building

**Development**:
\`\`\`bash
wails dev
\`\`\`

**Production**:
\`\`\`bash
wails build -clean -ldflags "-s -w"
\`\`\`

**Platform-Specific**:
- Windows: See \`BUILD_WINDOWS.md\`
- macOS: Universal binary (Intel + Apple Silicon)
- Linux: Built for amd64

## Common Tasks

### Adding a New Feature

1. **Plan**: Define the feature scope and requirements
2. **Backend** (if needed):
   - Add methods to handlers
   - Update database schema if needed
   - Add tests
3. **Frontend**:
   - Create/update components
   - Add i18n strings
   - Update store if needed
4. **Test**: Manual + automated testing
5. **Document**: Update README if user-facing

### Adding Translations

1. Edit \`frontend/src/i18n.js\`
2. Add keys to both \`en\` and \`zh\` sections
3. Use \`store.i18n.t('newKey')\` in templates
4. Test language switching

### Database Changes

1. Edit \`internal/database/sqlite.go\`
2. Update \`Init()\` function with new schema
3. No migrations (development phase)
4. Users should back up data before updates

### UI Changes

1. Use existing Tailwind classes
2. Follow dark mode pattern: \`class="bg-bg-primary text-text-primary"\`
3. Ensure responsive design
4. Test in both themes
5. Add icons from Phosphor Icons (\`ph ph-*\`)

## Important Conventions

### Naming

**Go**:
- Exported: \`PascalCase\` (e.g., \`FetchFeed\`)
- Unexported: \`camelCase\` (e.g., \`parseXML\`)
- Interfaces: Usually noun (e.g., \`Reader\`, \`Handler\`)

**Vue**:
- Components: \`PascalCase\` files (e.g., \`ArticleList.vue\`)
- Props/emits: \`camelCase\` (e.g., \`feedId\`, \`onUpdate\`)
- CSS classes: \`kebab-case\` (e.g., \`article-card\`)

### File Organization

- One component per file
- Co-locate tests with code
- Group related functionality
- Keep files under 500 lines

### Error Handling

**Go**:
\`\`\`go
if err != nil {
    return fmt.Errorf("context: %w", err)
}
\`\`\`

**Vue**:
\`\`\`javascript
try {
    // operation
} catch (e) {
    console.error(e);
    window.showToast(store.i18n.t('error'), 'error');
}
\`\`\`

## Security Considerations

1. **Input Validation**: Validate all user inputs
2. **SQL Injection**: Use parameterized queries
3. **XSS**: Vue escapes by default, don't use \`v-html\`
4. **API Keys**: Store locally, never commit
5. **Feed Content**: Displayed in iframe with sandbox

## Performance Tips

1. **Backend**:
   - Use goroutines for concurrent feed fetching
   - Cache feed data when appropriate
   - Optimize database queries

2. **Frontend**:
   - Virtual scrolling for large lists
   - Lazy load images
   - Debounce search inputs
   - Use \`v-show\` vs \`v-if\` appropriately

## Debugging

**Backend**:
- Use \`fmt.Println\` or \`log.Printf\`
- Check Wails console output
- Use Go debugger (Delve)

**Frontend**:
- Browser DevTools console
- Vue DevTools extension
- Check Network tab for API calls

**Common Issues**:
1. **CORS**: Not applicable (native app)
2. **Feed parsing**: Check feed URL format
3. **Translation**: Verify API key for DeepL
4. **Build errors**: Check Go/Node versions

## Version Management

- Version in \`wails.json\`: \`"version"\` and \`"productVersion"\`
- Version in frontend: \`SettingsModal.vue\` About tab
- Follow Semantic Versioning (MAJOR.MINOR.PATCH)
- Current: 1.1.0

## Resources

- **Wails Docs**: https://wails.io/docs/
- **Vue.js Docs**: https://vuejs.org/
- **Tailwind CSS**: https://tailwindcss.com/
- **Go Docs**: https://golang.org/doc/
- **gofeed**: https://github.com/mmcdole/gofeed

## Getting Help

1. Check existing issues on GitHub
2. Read documentation (README, CONTRIBUTING)
3. Search discussions
4. Create new issue with template

---

**Remember**: When in doubt, follow existing patterns in the codebase. Consistency is key!
