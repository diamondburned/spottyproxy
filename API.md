## API Draft

```
/api/v0:
	-> /login
	<- 200 { username, sessionToken }
	-> Authorization: sessionToken

/login
	/login { username, password, secret } POST - login to spotify

/search - search for anything
	/search?query=foo - search for foo
	/search?query=foo&suggest - search for foo and suggest

/track
	/track/{id} - get track
		/track/{id}/audio.(opus|mp3) - play track, medium cache
		/track/{id}/cover.(jpg|png) - get cover, aggressive cache

/artist
	/artist/{id} - get artist info
	/artist/{id}/tracks - get artist tracks
	/artist/{id}/albums - get artist albums

/album
	/album/{id} - get album info
	/album/{id}/tracks - get album tracks

/playlist - list user's playlists
	/playlist POST - create playlist
	/playlist/{id} - get specific playlist
		/playlist/{id} POST - update playlist
		/playlist/{id} DELETE - delete playlist
```

## Implementation

Each user gets no sessions or one session at most. Whenever the user logs into
the server, they will always get the same session and the same session token
back, provided that they had the correct credentials and that the server didn't
restart.

Some per-user mutex synchronization is needed: we'd want to lock the mutex
almost every time a request is invoked on the server. This becomes harder to do
when we acquire a lock to stream over a music file, since seeking implies it
would also acquire the same lock. However, we can't not lock while streaming a
music file, as that would break the file. Or can we?
