import { ref } from 'vue';
import type { Me } from '@/script/types';

let userInstance: ReturnType<typeof createUser> | null = null;

function createUser() {
  const user = ref<Me | null>(null);
  const baseUrl = ref(import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080');
  const proxyUrl = ref(import.meta.env.VITE_PROXY_URL || 'http://localhost:8081');

  function saveUser(u: Me | null) {
    localStorage.setItem('activeUser', JSON.stringify(u));
  }

  function loadUser(): Me | null {
    const u = localStorage.getItem('activeUser');
    return u ? JSON.parse(u) : null;
  }

  function setUser(u: Me | null) {
    user.value = u;
    saveUser(u);
  }

  async function fetchMe(): Promise<Me | {}> {
    try {
      const response = await fetch(`${proxyUrl.value}/me`, {
        method: 'GET',
        credentials: 'include',
      });

      if (response.redirected) {
        window.open(response.url, '_blank');
        return { redirected: true };
      }

      if (!response.ok) {
        console.error(`Fetch failed: ${response.status} ${response.statusText}`);
        return { id: -1 } as Me;
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Fetch error:', error);
      return {};
    }
  }

  setUser(loadUser());

  return {
    user,
    baseUrl,
    proxyUrl,
    setUser,
    fetchMe,
  };
}

export function useUser() {
  if (!userInstance) {
    userInstance = createUser();
  }
  return userInstance;
}
