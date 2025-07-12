const developerToken = "eyJhbGciOiJFUzI1NiIsImtpZCI6IjhUNDZDUERESDIifQ.eyJpc3MiOiIzVzlTVDNaUDYzIiwiaWF0IjoxNzUwNTgzNDA0LCJleHAiOjE3NjYzNjA0MDR9.i1M2Ylw67Jqj7hWhsYx4Urj9ZzuUtJgU4lPzfTQOlDXykayei_j7ExToJV0Q2zpMwqRvja1yYzZ2VzMh82tlbA";
const limit = 100;

async function authorizeAppleMusic() {
	await MusicKit.configure({
		developerToken,
		app: {
			name: "LiveRadar",
			build: "0.0.0",
		}
	});
	const music = MusicKit.getInstance();
	await music.authorize();
	const initialRes = await music.api.music("/v1/me/library/artists", {l: "ja-jp", limit});
	const progressBar = document.getElementById("apple-music-progress-bar");
	progressBar.max = initialRes.data.meta.total;
	const artists = initialRes.data.data.map(a => a.attributes.name);
	for (let val = limit; val < progressBar.max; val += limit) {
		progressBar.value = val;
		const res = await music.api.music("/v1/me/library/artists", {l: "ja-jp", limit, offset: val});
		artists.push(...res.data.data.map(a => a.attributes.name))
	}
	progressBar.value = progressBar.max;
	return artists;
}