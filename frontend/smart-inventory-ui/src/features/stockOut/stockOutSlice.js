import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { api } from "../../services/api";

export const createStockOut = createAsyncThunk(
  "stockOut/create",
  async (items) => {
    const res = await api.post("/stock-out", { items });
    return res.data;
  }
);

export const markInProgress = createAsyncThunk(
  "stockOut/inProgress",
  async (id) => {
    await api.post(`/stock-out/${id}/in-progress`);
    return id;
  }
);

export const completeStockOut = createAsyncThunk(
  "stockOut/complete",
  async (id) => {
    await api.post(`/stock-out/${id}/complete`);
    return id;
  }
);

export const cancelStockOut = createAsyncThunk(
  "stockOut/cancel",
  async (id) => {
    await api.post(`/stock-out/${id}/cancel`);
    return id;
  }
);
export const fetchStockOut = createAsyncThunk(
  "stockOut/fetch",
  async () => {
    const res = await api.get("/stock-out");
    return res.data;
  }
);

const stockOutSlice = createSlice({
  name: "stockOut",
  initialState: {
    currentId: null,
    loading: false,
  },
  extraReducers: (builder) => {
    builder
      .addCase(createStockOut.pending, (state) => {
        state.loading = true;
      })
      .addCase(createStockOut.fulfilled, (state, action) => {
        state.currentId = action.payload.id;
        state.loading = false;
      })
  .addCase(fetchStockOut.fulfilled, (state, action) => {
  state.orders = action.payload.map((o) => ({
    id: o.ID,
    status: o.Status,
      created_at: o.CreatedAt,
      done_at: o.DoneAt,
    items: o.Items.map((i) => ({
      inventory_item_id: i.InventoryItemID,
      product_name: i.ProductName,
      qty: i.Qty,
    
    })),
  }));
});
      
      
  },
});

export default stockOutSlice.reducer;