// Global type declarations

export interface ConfirmDialogOptions {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  isDanger?: boolean;
}

export interface InputDialogOptions {
  title: string;
  message: string;
  placeholder?: string;
  defaultValue?: string;
  confirmText?: string;
  cancelText?: string;
}

export type ToastType = 'success' | 'error' | 'info' | 'warning';

declare global {
  interface Window {
    showConfirm: (ConfirmDialogOptions) => Promise<boolean>;
    showInput: (InputDialogOptions) => Promise<string | null>;
    showToast: (string, ToastType?, number?) => void;
  }
}

export {};
