import { Toaster } from "react-hot-toast";
import Navigations from "./Navigations";

function App() {
  return (
    <>
      <Toaster position="top-center" reverseOrder={false} />
      <Navigations />
    </>
  );
}

export default App;
