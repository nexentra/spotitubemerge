// @ts-nocheck
import axios from "axios";
import jsCookie from "js-cookie";
import { useRouter } from "next/router";
import { useEffect,useState } from "react";
import { setCookie } from "cookies-next";
import { toast } from "react-toastify";
import getRandomGifUrl from "@/utils/gifSelector";

function Callback() {
  const router = useRouter();
  const handleCallback = async () => {
    const code = router.query.code; // Get the code query parameter from the URL
    if (!code){
    toast.error("Error logging in to Spotify!");
    router.push("/");
  }
    try {
      const response = await axios.post(`${(process.env.NODE_ENV !="production") ? "http://localhost:8080" : ""}/api/auth/spotify/callback`, { code });
      let token = response.data.token
      if (token) {
        const expiryTime = response.data.token.expiry;

        const expiryDate = new Date(expiryTime);
        const currentDate = new Date();

        const timeDiff = expiryDate - currentDate;
        const hoursDiff = timeDiff / (1000 * 60 * 60);
        jsCookie.set("spotify-token", JSON.stringify(response.data.token), {
          expires: 0.04,
        });
        toast.success("Successfully logged in to Spotify!");
        if (jsCookie.get("yt-token")) {
          router.push("/merger");
        }
        else{
          router.push("/auth");
        }
      }
    } catch (error) {
      console.error("Error:", error);
      toast.error("Error logging in to YouTube!");
      router.push("/");
    }
  };

  useEffect(() => {
    if (router?.isReady) {
      handleCallback();
    }
  }, [router?.isReady]);

  const randomGifUrl = getRandomGifUrl();

    return (
      <div className="container mx-auto">
      <div className=" flex flex-col items-center justify-center py-16 xl:py-24">
        <h1 className=" text-4xl font-bold text-primary">
          Please wait we are doing Science!
        </h1>
        <img className=" max-w-[600px] max-h-[600px] mt-6" src={randomGifUrl} alt="some gif" />
      </div>
    </div>
    );
  }
  
  export default Callback;