import BackgroundImage from "../assets/onboarding.png";
import { useState } from "react";
import { Button } from "../components/ui/button";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  FetchTrackerFile,
  ListenToPeers,
  MakeRequiredFile,
} from "../../wailsjs/go/main/App.js";
import { useNavigate } from "react-router-dom";

const schema = yup.object({
  url: yup.string().url("Invalid URL").required("URL is required"),
  port: yup
    .number()
    .required("Port is required")
    .typeError("Port must be a number")
    .min(6000, "Port must be greater than or equal to 6,000")
    .max(10000, "Port must be less than or equal to 10,000"),
});

function Onboarding() {
  const navigate = useNavigate();
  const [screen, setScreen] = useState("Home");
  const { handleSubmit, register, formState } = useForm({
    resolver: yupResolver(schema),
  });
  const handleLogin = async (data) => {
    try {
      localStorage.setItem("url", data.url);
      console.log("Login data:", data.port);
      await FetchTrackerFile(data.url);
      await ListenToPeers(data.port.toString());
      await MakeRequiredFile();
      console.log("after ", data);
      if (true) {
        navigate("/leech");
      }
    } catch (error) {
      console.error("Error fetching tracker content:", error);
    }
  };
  return (
    <>
      <div className="flex h-screen">
        <section className="w-full lg:w-2/3 flex items-center justify-center">
          <div className="flex items-center justify-center w-full h-full flex-col px-4">
            {screen === "Home" && (
              <>
                <h2 className="bg-clip-text text-transparent text-center bg-gradient-to-b from-neutral-900 to-neutral-700 dark:from-neutral-600 dark:to-white text-2xl md:text-4xl lg:text-7xl font-sans py-2 md:py-10 relative z-20 font-bold tracking-tight">
                  Share anything with our <span className="bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500 bg-clip-text text-transparent">Tunnel!</span>
                </h2>
                <p className="max-w-xl mx-auto text-sm md:text-lg text-neutral-700 dark:text-neutral-400 text-center">
                  A TCP and LAN based Peer to Peer file sharing application that
                  allows you to share files with your friends and family.
                </p>
                <div className="mt-4 z-50">
                  <Button
                    onClick={() => setScreen("Login")}
                    className="bg-gray-800 text-white cursor-pointer hover:bg-gray-900 transition duration-200 ease-in-out shadow-xl/30"
                  >
                    Get Started!
                  </Button>
                </div>
              </>
            )}

            {screen === "Login" && (
              <>
                <div className="w-full max-w-[496px] z-50">
                  <div className="font-semibold text-gray-400 text-2xl">
                    Tunnel Application
                  </div>
                  <div className="my-4">
                    <h4 className="text-4xl font-semibold">Connect 👋</h4>
                    <p className="text-gray-500 mt-3">
                      Connect to your prefered <span className="font-bold">tracker...</span> 
                    </p>
                  </div>
                  <form
                    onSubmit={handleSubmit(handleLogin)}
                    className="flex flex-col gap-2"
                  >
                    <div className="flex flex-col gap-4">
                      <label htmlFor="url" className="text-gray-600">
                        Tracker URL
                      </label>
                      <input
                        type="url"
                        className="border border-gray-500 rounded-md py-2 px-3"
                        placeholder="Enter your URL here..."
                        defaultValue=""
                        {...register("url")}
                      />
                      {formState.errors.url && (
                        <p className="text-sm text-red-500">
                          {formState.errors.url.message}
                        </p>
                      )}
                    </div>
                    <div className="flex relative flex-col gap-2">
                      <label htmlFor="password" className="text-gray-600">
                        Port
                      </label>
                      <input
                        className="border border-gray-500 rounded-md py-2 px-3"
                        placeholder="Enter your port number here..."
                        type="number"
                        defaultValue=""
                        {...register("port")}
                      />
                      {formState.errors.port && (
                        <p className="text-sm text-red-500">
                          {formState.errors.port.message}
                        </p>
                      )}
                    </div>
                    <Button
                      type="submit"
                      onClick={handleLogin}
                      className="bg-gray-800 text-white cursor-pointer my-4 hover:bg-gray-900 transition duration-200 ease-in-out shadow-xl/30"
                    >
                      Connect
                    </Button>
                  </form>
                </div>
              </>
            )}
          </div>
        </section>
        <section className="lg:flex w-full -scale-x-100 lg:w-1/3  h-full right-0">
          <img
            src={BackgroundImage}
            alt="Onboarding image"
            className="object-cover w-full h-full"
          />
        </section>
      </div>
    </>
  );
}

export default Onboarding;
