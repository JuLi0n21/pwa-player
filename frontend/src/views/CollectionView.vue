<script setup lang="ts">
import type { Song, CollectionPreview } from '../script/types'
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/userStore';
import CollectionListItem from '../components/CollectionListItem.vue'

const userStore = useUserStore();

const collections = ref<CollectionPreview[]>([]);

onMounted(async () => {
  const data = await userStore.fetchCollections();
  data.forEach(song => {
    song.previewimage = `${userStore.baseUrl}api/v1/images/${song.previewimage}`;
  })
  collections.value = data;

});
</script>

<template>
  <main class="flex-1 text-center flex flex-col h-full overflow-scroll">
    <div class="flex flex-col overflow-scroll">
      <CollectionListItem 
        v-for="(collection, index) in collections"
        :key="index"
        :collection="collection"
      />
    </div>
  </main>
</template>

