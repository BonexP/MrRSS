#!/bin/bash
# scripts/bump-version.sh - Bump version across all files

set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 <new-version>"
    echo "Example: $0 1.2.4"
    exit 1
fi

NEW_VERSION=$1

echo "🔄 Bumping version to $NEW_VERSION..."

# Update frontend/package.json
sed -i "s/\"version\": \"[^\"]*\"/\"version\": \"$NEW_VERSION\"/" frontend/package.json
echo "✅ Updated frontend/package.json"

# Update wails.json (version field)
sed -i "s/\"version\": \"[^\"]*\"/\"version\": \"$NEW_VERSION\"/" wails.json
echo "✅ Updated wails.json (version)"

# Update wails.json (productVersion field)
sed -i "s/\"productVersion\": \"[^\"]*\"/\"productVersion\": \"$NEW_VERSION\"/" wails.json
echo "✅ Updated wails.json (productVersion)"

# Update internal/version/version.go
sed -i "s/const Version = \"[^\"]*\"/const Version = \"$NEW_VERSION\"/" internal/version/version.go
echo "✅ Updated internal/version/version.go"

# Update frontend/src/components/modals/settings/AboutTab.vue
sed -i "s/appVersion: \"[^\"]*\"/appVersion: \"$NEW_VERSION\"/" frontend/src/components/modals/settings/AboutTab.vue
echo "✅ Updated AboutTab.vue"

# Update README.md badges
sed -i "s/version-[0-9]\+\.[0-9]\+\.[0-9]\+/version-$NEW_VERSION/" README.md
echo "✅ Updated README.md version badge"

# Update README_zh.md badges
sed -i "s/version-[0-9]\+\.[0-9]\+\.[0-9]\+/version-$NEW_VERSION/" README_zh.md
echo "✅ Updated README_zh.md version badge"

echo "🎉 Version bumped to $NEW_VERSION!"
echo "📝 Remember to update CHANGELOG.md"
