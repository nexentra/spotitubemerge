// @ts-nocheck
import axios from "axios";
import { postAPI, getAPI } from "../../utils/api";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

import dynamic from "next/dynamic";
import { toast } from "react-toastify";
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

  const fetchYoutubePlaylists = async () => {
    const response = await getAPI("youtube-playlist");
    setYoutubeData(response?.data?.playlists);
  };

  const fetchSpotifyPlaylists = async () => {
    const response = await getAPI("spotify-playlist");
    setSpotifyData(response?.data);
  };

  const startMerger = async () =>
    postAPI("merge-yt-spotify", {
      "spotify-playlists": selectedItemsSpotify,
      "youtube-playlists": selectedItemsYt,
    });

  useEffect(() => {
    if (!jsCookie.get("spotify-token" && "yt-token")){
      toast.error("Please login to continue!!");
      router.push("/auth");
      return;
    }
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
