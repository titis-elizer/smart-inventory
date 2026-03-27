import { configureStore } from "@reduxjs/toolkit";

import inventoryReducer from "../features/inventory/inventorySlice";
import stockInReducer from "../features/stockIn/stockInSlice";
import stockOutReducer from "../features/stockOut/stockOutSlice";
import reportReducer from "../features/report/reportSlice";

export const store = configureStore({
  reducer: {
    inventory: inventoryReducer,
    stockIn: stockInReducer,
    stockOut: stockOutReducer,
    report: reportReducer,
  },
});