<script setup lang="ts">
import type { CollectionPreview } from "@/script/types";
import { useUser } from "@/composables/useUser";

const userStore = useUser();
const props = defineProps<{ collection: CollectionPreview }>();
</script>

<template>
  <RouterLink :to="'/collection/' + props.collection.index" class="w-full">
    <div class="flex items-center border rounded-lg h-24 overflow-hidden bordercolor">
      <img
        class="p-1 rounded-lg w-24 h-24"
        :src="encodeURI(`${userStore.cloudflareUrl.value}${props.collection.previewimage}`)"
        loading="lazy"
      />
      <div class="flex flex-col overflow-hidden text-left">
        <h3 class="overflow-hidden text-ellipsis whitespace-nowrap info">{{ props.collection.name }}</h3>
        <h5 class="text-sm info">{{ props.collection.length }} Songs</h5>
      </div>
    </div>
  </RouterLink>
</template>
