export default function MainLayout({ children }) {
  const path = window.location.pathname;

  return (
    <div className="layout">

      {/* SIDEBAR */}
      <div className="sidebar">
        <div className="logo">📦 Inventory</div>

        <nav className="nav">
          <a
            href="/inventory"
            className={`nav-link ${
              path === "/inventory" ? "active" : ""
            }`}
          >
            📋 Inventory
          </a>

          <a
            href="/stock-in"
            className={`nav-link ${
              path === "/stock-in" ? "active" : ""
            }`}
          >
            📥 Stock In
          </a>

          <a
            href="/stock-out"
            className={`nav-link ${
              path === "/stock-out" ? "active" : ""
            }`}
          >
            📤 Stock Out
          </a>

          <a
            href="/report"
            className={`nav-link ${
              path === "/report" ? "active" : ""
            }`}
          >
            📊 Report
          </a>
        </nav>
      </div>

      {/* MAIN CONTENT */}
      <div className="content">

        {/* OPTIONAL TOPBAR */}
        <div className="topbar">
          Welcome 👋
        </div>

        {children}
      </div>
    </div>
  );
}