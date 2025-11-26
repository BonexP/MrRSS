export interface WebsiteMessages {
  nav: {
    features: string;
    download: string;
    github: string;
    getStarted: string;
  };
  hero: {
    badge: string;
    title1: string;
    title2: string;
    subtitle: string;
    downloadNow: string;
    starOnGithub: string;
  };
  features: {
    sectionTitle: string;
    heading: string;
    subtitle: string;
    autoTranslation: {
      title: string;
      description: string;
    };
    smartDiscovery: {
      title: string;
      description: string;
    };
    privacy: {
      title: string;
      description: string;
    };
    crossPlatform: {
      title: string;
      description: string;
    };
    keyboardShortcuts: {
      title: string;
      description: string;
    };
    automation: {
      title: string;
      description: string;
    };
  };
  download: {
    heading: string;
    subtitle: string;
    windows: {
      title: string;
      subtitle: string;
      button: string;
    };
    macos: {
      title: string;
      subtitle: string;
      button: string;
    };
    linux: {
      title: string;
      subtitle: string;
      button: string;
    };
  };
  footer: {
    copyright: string;
    madeWith: string;
  };
}

export type SupportedLocale = 'en' | 'zh';
