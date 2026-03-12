<template>
  <div class="w-full h-full text-xs animate-pulse">
    
    <div class="mb-4">
      <div class="my-2 ml-4 rounded w-24 h-6 action"></div>
      <ul class="space-y-1 px-2">
        <li v-for="i in 2" :key="'art-skel-' + i">
          <div class="flex items-center opacity-50 p-3 border rounded-xl bordercolor">
            <div class="mr-3 rounded-full w-10 h-10"></div>
            <div class="rounded w-32 h-4 info"></div>
          </div>
        </li>
      </ul>
    </div>

    <div>
      <div class="my-2 ml-4 rounded w-24 h-6 action"></div>
      <ul class="space-y-1 px-2">
        <li v-for="i in 6" :key="'song-skel-' + i">
          <div class="flex items-center p-2 rounded-xl w-full">
            <div class="rounded-lg w-12 h-12 shrink-0"></div>
            
            <div class="flex-1 space-y-2 ml-3">
              <div class="rounded w-3/4 h-4 info"></div>
              <div class="rounded w-1/2 h-3 b action"></div>
            </div>
          </div>
        </li>
      </ul>
    </div>
    
  </div>
</template>

<style scoped>
.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.info { background-color: currentColor; opacity: 0.1; }
.action { background-color: currentColor; opacity: 0.05; }
</style>