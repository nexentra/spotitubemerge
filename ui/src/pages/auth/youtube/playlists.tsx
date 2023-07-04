// @ts-nocheck
import axios from "axios";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

function Playlists() {
  const router = useRouter();
  const [data, setData] = useState<any>([]);
  async function fetcher() {
    try {
      let response = await axios.get("/api/youtube-playlist", {
        headers: {
          Authorization: `${jsCookie.get("yt-token")}`,
        },
      });
      setData(JSON.parse(response.data.playlists));
      console.log(JSON.parse(response.data.playlists));
    } catch (error: any) {
      // setError(error.message);
    }
  }

  useEffect(() => {
    if (router?.isReady) {
      fetcher();
    }
  }, [router?.isReady]);

  return (
    <div className="container ">
      <button
        onClick={async () => {
          let data1 = [];
          data1.push(data.items[10].id);
          let queryString = data1.join();
          console.log("data1", data1);
          console.log("queryString", queryString);
          let response1 = await axios.get("/api/youtube-items", {
            headers: {
              Authorization: `${jsCookie.get("yt-token")}`,
            },
            params: {
              strings: queryString,
            },
          });
          console.log("response1", response1.data);
        }}
      >
        test
      </button>
      {data?.items?.map((item: any, key) => {
        return (
          <div key={key} className="flex flex-col items-center justify-center">
            <img
              src={item.snippet.thumbnails.default.url}
              alt=""
              className={`w-[${item.snippet.thumbnails.default.width}] h-[${item.snippet.thumbnails.default.height}]`}
            />
            <h1 className="text-3xl font-bold text-primary">
              {item.snippet.localized.title}
            </h1>
          </div>
        );
      })}
    </div>
  );
}

export default Playlists;
