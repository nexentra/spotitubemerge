// @ts-nocheck
import axios from "axios";
import { postAPI, getAPI } from "../../utils/api";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

import dynamic from "next/dynamic";
import { toast } from "react-toastify";
import { Beforeunload } from "react-beforeunload";
import getRandomGifUrl from "@/utils/gifSelector";
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
  const [isMergerStarted, setIsMergerStarted] = useState<boolean>(false);

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
    console.log(response?.data, "asdqweasd");
    setSpotifyData(response?.data);
  };

  const startMerger = async () => {
    if (selectedItemsYt.length === 0 || selectedItemsSpotify.length === 0) {
      toast.error("Please select atleast one playlist from each side!!");
      return;
    }
    setIsMergerStarted(true);
    window.scrollTo(0, 0);
    document.body.style.overflow = "hidden";
    const merger = () => postAPI("merge-yt-spotify", {
      "spotify-playlists": selectedItemsSpotify,
      "youtube-playlists": selectedItemsYt,
    });

    merger()
      .then((res) => {
        setIsMergerStarted(false);
        // toast.success("Merger finished successfully!!");
        document.body.style.overflow = "auto";
      })
      .catch((err) => {
        setIsMergerStarted(false);
        // toast.error("Merger failed!! :(");
        document.body.style.overflow = "auto";
      });
    toast.promise(merger, {
      pending: "Promise is pending",
      success: "Promise resolved 👌",
      error: "Promise rejected 🤯",
    });
  };

  useEffect(() => {
    if (!jsCookie.get("spotify-token" && "yt-token")) {
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

  const randomGifUrl = getRandomGifUrl();

  return (
    <>
      {isMergerStarted ? (
        <div className="z-10 fixed w-full h-full flex justify-center bg-gradient-to-br from-sky-50 to-gray-200 dark:bg-trueGray-900 dark:from-gray-800 dark:to-gray-900 duration-700">
          <Beforeunload onBeforeunload={(event) => event.preventDefault()} />
          <div className="container flex flex-col items-center py-10">
            <h1 className=" text-center text-2xl font-bold text-primary ">
              Please wait we are merging your playlists! It might take some time
              depending on your playlist's size. You can continue browsing! 🚀
            </h1>
            <img
              className=" max-w-[600px] max-h-[600px] mt-6"
              src={randomGifUrl}
              alt="some gif"
            />
          </div>
        </div>
      ):
      <div className={`container mx-auto`}>
        <div className="use flex flex-col items-center justify-center">
          <button
            disabled={isMergerStarted}
            className="w-[250px] px-6 py-3 mt-3 text-center text-white bg-indigo-600 rounded-md lg:ml-5"
            onClick={startMerger}
          >
            Start Merger
          </button>
        </div>
        <div className="mx-5 md:mx-0 grid grid-cols-2 gap-4">
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
                data={spotifyData?.playlists}
                selectedItemsFunc={selectedItemsSpotifyFunc}
              />
            </div>
          )}
        </div>
      </div>}
      
    </>
  );
}

export default Playlists;
