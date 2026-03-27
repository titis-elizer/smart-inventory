import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import toast from "react-hot-toast";

import ItemSelect from "../components/ItemSelect";
import {
  createStockOut,
  markInProgress,
  completeStockOut,
  cancelStockOut,
  fetchStockOut,
} from "../features/stockOut/stockOutSlice";
import { fetchInventory } from "../features/inventory/inventorySlice";

export default function StockOutPage() {
  const dispatch = useDispatch();
const orders = useSelector((s) => s.stockOut.orders);
const safeOrders = orders || [];

  const [itemId, setItemId] = useState("");
  const [qty, setQty] = useState(0);

 const create = async () => {
  try {
    await dispatch(
      createStockOut([
        { inventory_item_id: itemId, qty: parseInt(qty) },
      ])
    ).unwrap(); // 🔥 ini kunci

    toast.success("Allocated successfully");
  } catch (err) {
    console.log("ERROR:", err);

    // ✅ custom message dari backend
    if (typeof err === "string") {
      toast.error(err);
    } else if (err?.error) {
      toast.error(err.error);
    } else {
      toast.error("Insufficient stock");
    }
  }
};

  const getStatusClass = (status) => {
    if (status === "pending") return "status pending";
    if (status === "in_progress") return "status progress";
    if (status === "complete") return "status done";
    if (status === "cancelled") return "status cancel";
    return "status";
  };

  const handleAction = async (action, id) => {
  try {
    await dispatch(action(id));

    // ⏱️ delay 300ms biar backend commit dulu
    await new Promise((res) => setTimeout(res, 300));

    dispatch(fetchStockOut());
  } catch (err) {
    console.error(err);
  }
};
const getActions = (status, id) => {
  switch (status) {
    case "allocated":
      return [
        {
          label: "Process",
          className: "btn btn-process",
          action: () => handleAction(markInProgress, id),
        },
        {
          label: "Cancel",
          className: "btn btn-cancel-red",
          action: () => handleAction(cancelStockOut, id),
        },
      ];

    case "in_progress":
      return [
        {
          label: "Done",
          className: "btn btn-done",
          action: () => handleAction(completeStockOut, id),
        },
        {
          label: "Cancel",
          className: "btn btn-cancel-red",
          action: () => handleAction(cancelStockOut, id),
        },
      ];

    default:
      return [];
  }
};
  useEffect(() => {
  dispatch(fetchInventory());
  dispatch(fetchStockOut());
}, []);

  return (
    <div className="container">

      {/* HEADER */}
      <div className="header">
        <div className="title">📤 Stock Out</div>
        <div className="subtitle">
          Manage outgoing inventory orders
        </div>
      </div>

      {/* FORM */}
      <div className="form-card">
        <div className="form-group">
          <ItemSelect value={itemId} onChange={setItemId} />
        </div>

        <div className="form-group">
          <input
            type="number"
            placeholder="Quantity"
            className="input"
            onChange={(e) => setQty(e.target.value)}
          />
        </div>

        <button className="btn btn-dark" onClick={create}>
          Create Order
        </button>
      </div>

      {/* TABLE */}
      <div className="card" style={{ marginTop: "20px" }}>
        <table className="table">
        <thead>
  <tr>
    <th>ID</th>
    <th>Product</th>
    <th>Qty</th>
    <th>Status</th>
    <th style={{ textAlign: "center" }}>Action</th>
  </tr>
</thead>

       <tbody>
  {safeOrders.map((o) => {
    const item = o.items?.[0]; // asumsi 1 item dulu

    return (
      <tr key={o.id}>
        {/* ID */}
        <td>{o.id}</td>

        {/* PRODUCT NAME */}
        <td>{item?.product_name || "-"}</td>

        {/* QTY */}
        <td>{item?.qty || 0}</td>

        {/* STATUS */}
        <td>
          <span className={getStatusClass(o.status)}>
            {o.status}
          </span>
        </td>

        {/* ACTION */}
      <td style={{ textAlign: "center" }}>
  <div className="actions">

    {getActions(o.status, o.id).length > 0 ? (
      getActions(o.status, o.id).map((btn, idx) => (
        <button
          key={idx}
          className={btn.className}
          onClick={btn.action}
        >
          {btn.label}
        </button>
      ))
    ) : (
      <span style={{ color: "#9ca3af", fontSize: "12px" }}>
        No Action
      </span>
    )}

  </div>
</td>
      </tr>
    );
  })}
</tbody>
        </table>
      </div>
    </div>
  );
}