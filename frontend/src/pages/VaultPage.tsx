// pantalla de Vault con acciones claras: Create, Read, Search/List, Update, Delete.
import { useEffect, useState } from "react";
import { api } from "../api";

type Secret = {
  id: number;
  username: string;
  url?: string;
  notes?: string;
  icon?: string;
  title?: string;
  created_at?: string;
  updated_at?: string;
};

export default function VaultPage() {
  const [q, setQ] = useState("");
  const [items, setItems] = useState<Secret[]>([]);
  const [total, setTotal] = useState(0);
  const [message, setMessage] = useState<string | null>(null);

  const [createForm, setCreateForm] = useState({
    username: "danie",
    password_plain: "p@ss",
    url: "https://github.com",
    notes: "mi cuenta",
    icon: "github",
    title: "GitHub",
  });

  const [readId, setReadId] = useState<number | "">("");
  const [readItem, setReadItem] = useState<Secret | null>(null);

  const [updateId, setUpdateId] = useState<number | "">("");
  const [updateForm, setUpdateForm] = useState<{ notes?: string; password_plain?: string }>({});

  const [deleteId, setDeleteId] = useState<number | "">("");

  const load = async () => {
    setMessage(null);
    try {
      const res = await api.listSecrets(q, 50, 0);
      setItems(res.items || []);
      setTotal(res.total || 0);
    } catch (e: any) {
      setMessage(e.message || "Error");
    }
  };

  useEffect(() => { load(); /* auto-load al entrar */ }, []);

  const onCreate = async () => {
    setMessage(null);
    try {
      const res = await api.createSecret(createForm);
      setMessage(`Create OK (id=${res.id})`);
      await load();
    } catch (e: any) {
      setMessage(e.message || "Error");
    }
  };

  const onRead = async () => {
    setMessage(null);
    setReadItem(null);
    if (readId === "") return;
    try {
      const res = await api.getSecret(Number(readId));
      setReadItem(res);
      setMessage("Read OK");
    } catch (e: any) {
      setMessage(e.message || "Error");
    }
  };

  const onUpdate = async () => {
    setMessage(null);
    if (updateId === "") return;
    try {
      await api.updateSecret(Number(updateId), updateForm);
      setMessage("Update OK");
      await load();
    } catch (e: any) {
      setMessage(e.message || "Error");
    }
  };

  const onDelete = async () => {
    setMessage(null);
    if (deleteId === "") return;
    try {
      await api.deleteSecret(Number(deleteId));
      setMessage("Delete OK");
      await load();
    } catch (e: any) {
      setMessage(e.message || "Error");
    }
  };

  return (
    <div style={{ display: "grid", gap: 24 }}>
      <h2>Vault</h2>
      {message && <div style={{ padding: 8, background: "#f3f3f3" }}>{message}</div>}

      {/* SEARCH / LIST */}
      <section>
        <h3>Search / List</h3>
        <div style={{ display: "flex", gap: 8, marginBottom: 8 }}>
          <input value={q} onChange={(e) => setQ(e.target.value)} placeholder="q (buscar por título/url/username)" />
          <button onClick={load}>Buscar</button>
        </div>
        <div style={{ fontSize: 12, marginBottom: 8 }}>Total: {total}</div>
        <table style={{ width: "100%", borderCollapse: "collapse" }}>
          <thead>
            <tr style={{ textAlign: "left", borderBottom: "1px solid #ddd" }}>
              <th>ID</th><th>Title</th><th>Username</th><th>URL</th><th>Updated</th>
            </tr>
          </thead>
          <tbody>
            {items.map(it => (
              <tr key={it.id} style={{ borderBottom: "1px solid #eee" }}>
                <td>{it.id}</td>
                <td>{it.title || "-"}</td>
                <td>{it.username}</td>
                <td>{it.url || "-"}</td>
                <td>{it.updated_at ? new Date(it.updated_at).toLocaleString() : "-"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>

      {/* CREATE */}
      <section>
        <h3>Create</h3>
        <div style={{ display: "grid", gap: 6, maxWidth: 600 }}>
          {Object.entries(createForm).map(([k, v]) => (
            <div key={k} style={{ display: "grid", gridTemplateColumns: "140px 1fr", gap: 8 }}>
              <label style={{ textTransform: "capitalize" }}>{k}</label>
              <input
                value={v as string}
                onChange={(e) => setCreateForm({ ...createForm, [k]: e.target.value })}
                placeholder={k}
                type={k === "password_plain" ? "password" : "text"}
              />
            </div>
          ))}
          <button onClick={onCreate}>Create</button>
        </div>
      </section>

      {/* READ */}
      <section>
        <h3>Read (by id)</h3>
        <div style={{ display: "flex", gap: 8 }}>
          <input
            value={readId}
            onChange={(e) => setReadId(e.target.value ? Number(e.target.value) : "")}
            placeholder="id" type="number"
          />
          <button onClick={onRead}>Read</button>
        </div>
        {readItem && (
          <pre style={{ background: "#f7f7f7", padding: 12, marginTop: 8 }}>
{JSON.stringify(readItem, null, 2)}
          </pre>
        )}
      </section>

      {/* UPDATE */}
      <section>
        <h3>Update</h3>
        <div style={{ display: "grid", gap: 6, maxWidth: 600 }}>
          <div style={{ display: "flex", gap: 8 }}>
            <input
              value={updateId}
              onChange={(e) => setUpdateId(e.target.value ? Number(e.target.value) : "")}
              placeholder="id" type="number"
            />
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "140px 1fr", gap: 8 }}>
            <label>notes</label>
            <input value={updateForm.notes || ""} onChange={(e) => setUpdateForm({ ...updateForm, notes: e.target.value })} placeholder="notes" />
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "140px 1fr", gap: 8 }}>
            <label>password_plain</label>
            <input
              type="password"
              value={updateForm.password_plain || ""}
              onChange={(e) => setUpdateForm({ ...updateForm, password_plain: e.target.value })}
              placeholder="password_plain"
            />
          </div>
          <button onClick={onUpdate}>Update</button>
        </div>
      </section>

      {/* DELETE */}
      <section>
        <h3>Delete</h3>
        <div style={{ display: "flex", gap: 8 }}>
          <input
            value={deleteId}
            onChange={(e) => setDeleteId(e.target.value ? Number(e.target.value) : "")}
            placeholder="id" type="number"
          />
          <button onClick={onDelete} style={{ background: "#ffe8e8" }}>Delete</button>
        </div>
      </section>
    </div>
  );
}
