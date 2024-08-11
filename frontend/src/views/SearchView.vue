<script setup lang="ts">
</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="flex justify-start my-2 mx-1 space-x-1">
        <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/"><i class="fa-solid fa-arrow-left"></i>
        </RouterLink>
        <h1 class="absolute left-0 right-0 text-center"> Search </h1>
      </nav>
      <hr>
    </div>
  </header>
  <main class="flex-1 flex-col">
    <h1> Search...</h1>
    <input class="flex-1 border border-pink-500 accent-pink-800 bg-yellow-300 bg-opacity-20 rounded-lg m-2 p-2" />
  </main>
</template>
