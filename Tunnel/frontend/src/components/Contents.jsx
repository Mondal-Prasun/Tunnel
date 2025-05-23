import { CloudDownload } from "lucide-react";
import ContentImage from "../assets/demo.jpg";
import { Card } from "./ui/card";
import { decode } from "@/utils/common";

function Contents({ item }) {
  return (
    <>
      <Card className="flex flex-col gap-2 items-center w-full h-full bg-white rounded-2xl shadow-md z-50 shadow-xl/30 pt-0">
        <img
          src={item.fileImage}
          alt={item.fileName}
          className="w-full rounded-t-2xl"
        />
        <div className="w-full flex px-4 py-2 flex-col gap-2">
          <h2 className="font-extrabold text-lg text-gray-900 truncate tracking-tight">
            {item.fileName}
          </h2>
          <div className="flex items-center justify-between text-sm text-gray-600">
            <span className="flex items-center gap-1">
              <CloudDownload className="w-4 h-4 text-blue-500" />
              <span className="font-medium">
                {(item.fileSize / 1024 / 1024).toFixed(2)} MB
              </span>
            </span>
            <span className="bg-blue-100 text-blue-700 px-2 py-0.5 rounded text-xs font-semibold">
              {item.fileExt?.toUpperCase() || ""}
            </span>
          </div>
        </div>
      </Card>

    </>
  );
}

export default Contents;
