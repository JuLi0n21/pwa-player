import { ref } from 'vue';

export function useThemeColors() {
  const bgColor = ref(localStorage.getItem('bgColor') || '#1c1719');
  const actionColor = ref(localStorage.getItem('actionColor') || '#eab308');
  const infoColor = ref(localStorage.getItem('infoColor') || '#ec4899');
  const borderColor = ref(localStorage.getItem('borderColor') || '#ec4899');

  function applyColors(bg: string, main: string, info: string, border: string) {
    document.documentElement.style.setProperty('--background-color', bg);
    document.documentElement.style.setProperty('--action-color', main);
    document.documentElement.style.setProperty('--information-color', info);
    document.documentElement.style.setProperty('--border-color', border);
  }

  function save(
    bg: string | null = null,
    main: string | null = null,
    info: string | null = null,
    border: string | null = null
  ) {
    bgColor.value = bg ?? bgColor.value;
    actionColor.value = main ?? actionColor.value;
    infoColor.value = info ?? infoColor.value;
    borderColor.value = border ?? borderColor.value;

    applyColors(bgColor.value, actionColor.value, infoColor.value, borderColor.value);

    localStorage.setItem('bgColor', bgColor.value);
    localStorage.setItem('actionColor', actionColor.value);
    localStorage.setItem('infoColor', infoColor.value);
    localStorage.setItem('borderColor', borderColor.value);
  }

  // Initialize colors on composable use
  applyColors(bgColor.value, actionColor.value, infoColor.value, borderColor.value);

  return {
    bgColor,
    actionColor,
    infoColor,
    borderColor,
    save,
  };
}
