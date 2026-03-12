import type { Apiv1Song, v1CollectionPreview, v1Collection } from "@/generated";

export type Song = {
  hash: string;
  name: string;
  artist: string;
  length: number;
  url: string;
  previewimage: string;
  mapper: string;
};

const basePath = import.meta.env.BACKEND_URL || "http://localhost:8080";

export function mapToSong(apiSong: Apiv1Song): Song {
  const image = apiSong.image;
  const imageIsMissing = !image || image === "404.png";

  return {
    hash: apiSong.md5Hash,
    name: apiSong.title,
    artist: apiSong.artist,
    length: Number(apiSong.totalTime),
    url: `${basePath}/api/v1/audio/${btoa(apiSong.folder + "/" + apiSong.audio).replace(/=+$/, "")}`,
    previewimage: imageIsMissing ? "/404.gif" : `${basePath}/api/v1/image/${btoa(image).replace(/=+$/, "")}`,
    mapper: apiSong.creator,
  };
}

export type CollectionPreview = {
  index: number;
  name: string;
  length: number;
  previewimage: string;
};

export type Collection = {
  name: string;
  items: number;
  songs: Song[];
};

export function mapApiToCollection(coll: v1Collection): Collection {
  return {
    name: coll.name,
    items: coll.items,
    songs: mapApiToSongs(coll.songs),
  };
}

export function mapToCollectionPreview(apiCollection: v1CollectionPreview, index: number): CollectionPreview {
  const image = apiCollection.image;
  const imageIsMissing = !image || image === "404.png";

  return {
    index: index,
    name: apiCollection.name,
    length: apiCollection.items,
    previewimage: imageIsMissing ? "/404.gif" : `${basePath}/api/v1/image/${btoa(image).replace(/=+$/, "")}`,
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

export function mapApiToCollectionPreview(apiCollections: v1CollectionPreview[], offset: number): CollectionPreview[] {
  return apiCollections.map((c, i) => mapToCollectionPreview(c, i + offset));
}
