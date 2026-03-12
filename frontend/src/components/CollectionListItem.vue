<script setup lang="ts">
import type { CollectionPreview } from "@/script/types";

const props = defineProps<{ collection: CollectionPreview }>();
</script>

<template>
  <RouterLink :to="'/collection/' + props.collection.index">
    <div class="flex border rounded-lg bordercolor">
      <img class="m-2 rounded-lg w-20 h-20" :src="props.collection.previewimage" loading="lazy" />
      <div class="flex flex-col">
        <h3 class="self-start info">{{ props.collection.name }}</h3>
        <h5 class="self-start text-sm info">{{ props.collection.length }} Songs</h5>
      </div>
    </div>
  </RouterLink>
</template>
