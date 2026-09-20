import type { ReactNode } from "react";

import { Table } from "@radix-ui/themes";
import { UnorderedListIcon } from "nfx-ui/icons";

import EmptyState from "./EmptyState";

export type DataTableColumn<T> = {
  key: string;
  header: string;
  width?: string;
  render?: (row: T) => ReactNode;
};

type DataTableProps<T> = {
  columns: DataTableColumn<T>[];
  rows: T[];
  rowKey: (row: T) => string;
  empty?: string;
  loading?: boolean;
  onRowClick?: (row: T) => void;
};

export function DataTable<T>({ columns, rows, rowKey, empty, loading, onRowClick }: DataTableProps<T>) {
  if (loading) {
    return <EmptyState icon={UnorderedListIcon} title={empty ?? "…"} />;
  }
  if (!rows.length) {
    return <EmptyState icon={UnorderedListIcon} title={empty ?? "—"} />;
  }
  return (
    <Table.Root size="2">
      <Table.Header>
        <Table.Row>
          {columns.map((column) => (
            <Table.ColumnHeaderCell key={column.key} style={column.width ? { width: column.width } : undefined}>
              {column.header}
            </Table.ColumnHeaderCell>
          ))}
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {rows.map((row) => (
          <Table.Row
            key={rowKey(row)}
            onClick={onRowClick ? () => onRowClick(row) : undefined}
            style={onRowClick ? { cursor: "pointer" } : undefined}
          >
            {columns.map((column) => (
              <Table.Cell key={column.key}>
                {column.render ? column.render(row) : String((row as Record<string, unknown>)[column.key] ?? "")}
              </Table.Cell>
            ))}
          </Table.Row>
        ))}
      </Table.Body>
    </Table.Root>
  );
}
