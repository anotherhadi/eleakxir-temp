export interface Result {
  DataleakName: string;
  Columns: string[];
  Content: string[];
}

export interface Dataleaks {
  Dataleaks: Dataleak[];
  TotalRows: number;
  TotalDataleaks: number;
  TotalSize: number; // In MB
}

export interface Dataleak {
  Name: string;
  Path: string;
  Columns: string[];
  Length: number;
  Size: number; // In MB
}

export type ResearchStatus = "idle" | "searching" | "complete" | "error";
