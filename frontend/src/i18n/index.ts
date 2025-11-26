import { createI18n } from 'vue-i18n';
import type { I18n } from 'vue-i18n';
import type { TranslationMessages, SupportedLocale } from './types';

// Import translation messages
import en from './locales/en';
import zh from './locales/zh';

const i18n: I18n<
  { en: TranslationMessages; zh: TranslationMessages },
  unknown,
  unknown,
  SupportedLocale,
  false
> = createI18n({
  legacy: false, // Use Composition API mode
  locale: (localStorage.getItem('locale') || 'en') as SupportedLocale,
  fallbackLocale: 'en' as SupportedLocale,
  messages: {
    en,
    zh,
  },
});

export default i18n;
