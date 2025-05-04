import {useMemo} from "react";
import { useLocation } from "react-router-dom";
import { CloudDownload, CloudUpload } from "lucide-react";

function useNavigation(){
    const location = useLocation();
    const paths = useMemo(()=> [{
        name: "Leech",
        href: "/leech",
        icon: <CloudDownload />,
        active: location.pathname.startsWith("/leech"), // Check if the path is active
      },
      {
        name: "Seeder",
        href: "/seeder",
        icon: <CloudUpload />,
        active: location.pathname.startsWith("/seeder"), // Check if the path is active
      },], [location.pathname]);

      return paths;
}

export default useNavigation;