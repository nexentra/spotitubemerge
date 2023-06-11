// @ts-nocheck
import axios from "axios";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect,useState } from "react";

function Playlists() {
  const router = useRouter();
  const [data, setData] = useState<any>([]);
  async function fetcher(){
    try {
      const response = await axios.get("/api/youtube-playlist",{
        headers: {
          'Authorization': `${jsCookie.get("yt-token") }`
        }
      })
      setData(JSON.parse(response.data.playlists));
      console.log(JSON.parse(response.data.playlists))
    } catch (error:any) {
      // setError(error.message);
    }
  }

  useEffect(() => {
    if (router?.isReady) {
      fetcher()
    }
  }, [router?.isReady]);

    return (
      <div className="container ">
        {
         data?.items?.map((item:any, key)=>{
            return (
              <div key={key} className="flex flex-col items-center justify-center">
                <img src={item.snippet.thumbnails.default.url}  alt="" className={`w-[${item.snippet.thumbnails.default.width}] h-[${item.snippet.thumbnails.default.height}]`}/>
                <h1 className="text-3xl font-bold text-primary">{item.snippet.localized.title}</h1>
              </div>
            )
          }
          )
        }
      </div>
    );
  }
  
  export default Playlists;