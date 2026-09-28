import type { ConnectionConfig, CreateTaskPayload } from '@/types/domain'

export function isValidDatabaseURL(value: string) {
  try {
    const url = new URL(value.trim())
    return ['mysql:', 'postgres:', 'postgresql:'].includes(url.protocol) && Boolean(url.hostname && url.username && url.pathname && url.pathname !== '/')
  } catch {
    return false
  }
}

export function isValidConnection(connection: ConnectionConfig) {
  return Boolean(connection.host.trim() && connection.database.trim() && connection.username.trim() && connection.port >= 1 && connection.port <= 65535)
}

// This value is safe for browser history only: the raw password and URL query
// are deliberately never retained. It is a human-readable label, not a URL to
// reconnect with.
export function maskedConnectionURL(mode: 'url' | 'fields', rawURL: string, connection: ConnectionConfig): string {
  if (mode === 'url') {
    const parsed = new URL(rawURL)
    const username = decodeURIComponent(parsed.username) || '***'
    return `${parsed.protocol}//${username}:••••••@${parsed.host}${parsed.pathname}`
  }
  const scheme = connection.dialect === 'postgres' ? 'postgres' : 'mysql'
  const username = encodeURIComponent(connection.username.trim() || '***')
  const database = encodeURIComponent(connection.database.trim())
  return `${scheme}://${username}:••••••@${connection.host.trim()}:${connection.port}/${database}`
}

export function connectionDatabaseName(mode: 'url' | 'fields', rawURL: string, connection: ConnectionConfig): string {
  if (mode === 'url') {
    const parsed = new URL(rawURL)
    return decodeURIComponent(parsed.pathname.replace(/^\//, '')) || '—'
  }
  return connection.database.trim() || '—'
}

export function connectionDialect(mode: 'url' | 'fields', rawURL: string, connection: ConnectionConfig): 'mysql' | 'postgres' {
  if (mode === 'url') {
    const protocol = new URL(rawURL).protocol
    return protocol === 'postgres:' || protocol === 'postgresql:' ? 'postgres' : 'mysql'
  }
  return connection.dialect
}

export function buildTaskPayload(mode: 'url' | 'fields', connectionURL: string, connection: ConnectionConfig, excludeTables: string, redactObjectNames: boolean): CreateTaskPayload {
  const payload: CreateTaskPayload = {
    excludeTables: excludeTables.split(',').map((name) => name.trim()).filter(Boolean),
    redactObjectNames,
  }
  if (mode === 'url') payload.connectionUrl = connectionURL.trim()
  else payload.connection = { ...connection }
  return payload
}
