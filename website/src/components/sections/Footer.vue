<template>
  <footer class="bg-gray-900 text-white border-t border-gray-800">
    <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row justify-between items-center">
        <div class="flex items-center mb-4 md:mb-0">
          <img src="/assets/logo.svg" alt="MrRSS Logo" class="h-8 w-8 mr-3" />
          <span class="font-bold text-xl">MrRSS</span>
        </div>
        <div class="flex space-x-6">
          <a
            href="https://github.com/WCY-dt/MrRSS"
            class="text-gray-400 hover:text-white transition-colors"
          >
            <span class="sr-only">GitHub</span>
            <IconGithub class="h-6 w-6" />
          </a>
        </div>
      </div>
      <div class="mt-8 border-t border-gray-800 pt-8 md:flex md:items-center md:justify-between">
        <p class="text-base text-gray-400">
          {{ $t('footer.copyright') }}
        </p>
        <p class="mt-4 md:mt-0 text-base text-gray-500">
          {{ $t('footer.madeWith') }}
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import IconGithub from '@/components/icons/IconGithub.vue';
</script>
