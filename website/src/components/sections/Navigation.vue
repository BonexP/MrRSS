<template>
  <nav
    class="fixed w-full z-50 transition-all duration-300 bg-white/80 backdrop-blur-md border-b border-gray-100"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <img
            src="/assets/logo.svg"
            alt="MrRSS Logo"
            class="h-8 w-8 mr-3 hover:rotate-12 transition-transform duration-300"
          />
          <span
            class="font-bold text-xl bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-sky-950"
            >MrRSS</span
          >
        </div>
        <div class="hidden md:flex items-center space-x-8">
          <a
            href="#features"
            class="text-gray-600 hover:text-blue-600 font-medium transition"
            data-aos="fade-in"
            >{{ $t('nav.features') }}</a
          >
          <a
            href="#download"
            class="text-gray-600 hover:text-blue-600 font-medium transition"
            data-aos="fade-in"
            >{{ $t('nav.download') }}</a
          >
          <a
            href="https://github.com/WCY-dt/MrRSS"
            target="_blank"
            class="text-gray-600 hover:text-blue-600 font-medium transition"
            rel="noopener noreferrer"
            >{{ $t('nav.github') }}</a
          >
          <button
            @click="switchLanguage"
            class="text-gray-600 hover:text-blue-600 font-medium transition cursor-pointer"
          >
            {{ currentLangText }}
          </button>
          <a
            href="https://github.com/WCY-dt/MrRSS/releases"
            class="bg-gradient-to-r from-blue-600 to-sky-950 text-white px-5 py-2 rounded-full hover:shadow-lg hover:scale-105 transition duration-300 font-medium"
            >{{ $t('nav.getStarted') }}</a
          >
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

interface Props {
  onSwitchLanguage: () => void;
}

const props = defineProps<Props>();

const { locale } = useI18n();

const currentLangText = computed(() => (locale.value === 'en' ? '中文' : 'English'));

function switchLanguage() {
  props.onSwitchLanguage();
}
</script>
