import { lazy } from "react";
import { Outlet } from "react-router-dom";
const MobileNav = lazy(() => import("@/components/MobileNav"));
const DesktopNav = lazy(() => import("@/components/DesktopNav"));

function DashBoard() {
  return (
    <>
      <div className="h-screen flex w-full gap-4">
        <DesktopNav />
        <div className="flex-1 overflow-auto no-scrollbar">
          <main className="h-full w-full no-scrollbar">
            <Outlet />
          </main>
        </div>
        <MobileNav />
      </div>
    </>
  );
}

export default DashBoard;
