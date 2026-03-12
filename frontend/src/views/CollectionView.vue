<script setup lang="ts">
import { type Song, type CollectionPreview, mapApiToCollectionPreview } from "../script/types";
import { ref, onMounted } from "vue";
import CollectionListItem from "../components/CollectionListItem.vue";
import { useApi } from "@/composables/useApi";

const { musicApi } = useApi();
const api = musicApi.value;

const collections = ref<CollectionPreview[]>([]);
const limit = ref(10);
const offset = ref(0);
const isLoading = ref(false);

const fetchCollections = async () => {
  if (isLoading.value) return;
  isLoading.value = true;

  const response = await api.musicBackendSearchCollections("", limit.value, offset.value);
  let songs = mapApiToCollectionPreview(response.data.collections || [], offset.value);
  collections.value = [...collections.value, ...songs];
  offset.value += limit.value;

  isLoading.value = false;
};

onMounted(async () => {
  await fetchCollections();

  const container = document.querySelector(".collection-container");
  if (container) {
    container.addEventListener("scroll", async () => {
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
  <main class="flex flex-col flex-1 h-full overflow-y-hidden text-center">
    <div class="flex flex-col overflow-y-scroll collection-container">
      <CollectionListItem v-for="(collection, index) in collections" :key="index" :collection="collection" />
    </div>
  </main>
</template>
