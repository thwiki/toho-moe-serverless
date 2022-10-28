import fetch, { Headers } from "node-fetch";
import { DateTime } from "luxon";
import fs from "fs/promises";
import { Result, ShortUrl } from "./types/shlink";
import * as dotenv from "dotenv";
dotenv.config();

async function* crawlShortUrls(source: string, apiKey: string) {
	let page = 1;

	while (true) {
		const url = new URL(source);
		url.searchParams.set("page", page.toString(10));

		const result = (await (await fetch(url.href, { method: "GET", headers: new Headers({ "X-Api-Key": apiKey }) })).json()) as Result;

		if ("detail" in result) throw new Error(result.detail);

		yield* result.shortUrls.data;

		const pagination = result.shortUrls.pagination;
		if (pagination.currentPage >= pagination.pagesCount || pagination.itemsInCurrentPage === 0) {
			break;
		}
		page++;
	}
}

(async () => {
	const buildDate = DateTime.now();

	const shortUrls: ShortUrl[] = [];
	for await (const url of crawlShortUrls(process.env.SOURCE_URL ?? "", process.env.API_KEY ?? "")) {
		shortUrls.push(url);
	}

	if (shortUrls.length === 0) {
		throw new Error("no short urls are found");
	}

	const shortUrlMap = Object.fromEntries(shortUrls.map((shortUrl) => [shortUrl.shortCode, { slug: shortUrl.shortCode, url: shortUrl.longUrl }]));
	const shortUrlMapJson = JSON.stringify(shortUrlMap);
	const shortUrlMapJsonBytes = new TextEncoder().encode(shortUrlMapJson);

	await fs.writeFile(
		"./utils/data.go",
		`package utils

var date = ${buildDate.toMillis().toString(10)}
var data = []byte\{${Array.from(shortUrlMapJsonBytes).join(", ")}\}
`,
		{ encoding: "utf-8" }
	);
})();
