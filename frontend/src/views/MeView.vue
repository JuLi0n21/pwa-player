<script setup lang="ts">
import { ref, onMounted } from "vue";
import SongItem from "../components/SongItem.vue";
import { useAudio } from "@/composables/useAudio";
import { useUser } from "@/composables/useUser";
import type { Me } from "@/script/types";

const audioStore = useAudio();
const userStore = useUser();

const bgColor = ref("");
const actionColor = ref("");
const infoColor = ref("");
const borderColor = ref("");

const loginStatus = ref("Login");

const isHealthy = ref<boolean | null>(null);

function save(bg: string | null, main: string | null, info: string | null, border: string | null) {
  document.documentElement.style.setProperty("--background-color", bg ?? bgColor.value);
  document.documentElement.style.setProperty("--action-color", main ?? actionColor.value);
  document.documentElement.style.setProperty("--information-color", info ?? infoColor.value);
  document.documentElement.style.setProperty("--border-color", border ?? borderColor.value);

  localStorage.setItem("bgColor", bg ?? bgColor.value);
  localStorage.setItem("actionColor", main ?? actionColor.value);
  localStorage.setItem("infoColor", info ?? infoColor.value);
  localStorage.setItem("borderColor", border ?? borderColor.value);
}

async function getMe() {
  const data = (await userStore.fetchMe()) as Me;
  if (data.redirected == true) {
    loginStatus.value = "waiting for login, click to refresh!";
    console.log("redirect detected");
  }

  if (data.id === null || data.id === undefined || Object.keys(data).length === 0) {
    return;
  }

  userStore.setUser(data);
  userStore.cloudflareUrl.value = data.endpoint;
}

async function pasteFromClipboard() {
  try {
    const text = await navigator.clipboard.readText();
    userStore.cloudflareUrl.value = text;
    await checkHealth(text);
  } catch (e) {
    console.error("Failed to read clipboard");
  }
}

async function copyToClipboard() {
  await navigator.clipboard.writeText(userStore.cloudflareUrl.value);
}

async function checkHealth(url: string) {
  if (!url) return;
  try {
    const response = await fetch(`${url}/ping`);
    isHealthy.value = response.ok;
  } catch (e) {
    isHealthy.value = false;
  }
}

function update(e: Event) {
  const val = (e.target as HTMLInputElement).value;
  userStore.cloudflareUrl.value = val;
  checkHealth(val);
}

onMounted(() => {
  reset();
  if (userStore.cloudflareUrl.value) {
    checkHealth(userStore.cloudflareUrl.value);
  }
});

function reset() {
  bgColor.value = localStorage.getItem("bgColor") || "#1c1719";
  actionColor.value = localStorage.getItem("actionColor") || "#eab308";
  infoColor.value = localStorage.getItem("infoColor") || "#ec4899";
  borderColor.value = localStorage.getItem("borderColor") || "#ec4899";

  document.documentElement.style.setProperty("--background-color", bgColor.value);
  document.documentElement.style.setProperty("--action-color", actionColor.value);
  document.documentElement.style.setProperty("--information-color", infoColor.value);
  document.documentElement.style.setProperty("--border-color", borderColor.value);
}
</script>

