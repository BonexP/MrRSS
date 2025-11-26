import type { WebsiteMessages } from '../types';

const zh: WebsiteMessages = {
  nav: {
    features: '功能特性',
    download: '下载',
    github: 'GitHub',
    getStarted: '开始使用',
  },
  hero: {
    badge: '新版本现已发布',
    title1: '您的订阅源，',
    title2: '焕然一新。',
    subtitle:
      '体验支持智能订阅源发现、自动翻译和完整隐私保护的 RSS 阅读。支持 Windows、macOS 和 Linux 平台。',
    downloadNow: '立即下载',
    starOnGithub: '在 GitHub 上 Star',
  },
  features: {
    sectionTitle: '功能特性',
    heading: '为什么选择 MrRSS？',
    subtitle: '简洁与强大的完美结合，让你完全掌控自己的新闻订阅源。',
    autoTranslation: {
      title: '自动翻译',
      description:
        '使用 Google 翻译或 DeepL API 即时翻译文章标题和内容。打破语言障碍，轻松访问全球内容。',
    },
    smartDiscovery: {
      title: '智能订阅源发现',
      description: '自动从友情链接和相关博客中发现新订阅源。一键智能扩展你的阅读列表。',
    },
    privacy: {
      title: '隐私至上',
      description:
        '你的数据完全存储在本地 SQLite 数据库中。无云同步、无追踪、无遥测。完全掌控你的阅读习惯和个人信息。',
    },
    crossPlatform: {
      title: '跨平台原生应用',
      description:
        '使用 Wails 和 Go 构建，在 Windows、macOS 和 Linux 上提供原生性能。数据本地存储，无云追踪。',
    },
    keyboardShortcuts: {
      title: '丰富的键盘快捷键',
      description:
        '使用全面的键盘快捷键浏览文章、标记已读/未读和管理订阅源。高级用户无需触碰鼠标即可保持高效。',
    },
    automation: {
      title: '自动化规则',
      description:
        '根据关键词、来源或自定义条件设置智能规则，自动过滤、分类和管理文章。让自动化处理日常工作。',
    },
  },
  download: {
    heading: '准备好开始了吗？',
    subtitle: '立即下载适用于你操作系统的 MrRSS，掌控你的新闻订阅源。',
    windows: {
      title: 'Windows',
      subtitle: 'Windows 10/11 (64位)',
      button: '下载 .exe',
    },
    macos: {
      title: 'macOS',
      subtitle: '通用版 (Intel 和 Apple Silicon)',
      button: '下载 .dmg',
    },
    linux: {
      title: 'Linux',
      subtitle: 'AppImage / .deb / .rpm',
      button: '下载',
    },
  },
  footer: {
    copyright: '© 2025 MrRSS 团队。GPL 3.0 许可证开源。',
    madeWith: '由 Ch3nyang 用 ❤️ 制作',
  },
};

export default zh;
