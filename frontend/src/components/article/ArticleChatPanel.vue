<script setup lang="ts">
import { ref, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import { PhChatCircleText, PhX, PhPaperPlaneRight, PhSpinner } from '@phosphor-icons/vue';
import type { Article } from '@/types/models';

interface ChatMessage {
  role: 'user' | 'assistant';
  content: string;
}

interface Props {
  article: Article;
  articleContent: string;
  settings: { ai_chat_enabled: boolean };
}

const props = defineProps<Props>();

const emit = defineEmits<{
  close: [];
}>();

const { t } = useI18n();

const isOpen = ref(true);
const isLoading = ref(false);
const inputMessage = ref('');
const messages = ref<ChatMessage[]>([]);
const chatContainer = ref<HTMLElement | null>(null);
const isFirstMessage = ref(true); // Track if this is the first message in the conversation

// Resize functionality
const isResizing = ref(false);
const startX = ref(0);
const startY = ref(0);
const startWidth = ref(0);
const startHeight = ref(0);
const panelElement = ref<HTMLElement | null>(null);

function startResize(e: MouseEvent) {
  isResizing.value = true;
  startX.value = e.clientX;
  startY.value = e.clientY;

  // Get the chat panel element
  const panel = panelElement.value;
  if (panel) {
    const rect = panel.getBoundingClientRect();
    startWidth.value = rect.width;
    startHeight.value = rect.height;
  }

  document.addEventListener('mousemove', resize);
  document.addEventListener('mouseup', stopResize);
  e.preventDefault();
  e.stopPropagation();
}

function resize(e: MouseEvent) {
  if (!isResizing.value) return;

  const deltaX = startX.value - e.clientX;
  const deltaY = startY.value - e.clientY;

  const newWidth = Math.max(300, startWidth.value + deltaX);
  const newHeight = Math.max(200, startHeight.value + deltaY);

  const panel = panelElement.value;
  if (panel) {
    // Remove Tailwind width/height classes and set custom size
    panel.classList.remove('w-96', 'h-96', 'w-[calc(100%-2rem)]', 'md:w-96');
    panel.style.width = `${newWidth}px`;
    panel.style.height = `${newHeight}px`;
    panel.style.maxWidth = 'none';
    panel.style.maxHeight = 'none';
  }
}

function stopResize() {
  isResizing.value = false;
  document.removeEventListener('mousemove', resize);
  document.removeEventListener('mouseup', stopResize);
}

async function sendMessage() {
  const message = inputMessage.value.trim();
  if (!message || isLoading.value) return;

  // Add user message
  messages.value.push({ role: 'user', content: message });
  inputMessage.value = '';
  isLoading.value = true;

  // Scroll to bottom
  await nextTick();
  scrollToBottom();

  try {
    const requestBody: any = {
      messages: messages.value.slice(-10), // Keep last 10 messages for context
      is_first_message: isFirstMessage.value,
    };

    // Only include article content for the first message
    if (isFirstMessage.value) {
      requestBody.article_title = props.article.title;
      requestBody.article_url = props.article.url;
      requestBody.article_content = props.articleContent.slice(0, 10000);
    }

    const response = await fetch('/api/ai-chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(requestBody),
    });

    if (response.ok) {
      const data = await response.json();
      messages.value.push({ role: 'assistant', content: data.response });
      isFirstMessage.value = false; // Mark that we've sent the first message
    } else {
      // Get error text from response
      const errorText = await response.text();
      console.error('AI chat error response:', response.status, errorText);

      // Try to parse as JSON first
      let errorMessage = t('aiChatError');
      try {
        const errorData = JSON.parse(errorText);
        errorMessage = errorData.error || errorData || t('aiChatError');
      } catch {
        // If not JSON, use the raw text
        errorMessage = errorText || t('aiChatError');
      }

      messages.value.push({
        role: 'assistant',
        content: errorMessage,
      });
    }
  } catch (e) {
    console.error('AI chat error:', e);
    messages.value.push({ role: 'assistant', content: t('aiChatError') });
  } finally {
    isLoading.value = false;
    await nextTick();
    scrollToBottom();
  }
}

