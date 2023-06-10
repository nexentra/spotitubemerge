"use client"
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

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get("/api/auth/spotify");
        setData(response.data);
      } catch (error:any) {
        setError(error.message);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    if (data) {
      console.log(data, typeof data);
      let newData = JSON.parse(data);
      newData.authUrl = newData.authUrl.replace(/\u0026/g, "&");
      console.log(newData?.authUrl);
      window.location = newData?.authUrl;
    }
  }, [data]);

  return (
    <div>
      <h1 className="bg-gray-100">Hello, world!</h1>
      <p>
        This is <code>pages/index.tsx</code>.
      </p>
      <p>
        Check out <Link href="/foo">foo</Link>.
      </p>

      <h2>Memory allocation stats from Go server</h2>
      {error && (
        <p>
          Error fetching profile: <strong>{error}</strong>
        </p>
      )}
      {!error && !data && <p>Loading ...</p>}
      {!error && data && <pre>{data}</pre>}
    </div>
  );
};
