import { Button } from "@/components/ui/button";
import { FolderSearch2 } from "lucide-react";
import ContentImage from "../assets/demo.jpg";
import { useParams } from "react-router-dom";
import { useState } from "react";
// import modal
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
//downloading content ...
import DownloadingGif from "../assets/searching.gif";
import CompleteGif from "../assets/success.gif";

function ContentDetails() {
  const [open, setOpen] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [downloaded, setDownloaded] = useState(false);
  const contentId = useParams();
  // need to call the API to get the content details using this contentId
  const content = {
    thumbnail: "https://www.w3schools.com/html/mov_bbb.mp4",
    title: "Content 1",
    description: "Description for content 1",
  };

  //demo json of file data
  const fileSegments = [
    {
      fileName: "file1.txt",
      location: "Location 1",
      fileHash: "1234567890abcdef",
    },
    {
      fileName: "file2.txt",
      location: "Location 2",
      fileHash: "abcdef1234567890",
    },
    {
      fileName: "file3.txt",
      location: "Location 3",
      fileHash: "fedcba0987654321",
    },
    {
      fileName: "file4.txt",
      location: "Location 4",
      fileHash: "0123456789abcdef",
    },
  ];


  const handleDownload = ()=> {
    console.log("Download clicked");
    // Handle the download logic here
    setTimeout(() => {
      setDownloaded(true);
    }, 5000); // Simulate a 2-second download time
  }
  return (
    <>
      <div className="flex flex-col lg:flex-row w-full h-full lg:w-[75%] lg:mx-auto p-2 gap-4 lg:justify-between lg:items-center ">
        <div className="flex flex-col lg:flex-row gap-4 w-full lg:mx-auto lg:justify-between ">
          <div className="flex flex-col w-full lg:w-[50%]">
            <img src={ContentImage} alt={content.title} className="rounded" />
          </div>
          <div className="flex flex-col gap-4 w-full lg:w-[50%]">
            <div>
              <h1 className="font-bold">Title</h1>
              <p className="font-semibold">Size: </p>
              <p>Parent file hash: </p>
              <p>
                Lorem ipsum dolor sit amet consectetur adipisicing elit.
                Quisquam, voluptatibus.
              </p>
            </div>
            <Button
              className="bg-gray-800 text-white cursor-pointer hover:bg-gray-900 transition duration-200 ease-in-out shadow-xl/30"
              onClick={() => setOpen(true)}
            >
              <FolderSearch2 />
              Request
            </Button>
          </div>
        </div>
      </div>
      {open && (
        <Dialog open={!!open} onOpenChange={setOpen}>
          <DialogTrigger>Open</DialogTrigger>
          <DialogContent className="bg-white text-black rounded-md shadow-lg border-none">
            {downloading ? (
              <>
                <DialogTitle className="flex flex-col items-center justify-center">
                  <img src={downloaded? CompleteGif :DownloadingGif} alt="downloading..." />
                  <p className="font-bold text-2xl text-gray-500">
                    {downloaded? "Download Completed!" : "Downloading..."}
                  </p>
                </DialogTitle>
                <DialogFooter>
                  <Button
                    className={`${downloaded? "bg-green-500": "bg-red-700"} text-white cursor-pointer ${downloaded?"hover:bg-green-600":"hover:bg-red-800"} transition duration-200 ease-in-out shadow-xl/30`}
                    onClick={() => {
                      setOpen(false); 
                      setDownloading(false); 
                      setDownloaded(false); 
                      window.location.href = "/leech";
                    }}
                  >
                    {downloaded? "Go to Dashboard": "Close"}
                  </Button>
                </DialogFooter>
              </>
            ) : (
              <>
                <DialogTitle>
                  <p className="font-bold">{content.title}</p>
                </DialogTitle>
                <DialogHeader>
                  {fileSegments.length > 0 ? (
                    fileSegments.map((fileSegment, index) => (
                      <div key={index} className="flex flex-col gap-2 p-4">
                        <p className="font-semibold ">
                          File Hash:{" "}
                          <span className="text-blue-700">
                            {" "}
                            {fileSegment.fileHash}{" "}
                          </span>
                        </p>
                        <p className="font-semibold">
                          Address:{" "}
                          <span className="text-gray-500">
                            {fileSegment.location}
                          </span>
                        </p>
                      </div>
                    ))
                  ) : (
                    <div className="flex flex-col gap-2 p-4">
                      <h1 className="font-bold">No file segments found</h1>
                    </div>
                  )}
                </DialogHeader>
                <DialogFooter>
                  <Button
                    className="bg-red-700 text-white cursor-pointer hover:bg-red-800 transition duration-200 ease-in-out shadow-xl/30"
                    onClick={() => setOpen(false)}
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={() => {setDownloading(true); handleDownload()}}
                    className="bg-green-500 text-white cursor-pointer hover:bg-green-600 transition duration-200 ease-in-out shadow-xl/30"
                  >
                    Download
                  </Button>
                </DialogFooter>
              </>
            )}
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}

export default ContentDetails;
