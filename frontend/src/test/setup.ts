// src/test/setup.ts
import { vi } from 'vitest';

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => {},
  }),
});

// Mock fetch
global.fetch = vi.fn();

// Mock window.showToast, window.showConfirm, etc.
Object.defineProperty(window, 'showToast', {
  writable: true,
  value: () => {},
});

Object.defineProperty(window, 'showConfirm', {
  writable: true,
  value: () => Promise.resolve(true),
});

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};
