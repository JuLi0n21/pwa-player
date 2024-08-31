<script setup lang="ts">
import { defineProps } from 'vue';
import type { CollectionPreview } from '@/script/types';
import { useUserStore } from '@/stores/userStore';
import { useAudioStore } from '@/stores/audioStore';
const userStore = useUserStore();
const audioStore = useAudioStore();
const props = defineProps<{ collection: CollectionPreview }>();


function hi() {
}

</script>

<template>

  <RouterLink @click.navtive="hi" :to="'/collection/' + props.collection.index">
    <div class=" border bordercolor rounded-lg flex">
      <img class="h-20 w-20 m-2 rounded-lg" :src="props.collection.previewimage" loading="lazy" />
      <div class="flex flex-col">
        <h3 class="self-start info">{{ props.collection.name }}</h3>
        <h5 class="self-start info text-sm">{{ props.collection.length }} Songs </h5>
      </div>
    </div>
  </RouterLink>

</template>
