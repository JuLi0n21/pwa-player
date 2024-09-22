<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/userStore';
import SongItem from '../components/SongItem.vue'
import { useAudioStore } from '@/stores/audioStore';

const audioStore = useAudioStore();
const userStore = useUserStore();

const bgColor = ref('');
const actionColor = ref('');
const infoColor = ref('');
const borderColor = ref('');

const loginStatus = ref('Login');

function update() {
  var input = document.getElementById("url-input") as HTMLInputElement;
  userStore.baseUrl = input.value;

}

function save(bg: string | null, main: string | null, info: string | null, border: string | null) {

  document.documentElement.style.setProperty('--background-color', bg ?? bgColor.value);
  document.documentElement.style.setProperty('--action-color', main ?? actionColor.value);
  document.documentElement.style.setProperty('--information-color', info ?? infoColor.value);
  document.documentElement.style.setProperty('--border-color', border ?? borderColor.value);

  localStorage.setItem('bgColor', bg ?? bgColor.value);
  localStorage.setItem('actionColor', main ?? actionColor.value);
  localStorage.setItem('infoColor', info ?? infoColor.value);
  localStorage.setItem('borderColor', border ?? borderColor.value);

  console.log("bg", bgColor.value, "action:", actionColor.value, "info", infoColor.value, "border", borderColor.value)
}

async function getMe() {

  const data = await userStore.fetchMe() as Me;
  if (data.redirected == true) {
    loginStatus.value = "waiting for login, click to refresh!"
    console.log("redirect detected");
  }

  console.log(data)
  if (data.id === null || data.id === undefined || Object.keys(data).length === 0) {
    return
  }

  console.log("active user: ", data.name)
  userStore.setUser(data);
  userStore.baseUrl = data.endpoint;

}

onMounted(() => {
  reset();
})

function reset() {

  bgColor.value = localStorage.getItem('bgColor') || '#1c1719';
  actionColor.value = localStorage.getItem('actionColor') || '#eab308';
  infoColor.value = localStorage.getItem('infoColor') || '#ec4899';
  borderColor.value = localStorage.getItem('borderColor') || '#ec4899';

  document.documentElement.style.setProperty('--background-color', bgColor.value);
  document.documentElement.style.setProperty('--action-color', actionColor.value);
  document.documentElement.style.setProperty('--information-color', infoColor.value);
  document.documentElement.style.setProperty('--border-color', borderColor.value);

}

</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="flex justify-start my-2 mx-1 space-x-1">
        <RouterLink class="p-1 rounded-full backdrop--light shadow-xl" to="/"><i class="fa-solid fa-arrow-left"></i>
        </RouterLink>
      </nav>
      <hr>
    </div>
  </header>

  <main class="flex-1 flex flex-col overflow-scroll">
    <h1> Meeeeee </h1>
    <input @change="update" type="text" id="url-input" :value="userStore.baseUrl" disabled />
    <br>
    <button v-if="!userStore.User" @click="getMe" class="border bordercolor rounded-lg p-0.5">{{ loginStatus }}</button>
    <div v-if="userStore.User" class="flex p-5 justify-between">
      <img :src="userStore.User.avatar_url" class="w-1/3">
      <div>
        <p>{{ userStore.User.name }}</p>
        <p>{{ userStore.User.endpoint == "" ? 'Not Connected' : 'Connected' }}</p>
        <p>Sharing: <button @click="share" class="border bordercolor rounded-lg p-0.5">{{ userStore.User.share
            }}</button></p>
        <button @click="getMe" class="border bordercolor rounded-lg p-0.5"> Refresh
        </button>

      </div>
    </div>

    <div class="flex flex-col w-full justify-around p-10">
      <div class="flex flex-1 justify-between">
        <p>Background:</p>
        <input type="color" id="bgPicker" v-model="bgColor" @input="save()"
          class="appearance-none w-8 h-8 border border-2 p-0 overflow-hidden cursor-pointer">
      </div>

      <div class="flex flex-1 justify-between">
        <p>Main:</p>
        <input type="color" id="actionPicker" v-model="actionColor" @input="save()"
          class="appearance-none w-8 h-8 border border-2 p-0 overflow-hidden cursor-pointer">
      </div>

      <div class="flex flex-1 justify-between">
        <p>Secondary:</p>
        <input type="color" id="infoPicker" v-model="infoColor" @input="save()"
          class="appearance-none w-8 h-8 border border-2 p-0 overflow-hidden cursor-pointer">
      </div>

      <div class="flex flex-1 justify-between">
        <p>Border:</p>
        <input type="color" id="borderPicker" v-model="borderColor" @input="save()"
          class="appearance-none w-8 h-8 border border-2 p-0 overflow-hidden cursor-pointer">
      </div>
    </div>

    <div class="w-full p-2">
      <p>Current</p>
      <SongItem :song="audioStore.currentSong" />
    </div>

    <div class="w-full p-2 bg-black">
      <p class="flex-1 flex justify-between" style=" color : #57db5d">StaryNight <button style="border-color : #b3002d"
          class="border rounded-lg p-0.5" @click="save('#000000', '#5e2d8f', '#57db5d', '#b3002d')">Choose
        </button>
      </p>
      <SongItem :song="audioStore.currentSong" :border="'#b3002d'" :action="'#5e2d8f'" :info="'#57db5d'" />
    </div>
    <div class="w-full p-2" style="background-color: #1c1719">
      <p class="flex-1 flex justify-between" style=" color : #ec4889">Default<button style="border-color : #ec4889"
          class="border rounded-lg p-0.5" @click="save('#1c1719', '#eab308', '#ec4889', '#ec4889')">Choose
        </button>
      </p>
      <SongItem :song="audioStore.currentSong" :border="'#ec4889'" :info="'#ec4889'" :action="'#eab308'" />
    </div>

    <div class="w-full p-2" style="background-color: #ff4c4c">
      <p class="flex-1 flex justify-between" style="color: #ffffff">
        Bright Sunset
        <button style="border-color: #ffffff" class="border rounded-lg p-0.5"
          @click="save('#ff4c4c', '#ffcc00', '#ffffff', '#ffffff')">Choose</button>
      </p>
      <SongItem :song="audioStore.currentSong" :border="'#ffffff'" :info="'#ffffff'" :action="'#ffcc00'" />
    </div>

    <div class="w-full p-2" style="background-color: #003d00">
      <p class="flex-1 flex justify-between" style="color: #e0f8d8">
        Forest Night
        <button style="border-color: #e0f8d8" class="border rounded-lg p-0.5"
          @click="save('#003d00', '#a8d5a2', '#e0f8d8', '#e0f8d8')">Choose</button>
      </p>
      <SongItem :song="audioStore.currentSong" :border="'#e0f8d8'" :info="'#e0f8d8'" :action="'#a8d5a2'" />
    </div>

    <div class="w-full p-2" style="background-color: #00274d">
      <p class="flex-1 flex justify-between" style="color: #00ffff">
        Electric Blue
        <button style="border-color: #00ffff" class="border rounded-lg p-0.5"
          @click="save('#00274d', '#0099ff', '#00ffff', '#00ffff')">Choose</button>
      </p>
      <SongItem :song="audioStore.currentSong" :border="'#00ffff'" :info="'#00ffff'" :action="'#0099ff'" />
    </div>

  </main>
</template>

<style scoped>
input[type='color']::-webkit-color-swatch-wrapper {
  padding: 0;
}

input[type='color']::-webkit-color-swatch {
  border: none;
}
</style>
