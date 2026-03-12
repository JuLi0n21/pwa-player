<template>
  <div class="flex items-center bg-white/5 border rounded-lg h-24 animate-pulse bordercolor">
    <div class="bg-white/10 m-2 rounded-lg w-20 h-20 shrink-0"></div>
    <div class="flex flex-col flex-1 justify-center gap-2 p-2 overflow-hidden">
      <div class="bg-white/10 rounded w-3/4 h-5"></div>
      <div class="bg-white/5 rounded w-1/4 h-3"></div>
    </div>
  </div>
</template>
<style scoped>
.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
