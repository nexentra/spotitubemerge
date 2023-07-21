// @ts-nocheck
import axios from "axios";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

import dynamic from "next/dynamic";
const PlaylistTable = dynamic(
  () => import("@/components/spotify/playlist_table"),
  { ssr: false }
);

function Playlists() {
  const router = useRouter();
  const [data, setData] = useState<any>(null);
  const [selectedItems, setSelectedItems] = useState<string[]>([]);
  function selectedItemsFunc(args: string[]) {
    setSelectedItems(args);
  }
  async function fetcher() {
    try {
      const response = await axios.get(
        `${(process.env.NODE_ENV !="production") ? "http://localhost:8080" : ""}/api/spotify-playlist`,
        {
          headers: {
            AuthorizationSpotify: `${jsCookie.get("spotify-token")}`,
          },
        }
      );
      console.log(response.data);
      setData(response.data);
    } catch (error: any) {
      console.log(error);
    }
  }

  useEffect(() => {
    if (router?.isReady) {
      fetcher();
    }
  }, [router?.isReady]);

  return (
    <div className="mx-5 md:mx-0">
      <PlaylistTable
        data={data?.playlists?.items}
        selectedItemsFunc={selectedItemsFunc}
      />
    </div>
  );
}

export default Playlists;
