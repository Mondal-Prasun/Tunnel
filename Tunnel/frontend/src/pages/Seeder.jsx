import { CloudUpload } from "lucide-react";
import BackgroundImage from "../assets/upload.jpg";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import { AnnounceCurrentFile } from "../../wailsjs/go/main/App.js";
import { useNavigate } from "react-router-dom";

const schema = yup.object({
  thumbnail: yup
    .mixed()
    .required("Thumbnail is required") // Ensure this is checked first
    .test("fileType", "Unsupported file format", (value) => {
      return (
        value &&
        value[0] &&
        ["image/jpeg", "image/png", "image/jpg"].includes(value[0].type)
      );
    }),
  title: yup.string().required("Title is required"),
  filePath: yup
    .string()
    .required("File path is required")
    .test("validPath", "Invalid file path format", (value) => {
      // Windows: C:\folder\file.txt or Unix: /home/user/file.txt
      return /^([a-zA-Z]:\\|\/).+/.test(value);
    }),
});

function Seeder() {
  const url = localStorage.getItem("url");
  const navigate = useNavigate();
  const [thumbnail, setThumbnail] = useState(null);
  const { handleSubmit, register, formState } = useForm({
    resolver: yupResolver(schema),
  });
  const handleUpload = async (data) => {
    try {
      if (!thumbnail) {
        alert("Please select a thumbnail");
        return;
      }
      const getBase64 = (file) => {
        return new Promise((resolve, reject) => {
          const reader = new FileReader();
          reader.readAsDataURL(file);
          reader.onload = () => resolve(reader.result);
          reader.onerror = (error) => reject(error);
        });
      };
      const base64Thumbnail = await getBase64(thumbnail);
      console.log("Base64 Thumbnail:", base64Thumbnail);
      const payload = {
        thumbnail: base64Thumbnail,
        title: data.title,
        filePath: data.filePath,
      };
      console.log(payload.filePath);
      console.log(payload.title);
      console.log(payload.base64Thumbnail);
      console.log("url: ",url);

      await AnnounceCurrentFile(
        data.filePath,
        base64Thumbnail,
        data.title,
        url
      );
      navigate("/leech");
    } catch (error) {
      console.error("Error uploading content:", error);
      alert("Error uploading content. Please try again.");
    }
    // Handle the upload logic here
  };
  return (
    <>
      <div className="flex h-screen">
        <section className="w-full lg:w-2/3 flex items-center justify-center">
          <div className="w-full max-w-[496px] p-4 z-50">
            <div className="my-4">
              <h4 className="text-4xl font-semibold">Seeder</h4>
              <p className="text-gray-500 mt-3">Upload your content here...</p>
            </div>
            <form
              onSubmit={handleSubmit(handleUpload)}
              className="flex flex-col gap-2"
            >
              <div className="flex flex-col gap-2">
                <label htmlFor="url" className="text-gray-600">
                  Thumbnail
                </label>
                <input
                  type="file"
                  accept="image/jpeg,image/png,image/jpg" // <-- restrict file types in the file picker
                  className="border border-dashed border-gray-500 rounded-md py-2 px-3 "
                  {...register("thumbnail")}
                  onChange={(e) => setThumbnail(e.target.files[0])}
                />
                {formState.errors.thumbnail && (
                  <p className="text-red-500 text-sm">
                    {formState.errors.thumbnail.message}
                  </p>
                )}
                {thumbnail && (
                  <img
                    src={URL.createObjectURL(thumbnail)}
                    alt="thumbnail"
                    className="w-full h-50 object-cover rounded-md mt-2"
                  />
                )}
              </div>
              <div className="flex relative flex-col gap-2">
                <label htmlFor="password" className="text-gray-600">
                  Title
                </label>
                <input
                  type="text"
                  placeholder="Enter your title here..."
                  className="border border-gray-500 rounded-md py-2 px-3 placeholder:text-gray-900"
                  defaultValue=""
                  {...register("title")}
                />
                {formState.errors.title && (
                  <p className="text-sm text-red-500">
                    {formState.errors.title.message}
                  </p>
                )}
              </div>
              <div className="flex relative flex-col gap-2">
                <label htmlFor="filePath" className="text-gray-600">
                  File Path
                </label>
                <input
                  type="text"
                  placeholder="Enter file path (e.g. C:\\Users\\... or /home/user/...)"
                  className="border border-gray-500 rounded-md py-2 px-3"
                  {...register("filePath")}
                />
                {formState.errors.filePath && (
                  <p className="text-sm text-red-500">
                    {formState.errors.filePath.message}
                  </p>
                )}
              </div>
              <Button
                type="submit"
                className="bg-gray-800 text-white cursor-pointer hover:bg-gray-900 transition duration-200 ease-in-out shadow-xl/30 my-4"
              >
                <CloudUpload />
                Upload
              </Button>
            </form>
          </div>
        </section>
        <section className="hidden lg:flex w-1/3 h-full">
          <img
            src={BackgroundImage}
            alt="seeder"
            className="object-cover w-full h-full"
          />
        </section>
      </div>
    </>
  );
}

export default Seeder;
