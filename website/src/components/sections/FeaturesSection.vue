<template>
  <section id="features" class="py-24 bg-white relative overflow-hidden">
    <div
      class="absolute top-0 left-0 w-full h-full bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMSIgY3k9IjEiIHI9IjEiIGZpbGw9IiNlNWU3ZWIiLz48L3N2Zz4=')] opacity-30"
    ></div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative">
      <div class="text-center max-w-3xl mx-auto mb-16" data-aos="fade-up">
        <h2 class="text-base text-blue-600 font-semibold tracking-wide uppercase">
          {{ $t('features.sectionTitle') }}
        </h2>
        <p class="mt-2 text-4xl font-extrabold tracking-tight text-gray-900 sm:text-5xl">
          {{ $t('features.heading') }}
        </p>
        <p class="mt-4 text-xl text-gray-500">
          {{ $t('features.subtitle') }}
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <!-- Feature 1 - Auto-Translation -->
        <div
          class="relative p-8 bg-white rounded-2xl shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group"
          data-aos="fade-up"
          data-aos-delay="100"
        >
          <div
            class="absolute top-0 right-0 -mt-4 -mr-4 w-24 h-24 bg-blue-50 rounded-full opacity-50 group-hover:scale-110 transition-transform duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="flex items-center justify-center h-14 w-14 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 text-white shadow-lg mb-6 group-hover:rotate-6 transition-transform duration-300"
            >
              <IconLanguages />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">
              {{ $t('features.autoTranslation.title') }}
            </h3>
            <p class="text-gray-500 leading-relaxed">
              {{ $t('features.autoTranslation.description') }}
            </p>
          </div>
        </div>

        <!-- Feature 2 - Smart Feed Discovery -->
        <div
          class="relative p-8 bg-white rounded-2xl shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group"
          data-aos="fade-up"
          data-aos-delay="200"
        >
          <div
            class="absolute top-0 right-0 -mt-4 -mr-4 w-24 h-24 bg-emerald-50 rounded-full opacity-50 group-hover:scale-110 transition-transform duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="flex items-center justify-center h-14 w-14 rounded-xl bg-gradient-to-br from-emerald-500 to-emerald-600 text-white shadow-lg mb-6 group-hover:rotate-6 transition-transform duration-300"
            >
              <IconSearch />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">
              {{ $t('features.smartDiscovery.title') }}
            </h3>
            <p class="text-gray-500 leading-relaxed">
              {{ $t('features.smartDiscovery.description') }}
            </p>
          </div>
        </div>

        <!-- Feature 3 - Privacy First -->
        <div
          class="relative p-8 bg-white rounded-2xl shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group"
          data-aos="fade-up"
          data-aos-delay="300"
        >
          <div
            class="absolute top-0 right-0 -mt-4 -mr-4 w-24 h-24 bg-purple-50 rounded-full opacity-50 group-hover:scale-110 transition-transform duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="flex items-center justify-center h-14 w-14 rounded-xl bg-gradient-to-br from-purple-500 to-purple-600 text-white shadow-lg mb-6 group-hover:rotate-6 transition-transform duration-300"
            >
              <IconShield />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">{{ $t('features.privacy.title') }}</h3>
            <p class="text-gray-500 leading-relaxed">
              {{ $t('features.privacy.description') }}
            </p>
          </div>
        </div>

        <!-- Feature 4 - Cross-Platform -->
        <div
          class="relative p-8 bg-white rounded-2xl shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group"
          data-aos="fade-up"
          data-aos-delay="400"
        >
          <div
            class="absolute top-0 right-0 -mt-4 -mr-4 w-24 h-24 bg-amber-50 rounded-full opacity-50 group-hover:scale-110 transition-transform duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="flex items-center justify-center h-14 w-14 rounded-xl bg-gradient-to-br from-amber-500 to-amber-600 text-white shadow-lg mb-6 group-hover:rotate-6 transition-transform duration-300"
            >
              <IconMonitor />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">
              {{ $t('features.crossPlatform.title') }}
            </h3>
            <p class="text-gray-500 leading-relaxed">
              {{ $t('features.crossPlatform.description') }}
            </p>
          </div>
        </div>

        <!-- Feature 5 - Keyboard Shortcuts -->
        <div
          class="relative p-8 bg-white rounded-2xl shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group"
          data-aos="fade-up"
          data-aos-delay="500"
        >
          <div
            class="absolute top-0 right-0 -mt-4 -mr-4 w-24 h-24 bg-rose-50 rounded-full opacity-50 group-hover:scale-110 transition-transform duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="flex items-center justify-center h-14 w-14 rounded-xl bg-gradient-to-br from-rose-500 to-rose-600 text-white shadow-lg mb-6 group-hover:rotate-6 transition-transform duration-300"
            >
              <IconKeyboard />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">
              {{ $t('features.keyboardShortcuts.title') }}
            </h3>
            <p class="text-gray-500 leading-relaxed">
              {{ $t('features.keyboardShortcuts.description') }}
            </p>
          </div>
        </div>

        <!-- Feature 6 - Automation Rules -->
        <div
          class="relative p-8 bg-white rounded-2xl shadow-lg border border-gray-100 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 group"
          data-aos="fade-up"
          data-aos-delay="600"
        >
          <div
            class="absolute top-0 right-0 -mt-4 -mr-4 w-24 h-24 bg-cyan-50 rounded-full opacity-50 group-hover:scale-110 transition-transform duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="flex items-center justify-center h-14 w-14 rounded-xl bg-gradient-to-br from-cyan-500 to-cyan-600 text-white shadow-lg mb-6 group-hover:rotate-6 transition-transform duration-300"
            >
              <IconSettings />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">
              {{ $t('features.automation.title') }}
            </h3>
            <p class="text-gray-500 leading-relaxed">
              {{ $t('features.automation.description') }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import IconLanguages from '@/components/icons/IconLanguages.vue';
import IconSearch from '@/components/icons/IconSearch.vue';
import IconShield from '@/components/icons/IconShield.vue';
import IconMonitor from '@/components/icons/IconMonitor.vue';
import IconKeyboard from '@/components/icons/IconKeyboard.vue';
import IconSettings from '@/components/icons/IconSettings.vue';
</script>
