// @ts-nocheck
import axios from "axios";
import { useRouter } from "next/router";
import { useEffect } from "react";
import jsCookie from "js-cookie";

function Callback() {
  const router = useRouter();
  const handleCallback = async () => {
    const code = router.query.code; // Get the code query parameter from the URL
    try {
      const response = await axios.post("/api/auth/youtube/callback", { code });
      let token = response.data.token
      if (token) {
        const expiryTime = response.data.token.expiry;

        const expiryDate = new Date(expiryTime);
        const currentDate = new Date();

        const timeDiff = expiryDate - currentDate;
        const hoursDiff = timeDiff / (1000 * 60 * 60);
        jsCookie.set("yt-token", JSON.stringify(response.data.token), {
          expires: hoursDiff,
        });
        router.push("/auth/youtube/playlists");
      }
    } catch (error) {
      console.error("Error:", error);
      // Handle the error appropriately
    }
  };

  useEffect(() => {
    if (router?.isReady) {
      handleCallback();
    }
  }, [router?.isReady]);

  return (
    <div className="container ">
      <div className=" flex h-[calc(100vh-160px)] items-center justify-center py-16 xl:py-24">
        <h1 className=" text-6xl font-bold text-primary">Coming Soon!</h1>
      </div>
    </div>
  );
}

export default Callback;
