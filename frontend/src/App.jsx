import React, { useState, useEffect } from 'react'

const API = 'http://localhost:8080'

export default function App() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [token, setToken] = useState('')
  const [items, setItems] = useState([])
  const [title, setTitle] = useState('')
  const [username, setUsername] = useState('')
  const [pwd, setPwd] = useState('')
  const [url, setUrl] = useState('')
  const [notes, setNotes] = useState('')
  const [detail, setDetail] = useState(null)
  const [query, setQuery] = useState('')

  const authHeaders = token ? { 'Authorization': `Bearer ${token}` } : {}

  async function register(e) {
    e.preventDefault()
    const r = await fetch(`${API}/auth/register`, {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ email, password })
    })
    if (r.ok) alert('Usuario registrado. Ahora haz login.')
    else alert('Registro falló')
  }

  async function login(e) {
    e.preventDefault()
    const r = await fetch(`${API}/auth/login`, {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ email, password })
    })
    if (!r.ok) return alert('Login falló')
    const data = await r.json()
    setToken(data.token)
    setEmail(''); setPassword('')
    loadItems()
  }

  async function loadItems() {
    if (!token) return
    const url = new URL(`${API}/api/vault`)
    if (query) url.searchParams.set('q', query)
    const r = await fetch(url, { headers: authHeaders })
    if (r.ok) setItems(await r.json())
  }

  async function createItem(e) {
    e.preventDefault()
    if (!token) return alert('Login primero')
    const r = await fetch(`${API}/api/vault`, {
      method: 'POST',
      headers: { 'Content-Type':'application/json', ...authHeaders },
      body: JSON.stringify({ title, username, password: pwd, url, notes, icon: '' })
    })
    if (r.ok) {
      setTitle(''); setUsername(''); setPwd(''); setUrl(''); setNotes('')
      loadItems()
    } else {
      alert('Error creando item')
    }
  }

  async function viewDetail(id) {
    const r = await fetch(`${API}/api/vault/${id}`, { headers: authHeaders })
    if (r.ok) setDetail(await r.json())
  }

  async function delItem(id) {
    if (!confirm('¿Eliminar item?')) return
    const r = await fetch(`${API}/api/vault/${id}`, { method: 'DELETE', headers: authHeaders })
    if (r.ok) {
      setDetail(null)
      loadItems()
    }
  }

  useEffect(() => { if (token) loadItems() }, [token])
  useEffect(() => { const t = setTimeout(loadItems, 300); return () => clearTimeout(t) }, [query])

  return (
    <div style={{maxWidth: 900, margin: '20px auto', fontFamily: 'system-ui, Arial'}}>
      <h1>password-danie</h1>

      {!token && (
        <div style={{display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16}}>
          <form onSubmit={register} style={card}>
            <h3>Registro</h3>
            <input placeholder="email" value={email} onChange={e=>setEmail(e.target.value)} style={inp} />
            <input placeholder="password" type="password" value={password} onChange={e=>setPassword(e.target.value)} style={inp}/>
            <button type="submit">Crear cuenta</button>
          </form>

          <form onSubmit={login} style={card}>
            <h3>Login</h3>
            <input placeholder="email" value={email} onChange={e=>setEmail(e.target.value)} style={inp} />
            <input placeholder="password" type="password" value={password} onChange={e=>setPassword(e.target.value)} style={inp}/>
            <button type="submit">Entrar</button>
          </form>
        </div>
      )}

      {token && (
        <>
          <div style={{display:'flex', alignItems:'center', gap:12, margin:'12px 0'}}>
            <input placeholder="Buscar (q)" value={query} onChange={e=>setQuery(e.target.value)} style={{...inp, flex:1}} />
            <button onClick={loadItems}>Buscar</button>
            <button onClick={()=>{setToken(''); setItems([]); setDetail(null)}}>Salir</button>
          </div>

          <form onSubmit={createItem} style={card}>
            <h3>Nuevo Item</h3>
            <input placeholder="title" value={title} onChange={e=>setTitle(e.target.value)} style={inp}/>
            <input placeholder="username" value={username} onChange={e=>setUsername(e.target.value)} style={inp}/>
            <input placeholder="password" type="password" value={pwd} onChange={e=>setPwd(e.target.value)} style={inp}/>
            <input placeholder="url" value={url} onChange={e=>setUrl(e.target.value)} style={inp}/>
            <textarea placeholder="notes" value={notes} onChange={e=>setNotes(e.target.value)} style={inp}/>
            <button type="submit">Crear</button>
          </form>

          <div style={{display:'grid', gridTemplateColumns:'1fr 1fr', gap:16}}>
            <div style={card}>
              <h3>Mis Items</h3>
              {items.length===0 && <p>No hay items</p>}
              <ul>
                {items.map(it=>(
                  <li key={it.id} style={{display:'flex', justifyContent:'space-between', margin:'6px 0'}}>
                    <span>{it.title} — {it.username}</span>
                    <div style={{display:'flex', gap:8}}>
                      <button onClick={()=>viewDetail(it.id)}>Ver</button>
                      <button onClick={()=>delItem(it.id)}>Borrar</button>
                    </div>
                  </li>
                ))}
              </ul>
            </div>

            <div style={card}>
              <h3>Detalle</h3>
              {!detail ? <p>Selecciona un item</p> : (
                <div>
                  <div><b>Title:</b> {detail.title}</div>
                  <div><b>Username:</b> {detail.username}</div>
                  <div><b>Password:</b> {detail.password}</div>
                  <div><b>URL:</b> {detail.url}</div>
                  <div><b>Notes:</b> {detail.notes}</div>
                  <div><b>id:</b> {detail.id}</div>
                </div>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  )
}

const card = { padding:16, border:'1px solid #ddd', borderRadius:8, boxShadow:'0 1px 3px rgba(0,0,0,.06)' }
const inp = { padding:10, border:'1px solid #ccc', borderRadius:6, margin:'6px 0', width:'100%' }
