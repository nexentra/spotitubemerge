// @ts-nocheck
import axios from "axios";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

function PlaylistItem() {
  const router = useRouter();

  useEffect(() => {
    async function fetcher() {
      let response1 = await axios.get(
        `${(process.env.NODE_ENV !="production") ? "http://localhost:8080" : ""}/api/spotify-items`,
        {
          headers: {
            AuthorizationSpotify: `${jsCookie.get("spotify-token")}`,
          },
          params: {
            strings: router.query.id,
          },
        }
      );
      console.log("response1", response1.data);
    }

    if (router?.isReady) fetcher();
  }, [router?.isReady]);
  return <></>;
}

export default PlaylistItem;
