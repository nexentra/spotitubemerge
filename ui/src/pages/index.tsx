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
        <title>Spotitubemerge - Merge your playlists in few clicks!!</title>
        <meta
          name="description"
          content="Spotitubemerge is a simple web app that allows you to merge your Spotify and Youtube playlists in few clicks!!"
        />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Hero/>
      
      <SectionTitle
        pretitle="Nextly Benefits"
        title=" Why should you use this landing page">
        Nextly is a free landing page & marketing website template for startups
        and indie projects. Its built with Next.js & TailwindCSS. And its
        completely open-source.
      </SectionTitle>

      <Benefits data={benefitOne} />
      <Benefits imgPos="right" data={benefitTwo} />

      <Cta />
    </>
  );
}