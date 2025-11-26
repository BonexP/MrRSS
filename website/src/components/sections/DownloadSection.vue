<template>
  <section id="download" class="py-24 bg-gray-50 relative">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center relative z-10">
      <h2 class="text-3xl font-extrabold text-gray-900 sm:text-4xl mb-4" data-aos="fade-up">
        {{ $t('download.heading') }}
      </h2>
      <p
        class="text-xl text-gray-500 mb-12 max-w-2xl mx-auto"
        data-aos="fade-up"
        data-aos-delay="100"
      >
        {{ $t('download.subtitle') }}
      </p>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto">
        <!-- Windows -->
        <a
          href="https://github.com/WCY-dt/MrRSS/releases"
          class="group relative bg-white p-8 rounded-2xl shadow-md hover:shadow-2xl transition-all duration-300 border border-gray-100 hover:-translate-y-2 flex flex-col items-center"
          data-aos="fade-up"
          data-aos-delay="200"
        >
          <div
            class="h-20 w-20 bg-blue-50 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-blue-600 transition-colors duration-300"
          >
            <IconWindows
              class="w-10 h-10 text-blue-600 group-hover:text-white transition-colors duration-300"
            />
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">{{ $t('download.windows.title') }}</h3>
          <p class="text-gray-500 text-sm mb-6">{{ $t('download.windows.subtitle') }}</p>
          <span
            class="mt-auto inline-flex items-center justify-center px-6 py-2 border border-transparent text-sm font-medium rounded-full text-blue-700 bg-blue-100 group-hover:bg-blue-600 group-hover:text-white transition-colors duration-300"
          >
            {{ $t('download.windows.button') }}
          </span>
        </a>

        <!-- macOS -->
        <a
          href="https://github.com/WCY-dt/MrRSS/releases"
          class="group relative bg-white p-8 rounded-2xl shadow-md hover:shadow-2xl transition-all duration-300 border border-gray-100 hover:-translate-y-2 flex flex-col items-center"
          data-aos="fade-up"
          data-aos-delay="300"
        >
          <div
            class="h-20 w-20 bg-gray-100 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-gray-800 transition-colors duration-300"
          >
            <IconMacos
              class="w-10 h-10 text-gray-800 group-hover:text-white transition-colors duration-300"
            />
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">{{ $t('download.macos.title') }}</h3>
          <p class="text-gray-500 text-sm mb-6">{{ $t('download.macos.subtitle') }}</p>
          <span
            class="mt-auto inline-flex items-center justify-center px-6 py-2 border border-transparent text-sm font-medium rounded-full text-gray-700 bg-gray-200 group-hover:bg-gray-800 group-hover:text-white transition-colors duration-300"
          >
            {{ $t('download.macos.button') }}
          </span>
        </a>

        <!-- Linux -->
        <a
          href="https://github.com/WCY-dt/MrRSS/releases"
          class="group relative bg-white p-8 rounded-2xl shadow-md hover:shadow-2xl transition-all duration-300 border border-gray-100 hover:-translate-y-2 flex flex-col items-center"
          data-aos="fade-up"
          data-aos-delay="400"
        >
          <div
            class="h-20 w-20 bg-yellow-50 rounded-2xl flex items-center justify-center mb-6 group-hover:bg-yellow-500 transition-colors duration-300"
          >
            <IconLinux
              class="w-10 h-10 text-yellow-600 group-hover:text-white transition-colors duration-300"
            />
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">{{ $t('download.linux.title') }}</h3>
          <p class="text-gray-500 text-sm mb-6">{{ $t('download.linux.subtitle') }}</p>
          <span
            class="mt-auto inline-flex items-center justify-center px-6 py-2 border border-transparent text-sm font-medium rounded-full text-yellow-700 bg-yellow-100 group-hover:bg-yellow-500 group-hover:text-white transition-colors duration-300"
          >
            {{ $t('download.linux.button') }}
          </span>
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import IconWindows from '@/components/icons/IconWindows.vue';
import IconMacos from '@/components/icons/IconMacos.vue';
import IconLinux from '@/components/icons/IconLinux.vue';
</script>
