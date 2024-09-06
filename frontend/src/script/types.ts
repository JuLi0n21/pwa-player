export type Song = {
	hash: string;
	name: string;
	artist: string;
	length: number;
	url: string;
	previewimage: string;
	mapper: string;
};

export type CollectionPreview = {
	index: number;
	name: string;
	length: number;
	previewimage: string;
};

export type Me = {
	id: number;
	name: string;
	avatar_url: string;
	endpoint: string;
	share: boolean;
};
