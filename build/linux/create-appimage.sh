#!/bin/bash
# Script to create a Linux AppImage for MrRSS
#
# Application Information:
# Name: MrRSS
# Description: A Modern, Cross-Platform Desktop RSS Reader
# Publisher: MrRSS Team
# URL: https://github.com/WCY-dt/MrRSS
# Copyright: Copyright © MrRSS Team

set -e

APP_NAME="MrRSS"
# Get version from wails.json if available, otherwise use default
VERSION=$(grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' wails.json 2>/dev/null | head -1 | sed 's/.*"\([^"]*\)".*/\1/' || echo "1.1.0")
APP_PUBLISHER="MrRSS Team"
APP_URL="https://github.com/WCY-dt/MrRSS"
APP_DESCRIPTION="A Modern, Cross-Platform Desktop RSS Reader"
BUILD_DIR="build/bin"
APPDIR="build/appimage/${APP_NAME}.AppDir"
APPIMAGE_NAME="${APP_NAME}-${VERSION}-linux-amd64.AppImage"

echo "Creating AppImage for ${APP_NAME} ${VERSION}..."
echo "Publisher: ${APP_PUBLISHER}"
echo "Description: ${APP_DESCRIPTION}"
echo ""

# Check if binary exists
if [ ! -f "${BUILD_DIR}/${APP_NAME}" ]; then
    echo "Error: Binary not found at ${BUILD_DIR}/${APP_NAME}"
    echo "Please build the application first with: wails build -platform linux/amd64"
    exit 1
fi

# Create AppDir structure
echo "Creating AppDir structure..."
rm -rf "build/appimage"
mkdir -p "${APPDIR}/usr/bin"
mkdir -p "${APPDIR}/usr/share/applications"
mkdir -p "${APPDIR}/usr/share/icons/hicolor/256x256/apps"

# Copy binary
echo "Copying binary..."
cp "${BUILD_DIR}/${APP_NAME}" "${APPDIR}/usr/bin/"
chmod +x "${APPDIR}/usr/bin/${APP_NAME}"

# Create desktop file
echo "Creating desktop file..."
cat > "${APPDIR}/usr/share/applications/${APP_NAME}.desktop" << EOF
[Desktop Entry]
Type=Application
Name=${APP_NAME}
GenericName=RSS Reader
Comment=${APP_DESCRIPTION}
Exec=${APP_NAME}
Icon=${APP_NAME}
Categories=Network;News;Feed;
Terminal=false
StartupWMClass=${APP_NAME}
Keywords=RSS;Atom;Feed;News;Reader;
X-GNOME-UsesNotifications=true
EOF

# Create AppRun script
echo "Creating AppRun script..."
cat > "${APPDIR}/AppRun" << 'EOF'
#!/bin/bash
SELF=$(readlink -f "$0")
HERE=${SELF%/*}
export PATH="${HERE}/usr/bin:${PATH}"
export LD_LIBRARY_PATH="${HERE}/usr/lib:${LD_LIBRARY_PATH}"
exec "${HERE}/usr/bin/MrRSS" "$@"
EOF
chmod +x "${APPDIR}/AppRun"

# Copy icon (if exists, otherwise create placeholder)
if [ -f "imgs/logo.svg" ] && (command -v inkscape &> /dev/null || command -v convert &> /dev/null); then
    echo "Converting icon..."
    # If inkscape is available, convert SVG to PNG
    if command -v inkscape &> /dev/null; then
        inkscape "imgs/logo.svg" -o "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png" -w 256 -h 256
        cp "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png" "${APPDIR}/${APP_NAME}.png"
    elif command -v convert &> /dev/null; then
        convert -background none -size 256x256 "imgs/logo.svg" "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png"
        cp "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png" "${APPDIR}/${APP_NAME}.png"
    fi
elif [ -f "build/appicon.png" ]; then
    # Fallback to pre-built PNG icon from Wails build process
    echo "Using existing PNG icon..."
    cp "build/appicon.png" "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png"
    cp "build/appicon.png" "${APPDIR}/${APP_NAME}.png"
else
    echo "Warning: Neither inkscape nor ImageMagick found, and no pre-built icon available."
    echo "To fix: Install imagemagick (sudo apt-get install imagemagick) or place a 256x256 PNG icon at build/appicon.png"
fi

# Copy desktop file to root
cp "${APPDIR}/usr/share/applications/${APP_NAME}.desktop" "${APPDIR}/"

# Download appimagetool if not present
APPIMAGETOOL="build/appimagetool-x86_64.AppImage"
if [ ! -f "${APPIMAGETOOL}" ]; then
    echo "Downloading appimagetool..."
    wget -q "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage" -O "${APPIMAGETOOL}"
    chmod +x "${APPIMAGETOOL}"
fi

# Create AppImage
echo "Creating AppImage..."
rm -f "${BUILD_DIR}/${APPIMAGE_NAME}"
# Use --appimage-extract-and-run if FUSE is not available (e.g., in CI environments)
if [ -n "${CI}" ] || ! [ -e /dev/fuse ]; then
    echo "FUSE not available, using --appimage-extract-and-run mode"
    ARCH=x86_64 "${APPIMAGETOOL}" --appimage-extract-and-run "${APPDIR}" "${BUILD_DIR}/${APPIMAGE_NAME}"
else
    ARCH=x86_64 "${APPIMAGETOOL}" "${APPDIR}" "${BUILD_DIR}/${APPIMAGE_NAME}"
fi

# Clean up
rm -rf "build/appimage"

echo "AppImage created successfully: ${BUILD_DIR}/${APPIMAGE_NAME}"
echo ""
echo "Installation instructions:"
echo "1. Make the AppImage executable: chmod +x ${APPIMAGE_NAME}"
echo "2. Run the AppImage: ./${APPIMAGE_NAME}"
echo ""
echo "User data will be stored in: ~/.local/share/MrRSS/"
