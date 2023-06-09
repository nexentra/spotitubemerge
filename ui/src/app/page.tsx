 // @ts-ignore 
 "use client";
import { useState } from 'react';
import axios from 'axios';
import useSWR from "swr";
import Link from "next/link";

export default function Home() {
  return (
    <>
    <Ping/>
    </>
  )
}



const Ping = () => {
  async function fetcher(url: string) {
    const resp = await fetch(url);
    return resp.text();
  }
  const { data, error } = useSWR("/api", fetcher, { refreshInterval: 1000 });
  return  <div>
  <h1 className='bg-gray-100'>Hello, world!</h1>
  <p>This is <code>pages/index.tsx</code>.</p>
  <p>Check out <Link href="/foo">foo</Link>.</p>

  <h2>Memory allocation stats from Go server</h2>
  {error && <p>Error fetching profile: <strong>{error}</strong></p>}
  {!error && !data && <p>Loading ...</p>}
  {!error && data && <pre>{data}</pre>}
</div>
};