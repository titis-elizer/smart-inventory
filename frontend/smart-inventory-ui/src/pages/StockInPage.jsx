import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import toast from "react-hot-toast";

import ItemSelect from "../components/ItemSelect";
import {
  createStockIn,
  updateStockInStatus,
} from "../features/stockIn/stockInSlice";
import { fetchInventory } from "../features/inventory/inventorySlice";
import { fetchStockIn } from "../features/stockIn/stockInSlice";

export default function StockInPage() {
  const dispatch = useDispatch();
  const orders = useSelector((s) => s.stockIn.orders);

  const [itemId, setItemId] = useState("");
  const [qty, setQty] = useState(0);

  useEffect(() => {
    dispatch(fetchInventory());
    dispatch(fetchStockIn());
  }, []);

  // 🔥 CREATE
  const create = async () => {
    if (!itemId) return toast.error("Select item");
    if (qty <= 0) return toast.error("Qty invalid");

    try {
      await dispatch(
        createStockIn([
          { inventory_item_id: itemId, qty: parseInt(qty) },
        ])
      ).unwrap();

      toast.success("Stock In created 🚀");

      dispatch(fetchStockIn());
    } catch (err) {
      toast.error(err || "Failed");
    }
  };

  // 🔥 ACTION HANDLER (DRY)
const handleAction = async (id, status) => {
  try {
    await dispatch(updateStockInStatus({ id, status })).unwrap();

    // delay kecil biar DB commit dulu
    await new Promise((r) => setTimeout(r, 300));

    dispatch(fetchStockIn());

    // ✅ toast success
    toast.success(`Status updated to ${status}`);
  } catch (err) {
    toast.error(err || "Failed to update status");
  }
};

  // 🔥 ACTION MAPPER
const getActions = (status, id) => {
  if (status === "in_progress") {
    return [
      {
        label: "Done",
        className: "btn btn-done",
        action: () => handleAction(id, "done"),
      },
      {
        label: "Cancel",
        className: "btn btn-cancel-red",
        action: () => handleAction(id, "canceled"),
      },
    ];
  }

  return [];
};

  return (
    <div className="container">

      {/* HEADER */}
      <div className="header">
        <div className="title">📥 Stock In</div>
        <div className="subtitle">
          Add incoming inventory
        </div>
      </div>

      {/* FORM */}
      <div className="form-card">
        <ItemSelect value={itemId} onChange={setItemId} />

        <input
          type="number"
          placeholder="Quantity"
          className="input"
          onChange={(e) => setQty(e.target.value)}
        />

        <button className="btn btn-dark" onClick={create}>
          Create Stock In
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
              <th>Action</th>
            </tr>
          </thead>

          <tbody>
            {(orders || []).map((o) => {
              const item = o.items?.[0];

              return (
                <tr key={o.id}>
                  <td>{o.id}</td>
                  <td>{item?.product_name}</td>
                  <td>{item?.qty}</td>
                  <td>{o.status}</td>

                  <td style={{ textAlign: "center" }}>
                    <div className="actions">
                      {getActions(o.status, o.id).length > 0 ? (
                        getActions(o.status, o.id).map((btn, i) => (
                          <button
                            key={i}
                            className={btn.className}
                            onClick={btn.action}
                          >
                            {btn.label}
                          </button>
                        ))
                      ) : (
                        <span style={{ color: "#9ca3af" }}>
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