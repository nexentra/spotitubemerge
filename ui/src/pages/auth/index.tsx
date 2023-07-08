import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import axios from "axios";
import jsCookie from "js-cookie";

const Auth = () => {
  const router = useRouter();
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const fetchData = async (route: String) => {
    try {
      const response = await axios.get("http://localhost:8080" + "/api/auth/" + route);
      setData(response.data);
    } catch (error: any) {
      setError(error.message);
    }
  };

  useEffect(() => {
    if (jsCookie.get("spotify-token") && jsCookie.get("yt-token")) {
      router.push("/");
    }

    if (data) {
      console.log(data, typeof data);
      let newData: any = data;
      newData.authUrl = newData.authUrl.replace(/\u0026/g, "&");
      console.log(newData?.authUrl);
      window.location = newData?.authUrl;
    }
  }, [data]);

  return (
    <div className="flex items-center justify-center h-screen">
      <div className="flex flex-col">
        <button
          onClick={() => fetchData("spotify")}
          type="button"
          className="text-white bg-black font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center mr-2 mb-2"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            fill="currentColor"
            className="w-4 h-4 mr-2"
            viewBox="0 0 16 16"
          >
            {" "}
            <path d="M8 0a8 8 0 1 0 0 16A8 8 0 0 0 8 0zm3.669 11.538a.498.498 0 0 1-.686.165c-1.879-1.147-4.243-1.407-7.028-.77a.499.499 0 0 1-.222-.973c3.048-.696 5.662-.397 7.77.892a.5.5 0 0 1 .166.686zm.979-2.178a.624.624 0 0 1-.858.205c-2.15-1.321-5.428-1.704-7.972-.932a.625.625 0 0 1-.362-1.194c2.905-.881 6.517-.454 8.986 1.063a.624.624 0 0 1 .206.858zm.084-2.268C10.154 5.56 5.9 5.419 3.438 6.166a.748.748 0 1 1-.434-1.432c2.825-.857 7.523-.692 10.492 1.07a.747.747 0 1 1-.764 1.288z" />{" "}
          </svg>
          Sign in with Spotify
        </button>
        <button
          onClick={() => fetchData("youtube")}
          type="button"
          className="text-white bg-[#4285F4] hover:bg-[#4285F4]/90 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center mr-2 mb-2"
        >
          <svg
            className="w-4 h-4 mr-2"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="currentColor"
            viewBox="0 0 18 19"
          >
            <path
              fill-rule="evenodd"
              d="M8.842 18.083a8.8 8.8 0 0 1-8.65-8.948 8.841 8.841 0 0 1 8.8-8.652h.153a8.464 8.464 0 0 1 5.7 2.257l-2.193 2.038A5.27 5.27 0 0 0 9.09 3.4a5.882 5.882 0 0 0-.2 11.76h.124a5.091 5.091 0 0 0 5.248-4.057L14.3 11H9V8h8.34c.066.543.095 1.09.088 1.636-.086 5.053-3.463 8.449-8.4 8.449l-.186-.002Z"
              clip-rule="evenodd"
            />
          </svg>
          Sign in with Google
        </button>
      </div>
    </div>
  );
};

export default Auth;
