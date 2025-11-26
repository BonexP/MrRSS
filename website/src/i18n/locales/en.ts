import type { WebsiteMessages } from '../types';

const en: WebsiteMessages = {
  nav: {
    features: 'Features',
    download: 'Download',
    github: 'GitHub',
    getStarted: 'Get Started',
  },
  hero: {
    badge: 'New version is now available',
    title1: 'Your Feed,',
    title2: 'Reimagined.',
    subtitle:
      'Experience RSS reading with intelligent feed discovery, automatic translation, and complete privacy protection. Built for Windows, macOS, and Linux.',
    downloadNow: 'Download Now',
    starOnGithub: 'Star on GitHub',
  },
  features: {
    sectionTitle: 'Features',
    heading: 'Why choose MrRSS?',
    subtitle:
      'Designed for simplicity and power, giving you control over your news feed with modern tools.',
    autoTranslation: {
      title: 'Auto-Translation',
      description:
        'Instantly translate article titles and content using Google Translate or DeepL API. Break down language barriers and access global content effortlessly.',
    },
    smartDiscovery: {
      title: 'Smart Feed Discovery',
      description:
        'Automatically discover new feeds from friend links and related blogs. Expand your reading list intelligently with one click.',
    },
    privacy: {
      title: 'Privacy First',
      description:
        'Your data stays completely local with SQLite storage. No cloud sync, no tracking, no telemetry. Full control over your reading habits and personal information.',
    },
    crossPlatform: {
      title: 'Cross-Platform Native',
      description:
        'Built with Wails and Go for native performance on Windows, macOS, and Linux. Your data stays local with SQLite - no cloud tracking.',
    },
    keyboardShortcuts: {
      title: 'Rich Keyboard Shortcuts',
      description:
        'Navigate through articles, mark as read/unread, and manage feeds using comprehensive keyboard shortcuts. Power users can stay productive without ever touching the mouse.',
    },
    automation: {
      title: 'Automation Rules',
      description:
        'Set up intelligent rules to auto-filter, categorize, and manage articles based on keywords, sources, or custom conditions. Let automation handle the routine work.',
    },
  },
  download: {
    heading: 'Ready to get started?',
    subtitle: 'Download MrRSS for your operating system today and take control of your news feed.',
    windows: {
      title: 'Windows',
      subtitle: 'Windows 10/11 (64-bit)',
      button: 'Download .exe',
    },
    macos: {
      title: 'macOS',
      subtitle: 'Universal (Intel & Apple Silicon)',
      button: 'Download .dmg',
    },
    linux: {
      title: 'Linux',
      subtitle: 'AppImage / .deb / .rpm',
      button: 'Download',
    },
  },
  footer: {
    copyright: '© 2025 MrRSS Team. Open Source under GPL 3.0 License.',
    madeWith: 'Made with ❤️ by Ch3nyang',
  },
};

export default en;
