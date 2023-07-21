import jsCookie from "js-cookie";
import axios from "axios";
import { toast } from "react-toastify";

interface Params {
  baseUrl: string;
  headers: any;
  method: string;
}

const postConfig: Params = {
  baseUrl: `${
    process.env.NODE_ENV != "production" ? "http://localhost:8080" : ""
  }/api`,
  headers: {
    AuthorizationSpotify: jsCookie.get("spotify-token"),
    AuthorizationYoutube: jsCookie.get("yt-token"),
  },
  method: "post",
};

export const postAPI = async (url: string, data: any): Promise<any> => {
  return await axios({
    ...postConfig,
    url: `${postConfig.baseUrl}/${url}`,
    data,
  })
    .then((response) => {
      console.log(response);
      return {
        status: response.status,
        data: response.data,
      };
    })
    .catch((error) => {
      console.log(error);
      toast.error("Something happened, but I'm not telling you what.");
      return {
        status: error.status,
        data: error.response,
      };
    });
};

const getConfig: Params = {
  baseUrl: `${
    process.env.NODE_ENV != "production" ? "http://localhost:8080" : ""
  }/api`,
  headers: {
    AuthorizationSpotify: jsCookie.get("spotify-token"),
    AuthorizationYoutube: jsCookie.get("yt-token"),
  },
  method: "get",
};

export const getAPI = async (url: string, data: any): Promise<any> => {
  return await axios({
    ...getConfig,
    url: `${getConfig.baseUrl}/${url}${data ? `/${data}` : ""}`,
  })
    .then((response) => {
      console.log(response);
      return {
        status: response.status,
        data: response.data,
      };
    })
    .catch((error) => {
      console.log(error);
      toast.error("Something happened, but I'm not telling you what.");
      return {
        status: error.status,
        data: error.response,
      };
    });
};