function scrollToBottom() {
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault();
    sendMessage();
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="chat-panel">
      <div
        v-if="isOpen"
        ref="panelElement"
        class="chat-panel fixed bottom-4 right-4 md:bottom-6 md:right-6 w-96 h-96 bg-bg-primary border border-border rounded-xl shadow-2xl flex flex-col z-50"
        :class="{ 'select-none': isResizing }"
      >
        <!-- Header -->
        <div
          class="flex items-center justify-between p-3 border-b border-border bg-bg-secondary rounded-t-xl relative"
        >
          <div class="flex items-center gap-2">
            <PhChatCircleText :size="20" class="text-accent" />
            <span class="font-medium text-sm">{{ t('aiChat') }}</span>
          </div>
          <button
            class="p-1 hover:bg-bg-tertiary rounded-lg transition-colors"
            :title="t('close')"
            @click="emit('close')"
          >
            <PhX :size="18" class="text-text-secondary" />
          </button>

          <!-- Resize handle -->
          <div
            class="absolute -top-1 -left-1 w-3 h-3 cursor-nw-resize opacity-0 hover:opacity-100 transition-opacity"
            :class="isResizing ? 'opacity-100' : ''"
            @mousedown="startResize"
          >
            <div class="w-full h-full bg-accent rounded-full border border-white shadow-sm"></div>
          </div>
        </div>

        <!-- Messages -->
        <div ref="chatContainer" class="flex-1 overflow-y-auto p-3 space-y-3">
          <div
            v-if="messages.length === 0"
            class="flex items-center justify-center h-full text-text-secondary text-sm"
          >
            {{ t('aiChatWelcome') }}
          </div>
          <div
            v-for="(msg, index) in messages"
            :key="index"
            class="flex"
            :class="msg.role === 'user' ? 'justify-end' : 'justify-start'"
          >
            <div
              class="max-w-[80%] rounded-lg px-3 py-2 text-sm select-text cursor-text"
              :class="
                msg.role === 'user' ? 'bg-accent text-white' : 'bg-bg-secondary text-text-primary'
              "
            >
              <div class="whitespace-pre-wrap break-words">{{ msg.content }}</div>
            </div>
          </div>
          <div v-if="isLoading" class="flex justify-start">
            <div class="bg-bg-secondary rounded-lg px-3 py-2 text-sm">
              <PhSpinner :size="16" class="animate-spin" />
            </div>
          </div>
        </div>

        <!-- Input -->
        <div class="p-3 border-t border-border bg-bg-secondary rounded-b-xl">
          <div class="flex gap-2">
            <input
              v-model="inputMessage"
              type="text"
              :placeholder="t('aiChatInputPlaceholder')"
              class="flex-1 px-3 py-2 bg-bg-tertiary border border-border rounded-lg text-sm focus:outline-none focus:border-accent"
              :disabled="isLoading"
              @keydown="handleKeydown"
            />
            <button
              :disabled="isLoading || !inputMessage.trim()"
              class="px-3 py-2 bg-accent text-white rounded-lg hover:bg-accent-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              @click="sendMessage"
            >
              <PhPaperPlaneRight :size="18" />
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style>
/* Chat panel text selection */
.chat-panel {
  user-select: text !important;
  -webkit-user-select: text !important;
  -moz-user-select: text !important;
  -ms-user-select: text !important;
}

.chat-panel.select-none {
  user-select: none !important;
  -webkit-user-select: none !important;
  -moz-user-select: none !important;
  -ms-user-select: none !important;
}

/* Message content selection */
.chat-panel .select-text {
  user-select: text !important;
  -webkit-user-select: text !important;
  -moz-user-select: text !important;
  -ms-user-select: text !important;
  cursor: text !important;
}

.chat-panel .select-text * {
  user-select: text !important;
  -webkit-user-select: text !important;
  -moz-user-select: text !important;
  -ms-user-select: text !important;
}

.chat-panel-enter-active,
.chat-panel-leave-active {
  transition: all 0.3s ease;
}

.chat-panel-enter-from,
.chat-panel-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

.chat-panel-enter-to,
.chat-panel-leave-from {
  opacity: 1;
  transform: translateY(0) scale(1);
}
</style>
