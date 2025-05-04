import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Suspense, lazy } from "react";
const Loader = lazy(() => import("@/components/Loader"));
const DashBoard = lazy(() => import("@/Layout/Dashboard"));
const Leecher = lazy(() => import("@/pages/Leecher"));
const Onboarding = lazy(() => import("@/pages/Onboarding"));
const Seeder = lazy(() => import("@/pages/Seeder"));
const ContentDetails = lazy(() => import("@/pages/ContentDetails"));

const Navigations = () => {
  return (
    <>
      <BrowserRouter basename="/">
        <Suspense fallback={<Loader />} />
        <Routes>
          <Route index path="/" element={<Onboarding />} />
          <Route element={<DashBoard />}>
            <Route path="/leech" element={<Leecher />} />
            <Route path="/leech/:leechId?" element={<ContentDetails />} />
            <Route path="/seeder" element={<Seeder />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </>
  );
};

export default Navigations;
