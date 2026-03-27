import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { api } from "../../services/api";

export const createStockIn = createAsyncThunk(
  "stockIn/create",
  async (items, { rejectWithValue }) => {
    try {
      const res = await api.post("/stock-in", { items });
      return res.data;
    } catch (err) {
      return rejectWithValue(
        err.response?.data || "Failed"
      );
    }
  }
);
export const fetchStockIn = createAsyncThunk(
  "stockIn/fetch",
  async () => {
    const res = await api.get("/stock-in");
    return res.data;
  }
);

export const updateStockInStatus = createAsyncThunk(
  "stockIn/updateStatus",
  async ({ id, status }, { rejectWithValue }) => {
    try {
      await api.put(`/stock-in/${id}/status`, {
        status, 
      });

      return { id, status };
    } catch (err) {
      return rejectWithValue(
        err.response?.data || "Failed"
      );
    }
  }
);


const stockInSlice = createSlice({
  name: "stockIn",
  initialState: {
    currentId: null,
  },
  extraReducers: (builder) => {
    builder.addCase(createStockIn.fulfilled, (state, action) => {
      state.currentId = action.payload.id;
    })
    .addCase(fetchStockIn.fulfilled, (state, action) => {
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

export default stockInSlice.reducer;