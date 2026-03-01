# Audio-Stream
This will be an audio server backend in order to allow for audio files to be submitted and then indexed and played back.


## Technologies used 
- SQLite for a database to hold information about the movies
- Gin as the Go HTTP framework powering the API
- GORM as the Go ORM framework for simplified database interactions and migrations
- FFmpeg for aduio streaming

## API Routes

- All endpoints are served at `http://localhost:8080`

### 1. **Create Movie**
- **POST** `/audios`
    - Form Data: title (string, Required), audio file(Required), cover_image (optional file)
    - Response: 201 Created with movie details

### 2 **List Movies**
- **GET** `/audios`
    - Response: 200 OK with list of audios

## 3 **Stream Movie**
- **GET** `/audios/:id/`
    - Response: 200 OK streaming the audio file

### 4 **Update Movie**
- **PUT** `/audios/:id`
    - Form Data: title (string), video (optional file), cover_image (optional file)
    - Response: 200 OK with updated movie details

### 5 **Delete Movie**
- **DELETE** `/audios/:id`
    - Response: 204 No Content

## 6 **Stream Movie**
- **GET** `/audios/:id/HLS`
    - 200 will get the segments.
    - The client then streams the video using HLS.

## FFmpeg Requirement

This project uses FFmpeg to generate HLS playlists and audio segments for streaming.
FFmpeg **must be installed** on your system for the server to work.

### Install FFmpeg

**macOS (Homebrew)**  
```bash
brew install ffmpeg
```
**Ubuntu / Debian**
sudo apt update
sudo apt install ffmpeg

**Windows** 
Download from: `https://ffmpeg.org`

## Depency checks 
- Check vulnerabilities by running 
```govulncheck ./...```
- If govulncheck not instaled
  ```bash
    go install golang.org/x/vuln/cmd/govulncheck@latest

    echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.zshrc
    source ~/.zshrc

    govulncheck -h
    ```


## Future 
- Add real time audio streaming using ws and webrtc to allow for audio to be able to to be streamed.
