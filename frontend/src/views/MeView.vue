<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/stores/userStore';

const userStore = useUserStore();

function update() {
  var input = document.getElementById("url-input") as HTMLAudioElement;
  console.log(input.value)
  userStore.baseUrl = input.value;

}
</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="flex justify-start my-2 mx-1 space-x-1">
        <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/"><i class="fa-solid fa-arrow-left"></i>
        </RouterLink>
      </nav>
      <hr>
    </div>
  </header>

  <main class="flex-1">
    <h1> Meeeeee </h1>
    <input @change="update" type="text" id="url-input" :value="userStore.baseUrl" />
  </main>
</template>
