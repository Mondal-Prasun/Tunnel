import useNavigation from "@/hooks/useNavigation";
import { Card } from "./ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { Link } from "react-router-dom";
import { Button } from "./ui/button";

function MobileNav() {
  const paths = useNavigation();
  return (
    <>
      <Card className="fixed bg-white bottom-4 w-[calc(100vw-16px)] flex lg:hidden justify-center items-center p-2 z-50">
        <nav className="w-full">
          <ul className="flex flex-row justify-between items-center">
            {paths.map((path, id) => {
              return (
                <li key={id} className="relative">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Link to={path.href}>
                        <Button
                          size="icon"
                          // variant={path.active ? "default" : "outline"}
                          className={`${
                            path.active
                              ? "bg-blue-500 text-white"
                              : "bg-transparent"
                          }`}
                        >
                          {path.icon}
                        </Button>
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

export default MobileNav;
