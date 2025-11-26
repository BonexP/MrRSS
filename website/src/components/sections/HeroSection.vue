<template>
  <header class="relative pt-32 pb-20 lg:pt-48 lg:pb-32 overflow-hidden">
    <div class="absolute inset-0 bg-gradient-to-br from-blue-50 via-white to-indigo-50 -z-10"></div>
    <div
      class="absolute top-0 right-0 -mr-20 -mt-20 w-96 h-96 rounded-full bg-blue-100 blur-3xl opacity-50 animate-pulse"
    ></div>
    <div
      class="absolute bottom-0 left-0 -ml-20 -mb-20 w-80 h-80 rounded-full bg-indigo-100 blur-3xl opacity-50 animate-pulse"
      style="animation-delay: 1s"
    ></div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative">
      <div class="lg:grid lg:grid-cols-12 lg:gap-16 items-center">
        <div
          class="lg:col-span-6 text-center lg:text-left"
          data-aos="fade-right"
          data-aos-duration="1000"
        >
          <div
            class="inline-flex items-center px-3 py-1 rounded-full bg-blue-100 text-blue-600 text-sm font-medium mb-6"
          >
            <span class="flex h-2 w-2 rounded-full bg-blue-600 mr-2"></span>
            <span>{{ $t('hero.badge') }}</span>
          </div>
          <h1
            class="text-5xl lg:text-6xl font-extrabold tracking-tight text-gray-900 mb-6 leading-tight"
          >
            <span>{{ $t('hero.title1') }}</span> <br />
            <span class="bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-sky-950">{{
              $t('hero.title2')
            }}</span>
          </h1>
          <p class="text-xl text-gray-600 mb-8 leading-relaxed">
            {{ $t('hero.subtitle') }}
          </p>
          <div class="flex flex-col sm:flex-row justify-center lg:justify-start gap-4">
            <a
              href="#download"
              class="inline-flex items-center justify-center px-8 py-4 text-base font-bold text-white transition-all duration-200 bg-blue-600 border border-transparent rounded-xl hover:bg-blue-700 hover:shadow-lg hover:-translate-y-1 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600"
            >
              {{ $t('hero.downloadNow') }}
            </a>
            <a
              href="https://github.com/WCY-dt/MrRSS"
              target="_blank"
              class="inline-flex items-center justify-center px-8 py-4 text-base font-bold text-gray-700 transition-all duration-200 bg-white border border-gray-200 rounded-xl hover:bg-gray-50 hover:border-gray-300 hover:shadow-md hover:-translate-y-1 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-200"
              rel="noopener noreferrer"
            >
              <IconGithub class="w-5 h-5 mr-2" />
              <span>{{ $t('hero.starOnGithub') }}</span>
            </a>
          </div>
        </div>
        <div
          class="lg:col-span-6 mt-12 lg:mt-0 relative"
          data-aos="fade-left"
          data-aos-duration="1000"
          data-aos-delay="200"
        >
          <div
            class="relative rounded-2xl shadow-2xl bg-gray-900 p-2 border border-gray-800 transform rotate-2 hover:rotate-0 transition-transform duration-500"
          >
            <div class="absolute top-0 left-0 w-full h-full rounded-2xl pointer-events-none"></div>
            <img
              class="rounded-xl w-full"
              src="/assets/example.png"
              alt="MrRSS Application Screenshot"
            />
          </div>
          <!-- Decorative elements -->
          <div
            class="absolute -bottom-10 -left-10 w-24 h-24 bg-yellow-400 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob"
          ></div>
          <div
            class="absolute -top-10 -right-10 w-24 h-24 bg-purple-400 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000"
          ></div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import IconGithub from '@/components/icons/IconGithub.vue';
</script>
