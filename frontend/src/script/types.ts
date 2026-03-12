import type { Apiv1Song, V1CollectionPreview, V1CollectionResponse } from "@/generated";

export type Song = {
  hash: string;
  name: string;
  artist: string;
  length: number;
  url: string;
  previewimage: string;
  mapper: string;
};
export function mapToSong(apiSong: Apiv1Song): Song {
  const image = apiSong.image;

  return {
    hash: apiSong.md5Hash,
    name: apiSong.title,
    artist: apiSong.artist,
    length: Number(apiSong.totalTime),
    url: `/api/v1/audio/${btoa(apiSong.folder + "/" + apiSong.audio).replace(/=+$/, "")}`,
    previewimage: image ? `/api/v1/image/${btoa(image).replace(/=+$/, "")}` : "",
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

export function mapApiToCollection(coll: V1CollectionResponse): Collection {
  return {
    name: coll.name,
    items: coll.items,
    songs: mapApiToSongs(coll.songs),
  };
}

export function mapToCollectionPreview(apiCollection: V1CollectionPreview, index: number): CollectionPreview {
  const image = apiCollection.image;

  return {
    index: index,
    name: apiCollection.name,
    length: apiCollection.items,
    previewimage: image ? `/api/v1/image/${btoa(image).replace(/=+$/, "")}` : "",
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

export function mapApiToCollectionPreview(apiCollections: V1CollectionPreview[], offset: number): CollectionPreview[] {
  return apiCollections.map((c, i) => mapToCollectionPreview(c, i + offset));
}
