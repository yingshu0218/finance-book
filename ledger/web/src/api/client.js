const BASE_URL = '/api'

export async function fetchBooks() {
  const res = await fetch(`${BASE_URL}/books`)
  return res.json()
}

export async function createBook(name) {
  const res = await fetch(`${BASE_URL}/books`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name })
  })
  return res.json()
}

export async function deleteBook(name) {
  const res = await fetch(`${BASE_URL}/books/${name}`, {
    method: 'DELETE'
  })
  return res.json()
}

export async function fetchEntries(params = {}) {
  const query = new URLSearchParams(params).toString()
  const res = await fetch(`${BASE_URL}/entries?${query}`)
  return res.json()
}

export async function createEntry(data) {
  const res = await fetch(`${BASE_URL}/entries`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  return res.json()
}

export async function updateEntry(id, data) {
  const res = await fetch(`${BASE_URL}/entries/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  })
  return res.json()
}

export async function deleteEntry(id, book) {
  const res = await fetch(`${BASE_URL}/entries/${id}?book=${book}`, {
    method: 'DELETE'
  })
  return res.json()
}

export async function fetchBalance(params = {}) {
  const query = new URLSearchParams(params).toString()
  const res = await fetch(`${BASE_URL}/stats/balance?${query}`)
  return res.json()
}

export async function fetchCategories() {
  const res = await fetch(`${BASE_URL}/categories`)
  return res.json()
}
