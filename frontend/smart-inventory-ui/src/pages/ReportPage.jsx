import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { fetchStockOut } from "../features/stockOut/stockOutSlice";
import { fetchStockIn } from "../features/stockIn/stockInSlice";

export default function ReportPage() {
  const dispatch = useDispatch();

  const stockOut = useSelector((s) => s.stockOut.orders);
  const stockIn = useSelector((s) => s.stockIn.orders);

  useEffect(() => {
    dispatch(fetchStockOut());
    dispatch(fetchStockIn());
  }, []);

  // ✅ filter done only
  const doneOut = (stockOut || []).filter((o) => o.status === "done");
  const doneIn = (stockIn || []).filter((o) => o.status === "done");

  const formatDate = (d) => {
    if (!d) return "-";
    return new Date(d).toLocaleString();
  };

  return (
    <div className="container">

      {/* HEADER */}
      <div className="header">
        <div className="title">📊 Report</div>
        <div className="subtitle">
          Completed stock movements
        </div>
      </div>

      {/* STOCK OUT */}
      <div className="card" style={{ marginBottom: "20px" }}>
        <h3 style={{ padding: "16px" }}>📤 Stock Out (Done)</h3>

        <table className="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Product</th>
              <th>Qty</th>
              <th>Created</th>
              <th>Done</th>
            </tr>
          </thead>

          <tbody>
            {doneOut.map((o) => {
              const item = o.items?.[0];

              return (
                <tr key={o.id}>
                  <td>{o.id}</td>
                  <td>{item?.product_name}</td>
                  <td>{item?.qty}</td>
                  <td>{formatDate(o.created_at)}</td>
                  <td>{formatDate(o.done_at)}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* STOCK IN */}
      <div className="card">
        <h3 style={{ padding: "16px" }}>📥 Stock In (Done)</h3>

        <table className="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Product</th>
              <th>Qty</th>
              <th>Created</th>
              <th>Done</th>
            </tr>
          </thead>

          <tbody>
            {doneIn.map((o) => {
              const item = o.items?.[0];

              return (
                <tr key={o.id}>
                  <td>{o.id}</td>
                  <td>{item?.product_name}</td>
                  <td>{item?.qty}</td>
                  <td>{formatDate(o.created_at)}</td>
                  <td>{formatDate(o.done_at)}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

    </div>
  );
}