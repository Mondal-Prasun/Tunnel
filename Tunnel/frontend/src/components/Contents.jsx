import { Link } from "react-router-dom";
import ContentImage from "../assets/demo.jpg";
import DownloadingGif from "../assets/searching.gif";
import CompleteGif from "../assets/complete.gif";

function Contents({ item }) {
  const [open, setOpen] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [downloaded, setDownloaded] = useState(false);
  const handleDownload = () => {
    console.log("Download clicked");
    // Handle the download logic here
    setTimeout(() => {
      setDownloaded(true);
    }, 5000); // Simulate a 2-second download time
  };

  return (
    <>
      <div
        onClick={() => setOpen(true)}
        className="flex flex-col gap-2 items-center w-full h-full bg-white rounded shadow-md"
      >
        <img
          src={ContentImage}
          alt={item.fileName}
          className="w-full rounded-t"
        />
        <div className="w-full px-4 pb-4">
          <h2 className="font-bold">{item.fileDescription}</h2>
          {/* <p>{item.description}</p> */}
        </div>
      </div>
      {open && (
        <Dialog open={!!open} onOpenChange={setOpen}>
          <DialogTrigger>Open</DialogTrigger>
          <DialogContent className="bg-white text-black rounded-md shadow-lg border-none">
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
                      window.location.href = "/leech";
                    }}
                  >
                    {downloaded ? "Go to Dashboard" : "Close"}
                  </Button>
                </DialogFooter>
              </>
            ) : (
              <>
                <DialogTitle>
                  <p className="font-bold">{item.fileDescription}</p>
                </DialogTitle>
                <DialogHeader>
                  {item.fileSegments.length > 0 ? (
                    item.fileSegments.map((fileSegment, index) => (
                      <div key={index} className="flex flex-col gap-2 p-4">
                        <p className="font-semibold ">
                          File Hash:{" "}
                          <span className="text-blue-700">
                            {" "}
                            {fileSegment.fileSegmentHash}{" "}
                          </span>
                        </p>
                        <p className="font-semibold">
                          File Size:{" "}
                          <span className="text-gray-500">
                            {fileSegment.segContentSize}
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
                  <Button
                    onClick={() => {
                      setDownloading(true);
                      handleDownload();
                    }}
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

export default Contents;
