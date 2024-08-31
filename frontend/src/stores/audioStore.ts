import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Song, CollectionPreview } from '@/script/types';

export const useAudioStore = defineStore('audioStore', () => {

  const songSrc = ref('https://cdn.pixabay.com/audio/2024/05/24/audio_46382ae035.mp3')

  const artist = ref('Artist ');
  const title = ref('Title ');
  const bgimg = ref('https://assets.ppy.sh/beatmaps/2197744/covers/cover@2x.jpg?1722207959');
  const hash = ref('0000');

  const isPlaying = ref(true)
  const duration = ref('0:00')
  const currentTime = ref('0:00')
  const percentDone = ref(0)

  const shuffle = ref(false);
  const repeat = ref(false);

  const activeCollection = ref<Song[]>([]);
  const currentSong = ref<Song>(null);

  function saveSongToLocalStorage(song: Song) {
    localStorage.setItem('lastPlayedSong', JSON.stringify(song));
  }

  function loadSongFromLocalStorage(): Song | null {
    const song = localStorage.getItem('lastPlayedSong');
    return song ? JSON.parse(song) : null;
  }

  function saveCollectionToLocalStorage(collection: Song[]) {
    localStorage.setItem('lastActiveCollection', JSON.stringify(collection));
  }

  function loadCollectionFromLocalStorage(): Song[] | null {
    const collection = localStorage.getItem('lastActiveCollection');
    return collection ? JSON.parse(collection) : null;
  }

  function togglePlay() {
    var audio = document.getElementById("audio-player") as HTMLAudioElement;

    if (audio.paused) {
      audio.play();
    } else {
      audio.pause();
    }
  }


  function update() {
    var audio = document.getElementById("audio-player") as HTMLAudioElement;

    isPlaying.value = audio.paused;

    let current_min = Math.floor(audio.currentTime / 60);
    let current_sec = Math.round(audio.currentTime % 60);

    if (!isNaN(current_sec) && !isNaN(current_min)) {
      currentTime.value = current_min + ':' + current_sec.toString().padStart(2, '0');
    }

    let duration_min = Math.floor(audio.duration / 60);
    let duration_sec = Math.round(audio.duration % 60);

    if (!isNaN(duration_sec) && !isNaN(duration_min)) {
      duration.value = duration_min + ':' + duration_sec.toString().padStart(2, '0');
    }


    let percent = (audio.currentTime / audio.duration) * 100

    if (!isNaN(percent)) {
      percentDone.value = percent;
    }

    if (audio.ended) {
      next();
    }
  }

  function next() {
    var audio = document.getElementById("audio-player") as HTMLAudioElement;

    if (repeat.value) {
      audio.pause()
      audio.currentTime = 0;
      audio.play()
      return;
    }

    if (shuffle.value) {
      audio.pause()

      setSong(activeCollection.value[Math.floor(activeCollection.value.length * Math.random())])
      audio.play()
      return;
    }

    toggleNext()
  }

  function updateTime() {

    var audioslider = document.getElementById("audio-slider") as HTMLInputElement;

    var audio = document.getElementById("audio-player") as HTMLAudioElement;

    audio.currentTime = Math.round((Number(audioslider.value) / 100) * audio.duration)
  }

  function togglePrev() {
    let index = activeCollection.value.findIndex(s => s.hash == hash.value);
    setSong(activeCollection.value[Math.abs((index - 1 + activeCollection.value.length) % activeCollection.value.length)])

    console.log(Math.abs((index - 1 + activeCollection.value.length) % activeCollection.value.length))

  }

  function toggleNext() {
    let index = 0;
    if (shuffle.value) {
      index = Math.floor(activeCollection.value.length * Math.random())

    } else {
      index = activeCollection.value.findIndex(s => s.hash == hash.value);

    }
    setSong(activeCollection.value[Math.abs((index + 1) % activeCollection.value.length)])
  }

  function toggleShuffle() {
    console.log('shuffle', !shuffle.value)
    shuffle.value = !shuffle.value;
  }

  function toggleRepeat() {
    console.log('repeat', !repeat.value)
    repeat.value = !repeat.value;
  }

  function setSong(song: Song | null) {

    if (song === null) {
      return;
    }

    console.log('setSong', song)
    var audio = document.getElementById("audio-player") as HTMLAudioElement;

    songSrc.value = song.url;
    artist.value = song.artist;
    title.value = song.name;
    bgimg.value = song.previewimage;
    hash.value = song.hash;

    currentSong.value = song;
    saveSongToLocalStorage(song);

    if (audio === null) { return; }

    if (!audio.paused) {
      audio.pause
    }

    audio.src = song.url;


    audio.addEventListener('canplaythrough', () => {
      audio.play().catch(error => {
        console.error('Playback error:', error);
      });
    });
  }

  function setCollection(songs: Song[] | null) {
    if (songs === null) { return; }
    activeCollection.value = songs;
    saveCollectionToLocalStorage(songs);
  }

  setSong(loadSongFromLocalStorage());
  setCollection(loadCollectionFromLocalStorage());

  return { setCollection, songSrc, updateTime, artist, title, bgimg, shuffle, repeat, setSong, togglePlay, togglePrev, toggleNext, toggleRepeat, toggleShuffle, isPlaying, currentTime, duration, update, percentDone, currentSong }
})
