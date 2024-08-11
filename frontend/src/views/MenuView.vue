<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useHeaderStore } from '@/stores/headerStore';

const route = useRoute();

function isActive(path: string) {
  return route.path === path ? 'bg-blue-500 text-white' : '';
};

const headerStore = useHeaderStore();
</script>

<template>
  <main class="flex-1 flex flex-col h-full overflow-scroll">
    <header v-show="true">
      <div class="wrapper">
        <nav class="flex justify-start my-2 mx-1 space-x-1 overflow-x-scroll flex-nowrap text-nowrap">
          <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/"><i class="fa-solid fa-arrow-left"></i>
          </RouterLink>
          <RouterLink :class="`p-1 rounded-full backdrop--light shadow-xl ${isActive('/')}`" to="/menu/recent">Recently
            added</RouterLink>
          <RouterLink :class="`p-1 rounded-full backdrop--light shadow-xl ${isActive('/')}`" to="/menu/favourites">
            Favorites</RouterLink>
          <RouterLink :class="`p-1 rounded-full backdrop--light shadow-xl ${isActive('/')}`" to="/menu/collections">
            Collections</RouterLink>
        </nav>
        <hr>
      </div>
    </header>

    <router-view />
  </main>
</template>
