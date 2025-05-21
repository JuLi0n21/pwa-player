import { MusicBackendApi, Configuration, } from '@/generated';
import type { ConfigurationParameters } from '@/generated';
import { ref } from 'vue';

export function useApi() {
    const basePath = ref(import.meta.env.BACKEND_URL || 'http://localhost:8080');
    const musicApi = (): MusicBackendApi => {
    const configParams: ConfigurationParameters = {
        basePath: basePath.value,
    };

    const configuration = new Configuration(configParams);
    return new MusicBackendApi(configuration);
    };

    return {
        musicApi,
    };
}
