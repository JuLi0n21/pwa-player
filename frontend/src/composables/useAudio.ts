import { ref, onMounted } from "vue";
import { mapApiToSongs, type Song } from "@/script/types";
import { useApi } from "./useApi";
import { useUser } from "./useUser";

let audioInstance: ReturnType<typeof createAudio> | null = null;

function createAudio() {
  const { musicApi } = useApi();
  const userStore = useUser();

  const audioElement = document.createElement("audio");
  audioElement.setAttribute("id", "global-audio");
  audioElement.setAttribute("controls", "");
  audioElement.classList.add("hidden");
  document.body.appendChild(audioElement);

  const state = {
    isPlaying: ref(false),
    duration: ref("0:00"),
    currentTime: ref("0:00"),
    percentDone: ref(0),
    shuffle: ref(false),
    repeat: ref(false),
    activeSongs: ref<Song[] | null>([]),
    currentSong: ref<Song | null>(null),
    recentlyPlayed: ref(new Map<string, Song>()),
  };

  function saveToLocalStorage(key: string, data: any) {
    localStorage.setItem(key, JSON.stringify(data));
  }

  function loadFromLocalStorage<T>(key: string): T | null {
    const item = localStorage.getItem(key);
    return item ? (JSON.parse(item) as T) : null;
  }

  function setSong(song: Song | null) {
    if (!song) return;
    state.currentSong.value = song;
    const map = state.recentlyPlayed.value;

    saveToLocalStorage("lastPlayedSong", song);
    if (map.has(song.hash)) {
      map.delete(song.hash);
    }
    map.set(song.hash, song);

    audioElement.pause();
    audioElement.src = userStore.cloudflareUrl.value + "/" + song.url;
    audioElement.addEventListener("canplaythrough", () => audioElement.play().catch(console.error), { once: true });
  }

  function setCollection(song: Song[] | null) {
    if (!song) return;
    state.activeSongs.value = song;
    saveToLocalStorage("activeCollection", song);
  }

  function togglePlay() {
    if (audioElement.paused) {
      audioElement.play().catch(console.error);
    } else {
      audioElement.pause();
    }
  }

  function toggleNext() {
    if (!state.activeSongs.value || !state.currentSong.value) return;

    const songs = state.activeSongs.value;

    if (state.shuffle.value) {
      setSong(songs[Math.floor(Math.random() * songs.length)]);
    }

    const currentHash = state.currentSong.value.hash;
    const currentIndex = songs.findIndex((song) => song.hash === currentHash);

    if (currentIndex === -1) return;

    const nextIndex = (currentIndex + 1) % songs.length;
    setSong(songs[nextIndex]);
  }

  function togglePrevious() {
    if (!state.activeSongs.value || !state.currentSong.value) return;

    const songs = state.activeSongs.value;
    const currentHash = state.currentSong.value.hash;
    const currentIndex = songs.findIndex((song) => song.hash === currentHash);

    if (currentIndex === -1) return;

    const prevIndex = (currentIndex - 1 + songs.length) % songs.length;
    setSong(songs[prevIndex]);
  }

  function update() {
    const { currentTime: ct, duration: dur } = audioElement;
    state.isPlaying.value = !audioElement.paused;
    state.currentTime.value = formatTime(ct);
    state.duration.value = formatTime(dur);
    state.percentDone.value = isNaN(dur) ? 0 : (ct / dur) * 100;

    if (audioElement.ended) {
      if (state.repeat.value) {
        audioElement.currentTime = 0;
        audioElement.play();
        return;
      }

      toggleNext();
    }
  }

  function formatTime(seconds: number): string {
    const min = Math.floor(seconds / 60);
    const sec = Math.floor(seconds % 60)
      .toString()
      .padStart(2, "0");
    return `${min}:${sec}`;
  }

  function updateTime(value: number) {
    if (value) audioElement.currentTime = (Number(value) / 100) * audioElement.duration;
  }

  async function loadInitialSong() {
    try {
      const api = musicApi.value;
      const res = await api.musicBackendRecent();
      let songs = mapApiToSongs(res.data.songs);
      if (res.data?.songs?.length) {
        setSong(songs[0]);
      }
    } catch (err) {
      console.error("Failed to load song:", err);
    }
  }

  function init() {
    setSong(loadFromLocalStorage<Song>("lastPlayedSong"));
    setCollection(loadFromLocalStorage<Song[]>("activeCollection"));

    if (!state.currentSong.value) {
      loadInitialSong();
    }

    audioElement.addEventListener("timeupdate", update);
  }

  onMounted(init);

  return {
    ...state,
    togglePlay,
    toggleShuffle: () => (state.shuffle.value = !state.shuffle.value),
    toggleRepeat: () => (state.repeat.value = !state.repeat.value),
    toggleNext,
    togglePrevious,
    updateTime,
    setSong,
    setCollection,
    init,
  };
}

export function useAudio() {
  if (!audioInstance) {
    audioInstance = createAudio();
  }
  return audioInstance;
}
