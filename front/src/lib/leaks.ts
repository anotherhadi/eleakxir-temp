import type { Result } from "./types";

export function GetIndexOrFirstNonEmpty(row: string[], index: number): string {
  const label = row[index];
  if (label !== undefined && label !== "") {
    return row[index];
  }

  for (let i = 0; i < row.length; i++) {
    if (row[i] !== undefined && row[i] !== "") {
      return row[i];
    }
  }
  return "";
}

export function GetIndexToPrioritise(result: Result): number {
  if (result.Columns.includes("email")) {
    return result.Columns.indexOf("email");
  }
  if (result.Columns.includes("username")) {
    return result.Columns.indexOf("username");
  }
  if (result.Columns.includes("name")) {
    return result.Columns.indexOf("name");
  }
  if (result.Columns.includes("last_name")) {
    return result.Columns.indexOf("lastname");
  }
  if (result.Columns.includes("first_name")) {
    return result.Columns.indexOf("firstname");
  }
  if (result.Columns.includes("password")) {
    return result.Columns.indexOf("password");
  }
  if (result.Columns.includes("phone")) {
    return result.Columns.indexOf("phone");
  }
  if (result.Columns.includes("address")) {
    return result.Columns.indexOf("address");
  }
  return 0;
}

export function FormatSize(mb: number): string {
  const gb = mb / 1024;
  if (gb >= 1) {
    return `${gb.toFixed(gb < 10 ? 1 : 0)} GB`;
  }
  return `${Math.round(mb)} MB`;
}
