<template>
  <hr>
  <nav class="flex justify-around my-4 text-xl">
    <RouterLink class="p-1 rounded-full backdrop--light shadow-xl hover:text-pink-500" to="/me">
      <img src="https://a.ppy.sh/14100399" class="h-12 rounded-full">
    </RouterLink>

    <RouterLink
      class="flex flex-col justify-center text-2xl p-4 rounded-full backdrop--light shadow-xl hover:text-pink-500"
      to="/menu"><i class="fa-solid fa-house"></i>
    </RouterLink>

    <RouterLink
      class="flex flex-col justify-center text-2xl p-4 rounded-full backdrop--light shadow-xl hover:text-pink-500"
      to="/search">
      <i class="fa-solid fa-magnifying-glass"></i>
    </RouterLink>

  </nav>
</template>
