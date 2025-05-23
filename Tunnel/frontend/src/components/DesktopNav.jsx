import useNavigation from "@/hooks/useNavigation";
import { Card } from "./ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { Link } from "react-router-dom";
import { Button } from "./ui/button";

function DesktopNav() {
  const paths = useNavigation();
  return (
    <>
      <Card className="hidden lg:flex lg:flex-col lg:justify-between lg:items-center lg:h-screen lg:w-16 lg:px-2 lg:py-4 bg-white z-50">
        <nav>
          <ul className="flex lg:flex-col items-center gap-10">
            {paths.map((path, id) => {
              return (
                <li key={id} className="relative">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Link to={path.href} className="flex flex-col items-center gap-2">
                        <Button
                          size="icon"
                          // variant={path.active ? "outline" : "default"}
                          className={`${path.active ? "bg-blue-500 text-white": "bg-transparent"}`}
                        >
                          {path.icon}
                        </Button>
                        <span className={`${path.active?"text-[15px] font-semibold": "text-[15px]"}`}>{path.name}</span>
                      </Link>
                    </TooltipTrigger>
                    <TooltipContent>{path.name}</TooltipContent>
                  </Tooltip>
                </li>
              );
            })}
          </ul>
        </nav>
      </Card>
    </>
  );
}

export default DesktopNav;
