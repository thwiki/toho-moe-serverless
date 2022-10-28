export interface ErrorResult {
	type: string;
	title: string;
	detail: string;
	status: number;
}

export interface ShortUrlsResult {
	shortUrls: ShortUrls;
}

export interface ShortUrls {
	data: ShortUrl[];
	pagination: Pagination;
}

export interface ShortUrl {
	shortCode: string;
	shortUrl: string;
	longUrl: string;
	dateCreated: string;
	visitsCount: number;
	tags: string[];
	meta: ShortUrlMeta;
	domain?: string;
	title?: string;
	crawlable: boolean;
	forwardQuery: boolean;
}

export interface ShortUrlMeta {
	validSince?: string;
	validUntil?: string;
	maxVisits?: string;
}

export interface Pagination {
	currentPage: number;
	pagesCount: number;
	itemsPerPage: number;
	itemsInCurrentPage: number;
	totalItems: number;
}

export type Result = ShortUrlsResult | ErrorResult;
