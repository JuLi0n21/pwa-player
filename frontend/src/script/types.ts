import type { Apiv1Song, v1CollectionPreview } from '@/generated';

export type Song = {
	hash: string;
	name: string;
	artist: string;
	length: number;
	url: string;
	previewimage: string;
	mapper: string;
};

const basePath = import.meta.env.BACKEND_URL || 'http://localhost:8080';

export function mapToSong(apiSong: Apiv1Song): Song {
  return {
    hash: apiSong.md5Hash,
    name: apiSong.title,
    artist: apiSong.artist,
    length: Number(apiSong.totalTime),
    url: `${basePath}/api/v1/audio/${btoa(apiSong.folder + "/" + apiSong.audio).replace(/=+$/, '')}`,
    previewimage: `${basePath}/api/v1/image/${btoa(apiSong.image === "" || apiSong.image === undefined ? "404.png" : apiSong.image).replace(/=+$/, '')}`,
    mapper: apiSong.creator,
  };
}

export type CollectionPreview = {
	index: number;
	name: string;
	length: number;
	previewimage: string;
};


export function mapToCollectionPreview(
  apiCollection: v1CollectionPreview,
  index: number
): CollectionPreview {
  return {
    index,
    name: apiCollection.name,
    length: apiCollection.items,
    previewimage: `${basePath}/api/v1/images/${apiCollection.image}`,
  };
}

export type Me = {
	id: number;
	name: string;
	avatar_url: string;
	endpoint: string;
	share: boolean;
};

export function mapApiToSongs(apiSongs: Apiv1Song[]): Song[] {
  return apiSongs.map(mapToSong);
}

export function mapApiToCollectionPreview(
  apiCollections: v1CollectionPreview[]
): CollectionPreview[] {
  return apiCollections.map((c, i) => mapToCollectionPreview(c, i));
}