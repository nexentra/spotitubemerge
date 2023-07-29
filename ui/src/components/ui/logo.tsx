import { useTheme } from "next-themes";
import Image from "next/image";
import { useEffect, useState } from "react";

const LogoImage = () => {
  const [mounted, setMounted] = useState(false);
  const { theme } = useTheme();

  useEffect(() => setMounted(true), []);

  return (
    <>
    {
      mounted && 
        <Image
      src={theme === "dark" ? "/logo-white.svg" : "/logo-black.svg"}
      alt={"logo"}
      width={32}
      height={32}
      className={"w-[210px] h-[50px]"}
    />
        
    }
    </>
  );
};

export default LogoImage;
