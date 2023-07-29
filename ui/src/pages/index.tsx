import { useState, useEffect } from "react";
import axios from "axios";
import Link from "next/link";
import Head from "next/head";
import Navbar from "@/components/ui/navbar";
import Hero from "@/components/ui/hero";
import SectionTitle from "@/components/ui/sectionTitle";
import Benefits from "@/components/ui/benefits";
import { benefitOne, benefitTwo } from "@/components/ui/data";
import Cta from "@/components/ui/cta";
import Footer from "@/components/ui/footer";

export default function Home() {
  return (
    <>
      <Head>
        <title>SpotiTubeMerge - Merge your playlists in few clicks!!</title>
        <meta
          name="description"
          content="SpotiTubeMerge is a simple web app that allows you to merge your Spotify and Youtube playlists in few clicks!!"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Hero />

      <SectionTitle
        pretitle="SpotiTubeMerge Benefits"
        title="Why should you use our service"
      >
        SpotiTubeMerge is an online platform that lets users merge their Spotify
        and YouTube playlists.
      </SectionTitle>

      <Benefits data={benefitOne} />
      <Benefits imgPos="right" data={benefitTwo} />

      <Cta />
    </>
  );
}
