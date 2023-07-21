// @ts-nocheck
import axios from "axios";
import { useRouter } from "next/router";
import { useEffect } from "react";
import jsCookie from "js-cookie";
import { toast } from "react-toastify";

function Callback() {
  const router = useRouter();
  const handleCallback = async () => {
    const code = router.query.code; // Get the code query parameter from the URL
    if (!code) {
      toast.error("Error logging in to YouTube!");
      router.push("/");
    }
    try {
      const response = await axios.post(
        `${
          process.env.NODE_ENV != "production" ? "http://localhost:8080" : ""
        }/api/auth/youtube/callback`,
        { code }
      );
      let token = response.data.token;
      if (token) {
        const expiryTime = response.data.token.expiry;

        const expiryDate = new Date(expiryTime);
        const currentDate = new Date();

        const timeDiff = expiryDate - currentDate;
        const hoursDiff = timeDiff / (1000 * 60 * 60);
        jsCookie.set("yt-token", JSON.stringify(response.data.token), {
          expires: 0.04,
        });
        console.log("yt-token", JSON.stringify(response.data));
        toast.success("Successfully logged in to YouTube!");
        if (jsCookie.get("spotify-token")) {
          router.push("/merger");
        } else {
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

  const gifs = [
    "https://gifdb.com/images/high/comedian-jim-carrey-as-hackerman-vmf9qnz7nx5p9grz.gif",
    "/gifs/1.gif",
    "/gifs/2.gif",
    "/gifs/3.gif",
  ];

  const getRandomGifUrl = () => {
    const randomIndex = Math.floor(Math.random() * gifs.length);
    return gifs[randomIndex];
  };

  const randomGifUrl = getRandomGifUrl();

  return (
    <div className="container mx-auto">
      <div className=" flex flex-col items-center justify-center py-16 xl:py-24">
        <h1 className=" text-4xl font-bold text-primary">
          Please wait we are doing Science!
        </h1>
        <img className=" w-[600px] mt-6" src={randomGifUrl} alt="some gif" />
      </div>
    </div>
  );
}

export default Callback;