<template>
  <header>
    <div class="wrapper">
      <nav class="flex justify-start space-x-1 mx-1 my-2">
        <RouterLink class="shadow-xl backdrop--light p-1 rounded-full" to="/">
          <i class="fa-arrow-left fa-solid"></i>
        </RouterLink>
      </nav>
      <hr />
    </div>
  </header>

  <main class="flex flex-col flex-1 h-full overflow-y-scroll">
    <div class="flex flex-col gap-2 p-4">      
      <div class="flex gap-1 overflow-hidden">
        <input 
          @change="update" 
          type="text" 
          id="url-input" 
          :value="userStore.cloudflareUrl.value" 
          class="flex-1 bg-white/5 p-2 border-y border-l rounded-l-lg outline-none bordercolor"
          placeholder="https://..."
        />
        <button @click="copyToClipboard" class="bg-white/5 hover:bg-white/10 p-2 border bordercolor" title="Copy URL">
          <i class="fa-solid fa-copy"></i>
        </button>
        <button @click="pasteFromClipboard" class="bg-white/5 hover:bg-white/10 p-2 border rounded-r-lg text-yellow-500 bordercolor" title="Paste URL">
          <i class="fa-solid fa-paste"></i>
        </button>
        <div 
          v-if="isHealthy !== null"
          class="self-center rounded-full w-2 h-2"
          :class="isHealthy ? 'bg-green-500 shadow-[0_0_8px_green]' : 'bg-red-500 shadow-[0_0_8px_red]'"
        >
      </div>
      </div>
    </div>

    <br />
    
    <button v-if="!userStore.user.value" @click="getMe" class="mx-4 p-0.5 border rounded-lg bordercolor">
      {{ loginStatus }}
    </button>

    <div class="flex flex-col justify-around p-10 w-full">
      <div class="flex flex-1 justify-between">
        <p>Background:</p>
        <input
          type="color"
          id="bgPicker"
          v-model="bgColor"
          @input="save(null, null, null, null)"
          class="p-0 border-2 w-8 h-8 overflow-hidden appearance-none cursor-pointer"
        />
      </div>

      <div class="flex flex-1 justify-between">
        <p>Main:</p>
        <input
          type="color"
          id="actionPicker"
          v-model="actionColor"
          @input="save(null, null, null, null)"
          class="p-0 border-2 w-8 h-8 overflow-hidden appearance-none cursor-pointer"
        />
      </div>

      <div class="flex flex-1 justify-between">
        <p>Secondary:</p>
        <input
          type="color"
          id="infoPicker"
          v-model="infoColor"
          @input="save(null, null, null, null)"
          class="p-0 border-2 w-8 h-8 overflow-hidden appearance-none cursor-pointer"
        />
      </div>

      <div class="flex flex-1 justify-between">
        <p>Border:</p>
        <input
          type="color"
          id="borderPicker"
          v-model="borderColor"
          @input="save(null, null, null, null)"
          class="p-0 border-2 w-8 h-8 overflow-hidden appearance-none cursor-pointer"
        />
      </div>
    </div>
    <div class="p-2 w-full">
      <p>Current</p>
      <SongItem :song="audioStore.currentSong.value" />
    </div>

    <div class="bg-black p-2 w-full">
      <p class="flex flex-1 justify-between" style="color: #57db5d">
        StaryNight
        <button
          style="border-color: #b3002d"
          class="p-0.5 border rounded-lg"
          @click="save('#000000', '#5e2d8f', '#57db5d', '#b3002d')"
        >
          Choose
        </button>
      </p>
      <SongItem :song="audioStore.currentSong.value" :border="'#b3002d'" :action="'#5e2d8f'" :info="'#57db5d'" />
    </div>
    <div class="p-2 w-full" style="background-color: #1c1719">
      <p class="flex flex-1 justify-between" style="color: #ec4889">
        Default<button
          style="border-color: #ec4889"
          class="p-0.5 border rounded-lg"
          @click="save('#1c1719', '#eab308', '#ec4889', '#ec4889')"
        >
          Choose
        </button>
      </p>
      <SongItem :song="audioStore.currentSong.value" :border="'#ec4889'" :info="'#ec4889'" :action="'#eab308'" />
    </div>

    <div class="p-2 w-full" style="background-color: #ff4c4c">
      <p class="flex flex-1 justify-between" style="color: #ffffff">
        Bright Sunset
        <button
          style="border-color: #ffffff"
          class="p-0.5 border rounded-lg"
          @click="save('#ff4c4c', '#ffcc00', '#ffffff', '#ffffff')"
        >
          Choose
        </button>
      </p>
      <SongItem :song="audioStore.currentSong.value" :border="'#ffffff'" :info="'#ffffff'" :action="'#ffcc00'" />
    </div>

    <div class="p-2 w-full" style="background-color: #003d00">
      <p class="flex flex-1 justify-between" style="color: #e0f8d8">
        Forest Night
        <button
          style="border-color: #e0f8d8"
          class="p-0.5 border rounded-lg"
          @click="save('#003d00', '#a8d5a2', '#e0f8d8', '#e0f8d8')"
        >
          Choose
        </button>
      </p>
      <SongItem :song="audioStore.currentSong.value" :border="'#e0f8d8'" :info="'#e0f8d8'" :action="'#a8d5a2'" />
    </div>

    <div class="p-2 w-full" style="background-color: #00274d">
      <p class="flex flex-1 justify-between" style="color: #00ffff">
        Electric Blue
        <button
          style="border-color: #00ffff"
          class="p-0.5 border rounded-lg"
          @click="save('#00274d', '#0099ff', '#00ffff', '#00ffff')"
        >
          Choose
        </button>
      </p>
      <SongItem :song="audioStore.currentSong.value" :border="'#00ffff'" :info="'#00ffff'" :action="'#0099ff'" />
    </div>
  </main>
</template>

<style scoped>
input[type="color"]::-webkit-color-swatch-wrapper {
  padding: 0;
}

input[type="color"]::-webkit-color-swatch {
  border: none;
}
</style>
