export function matchesSearch(query: string, values: Array<string | number | undefined | null>) {
  const normalized = query.trim().toLowerCase()
  if (!normalized) {
    return true
  }

  return values.some((value) => String(value ?? '').toLowerCase().includes(normalized))
}
