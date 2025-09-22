//  pantalla de autenticación con pestañas Login/Registro.
import { useState } from "react";
import { api, setToken } from "../api";
import { useNavigate } from "react-router-dom";

export default function AuthPage() {
  const [tab, setTab] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("test@example.com");
  const [password, setPassword] = useState("Secret123!");
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState<string | null>(null);
  const nav = useNavigate();

  const submit = async () => {
    setMsg(null);
    setLoading(true);
    try {
      if (tab === "register") {
        await api.register(email, password);
        setMsg("Registro OK. Ahora inicia sesión.");
        setTab("login");
      } else {
        const res = await api.login(email, password);
        setToken(res.access_token);
        setMsg("Login OK");
        nav("/vault");
      }
    } catch (e: any) {
      setMsg(e.message || "Error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div style={{ display: "flex", gap: 12, marginBottom: 16 }}>
        <button onClick={() => setTab("login")} disabled={tab === "login"}>Login</button>
        <button onClick={() => setTab("register")} disabled={tab === "register"}>Registro</button>
      </div>

      <div style={{ display: "grid", gap: 8, maxWidth: 400 }}>
        <label>Email</label>
        <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="email" />

        <label>Password</label>
        <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="password" />

        <button onClick={submit} disabled={loading} style={{ marginTop: 8 }}>
          {loading ? "..." : tab === "register" ? "Crear cuenta" : "Entrar"}
        </button>

        {msg && <div style={{ color: "#444", marginTop: 8 }}>{msg}</div>}
      </div>

      <hr style={{ margin: "24px 0" }} />

      <ResetBlock />
    </div>
  );
}

function ResetBlock() {
  const [email, setEmail] = useState("test@example.com");
  const [token, setTok] = useState("");
  const [newPass, setNewPass] = useState("");
  const [msg, setMsg] = useState<string | null>(null);

  const request = async () => {
    setMsg(null);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/auth/reset/request`, {
        method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ email })
      });
      const j = await res.json();
      if (!res.ok) throw new Error(j.error || "error");
      setTok(j.reset_token || "");
      setMsg("Reset token generado (visible solo en dev).");
    } catch (e: any) {
      setMsg(e.message || "Error");
    }
  };

  const confirm = async () => {
    setMsg(null);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_BASE_URL}/api/v1/auth/reset/confirm`, {
        method: "POST", headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ token, new_password: newPass })
      });
      const j = await res.json();
      if (!res.ok) throw new Error(j.error || "error");
      setMsg("Contraseña actualizada, ya puedes loguearte.");
    } catch (e: any) {
      setMsg(e.message || "Error");
    }
  };

  return (
    <div>
      <h3 style={{ marginBottom: 8 }}>Reset de contraseña (dev)</h3>
      <div style={{ display: "grid", gap: 6, maxWidth: 480 }}>
        <div style={{ display: "grid", gap: 6, gridTemplateColumns: "1fr auto" }}>
          <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="email" />
          <button onClick={request}>Request</button>
        </div>
        <input value={token} onChange={(e) => setTok(e.target.value)} placeholder="reset token" />
        <div style={{ display: "grid", gap: 6, gridTemplateColumns: "1fr auto" }}>
          <input value={newPass} onChange={(e) => setNewPass(e.target.value)} type="password" placeholder="nueva contraseña" />
          <button onClick={confirm}>Confirm</button>
        </div>
        {msg && <div style={{ color: "#444" }}>{msg}</div>}
      </div>
    </div>
  );
}
