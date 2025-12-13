import { ref, onBeforeUnmount } from 'vue';

export function useResizablePanels() {
  const sidebarWidth = ref<number>(256);
  const articleListWidth = ref<number>(400);
  const isResizingSidebar = ref<boolean>(false);
  const isResizingArticleList = ref<boolean>(false);

  // Sidebar resize handlers
  function startResizeSidebar(): void {
    isResizingSidebar.value = true;
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
    window.addEventListener('mousemove', handleResizeSidebar);
    window.addEventListener('mouseup', stopResizeSidebar);
  }

  function handleResizeSidebar(): void {
    if (!isResizingSidebar.value) return;
    const newWidth = (window.event as MouseEvent).clientX;
    if (newWidth >= 180 && newWidth <= 450) {
      sidebarWidth.value = newWidth;
    }
  }

  function stopResizeSidebar(): void {
    isResizingSidebar.value = false;
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
    window.removeEventListener('mousemove', handleResizeSidebar);
    window.removeEventListener('mouseup', stopResizeSidebar);
  }

  // Article list resize handlers
  function startResizeArticleList(): void {
    isResizingArticleList.value = true;
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
    window.addEventListener('mousemove', handleResizeArticleList);
    window.addEventListener('mouseup', stopResizeArticleList);
  }

  function handleResizeArticleList(): void {
    if (!isResizingArticleList.value) return;
    const newWidth = (window.event as MouseEvent).clientX - sidebarWidth.value;
    if (newWidth >= 250 && newWidth <= 600) {
      articleListWidth.value = newWidth;
    }
  }

  function stopResizeArticleList(): void {
    isResizingArticleList.value = false;
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
    window.removeEventListener('mousemove', handleResizeArticleList);
    window.removeEventListener('mouseup', stopResizeArticleList);
  }

  // Cleanup
  onBeforeUnmount(() => {
    window.removeEventListener('mousemove', handleResizeSidebar);
    window.removeEventListener('mouseup', stopResizeSidebar);
    window.removeEventListener('mousemove', handleResizeArticleList);
    window.removeEventListener('mouseup', stopResizeArticleList);
  });

  return {
    sidebarWidth,
    articleListWidth,
    startResizeSidebar,
    startResizeArticleList,
  };
}
