import type { ReactNode } from 'react'

interface FiltersBarProps {
  search: string
  onSearchChange: (value: string) => void
  searchPlaceholder?: string
  children?: ReactNode
}

export function FiltersBar({
  search,
  onSearchChange,
  searchPlaceholder = 'Поиск…',
  children,
}: FiltersBarProps) {
  return (
    <div className="filters-bar card">
      <label className="search-field">
        <span>Поиск</span>
        <input
          type="search"
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder={searchPlaceholder}
        />
      </label>
      {children}
    </div>
  )
}
