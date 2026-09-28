import { describe, expect, it } from 'vitest'
import { buildTaskPayload, connectionDatabaseName, connectionDialect, isValidConnection, isValidDatabaseURL, maskedConnectionURL } from './connection'

const connection = { dialect: 'mysql' as const, host: '127.0.0.1', port: 3306, database: 'app', username: 'readonly', password: 'secret' }

describe('connection helpers', () => {
  it('accepts a complete mysql URL and rejects other protocols', () => {
    expect(isValidDatabaseURL('mysql://readonly:secret@127.0.0.1:3306/app')).toBe(true)
    expect(isValidDatabaseURL('postgres://readonly:secret@127.0.0.1/app')).toBe(true)
    expect(isValidDatabaseURL('mysql://readonly@127.0.0.1')).toBe(false)
  })

  it('builds exactly one connection input mode', () => {
    const byURL = buildTaskPayload('url', 'mysql://readonly:secret@127.0.0.1:3306/app', connection, 'audit_log, temp', true)
    expect(byURL.connectionUrl).toContain('mysql://')
    expect(byURL.connection).toBeUndefined()
    expect(byURL.excludeTables).toEqual(['audit_log', 'temp'])

    const byFields = buildTaskPayload('fields', '', connection, '', false)
    expect(byFields.connection).toEqual(connection)
    expect(byFields.connectionUrl).toBeUndefined()
    expect(isValidConnection(connection)).toBe(true)
    expect(isValidConnection({ ...connection, database: '' })).toBe(false)
  })

  it('stores a password-masked URL label only', () => {
    const masked = maskedConnectionURL('url', 'mysql://readonly:secret@db.example:3306/app_db?tls=skip', connection)
    expect(masked).toBe('mysql://readonly:••••••@db.example:3306/app_db')
    expect(masked).not.toContain('secret')
    expect(connectionDatabaseName('url', 'mysql://readonly:secret@db.example:3306/app_db?tls=skip', connection)).toBe('app_db')
    expect(connectionDialect('url', 'postgres://readonly:secret@db.example/app_db', connection)).toBe('postgres')
  })
})
