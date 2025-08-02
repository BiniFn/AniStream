# github.com/coeeter/aniways

# Aniways API routes (auto-generated)

## Routes

<details>
<summary>`/`</summary>

- **/**
	- _GET_
		- [RegisterRoutes.func1]()

</details>
<details>
<summary>`/anime/genres`</summary>

- **/anime**
	- **/genres**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).listGenres-fm]()

</details>
<details>
<summary>`/anime/genres/{genre}`</summary>

- **/anime**
	- **/genres/{genre}**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).animeByGenre-fm]()

</details>
<details>
<summary>`/anime/popular`</summary>

- **/anime**
	- **/popular**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).popularAnimes-fm]()

</details>
<details>
<summary>`/anime/random`</summary>

- **/anime**
	- **/random**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).randomAnime-fm]()

</details>
<details>
<summary>`/anime/recently-updated`</summary>

- **/anime**
	- **/recently-updated**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).listRecentlyUpdated-fm]()

</details>
<details>
<summary>`/anime/search`</summary>

- **/anime**
	- **/search**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).searchAnimes-fm]()

</details>
<details>
<summary>`/anime/seasonal`</summary>

- **/anime**
	- **/seasonal**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).seasonalAnimes-fm]()

</details>
<details>
<summary>`/anime/trending`</summary>

- **/anime**
	- **/trending**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).trendingAnimes-fm]()

</details>
<details>
<summary>`/anime/{id}`</summary>

- **/anime/{id}**
	- **/**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getAnimeByID-fm]()

</details>
<details>
<summary>`/anime/{id}/banner`</summary>

- **/anime/{id}**
	- **/banner**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getAnimeBanner-fm]()

</details>
<details>
<summary>`/anime/{id}/franchise`</summary>

- **/anime/{id}**
	- **/franchise**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getAnimeFranchise-fm]()

</details>
<details>
<summary>`/anime/{id}/trailer`</summary>

- **/anime/{id}**
	- **/trailer**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getAnimeTrailer-fm]()

</details>
<details>
<summary>`/anime/{id}/episodes`</summary>

- **/anime/{id}/episodes**
	- **/**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getAnimeEpisodes-fm]()

</details>
<details>
<summary>`/anime/{id}/episodes/{episodeID}/langs`</summary>

- **/anime/{id}/episodes**
	- **/{episodeID}/langs**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getEpisodeLangs-fm]()

</details>
<details>
<summary>`/anime/{id}/episodes/{episodeID}/stream/{type}`</summary>

- **/anime/{id}/episodes**
	- **/{episodeID}/stream/{type}**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getEpisodeStream-fm]()

</details>
<details>
<summary>`/anime/{id}/episodes/{episodeID}/stream/{type}/metadata`</summary>

- **/anime/{id}/episodes**
	- **/{episodeID}/stream/{type}/metadata**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getEpisodeStreamMetadata-fm]()

</details>
<details>
<summary>`/auth/forget-password`</summary>

- **/auth**
	- **/forget-password**
		- _POST_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).forgetPassword-fm]()

</details>
<details>
<summary>`/auth/login`</summary>

- **/auth**
	- **/login**
		- _POST_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).login-fm]()

</details>
<details>
<summary>`/auth/logout`</summary>

- **/auth**
	- **/logout**
		- _POST_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).logout-fm]()

</details>
<details>
<summary>`/auth/me`</summary>

- **/auth**
	- **/me**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).me-fm]()

</details>
<details>
<summary>`/auth/providers`</summary>

- **/auth**
	- **/providers**
		- _GET_
			- [RequireUser]()
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getProviders-fm]()

</details>
<details>
<summary>`/auth/providers/{provider}`</summary>

- **/auth**
	- **/providers/{provider}**
		- _DELETE_
			- [RequireUser]()
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).deleteProvider-fm]()

</details>
<details>
<summary>`/auth/reset-password/{token}`</summary>

- **/auth**
	- **/reset-password/{token}**
		- _PUT_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).resetPassword-fm]()

</details>
<details>
<summary>`/auth/u/{token}`</summary>

- **/auth**
	- **/u/{token}**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getUser-fm]()

</details>
<details>
<summary>`/auth/oauth/{provider}`</summary>

- **/auth/oauth**
	- **/{provider}**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).beginAuthHandler-fm]()

</details>
<details>
<summary>`/auth/oauth/{provider}/callback`</summary>

- **/auth/oauth**
	- **/{provider}/callback**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).callbackHandler-fm]()

</details>
<details>
<summary>`/healthz`</summary>

- **/healthz**
	- _GET_
		- [RegisterRoutes.func2]()

</details>
<details>
<summary>`/library`</summary>

- **/library**
	- **/**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getLibrary-fm]()

</details>
<details>
<summary>`/library/continue-watching`</summary>

- **/library**
	- **/continue-watching**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getContinueWatching-fm]()

</details>
<details>
<summary>`/library/import`</summary>

- **/library**
	- **/import**
		- _POST_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).importLibrary-fm]()

</details>
<details>
<summary>`/library/import/{id}`</summary>

- **/library**
	- **/import/{id}**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getLibraryImportStatus-fm]()

</details>
<details>
<summary>`/library/planning`</summary>

- **/library**
	- **/planning**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getPlanning-fm]()

</details>
<details>
<summary>`/library/{animeID}`</summary>

- **/library**
	- **/{animeID}**
		- _POST_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).createLibrary-fm]()
		- _PUT_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).updateLibrary-fm]()
		- _DELETE_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).deleteAnimeFromLib-fm]()
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getAnimeStatus-fm]()

</details>
<details>
<summary>`/settings`</summary>

- **/settings**
	- **/**
		- _GET_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).getSettings-fm]()
		- _POST_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).saveSettings-fm]()

</details>
<details>
<summary>`/users`</summary>

- **/users**
	- **/**
		- _DELETE_
			- [RequireUser]()
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).deleteUser-fm]()
		- _POST_
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).createUser-fm]()
		- _PUT_
			- [RequireUser]()
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).updateUser-fm]()

</details>
<details>
<summary>`/users/image`</summary>

- **/users**
	- **/image**
		- _PUT_
			- [RequireUser]()
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).updateImage-fm]()

</details>
<details>
<summary>`/users/password`</summary>

- **/users**
	- **/password**
		- _PUT_
			- [RequireUser]()
			- [oeeter/aniways/internal/transport/http/handlers.(*Handler).updatePassword-fm]()

</details>

Total # of routes: 38
