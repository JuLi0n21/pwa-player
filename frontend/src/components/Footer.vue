<template>
  <hr>
  <nav class="flex justify-around my-2 text-sm">
    <RouterLink class="p-1 rounded-full backdrop--light shadow-xl hover:text-pink-500" to="/me">
      <img src="https://a.ppy.sh/14100399" class="h-6 rounded-full">
    </RouterLink>

    <RouterLink class="flex flex-col justify-center p-2 rounded-full backdrop--light shadow-xl info" to="/menu"><i
        class="fa-solid fa-house"></i>
    </RouterLink>

    <RouterLink class="flex flex-col justify-center p-2 rounded-full backdrop--light shadow-xl info" to="/search">
      <i class="fa-solid fa-magnifying-glass"></i>
    </RouterLink>

  </nav>
</template>
