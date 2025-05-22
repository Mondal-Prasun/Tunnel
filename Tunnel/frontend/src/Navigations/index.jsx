import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Suspense, lazy } from "react";
const Loader = lazy(() => import("@/components/Loader"));
const DashBoard = lazy(() => import("@/Layout/Dashboard"));
const Leecher = lazy(() => import("@/pages/Leecher"));
const Onboarding = lazy(() => import("@/pages/Onboarding"));
const Seeder = lazy(() => import("@/pages/Seeder"));
const Assembler = lazy(() => import("@/pages/Assembler"));

const Navigations = () => {
  return (
    <>
      <BrowserRouter basename="/">
        <Suspense fallback={<Loader />} />
        <Routes>
          <Route index path="/" element={<Onboarding />} />
          <Route element={<DashBoard />}>
            <Route path="/leech" element={<Leecher />} />
            <Route path="/seeder" element={<Seeder />} />
            <Route path="/assembler" element={<Assembler />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </>
  );
};

export default Navigations;
