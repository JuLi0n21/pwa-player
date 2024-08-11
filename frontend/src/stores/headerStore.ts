import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useHeaderStore = defineStore('headerStore', () => {
  const showheader = ref(true)
  function show() {
    showheader.value = true;
  }

  function hide() {
    showheader.value = false;
  }

  return { showheader, show, hide }
})
