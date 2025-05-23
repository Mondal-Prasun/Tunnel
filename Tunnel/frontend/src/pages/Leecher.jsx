import Contents from "@/components/Contents";
import BackgroundImage from "../assets/background.png";
import DownloadingGif from "../assets/searching.gif";
import CompleteGif from "../assets/success.gif";
import { useEffect, useState } from "react";
import {
  FetchTrackerFile,
  ListenToPeers,
  GetRequiredContent,
} from "../../wailsjs/go/main/App.js";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import { set } from "react-hook-form";

function Leecher() {
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [downloaded, setDownloaded] = useState(false);
  const url = localStorage.getItem("url");
  let [contents, setContents] = useState([]);
  let [clickedContent, setClickedContent] = useState(null);
  // contents = [
  //   {
  //     fileHash:
  //       "a4406c90d99c8c3be19fa8cf6c5d34c7bd462d61edb98a00f01324cafa21c052",
  //     fileName: "blossom.mp4",
  //     fileSize: 62781236,
  //     allSegmentCount: 6,
  //     fileExt: ".mp4",
  //     fileSegments: [
  //       {
  //         fileSegmentHash:
  //           "5ee442944bd3fbb9fe06d494c0585520c666b51da0e0b95e9a857ae6fb54327a",
  //         segContentSize: 10464256,
  //         segFileSize: 10464508,
  //         SegmentNumber: 0,
  //         segAddress: ["127.0.0.1:6000"],
  //       },
  //       {
  //         fileSegmentHash:
  //           "21f4552784245410d946d7fed391bf7e68df9cafbda210fbd03e505b5ffb8a3b",
  //         segContentSize: 10464256,
  //         segFileSize: 10464508,
  //         SegmentNumber: 1,
  //         segAddress: ["127.0.0.1:6000"],
  //       },
  //       {
  //         fileSegmentHash:
  //           "118ebfb1bfe904eb6d98fa71ba47f300894ea9e73dac724005c8818e6612b192",
  //         segContentSize: 10464256,
  //         segFileSize: 10464508,
  //         SegmentNumber: 2,
  //         segAddress: ["127.0.0.1:6000"],
  //       },
  //       {
  //         fileSegmentHash:
  //           "082657d9ee96d7fa232a7867477ebb8e20af4eaef938e5d78f8bcbd9e2c277f8",
  //         segContentSize: 10464256,
  //         segFileSize: 10464508,
  //         SegmentNumber: 3,
  //         segAddress: ["127.0.0.1:6000"],
  //       },
  //       {
  //         fileSegmentHash:
  //           "f7dfea510b7a18aed71835067027967c32e7611ea592fc8b4ab6a39c8ac5e9c8",
  //         segContentSize: 10464256,
  //         segFileSize: 10464508,
  //         SegmentNumber: 4,
  //         segAddress: ["127.0.0.1:6000"],
  //       },
  //       {
  //         fileSegmentHash:
  //           "3701cd37e60da0b11853b04c33fc135b67f6e789c5dca4d9598fb86c09e0f077",
  //         segContentSize: 10459956,
  //         segFileSize: 10464508,
  //         SegmentNumber: 5,
  //         segAddress: ["127.0.0.1:6000"],
  //       },
  //     ],
  //   },
  // ];
  const handleCall = async () => {
    try {
      const data = await GetRequiredContent();
      await FetchTrackerFile(url);
      console.log(data);
      setContents(data);
      toast.success(`Successfully updated contents: ${new Date().toLocaleTimeString()}`);
    } catch (e) {
      console.log(e);
    }
  };
  useEffect(() => {
    handleCall();
    const interval = setInterval(() => {
      handleCall();
    }, 5000);
    return () => clearInterval(interval);
  }, []);
  
  const handleDownload = async() => {
    try {
      console.log("Download clicked");
      // Handle the download logic here
      setDownloading(true);
      await RequestRequiredSegments(clickedContent.fileHash);
      setDownloaded(true);
      setDownloading(false);
      toast.success("Content successfully downloaded!")
      // Simulate a 2-second download time
    } catch (error) {
      console.error("Error downloading content:", error);
      toast.error("Error downloading content. Please try again.");
    }
  };
  return (
    <>
      <section className="flex flex-col gap-4 w-full h-full p-4 relative">
        <div className="absolute top-4 right-4 z-10">
          <button
            className="bg-white/40 hover:bg-white/60 text-gray-800 font-semibold py-2 px-4 rounded shadow transition"
            onClick={handleCall}
            type="button"
          >
            &#x21bb; Refresh
          </button>
        </div>
        <section className="w-full md:w-[70%] lg:w-[50%] mx-auto md:h-[400px] rounded-lg md:block relative">
          <img
            className="h-full w-full rounded-2xl"
            src={BackgroundImage}
            alt="background-image"
          />
          <div
            className="absolute inset-0 bg-gradient-to-t from-white via-transparent to-transparent rounded-2xl"
            style={{ top: "-40px", height: "calc(100% + 40px)" }}
          ></div>
        </section>
        <h1 className="font-semibold text-xl">Contents</h1>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 w-full">
          {contents.map((content, index) => (
            <div
              key={index}
              onClick={() => {
                setClickedContent(content);
                setOpen(true);
              }}
            >
              <Contents item={content} />
            </div>
          ))}
        </div>
      </section>
      {open && (
        <Dialog open={!!open} onOpenChange={setOpen}>
          <DialogTrigger>Open</DialogTrigger>
          <DialogContent
            className="bg-white p-4 text-black rounded-md shadow-lg border-none"
            style={{
              width: "90vw",
              maxWidth: "600px",
              overflowY: "auto",
              overflowX: "hidden",
              maxHeight: "80vh",
              msOverflowStyle: "none",
              scrollbarWidth: "none",
            }}
          >
            {downloading ? (
              <>
                <DialogTitle className="flex flex-col items-center justify-center">
                  <img
                    src={downloaded ? CompleteGif : DownloadingGif}
                    alt="downloading..."
                  />
                  <p className="font-bold text-2xl text-gray-500">
                    {downloaded ? "Download Completed!" : "Downloading..."}
                  </p>
                </DialogTitle>
                <DialogFooter>
                  <Button
                    className={`${
                      downloaded ? "bg-green-500" : "bg-red-700"
                    } text-white cursor-pointer ${
                      downloaded ? "hover:bg-green-600" : "hover:bg-red-800"
                    } transition duration-200 ease-in-out shadow-xl/30`}
                    onClick={() => {
                      setOpen(false);
                      setDownloading(false);
                      setDownloaded(false);
                      navigate("/leech");
                    }}
                  >
                    {downloaded ? "Go to Dashboard" : "Close"}
                  </Button>
                </DialogFooter>
              </>
            ) : (
              <>
                <DialogTitle>
                  <p className="font-bold">{clickedContent.fileDescription}</p>
                </DialogTitle>
                <DialogHeader>
                  {clickedContent.fileSegments.length > 0 ? (
                    clickedContent.fileSegments.map((fileSegment, index) => (
                      <div key={index} className="flex flex-col gap-2 p-4">
                        <p className="font-semibold break-all w-full">
                          File Hash:{" "}
                          <span className="text-blue-700">
                            {" "}
                            {fileSegment.fileSegmentHash}{" "}
                          </span>
                        </p>
                        <p className="font-semibold">
                          File Size:{" "}
                          <span className="text-gray-500">
                            {parseInt((fileSegment.segContentSize/1024)/1024)} MB
                          </span>
                        </p>
                        <p className="font-semibold">
                          Segment Number:{" "}
                          <span className="text-gray-500">
                            {fileSegment.SegmentNumber}
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
                  {clickedContent.fileSegments.length > 0 && (
                    <Button
                      onClick={() => {
                        handleDownload();
                      }}
                      className="bg-green-500 text-white cursor-pointer hover:bg-green-600 transition duration-200 ease-in-out shadow-xl/30"
                    >
                      Download
                    </Button>
                  )}
                </DialogFooter>
              </>
            )}
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}

export default Leecher;
