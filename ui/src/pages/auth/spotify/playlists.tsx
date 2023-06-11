// @ts-nocheck
import axios from "axios";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect,useState } from "react";

function Playlists() {
  const router = useRouter();
  const [data, setData] = useState<any>(null);
  async function fetcher(){
    try {
      const response = await axios.get("/api/spotify-playlist",{
        headers: {
          'Authorization': `${jsCookie.get("spotify-token") }`
        }
      })
      console.log(response.data)
      setData(response.data);
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
          data?.playlists?.items?.map((item:any, key)=>{
            return (
              <div key={key} className="flex flex-col items-center justify-center">
                <img src={item.images[0].url}  alt="" className="w-96 h-96"/>
                <h1 className="text-3xl font-bold text-primary">{item.name}</h1>
              </div>
            )
          }
          )
        }
      </div>
    );
  }
  
  export default Playlists;