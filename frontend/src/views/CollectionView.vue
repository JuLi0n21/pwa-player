<script setup lang="ts">
import { type Song, type CollectionPreview, mapApiToCollectionPreview } from "../script/types";
import { ref, onMounted, nextTick } from "vue";
import CollectionListItem from "../components/CollectionListItem.vue";
import CollectionListItemSkeleton from "../components/CollectionListItemSkeleton.vue";
import { useApi } from "@/composables/useApi";

const { musicApi } = useApi();
const api = musicApi.value;

const containerRef = ref<HTMLElement | null>(null);
const collections = ref<CollectionPreview[]>([]);
const limit = ref(12);
const offset = ref(0);
const isLoading = ref(false);

const fetchCollections = async () => {
  if (isLoading.value) return;
  isLoading.value = true;

  try {
    const response = await api.musicBackendSearchCollections("", limit.value, offset.value);
    const newItems = response.data.collections || [];

    if (newItems.length > 0) {
      let mapped = mapApiToCollectionPreview(newItems, offset.value);
      collections.value = [...collections.value, ...mapped];
      offset.value += limit.value;

      await nextTick();

      const container = containerRef.value;
      if (container && container.scrollHeight <= container.clientHeight) {
        isLoading.value = false; 
        await fetchCollections(); 
      }
    }
  } catch (error) {
    console.error("Fetch failed:", error);
  } finally {
    isLoading.value = false;
  }
};

onMounted(async () => {
  await fetchCollections();

  const container = containerRef.value;
  if (container) {
    container.addEventListener("scroll", () => {
      const { scrollTop, scrollHeight, clientHeight } = container;
      if (scrollTop + clientHeight >= scrollHeight - 100 && !isLoading.value) {
        fetchCollections();
      }
    });
  }
});
</script>

<template>
  <main class="flex flex-col flex-1 h-full overflow-y-hidden text-center">
    <div 
      ref="containerRef" 
      class="gap-2 grid grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 p-2 overflow-y-auto collection-container"
    >
      <CollectionListItem 
        v-for="(collection, index) in collections" 
        :key="index" 
        :collection="collection" 
      />

      <template v-if="isLoading">
        <CollectionListItemSkeleton v-for="i in 12" :key="'skeleton-' + i" />
      </template>
    </div>
  </main>
</template>