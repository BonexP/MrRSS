import { createI18n } from 'vue-i18n';
import type { SupportedLocale } from './types';

// Import translation messages
import en from './locales/en';
import zh from './locales/zh';

const i18n = createI18n({
  legacy: false, // Use Composition API mode
  locale: (navigator.language.startsWith('zh') ? 'zh' : 'en') as SupportedLocale,
  fallbackLocale: 'en' as SupportedLocale,
  messages: {
    en,
    zh,
  },
});

export default i18n;
