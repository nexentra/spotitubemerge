import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import axios from "axios";
import jsCookie from "js-cookie";
import Link from "next/link";
import Image from "next/image";

const Auth = () => {
  const router = useRouter();
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const fetchData = async (route: String) => {
    try {
      console.log(
        process.env.NODE_ENV != "production",
        typeof (process.env.NODE_ENV != "production"),
        "fetching data"
      );
      const response = await axios.get(
        `${
          process.env.NODE_ENV != "production" ? "http://localhost:8080" : ""
        }/api/auth/` + route
      );
      setData(response.data);
    } catch (error: any) {
      setError(error.message);
    }
  };

  useEffect(() => {
    if (jsCookie.get("spotify-token") && jsCookie.get("yt-token")) {
      router.push("/merger");
    }

    if (data) {
      console.log(data, typeof data);
      let newData: any = data;
      newData.authUrl = newData.authUrl.replace(/\u0026/g, "&");
      console.log(newData?.authUrl);
      window.location = newData?.authUrl;
    }
  }, [router, data]);

  return (
    <div className="relative py-16 h-full">
      <div className="relative container max-h-[500px] m-auto px-6 text-gray-500 md:px-12 xl:px-40">
        <div className="m-auto md:w-8/12 lg:w-6/12 xl:w-6/12">
          <div className="rounded-xl bg-white shadow-xl">
            <div className="p-6 sm:p-16">
              <div className="space-y-4">
                <h2 className="mb-8 text-2xl text-cyan-900 font-bold">
                  Sign in to unlock the <br /> best of{" "}
                  <Image
                    src={
                      "/logo-black.svg"
                    }
                    alt="logo"
                    width="32"
                    height="32"
                    className="w-[210px] inline-block"
                  />
                </h2>
              </div>
              <div className="mt-16 grid space-y-4">
                <button
                  onClick={() => fetchData("youtube")}
                  className="group h-12 px-6 border-2 border-gray-300 rounded-full transition duration-300 hover:border-blue-400 focus:bg-blue-50 active:bg-blue-100"
                >
                  <div className="relative flex items-center space-x-4 justify-center">
                    <svg
                      className="absolute left-0 w-5"
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
                    <span className="block w-max font-semibold tracking-wide text-gray-700 text-sm transition duration-300 group-hover:text-blue-600 sm:text-base">
                      Connect with Google
                    </span>
                  </div>
                </button>
                <button
                  onClick={() => fetchData("spotify")}
                  className="group h-12 px-6 border-2 border-gray-300 rounded-full transition duration-300 hover:border-blue-400 focus:bg-blue-50 active:bg-blue-100"
                >
                  <div className="relative flex items-center space-x-4 justify-center">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      fill="currentColor"
                      className="absolute left-0 w-5"
                      viewBox="0 0 16 16"
                    >
                      {" "}
                      <path d="M8 0a8 8 0 1 0 0 16A8 8 0 0 0 8 0zm3.669 11.538a.498.498 0 0 1-.686.165c-1.879-1.147-4.243-1.407-7.028-.77a.499.499 0 0 1-.222-.973c3.048-.696 5.662-.397 7.77.892a.5.5 0 0 1 .166.686zm.979-2.178a.624.624 0 0 1-.858.205c-2.15-1.321-5.428-1.704-7.972-.932a.625.625 0 0 1-.362-1.194c2.905-.881 6.517-.454 8.986 1.063a.624.624 0 0 1 .206.858zm.084-2.268C10.154 5.56 5.9 5.419 3.438 6.166a.748.748 0 1 1-.434-1.432c2.825-.857 7.523-.692 10.492 1.07a.747.747 0 1 1-.764 1.288z" />{" "}
                    </svg>
                    <span className="block w-max font-semibold tracking-wide text-gray-700 text-sm transition duration-300 group-hover:text-blue-600 sm:text-base">
                      Connect with Spotify
                    </span>
                  </div>
                </button>
                {/* <button className="group h-12 px-6 border-2 border-gray-300 rounded-full transition duration-300 hover:border-blue-400 focus:bg-blue-50 active:bg-blue-100">
                  <div className="relative flex items-center space-x-4 justify-center">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="currentColor"
                      className="absolute left-0 w-5 text-gray-700"
                      viewBox="0 0 16 16"
                    >
                      <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z" />
                    </svg>
                    <span className="block w-max font-semibold tracking-wide text-gray-700 text-sm transition duration-300 group-hover:text-blue-600 sm:text-base">
                      Continue with Github
                    </span>
                  </div>
                </button> */}
              </div>

              <div className="mt-32 space-y-4 text-gray-600 text-center sm:-mb-8">
                <p className="text-xs">
                  By clicking connect with google, connect with spotify you
                  agree to our{" "}
                  <Link href={"terms-n-conditions"} className="underline">
                    Terms and Conditions
                  </Link>{" "}
                  {/* and confirm you have read our{" "}
                  <a href="#" className="underline">
                    Privacy and Cookie Statement
                  </a>
                  .
                </p>
                <p className="text-xs">
                  This site is protected by reCAPTCHA and the{" "}
                  <a href="#" className="underline">
                    Google Privacy Policy
                  </a>{" "}
                  and{" "}
                  <a href="#" className="underline">
                    Terms of Service
                  </a>{" "}
                  apply. */}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Auth;
