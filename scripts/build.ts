import fetch from "node-fetch";
import { DateTime } from "luxon";
import fs from "fs/promises";
import { LinksRecord } from "./types/pb";
import PocketBase from "pocketbase";
import * as dotenv from "dotenv";
dotenv.config();

globalThis.fetch = fetch as any;

(async () => {
	const buildDate = DateTime.now();

	const pb = new PocketBase(process.env.BASE_URL ?? "");

	await pb.collection("users").authWithPassword(process.env.API_USER ?? "", process.env.API_PASSWORD ?? "");

	if (!pb.authStore.isValid) {
		throw new Error("incorrect credentials");
	}

	const shortUrls = await pb.collection("links").getFullList<LinksRecord>({
		filter: "enabled = true",
	});

	if (shortUrls.length === 0) {
		throw new Error("no short urls are found");
	}

	const shortUrlMap = Object.fromEntries(shortUrls.map((shortUrl) => [shortUrl.slug, { slug: shortUrl.slug, url: shortUrl.url }]));
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
