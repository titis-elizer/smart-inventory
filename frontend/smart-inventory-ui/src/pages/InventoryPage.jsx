import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { fetchInventory, adjustStock } from "../features/inventory/inventorySlice";

export default function InventoryPage() {
  const dispatch = useDispatch();
  const { items } = useSelector((s) => s.inventory);

  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [editId, setEditId] = useState(null);
  const [newStock, setNewStock] = useState(0);

  const pageSize = 5;

  useEffect(() => {
    dispatch(fetchInventory(search));
  }, [search]);

  // 🔍 Filter + Pagination
  const filtered = items.filter((i) =>
    i.name.toLowerCase().includes(search.toLowerCase())
  );

  const totalPages = Math.ceil(filtered.length / pageSize);
  const paginated = filtered.slice(
    (page - 1) * pageSize,
    page * pageSize
  );

const handleSave = (item) => {
  const newVal = Number(newStock);

  // validasi
  if (isNaN(newVal)) {
    alert("Stock harus angka");
    return;
  }

  const qty = newVal - item.physical_stock;

  // optional: skip kalau tidak berubah
  if (qty === 0) {
    setEditId(null);
    return;
  }

  dispatch(
    adjustStock({
      item_id: item.id, // pastikan ini ada!
      qty: Number(qty),
    })
  );

  setEditId(null);
};

  return (
  <div className="page">
  <div className="container">

    {/* HEADER */}
    <div className="header">
      <div className="title">📦 Inventory</div>
      <div className="subtitle">
        Manage your product stock easily
      </div>
    </div>

    {/* SEARCH */}
    <div className="search-wrapper">
      <div className="search-box">
        <span className="search-icon">🔍</span>
        <input
          className="search-input"
          placeholder="Search product..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>
    </div>

    {/* TABLE */}
    <div className="card">
      <table className="table">
        <thead>
          <tr>
            <th>Product</th>
            <th>SKU</th>
            <th>Stock</th>
            <th style={{ textAlign: "center" }}>Action</th>
          </tr>
        </thead>

        <tbody>
          {paginated.map((item) => {
            const available =
              item.physical_stock - item.reserved_stock;

            return (
              <tr key={item.id}>
                <td><b>{item.name}</b></td>
                <td>{item.sku}</td>

                <td>
                  {editId === item.id ? (
                    <input
                      type="number"
                      className="input-stock"
                      value={newStock}
                      onChange={(e) => setNewStock(e.target.value)}
                    />
                  ) : (
                    <div className="badges">
                      <span className="badge blue">
                        P: {item.physical_stock}
                      </span>
                      <span className="badge yellow">
                        R: {item.reserved_stock}
                      </span>
                      <span
                        className={`badge ${
                          available > 10 ? "green" : "red"
                        }`}
                      >
                        A: {available}
                      </span>
                    </div>
                  )}
                </td>

                <td style={{ textAlign: "center" }}>
                  {editId === item.id ? (
                    <>
                      <button
                        className="btn btn-save"
                        onClick={() => handleSave(item)}
                      >
                        Save
                      </button>
                      <button
                        className="btn btn-cancel"
                        onClick={() => setEditId(null)}
                      >
                        Cancel
                      </button>
                    </>
                  ) : (
                    <button
                      className="btn btn-edit"
                      onClick={() => {
                        setEditId(item.id);
                        setNewStock(item.physical_stock);
                      }}
                    >
                      Edit
                    </button>
                  )}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>

    {/* PAGINATION */}
    <div className="pagination">
      <div>
        Page {page} of {totalPages}
      </div>

      <div className="page-buttons">
        {Array.from({ length: totalPages }, (_, i) => (
          <button
            key={i}
            onClick={() => setPage(i + 1)}
            className={`page-btn ${
              page === i + 1 ? "active" : ""
            }`}
          >
            {i + 1}
          </button>
        ))}
      </div>
    </div>

  </div>
</div>
  );
}