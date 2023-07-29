import {
  FaceSmileIcon,
  ChartBarSquareIcon,
  CursorArrowRaysIcon,
  DevicePhoneMobileIcon,
  AdjustmentsHorizontalIcon,
  SunIcon,
} from "@heroicons/react/24/solid";

import benefitOneImg from "/public/img/benefit-one.png";
import benefitTwoImg from "/public/img/benefit-two.png";

const benefitOne = {
  title: "Some of the benefits of using SpotiTubeMerge",
  desc: "Users can access SpotiTubeMerge anytime, anywhere, as it is an online platform. This flexibility allows for merging playlists and enjoying music across devices and locations.",
  image: benefitOneImg,
  bullets: [
    {
      title: "Merge Spotify and YouTube Playlists",
      desc: "SpotiTubeMerge allows users to effortlessly combine and synchronize their Spotify and YouTube playlists. This feature provides a unified music experience by consolidating playlists from both platforms.",
      icon: <FaceSmileIcon />,
    },
    {
      title: "Seamless Music Experience",
      desc: "By merging playlists, users can seamlessly switch between their favorite songs from Spotify and YouTube, enhancing their music enjoyment without the hassle of switching between apps.",
      icon: <ChartBarSquareIcon />,
    },
    {
      title: "User Privacy and Data Security",
      desc: "SpotiTubeMerge takes user privacy and data security seriously, ensuring that personal information is handled responsibly. Users can trust the platform with their account information and data.",
      icon: <CursorArrowRaysIcon />,
    },
  ],
};

const benefitTwo = {
  title: "Few more benefits of using SpotiTubeMerge",
  desc: "Overall, SpotiTubeMerge offers a valuable service that simplifies playlist management, enhances music listening experiences, and prioritizes user privacy and data protection.",
  image: benefitTwoImg,
  bullets: [
    {
      title: "Easy Account Creation",
      desc: "The platform offers a straightforward process for creating an account by connecting Google and Spotify accounts and granting necessary authentication tokens, making it convenient for users to get started.",
      icon: <DevicePhoneMobileIcon />,
    },
    {
      title: "Protected Intellectual Property",
      desc: "SpotiTubeMerge safeguards its content, website, software, graphics, logos, and trademarks, ensuring the protection of its intellectual property and maintaining a high-quality service.",
      icon: <AdjustmentsHorizontalIcon />,
    },
    // {
    //   title: "No Disruptions or Harmful Activities",
    //   desc: "By agreeing to the Terms of Use, users commit to not disrupt or harm the integrity of SpotiTubeMerge or interfere with other users' experiences, ensuring a positive and respectful community. ",
    //   icon: <SunIcon />,
    // },
  ],
};

export { benefitOne, benefitTwo };
