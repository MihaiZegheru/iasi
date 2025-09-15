export interface Problem {
  name: string;
  url: string;
  time: string;
  id: string;
  locks?: {
    hints?: boolean[]; // true = locked, false = unlocked
    editorial?: boolean; // true = locked, false = unlocked
  };
}
