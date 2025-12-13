import { ref } from 'vue';
import type { ConfirmDialogOptions, InputDialogOptions, ToastType } from '@/types/global';

export interface Toast {
  id: number;
  message: string;
  type: ToastType;
  duration: number;
}

export interface ConfirmDialogState extends ConfirmDialogOptions {
  onConfirm: () => void;
  onCancel: () => void;
}

export interface InputDialogState extends InputDialogOptions {
  onConfirm: (string) => void;
  onCancel: () => void;
}

export function useNotifications() {
  const confirmDialog = ref<ConfirmDialogState | null>(null);
  const inputDialog = ref<InputDialogState | null>(null);
  const toasts = ref<Toast[]>([]);

  function showConfirm(options: ConfirmDialogOptions): Promise<boolean> {
    return new Promise((resolve) => {
      confirmDialog.value = {
        ...options,
        onConfirm: () => {
          confirmDialog.value = null;
          resolve(true);
        },
        onCancel: () => {
          confirmDialog.value = null;
          resolve(false);
        },
      };
    });
  }

  function showInput(options: InputDialogOptions): Promise<string | null> {
    return new Promise((resolve) => {
      inputDialog.value = {
        ...options,
        onConfirm: (value: string) => {
          inputDialog.value = null;
          resolve(value);
        },
        onCancel: () => {
          inputDialog.value = null;
          resolve(null);
        },
      };
    });
  }

  function showToast(message: string, type: ToastType = 'info', duration: number = 3000): void {
    const id = Date.now();
    toasts.value.push({ id, message, type, duration });
  }

  function removeToast(id: number): void {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  }

  // Make these available globally
  function installGlobalHandlers(): void {
    window.showConfirm = showConfirm;
    window.showInput = showInput;
    window.showToast = showToast;
  }

  return {
    confirmDialog,
    inputDialog,
    toasts,
    showConfirm,
    showInput,
    showToast,
    removeToast,
    installGlobalHandlers,
  };
}
