<script setup lang="ts">
import { useUser } from "@/composables/useUser";

const userStore = useUser();
</script>

<template>
  <hr />
  <nav class="flex justify-around my-2 text-sm md:text-2xl">
    <RouterLink class="shadow-xl backdrop--light p-1 rounded-full hover:text-pink-500" to="/me">
      <img
        :src="userStore.user.value?.avatar_url ?? 'https://osu.ppy.sh/images/layout/avatar-guest.png'"
        class="rounded-full h-6 md:h-12"
      />
    </RouterLink>

    <RouterLink class="flex flex-col justify-center shadow-xl backdrop--light p-2 rounded-full info" to="/menu"
      ><i class="fa-solid fa-house"></i>
    </RouterLink>

    <RouterLink class="flex flex-col justify-center shadow-xl backdrop--light p-2 rounded-full info" to="/search">
      <i class="fa-solid fa-magnifying-glass"></i>
    </RouterLink>
  </nav>
</template>
