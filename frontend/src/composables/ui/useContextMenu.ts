import { ref } from 'vue';

export interface ContextMenuItem {
  label: string;
  icon?: string;
  action: string;
  danger?: boolean;
}

export interface ContextMenuState {
  show: boolean;
  x: number;
  y: number;
  items: ContextMenuItem[];
  data: unknown;

  callback?: (string, unknown) => void;
}

export function useContextMenu() {
  const contextMenu = ref<ContextMenuState>({
    show: false,
    x: 0,
    y: 0,
    items: [],
    data: null,
  });

  function openContextMenu(event: CustomEvent): void {
    contextMenu.value = {
      show: true,
      x: event.detail.x,
      y: event.detail.y,
      items: event.detail.items,
      data: event.detail.data,
      callback: event.detail.callback,
    };
  }

  function closeContextMenu(): void {
    contextMenu.value.show = false;
  }

  function handleContextMenuAction(action: string): void {
    if (contextMenu.value.callback) {
      contextMenu.value.callback(action, contextMenu.value.data);
    }
  }

  return {
    contextMenu,
    openContextMenu,
    closeContextMenu,
    handleContextMenuAction,
  };
}
