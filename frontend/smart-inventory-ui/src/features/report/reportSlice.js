import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { api } from "../../services/api";

export const fetchReports = createAsyncThunk(
  "report/fetch",
  async () => {
    const res = await api.get("/reports");
    return res.data;
  }
);

const reportSlice = createSlice({
  name: "report",
  initialState: {
    data: [],
  },
  extraReducers: (builder) => {
    builder.addCase(fetchReports.fulfilled, (state, action) => {
      state.data = action.payload;
    });
  },
});

export default reportSlice.reducer;