import { useState, useEffect } from "react";
import axios from "axios";
import Link from "next/link";

export default function Home() {
  return (
    <>
      <Ping />
    </>
  );
}

const Ping = () => {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const fetchData = async (route:String) => {
    try {
      const response = await axios.get("/api/auth/"+route);
      setData(response.data);
    } catch (error:any) {
      setError(error.message);
    }
  };

  useEffect(() => {
    if (data) {
      console.log(data, typeof data);
      let newData:any = data
      newData.authUrl = newData.authUrl.replace(/\u0026/g, "&");
      console.log(newData?.authUrl);
      window.location = newData?.authUrl;
    }
  }, [data]);

  return (
    <div>
      <button onClick={()=>fetchData("spotify")} className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">
        Login with Spotify
      </button>
      <button onClick={()=>fetchData("youtube")} className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">
        Login with youtube
      </button>
     </div>
  );
};
