import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Song, CollectionPreview } from '@/script/types';

export const useUserStore = defineStore('userStore', () => {
  const userId = ref(null)
  const baseUrl = ref('https://service.illegalesachen.download/')

  async function fetchSong(hash: string): Promise<Song> {
    try {

      const response = await fetch(`${baseUrl}api/v1/songs/${hash}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: Song = await response.json();
      return data;
    } catch (error) {
      console.error('Failed to fetch songs:', error);
      return {
        hash: "-1",
        name: "song name",
        artist: "artist name",
        length: 0,
        url: "song.mp3",
        previewimage: "404.im5",
        mapper: "map",
      } as Song;
    }
  }

  async function fetchWithCache<T>(cacheKey: string, url: string, cacheDuration: number = 24 * 60 * 60 * 1000): Promise<T> {
    const cacheTimestampKey = `${cacheKey}_timestamp`;

    const cachedData = localStorage.getItem(cacheKey);
    const cachedTimestamp = localStorage.getItem(cacheTimestampKey);

    if (cachedData && cachedTimestamp && (Date.now() - parseInt(cachedTimestamp)) < cacheDuration) {
      console.log(`Returning cached data for ${cacheKey}`);
      return JSON.parse(cachedData);
    }

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data: T = await response.json();

      localStorage.setItem(cacheKey, JSON.stringify(data));
      localStorage.setItem(cacheTimestampKey, Date.now().toString());

      return data;
    } catch (error) {
      console.error(`Failed to fetch data from ${url}:`, error);
      throw error;
    }
  }


  async function fetchCollections(): Promise<CollectionPreview[]> {
    const cacheKey = 'collections_cache';
    const url = `${baseUrl.value}api/v1/collections/`;
    return fetchWithCache<CollectionPreview[]>(cacheKey, url);
  }

  async function fetchCollection(index: number): Promise<Song[]> {
    const cacheKey = `collection_${index}_cache`;
    const url = `${baseUrl.value}api/v1/collection/${index}`;
    return fetchWithCache<Song[]>(cacheKey, url);
  }

  async function fetchFavorites(limit: number, offset: number): Promise<Song[]> {
    const cacheKey = `favorites_${limit}_${offset}_cache`;
    const url = `${baseUrl.value}api/v1/recent?limit=${limit}&offset=${offset}`;
    return fetchWithCache<Song[]>(cacheKey, url);
  }

  async function fetchRecent(limit: number, offset: number): Promise<Song[]> {
    const cacheKey = `recent_${limit}_${offset}_cache`;
    const url = `${baseUrl.value}api/v1/songs/recent?limit=${limit}&offset=${offset}`;
    return fetchWithCache<Song[]>(cacheKey, url);
  }



  return { fetchSong, fetchCollections, fetchCollection, fetchRecent, fetchFavorites, userId, baseUrl }
})
