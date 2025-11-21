# GitHub Copilot Instructions for MrRSS

## Project Context

MrRSS is a cross-platform desktop RSS reader built with Wails (Go + Vue.js). It emphasizes simplicity, privacy, and modern UI design.

## Tech Stack

- **Backend**: Go 1.21+, Wails v2, SQLite
- **Frontend**: Vue 3 (Composition API), Tailwind CSS
- **Tools**: npm, Wails CLI

## Code Patterns

### Backend (Go)

When writing Go code:

\`\`\`go
// Always use context for exported methods
func (h *Handler) MethodName(ctx context.Context, param string) (Result, error) {
    if param == "" {
        return Result{}, errors.New("param is required")
    }
    
    // Implementation
    return result, nil
}

// Use prepared statements for SQL
stmt, err := db.Prepare("SELECT * FROM table WHERE id = ?")
if err != nil {
    return nil, fmt.Errorf("prepare: %w", err)
}
defer stmt.Close()

// Handle errors explicitly
if err != nil {
    return nil, fmt.Errorf("operation failed: %w", err)
}
\`\`\`

### Frontend (Vue)

When writing Vue components:

\`\`\`vue
<script setup>
import { ref, computed, onMounted } from 'vue';
import { store } from '../store.js';

// Props with validation
const props = defineProps({
    item: { type: Object, required: true },
    isActive: { type: Boolean, default: false }
});

// Emits declaration
const emit = defineEmits(['update', 'delete']);

// Reactive state
const isLoading = ref(false);
const items = ref([]);

// Computed properties
const filteredItems = computed(() => {
    return items.value.filter(item => /* condition */);
});

// Methods
async function loadData() {
    isLoading.value = true;
    try {
        // API call via Wails bindings
        const data = await SomeBackendMethod();
        items.value = data;
    } catch (e) {
        console.error(e);
        window.showToast(store.i18n.t('errorLoadingData'), 'error');
    } finally {
        isLoading.value = false;
    }
}

// Lifecycle
onMounted(() => {
    loadData();
});
</script>

<template>
    <div class="container">
        <h2 class="text-lg font-semibold">{{ store.i18n.t('title') }}</h2>
        
        <div v-if="isLoading">{{ store.i18n.t('loading') }}</div>
        <div v-else-if="items.length === 0">{{ store.i18n.t('noItems') }}</div>
        <div v-else v-for="item in filteredItems" :key="item.id">
            {{ item.name }}
        </div>
    </div>
</template>

<style scoped>
.container {
    @apply p-4 bg-bg-primary rounded-lg;
}
</style>
\`\`\`

## Styling Guidelines

### Tailwind CSS Patterns

Use these semantic class combinations:

\`\`\`html
<!-- Buttons -->
<button class="btn-primary">Primary Action</button>
<button class="btn-secondary">Secondary Action</button>
<button class="btn-danger">Dangerous Action</button>

<!-- Cards -->
<div class="bg-bg-primary border border-border rounded-lg p-4">
    <h3 class="text-text-primary font-semibold">Title</h3>
    <p class="text-text-secondary text-sm">Description</p>
</div>

<!-- Inputs -->
<input class="input-field" type="text" />

<!-- Modal -->
<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
    <div class="bg-bg-primary w-full max-w-2xl rounded-2xl shadow-2xl border border-border">
        <!-- Content -->
    </div>
</div>
\`\`\`

### Theme Variables

Use CSS variables for theming:

\`\`\`css
/* Colors follow theme */
background-color: var(--color-bg-primary);
color: var(--color-text-primary);

/* Or use Tailwind classes */
class="bg-bg-primary text-text-primary"
\`\`\`

## Internationalization

Always use i18n for user-facing strings:

\`\`\`vue
<!-- Template -->
<h1>{{ store.i18n.t('welcome') }}</h1>
<button :title="store.i18n.t('clickToOpen')">

<!-- Script -->
window.showToast(store.i18n.t('successMessage'), 'success');
\`\`\`

To add new strings, edit \`frontend/src/i18n.js\`:

\`\`\`javascript
export const translations = {
    en: {
        // English translations
        newKey: 'New String',
    },
    zh: {
        // Chinese translations
        newKey: '新字符串',
    }
};
\`\`\`

## Common Patterns

### API Calls (Frontend → Backend)

\`\`\`javascript
// Import from generated bindings
import { MethodName } from './wailsjs/go/internal/handlers/Handler.js';

// Call backend method
try {
    const result = await MethodName(param1, param2);
    // Handle result
} catch (error) {
    console.error('Error:', error);
    window.showToast(store.i18n.t('error'), 'error');
}
\`\`\`

### Custom Events

\`\`\`javascript
// Dispatch event
window.dispatchEvent(new CustomEvent('event-name', {
    detail: { data: value }
}));

// Listen for event
window.addEventListener('event-name', (e) => {
    const data = e.detail.data;
    // Handle event
});
\`\`\`

### Toast Notifications

\`\`\`javascript
window.showToast(message, type);  // type: 'success' | 'error' | 'info' | 'warning'
\`\`\`

### Confirm Dialogs

\`\`\`javascript
const confirmed = await window.showConfirm(
    store.i18n.t('confirmTitle'),
    store.i18n.t('confirmMessage'),
    true  // isDanger (shows red button)
);

if (confirmed) {
    // Proceed with action
}
\`\`\`

## Database Operations

### Query Pattern

\`\`\`go
// Use prepared statements
func (db *Database) GetArticles(feedID int) ([]models.Article, error) {
    rows, err := db.conn.Query(\`
        SELECT id, title, url, content, published_at
        FROM articles
        WHERE feed_id = ?
        ORDER BY published_at DESC
    \`, feedID)
    if err != nil {
        return nil, fmt.Errorf("query: %w", err)
    }
    defer rows.Close()
    
    var articles []models.Article
    for rows.Next() {
        var article models.Article
        err := rows.Scan(&article.ID, &article.Title, &article.URL, &article.Content, &article.PublishedAt)
        if err != nil {
            return nil, fmt.Errorf("scan: %w", err)
        }
        articles = append(articles, article)
    }
    
    return articles, rows.Err()
}
\`\`\`

## Testing

### Backend Tests

\`\`\`go
func TestFunctionName(t *testing.T) {
    // Setup
    input := "test"
    expected := "result"
    
    // Execute
    result, err := FunctionName(input)
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
\`\`\`

### Frontend Tests

\`\`\`javascript
describe('Component', () => {
    it('should render correctly', () => {
        // Test implementation
    });
});
\`\`\`

## Documentation

When adding new features:

1. **Code Comments**: Document exported functions
   \`\`\`go
   // FetchFeed retrieves and parses an RSS feed from the given URL.
   // Returns an error if the URL is invalid or the feed cannot be parsed.
   func FetchFeed(url string) (*Feed, error) {
   \`\`\`

2. **README Updates**: Document user-facing features

3. **Changelog**: Update for releases

## Don'ts

❌ **Don't**:
- Use \`var\` in Vue (use \`ref\` or \`reactive\`)
- Hardcode strings (use i18n)
- Use inline styles (use Tailwind classes)
- Forget error handling
- Use \`any\` type without good reason
- Commit API keys or secrets
- Use \`v-html\` (XSS risk)
- Make breaking changes without discussion

## Do's

✅ **Do**:
- Use TypeScript-style JSDoc for better IDE support
- Follow existing code patterns
- Write tests for new features
- Keep functions small and focused
- Use meaningful variable names
- Handle edge cases
- Validate inputs
- Log errors appropriately
- Use semantic HTML

## File Naming

- Components: \`PascalCase.vue\` (e.g., \`ArticleList.vue\`)
- Go files: \`lowercase.go\` (e.g., \`fetcher.go\`)
- Test files: \`*_test.go\`
- Utilities: \`kebab-case.js\` (e.g., \`date-utils.js\`)

## Useful Commands

\`\`\`bash
# Development
wails dev                          # Run in dev mode
wails build                        # Build for production
wails build -clean                 # Clean build

# Testing
go test ./...                      # Run Go tests
cd frontend && npm test            # Run frontend tests

# Linting
go vet ./...                       # Go linter
cd frontend && npm run lint        # Frontend linter

# Dependencies
go mod tidy                        # Clean Go dependencies
cd frontend && npm install         # Install npm packages
\`\`\`

## Quick Reference

**Get current theme**: \`store.theme\` (returns \`'light'\` or \`'dark'\`)
**Get current language**: \`store.i18n.locale.value\` (returns \`'en'\` or \`'zh'\`)
**Global store**: \`import { store } from './store.js'\`
**Show notification**: \`window.showToast(message, type)\`
**Confirm action**: \`await window.showConfirm(title, message, isDanger)\`

---

When generating code, prioritize:
1. **Correctness**: Code that works
2. **Consistency**: Follow existing patterns
3. **Clarity**: Easy to understand
4. **Maintainability**: Easy to modify later
