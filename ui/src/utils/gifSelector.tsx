const gifs = [
    "https://gifdb.com/images/high/comedian-jim-carrey-as-hackerman-vmf9qnz7nx5p9grz.gif",
    "/gifs/1.gif",
    "/gifs/2.gif",
    "/gifs/3.gif",
    "/gifs/4.gif",
  ];

  const getRandomGifUrl = () => {
    const randomIndex = Math.floor(Math.random() * gifs.length);
    return gifs[randomIndex];
  };

    export default getRandomGifUrl;