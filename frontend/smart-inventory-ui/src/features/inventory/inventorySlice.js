import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { api } from "../../services/api";

// 🔍 FETCH INVENTORY
export const fetchInventory = createAsyncThunk(
  "inventory/fetch",
  async (search = "") => {
    const res = await api.get(`/inventory?search=${search}`);
    return res.data;
  }
);

// ✏️ ADJUST STOCK
export const adjustStock = createAsyncThunk(
  "inventory/adjust",
  async ({ item_id, qty }) => {
    await api.post("/inventory/adjust", {
      item_id,
      qty,
    });

    // return manual (karena backend cuma return "success")
    return { item_id, qty };
  }
);
const inventorySlice = createSlice({
  name: "inventory",
  initialState: {
    items: [],
    loading: false,
  },
  extraReducers: (builder) => {
    builder
      // 🔄 FETCH
      .addCase(fetchInventory.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchInventory.fulfilled, (state, action) => {
        // ✅ Mapping dari PascalCase → camelCase
        state.items = action.payload.map((item) => ({
          id: item.ID,
          name: item.Name,
          sku: item.SKU,
          physical_stock: item.PhysicalStock,
          reserved_stock: item.ReservedStock,
        }));
        state.loading = false;
      })

      // ✏️ ADJUST STOCK (optimistic update)
     .addCase(adjustStock.fulfilled, (state, action) => {
  const { item_id, qty } = action.payload;

  const index = state.items.findIndex(
    (item) => item.id === item_id
  );

  if (index !== -1) {
    state.items[index].physical_stock += qty;
  }
});
  },
});

export default inventorySlice.reducer;