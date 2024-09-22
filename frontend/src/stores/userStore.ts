import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Song, CollectionPreview, Me } from '@/script/types';

export const useUserStore = defineStore('userStore', () => {
  const userId = ref(null)
  const baseUrl = ref('https://service.illegalesachen.download')
  const proxyUrl = ref('https://proxy.illegalesachen.download')

  const User = ref<Me | null>(null)

  function saveUser(user: Me | null) {
    localStorage.setItem('activeUser', JSON.stringify(user));
  }

  function loadUser(): Me | null {
    const user = localStorage.getItem('activeUser');
    return user ? JSON.parse(user) : null;
  }

  function setUser(user: Me | null) {
    User.value = user;
    saveUser(user)
  }

  async function fetchSong(hash: string): Promise<Song> {
    try {

      const response = await fetch(`${baseUrl}/api/v1/songs/${hash}`);
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

  async function fetchWithCache<T>(cacheKey: string, url: string, cacheDuration: number = 24 * 60 * 60 * 1): Promise<T> {
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


  async function fetchCollections(offset: number, limit: number): Promise<CollectionPreview[]> {
    const cacheKey = `collections_cache_${offset}_${limit}`;
    const url = `${baseUrl.value}/api/v1/collections/?offset=${offset}&limit=${limit}`;
    return fetchWithCache<CollectionPreview[]>(cacheKey, url);
  }

  async function fetchCollection(index: number): Promise<Song[]> {
    const cacheKey = `collection_${index}_cache`;
    const url = `${baseUrl.value}/api/v1/collection/${index}`;
    return fetchWithCache<Song[]>(cacheKey, url);
  }

  async function fetchFavorites(limit: number, offset: number): Promise<Song[]> {
    const cacheKey = `favorites_${limit}_${offset}_cache`;
    const url = `${baseUrl.value}/api/v1/recent?limit=${limit}&offset=${offset}`;
    return fetchWithCache<Song[]>(cacheKey, url);
  }

  async function fetchRecent(limit: number, offset: number): Promise<Song[]> {
    const cacheKey = `recent_${limit}_${offset}_cache`;
    const url = `${baseUrl.value}/api/v1/songs/recent?limit=${limit}&offset=${offset}`;
    return fetchWithCache<Song[]>(cacheKey, url);
  }

  async function fetchActiveSearch(query: string): Promise<{}> {
    const cacheKey = `collections_activeSearch_${query}`;
    const url = `${baseUrl.value}/api/v1/search/active?q=${query}`;
    return fetchWithCache(cacheKey, url);
  }

  async function fetchSearchArtist(query: string): Promise<Song[]> {
    const cacheKey = `collections_artist_${query}`;
    const url = `${baseUrl.value}/api/v1/search/artist?q=${query}`;
    return fetchWithCache<Song[]>(cacheKey, url);
  }

  async function fetchMe(): Promise<Me | {}> {
    const url = `${proxyUrl.value}/me`;

    try {
      const response = await fetch(url, {
        method: 'GET',
        credentials: 'include'
      });
      console.log(response);

      if (response.redirected) {
        window.open(response.url, '_blank');
        return { "redirected": true };
      }

      if (!response.ok) {
        console.error(`Fetch failed with status: ${response.status} ${response.statusText}`);
        return { id: -1 } as Me;
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Fetch error:', error);
      return {} as Me;
    }
  }

  setUser(loadUser());

  return { fetchSong, fetchActiveSearch, fetchSearchArtist, fetchCollections, fetchCollection, fetchRecent, fetchFavorites, fetchMe, userId, baseUrl, proxyUrl, User, setUser }
})
