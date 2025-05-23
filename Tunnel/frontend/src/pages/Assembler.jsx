import { lazy, useEffect, useState } from "react";
import {
  GetRequiredContent,
  CheckIfAllSegmentAreAvaliable,
} from "../../wailsjs/go/main/App.js";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.jsx";
import { Button } from "@/components/ui/button.jsx";
import CompleteGif from "../assets/success.gif";
import DownloadingGif from "../assets/searching.gif";
import toast from "react-hot-toast";

const Contents = lazy(() => import("@/components/Contents"));

function Assembler() {
  const [building, setBuilding] = useState(false);
  const [buildingDone, setBuildingDone] = useState(false);
  const [open, setOpen] = useState(false);
  const [neededFileSegments, setNeededFileSegments] = useState([]);
  let [contents, setContents] = useState([]);
  let[clickedContent, setClickedContent] = useState(null);

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

  const handleGetRequitredContents = async () => {
    try {
      const data = await GetRequiredContent();
      setContents(data);
    } catch (error) {
      console.error("Error fetching required content:", error);
      toast.error("Error fetching required content");
    }
  };

  const handleCheckAllSegmentsAreAvaliable = async (clickedContent) => {
    try {
      const neededSegments = await CheckIfAllSegmentAreAvaliable(
        clickedContent.fileSegments
      );
      setClickedContent(clickedContent);
      console.log("Needed Segments:", neededSegments);
      if (neededSegments === null) {
        setNeededFileSegments([]);
      } else {
        setNeededFileSegments(neededSegments);
      }
      setOpen(true);
    } catch (error) {
      console.error("Error checking segments:", error);
      toast.error("Error while checking segments");
    }
  };

  const handleBuilding = async () => {
    try {
      console.log("Building...");
      setBuilding(true);
      await MakeOriginaleFile(clickedContent?.fileHash);
      setBuildingDone(true);
      setBuilding(false);
      toast.success("Content is ready!");
    } catch (error) {
      console.error("Error building:", error);
      toast.error("Error while building");
    }
  };

  useEffect(() => {
    handleGetRequitredContents();
  }, []);

  return (
    <>
      <div className="relative flex flex-col gap-4 w-full h-full p-4">
        <h1 className="font-semibold text-xl">Contents</h1>
        <div className="absolute top-4 right-4">
          <button
            className="bg-white/40 hover:bg-white/60 text-gray-800 font-semibold py-2 px-4 rounded shadow transition"
            onClick={handleGetRequitredContents}
            type="button"
          >
            &#x21bb; Refresh
          </button>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 w-full z-50">
          {contents.map((content, index) => (
            <div
              key={index}
              onClick={() => {
                console.log("Clicked content:");
                handleCheckAllSegmentsAreAvaliable(content);
              }}
            >
              <Contents item={content} />
            </div>
          ))}
        </div>
      </div>
      {open && (
        <Dialog open={!!open} onOpenChange={setOpen}>
          <DialogContent
            className="bg-white text-black rounded-md shadow-lg border-none"
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
            <DialogTitle>
              <h1>Download all segments</h1>
            </DialogTitle>
            <DialogHeader>
              {neededFileSegments && neededFileSegments.length > 0 ? (
                neededFileSegments.map((neededFileSegment, index) => (
                  <div key={index} className="flex flex-col gap-2 p-4">
                    <p className="font-semibold break-all w-full">
                      File Hash:{" "}
                      <span className="text-blue-700">
                        {neededFileSegment.fileSegmentHash}
                      </span>
                    </p>
                    <p className="font-semibold">
                      File Size:{" "}
                      <span className="text-gray-500">
                        {neededFileSegment.segFileSize / 1024 / 1024} MB
                      </span>
                    </p>
                    <p className="font-semibold">
                      Segment Number:{" "}
                      <span className="text-gray-500">
                        {neededFileSegment.SegmentNumber + 1}
                      </span>
                    </p>
                  </div>
                ))
              ) : (
                <div className="flex flex-col gap-2 py-4">
                  {buildingDone ? (
                    <img src={CompleteGif} alt="Completed" />
                  ) : (
                    <h1 className="font-semibold text-gray-600">
                      All segments are present!
                    </h1>
                  )}
                </div>
              )}
            </DialogHeader>
            <DialogFooter>
              <Button
                className="bg-red-700 text-white cursor-pointer hover:bg-red-800 transition duration-200 ease-in-out shadow-xl/30"
                onClick={() => {
                  setBuildingDone(false);
                  setBuilding(false);
                  setOpen(false);
                }}
              >
                Cancel
              </Button>

              <Button
                onClick={() => handleBuilding()}
                className="bg-green-500 text-white cursor-pointer hover:bg-green-600 transition duration-200 ease-in-out shadow-xl/30 px-4 rounded-md"
                disabled={
                  building || buildingDone || neededFileSegments.length > 0
                }
              >
                {building ? "Building..." : "Build"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}

export default Assembler;
