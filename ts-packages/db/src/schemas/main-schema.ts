import type { ColumnType } from "kysely";
export type Generated<T> = T extends ColumnType<infer S, infer I, infer U>
  ? ColumnType<S, I | undefined, U>
  : ColumnType<T, T | undefined, T>;
export type Timestamp = ColumnType<Date, Date | string, Date | string>;

export type Deal = {
    id: Generated<string>;
    title: string;
    amount: string;
    status: string;
    created_at: Generated<Timestamp>;
    updated_at: Timestamp;
};
export type DB = {
    deals: Deal;
};
