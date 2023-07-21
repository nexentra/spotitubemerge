import dynamic from "next/dynamic";
import Script from "next/script";
import Head from "next/head";
import { ThemeProvider } from "next-themes";
import "../globals.css";
import Navbar from "@/components/ui/navbar";
import Footer from "@/components/ui/footer";
import { ToastContainer } from 'react-toastify';
import { Portal } from '@headlessui/react';
import 'react-toastify/dist/ReactToastify.css';

function App({ Component, pageProps }: any) {
  return (
    <>
      <Head>
        <meta
          name="viewport"
          content="minimum-scale=1, initial-scale=1, width=device-width, shrink-to-fit=no, user-scalable=no, viewport-fit=cover"
        />
      </Head>
      <Portal>
          <ToastContainer autoClose={2000} theme="light" />
        </Portal>
      <ThemeProvider attribute="class">
        <Navbar/>
      <Component {...pageProps} />
      <Footer />
    </ThemeProvider>
      {/* <Script
        strategy="afterInteractive"
        src="https://www.googletagmanager.com/gtag/js?id=test_id"
      />
      <Script id="google-analytics" strategy="afterInteractive">
        {`
          window.dataLayer = window.dataLayer || [];
          function gtag(){window.dataLayer.push(arguments);}
          gtag('js', new Date());
          gtag('config', 'test_id', {
            page_path: window.location.pathname,
          });
          window.gtag = gtag;
        `}
      </Script> */}
    </>
  );
}

export default App;
