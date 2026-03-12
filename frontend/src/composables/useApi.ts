import { MusicBackendApi, Configuration } from "@/generated";
import type { ConfigurationParameters } from "@/generated";
import { useUser } from "./useUser";
import { computed } from "vue";

export function useApi() {
  const userStore = useUser();

  const musicApi = computed(() => {
    const configParams: ConfigurationParameters = {
      basePath: userStore.cloudflareUrl.value,
    };

    const configuration = new Configuration(configParams);
    return new MusicBackendApi(configuration);
  });

  return {
    musicApi,
  };
}
