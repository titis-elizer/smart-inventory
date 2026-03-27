import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Toaster } from "react-hot-toast";

import MainLayout from "./layouts/MainLayout";
import InventoryPage from "./pages/InventoryPage";
import StockInPage from "./pages/StockInPage";
import StockOutPage from "./pages/StockOutPage";
import ReportPage from "./pages/ReportPage";

export default function App() {
  return (
    <BrowserRouter>
      <Toaster position="top-right" />

      <MainLayout>
        <Routes>
          <Route path="/inventory" element={<InventoryPage />} />
          <Route path="/stock-in" element={<StockInPage />} />
          <Route path="/stock-out" element={<StockOutPage />} />
          <Route path="/report" element={<ReportPage />} />
        </Routes>
      </MainLayout>
    </BrowserRouter>
  );
}