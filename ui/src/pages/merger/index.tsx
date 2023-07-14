// @ts-nocheck
import axios from "axios";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

import dynamic from "next/dynamic";
const YoutubePlaylistTable = dynamic(
  () => import("@/components/youtube/playlist_table"),
  { ssr: false }
);
const SpotifyPlaylistTable = dynamic(
  () => import("@/components/spotify/playlist_table"),
  { ssr: false }
);

function Playlists() {
  const router = useRouter();
  const [youtubeData, setYoutubeData] = useState<any>(null);
  const [spotifyData, setSpotifyData] = useState<any>(null);
  const [selectedItemsYt, setSelectedItemsYt] = useState<string[]>([]);
  const [selectedItemsSpotify, setSelectedItemsSpotify] = useState<string[]>(
    []
  );

  function selectedItemsYtFunc(args: string[]) {
    setSelectedItemsYt(args);
  }

  function selectedItemsSpotifyFunc(args: string[]) {
    setSelectedItemsSpotify(args);
  }

  async function fetchYoutubePlaylists() {
    try {
      const response = await axios.get(
        "http://localhost:8080/api/youtube-playlist",
        {
          headers: {
            AuthorizationYoutube: `${jsCookie.get("yt-token")}`,
          },
        }
      );
      console.log(response?.data?.playlists);
      setYoutubeData(response?.data?.playlists);
    } catch (error: any) {
      console.log(error);
    }
  }

  async function fetchSpotifyPlaylists() {
    try {
      const response = await axios.get(
        "http://localhost:8080/api/spotify-playlist",
        {
          headers: {
            Authorization: `${jsCookie.get("spotify-token")}`,
          },
        }
      );
      console.log(response?.data);
      setSpotifyData(response?.data);
    } catch (error: any) {
      console.log(error);
    }
  }

  const startMerger = async () => {
    try {
      const response = await axios.post("http://localhost:8080" + "/api/merge-yt-spotify", {
        "spotify-playlists": selectedItemsSpotify,
        "youtube-playlists": selectedItemsYt
       },
      {
        headers: {
          AuthorizationSpotify: `${jsCookie.get("spotify-token")}`,
          AuthorizationYoutube: `${jsCookie.get("yt-token")}`,
        }
      });
      console.log(response,"response");
    } catch (error) {
      console.error("Error:", error);
      // Handle the error appropriately
    }
  };

  useEffect(() => {
    if (router?.isReady) {
      fetchYoutubePlaylists();
      fetchSpotifyPlaylists();
    }
  }, [router?.isReady]);

  useEffect(() => {
    console.log(selectedItemsYt, "yt");
    console.log(selectedItemsSpotify, "spotify");
  }, [selectedItemsYt, selectedItemsSpotify]);

  return (
    <div className="mx-5 md:mx-0 grid grid-cols-2 gap-4">
      <button onClick={startMerger}>Start Merger</button>
      {youtubeData && (
        <div>
          <h2>YouTube Playlists</h2>
          <YoutubePlaylistTable
            mergerPage={true}
            data={youtubeData}
            selectedItemsFunc={selectedItemsYtFunc}
          />
        </div>
      )}
      {spotifyData && (
        <div>
          <h2>Spotify Playlists</h2>
          <SpotifyPlaylistTable
            mergerPage={true}
            data={spotifyData?.playlists?.items}
            selectedItemsFunc={selectedItemsSpotifyFunc}
          />
        </div>
      )}
    </div>
  );
}

export default Playlists;
