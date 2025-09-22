// layout básico + protección de ruta por token (redirige si no hay token).
import { Outlet, useNavigate, useLocation, Link } from "react-router-dom";
import { useEffect, useState } from "react";
import { getToken, setToken } from "./api";

export default function App() {
  const nav = useNavigate();
  const loc = useLocation();
  const [token, setTok] = useState<string | null>(getToken());

  useEffect(() => {
    const t = getToken();
    setTok(t);
    if (!t && loc.pathname.startsWith("/vault")) nav("/");
  }, [loc.pathname, nav]);

  const logout = () => {
    setToken(null);
    setTok(null);
    nav("/");
  };

  return (
    <div style={{ fontFamily: "system-ui, sans-serif", padding: 24, maxWidth: 960, margin: "0 auto" }}>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
        <h1 style={{ margin: 0, fontSize: 20 }}>password-danie</h1>
        <nav style={{ display: "flex", gap: 12 }}>
          <Link to="/">Auth</Link>
          <Link to="/vault">Vault</Link>
          {token && <button onClick={logout} style={{ padding: "6px 10px" }}>Logout</button>}
        </nav>
      </header>
      <Outlet />
    </div>
  );
}
