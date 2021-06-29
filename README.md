# Visio

##### Secure, limitless, decentralised video storage

Visio is a service which allows you to store MP4 video of unlimited size on the Filecoin network via file.video for
free. It contains a web server which accepts HTTP requests for this purpose.

### Functionality

When accepting a video upload, Visio will first split the MP4 video into 30MB chunks if required. Then it will upload
each chunk to Filecoin, receiving a unique link to a HLS master playlist for each. Then it will parse the master
playlist and it's child sub-playlists to extract the links to each MPEG-TS file.

Using the links to the .ts files, it will form a JSON file mapping a pseudo .ts file to the actual link. This will be
required when the user requests a video such as `/x/video_id/1080p/4.ts` because it will lookup the resolution and the
index of the .ts file to request the actual file and relay that to the user. These JSON mappings are stored on Storj.io
in a decentralised fashion.

Should the user request the master playlist (e.g. `/x/video_id/root.m3u8`), Visio will compute the required M3U8
playlist file based on the entries for the video on Storj. The same is done if a sub-playlist file is requested (
e.g. `/x/video_id/1080p.m3u8`).

### Endpoints

- `/upload` - POST a video here with a header named "ID" to upload the video. TODO Compute a hash of each video (e.g.
  SHA512) to act as the video ID.
- `/x/{video_id}/root.m3u8` - GET the master M3U8 playlist for the specified video
- `/x/{video_id}/{resolution}.m3u8` - GET the resolution-specific sub-playlist for the video
- `/x/{video_id}/{resolution}/{n}.ts` - GET the n<sup>th</sup> MPEG-TS file associated for the video with the specified
  resolution.