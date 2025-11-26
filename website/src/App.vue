<template>
  <div id="app" class="bg-gray-50 text-gray-900 font-sans antialiased overflow-x-hidden">
    <!-- Navigation -->
    <Navigation :onSwitchLanguage="switchLanguage" />

    <!-- Hero Section -->
    <HeroSection />

    <!-- Features Section -->
    <FeaturesSection />

    <!-- Download Section -->
    <DownloadSection />

    <!-- Footer -->
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import Navigation from './components/sections/Navigation.vue';
import HeroSection from './components/sections/HeroSection.vue';
import FeaturesSection from './components/sections/FeaturesSection.vue';
import DownloadSection from './components/sections/DownloadSection.vue';
import Footer from './components/sections/Footer.vue';

const { locale } = useI18n();

function switchLanguage() {
  locale.value = locale.value === 'en' ? 'zh' : 'en';
}

onMounted(() => {
  // Initialize AOS
  if (typeof window !== 'undefined' && (window as any).AOS) {
    (window as any).AOS.init({
      once: true,
      offset: 50,
      duration: 800,
      easing: 'ease-out-cubic',
    });
  }
});
</script>

<style scoped>
/* Custom Scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}

html {
  scroll-behavior: smooth;
}

@keyframes blob {
  0% {
    transform: translate(0px, 0px) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
  100% {
    transform: translate(0px, 0px) scale(1);
  }
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animation-delay-4000 {
  animation-delay: 4s;
}
</style>
