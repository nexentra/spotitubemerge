import { useState, useEffect } from "react";
import axios from "axios";
import Link from "next/link";
import Head from "next/head";
import Navbar from "@/components/ui/navbar";

export default function Home() {
  return (
    <>
    <Head>
        <title>Spotitubemerge - Merge your playlists in few clicks!!</title>
        <meta
          name="description"
          content="Spotitubemerge is a simple web app that allows you to merge your Spotify and Youtube playlists in few clicks!!"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Navbar/>
      
    </>
  );
}