<script setup lang="ts">
import type { Song, CollectionPreview } from '../script/types'
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/userStore';
import CollectionListItem from '../components/CollectionListItem.vue'

const userStore = useUserStore();

const collections = ref<CollectionPreview[]>([]);
const limit = ref(10);
const offset = ref(0);
const isLoading = ref(false);

const fetchCollections = async () => {
  if (isLoading.value) return;
  isLoading.value = true;

  const data = await userStore.fetchCollections(offset.value, limit.value);
  data.forEach(song => {
    song.previewimage = `${userStore.baseUrl}/api/v1/images/${song.previewimage}?h=80&w=80`;
  });

  collections.value = [...collections.value, ...data];
  offset.value += limit.value;

  isLoading.value = false;
};

onMounted(async () => {
  await fetchCollections();

  const container = document.querySelector('.collection-container');
  if (container) {
    container.addEventListener('scroll', async () => {
      const scrollTop = container.scrollTop;
      const scrollHeight = container.scrollHeight;
      const clientHeight = container.clientHeight;

      if (scrollTop + clientHeight >= scrollHeight * 0.9 && !isLoading.value) {
        await fetchCollections();
      }
    });
  }
});

</script>

<template>
  <main class="flex-1 text-center flex flex-col h-full overflow-scroll">
    <div class="flex flex-col overflow-scroll collection-container">
      <CollectionListItem v-for="(collection, index) in collections" :key="index" :collection="collection" />
    </div>
  </main>
</template>
