<template>
  <div class="flex bg-white/5 m-1 border rounded-lg md:text-xl animate-pulse bordercolor">
    <div class="flex-shrink-0 bg-white/10 m-1 rounded-lg w-14 md:w-24 h-14 md:h-24"></div>
    <div class="flex flex-col flex-1 justify-center gap-2 px-2 overflow-hidden text-left">
      <div class="bg-white/10 rounded w-3/4 h-4 info"></div>
      <div class="bg-white/5 rounded w-1/2 h-3 action"></div>
      <div class="bg-white/5 rounded w-1/4 h-3 action"></div>
    </div>
  </div>
</template>

<style scoped>
.info {
  background-color: currentColor;
  opacity: 0.15;
}
.action {
  background-color: currentColor;
  opacity: 0.1;
}
</style>
