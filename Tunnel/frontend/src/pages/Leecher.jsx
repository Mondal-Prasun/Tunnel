import Contents from "@/components/Contents";
import BackgroundImage from "../assets/leech-background.jpg";
import { useEffect } from "react";

function Leecher() {
  const contents = [
    {
      id: 1,
      thumbnail: "https://via.placeholder.com/150/0000FF/808080?text=Content+1",
      title: "Content 1",
      description: "Description for content 1",
    },
    {
      id: 2,
      thumbnail: "https://via.placeholder.com/150/FF0000/FFFFFF?text=Content+2",
      title: "Content 2",
      description: "Description for content 2",
    },
    {
      id: 3,
      thumbnail: "https://via.placeholder.com/150/00FF00/000000?text=Content+3",
      title: "Content 3",
      description: "Description for content 3",
    },
    {
      id: 4,
      thumbnail: "https://via.placeholder.com/150/FFFF00/000000?text=Content+4",
      title: "Content 4",
      description: "Description for content 4",
    },
    {
      id: 5,
      thumbnail: "https://via.placeholder.com/150/FF00FF/FFFFFF?text=Content+5",
      title: "Content 5",
      description: "Description for content 5",
    },
    {
      id: 6,
      thumbnail: "https://via.placeholder.com/150/00FFFF/000000?text=Content+6",
      title: "Content 6",
      description: "Description for content 6",
    },
    {
      id: 7,
      thumbnail: "https://via.placeholder.com/150/800000/FFFFFF?text=Content+7",
      title: "Content 7",
      description: "Description for content 7",
    },
    {
      id: 8,
      thumbnail: "https://via.placeholder.com/150/808000/FFFFFF?text=Content+8",
      title: "Content 8",
      description: "Description for content 8",
    },
    {
      id: 9,
      thumbnail: "https://via.placeholder.com/150/008080/FFFFFF?text=Content+9",
      title: "Content 9",
      description: "Description for content 9",
    },
    {
      id: 10,
      thumbnail:
        "https://via.placeholder.com/150/800080/FFFFFF?text=Content+10",
      title: "Content 10",
      description: "Description for content 10",
    },
    {
      id: 11,
      thumbnail:
        "https://via.placeholder.com/150/000080/FFFFFF?text=Content+11",
      title: "Content 11",
      description: "Description for content 11",
    },
    {
      id: 12,
      thumbnail:
        "https://via.placeholder.com/150/808080/FFFFFF?text=Content+12",
      title: "Content 12",
      description: "Description for content 12",
    },
    {
      id: 13,
      thumbnail:
        "https://via.placeholder.com/150/FFA500/FFFFFF?text=Content+13",
      title: "Content 13",
      description: "Description for content 13",
    },
    {
      id: 14,
      thumbnail:
        "https://via.placeholder.com/150/FFC0CB/FFFFFF?text=Content+14",
      title: "Content 14",
      description: "Description for content 14",
    },
    {
      id: 15,
      thumbnail:
        "https://via.placeholder.com/150/ADD8E6/FFFFFF?text=Content+15",
      title: "Content 15",
      description: "Description for content 15",
    },
    {
      id: 16,
      thumbnail:
        "https://via.placeholder.com/150/90EE90/FFFFFF?text=Content+16",
      title: "Content 16",
      description: "Description for content 16",
    },
  ];
  useEffect(()=> {
    setTimeout(()=> {
      Promise.resolve(getTrackerFile()).catch((err)=> console.log("Error has occured", err))
    }, 10000)
  },[])
  return (
    <section className="flex flex-col gap-4 w-full h-full p-4">
      <section className="w-full mx-auto h-[250px] rounded-lg hidden md:block relative">
        <img
          className="h-full w-full rounded-2xl"
          src={BackgroundImage}
          alt="background-image"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-white via-transparent to-transparent rounded-2xl"></div>
      </section>
      <h1 className="font-semibold text-xl">Contents</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 w-full">
        {contents.map((content, index) => (
          <Contents key={index} item={content} />
        ))}
      </div>
    </section>
  );
}

export default Leecher;
