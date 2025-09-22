// App.tsx: UI mínima — login (JWT) y CRUD básico de secrets (lista + crear).
import { useEffect, useState } from 'react'

type Secret = {
  id: number
  username: string
  url: string
  notes: string
  icon: string
  title: string
  created_at: string
  updated_at: string
}

export default function App() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [token, setToken] = useState<string | null>(null)
  const [items, setItems] = useState<Secret[]>([])
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState({ username: '', password_plain: '', url: '', notes: '', icon: '', title: '' })
  const [q, setQ] = useState('')

  useEffect(() => {
    const t = localStorage.getItem('jwt')
    if (t) setToken(t)
  }, [])

  async function login() {
    const res = await fetch('/api/v1/auth/login', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    })
    const data = await res.json()
    if (!res.ok) return alert(data.error || 'login failed')
    localStorage.setItem('jwt', data.access_token)
    setToken(data.access_token)
  }

  async function load() {
    if (!token) return
    setLoading(true)
    const res = await fetch(`/api/v1/vault/entries?q=${encodeURIComponent(q)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    const data = await res.json()
    setLoading(false)
    if (!res.ok) return alert(data.error || 'fetch failed')
    setItems(data.items || [])
  }

  async function create() {
    if (!token) return
    const res = await fetch('/api/v1/vault/entries', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(form)
    })
    const data = await res.json()
    if (!res.ok) return alert(data.error || 'create failed')
    setForm({ username: '', password_plain: '', url: '', notes: '', icon: '', title: '' })
    load()
  }

  function logout() {
    localStorage.removeItem('jwt')
    setToken(null)
    setItems([])
  }

  useEffect(() => { if (token) load() }, [token])

  if (!token) {
    return (
      <div style={{ maxWidth: 420, margin: '3rem auto', fontFamily: 'system-ui' }}>
        <h1>password-danie</h1>
        <h2>Login</h2>
        <input placeholder="email" value={email} onChange={e => setEmail(e.target.value)} /><br />
        <input placeholder="password" type="password" value={password} onChange={e => setPassword(e.target.value)} /><br />
        <button onClick={login}>Entrar</button>
      </div>
    )
  }

  return (
    <div style={{ maxWidth: 900, margin: '2rem auto', fontFamily: 'system-ui' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>Vault</h1>
        <button onClick={logout}>Cerrar sesión</button>
      </header>

      <section style={{ margin: '1rem 0' }}>
        <input placeholder="buscar..." value={q} onChange={e => setQ(e.target.value)} />
        <button onClick={load} disabled={loading}>{loading ? 'Cargando...' : 'Buscar'}</button>
      </section>

      <section style={{ border: '1px solid #ddd', padding: 12, borderRadius: 8 }}>
        <h3>Nuevo secreto</h3>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 8 }}>
          <input placeholder="username" value={form.username} onChange={e => setForm({ ...form, username: e.target.value })} />
          <input placeholder="password" type="password" value={form.password_plain} onChange={e => setForm({ ...form, password_plain: e.target.value })} />
          <input placeholder="url" value={form.url} onChange={e => setForm({ ...form, url: e.target.value })} />
          <input placeholder="notes" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} />
          <input placeholder="icon" value={form.icon} onChange={e => setForm({ ...form, icon: e.target.value })} />
          <input placeholder="title (opcional)" value={form.title} onChange={e => setForm({ ...form, title: e.target.value })} />
        </div>
        <button style={{ marginTop: 8 }} onClick={create}>Crear</button>
      </section>

      <section style={{ marginTop: 16 }}>
        <h3>Resultados</h3>
        {items.length === 0 ? <p>No hay elementos.</p> : (
          <ul>
            {items.map(s => (
              <li key={s.id} style={{ padding: 8, borderBottom: '1px solid #eee' }}>
                <strong>{s.title || s.username}</strong> — <a href={s.url} target="_blank">{s.url}</a>
                <div style={{ fontSize: 12, color: '#555' }}>{s.notes}</div>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}
