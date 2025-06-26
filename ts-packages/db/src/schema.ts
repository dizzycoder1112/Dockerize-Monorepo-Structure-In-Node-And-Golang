import type { ColumnType, Generated } from "kysely";
import type { Buffer } from "node:buffer";

import type { BLOCKCHAIN_ID, USER_WALLET_TRANSACTIONS_ACTION } from "@ts-packages/shared/constants";

export interface Database {
  users: UsersTable;
  wallets: WalletsTable;
  vaultHistories: VaultHistoriesTable;
  vaults: VaultsTable;
  referralCodes: ReferralCodesTable;
  userBalances: UserBalancesTable;
  userWalletTransactions: UserWalletTransactionsTable;
}

export type Timestamp = ColumnType<Date, string | Date | undefined, Date>;

export interface UsersTable {
  id: string;
  nickname: string;
  referralCodeId: string | null;
  referralCodeCheck: ColumnType<boolean, undefined, boolean>;
  createdAt: Generated<Timestamp>;
  updatedAt: Generated<Timestamp>;
}

export interface UserBalancesTable {
  userId: string;
  balance: ColumnType<string, undefined, string>; // numeric(64)
  version: ColumnType<number, undefined, number>;
  updatedAt: Timestamp;
}

export interface UserWalletTransactionsTable {
  hash: Buffer;
  blockchainId: BLOCKCHAIN_ID;
  walletAddress: Buffer;
  action: USER_WALLET_TRANSACTIONS_ACTION;
  userId: string;
  createdAt: Timestamp;
}

export interface WalletsTable {
  address: Buffer;
  blockchainId: BLOCKCHAIN_ID;
  publicKey: string;
  userId: string | null;
  createdAt: Generated<Timestamp>;
  updatedAt: Generated<Timestamp>;
}

export interface VaultsTable {
  id: Generated<number>;
  address: Buffer;
  blockchainId: BLOCKCHAIN_ID;
  note: ColumnType<string, string | null, string | null>;
  createdAt: Timestamp;
}

export interface VaultHistoriesTable {
  hash: Buffer;
  blockchainId: BLOCKCHAIN_ID;
  symbol: string;
  source: Buffer;
  dest: Buffer;
  amount: ColumnType<string, string | bigint, string | never>;
  logicTime: ColumnType<string, string | bigint, never>;
  isHandled: ColumnType<boolean, boolean | undefined, boolean>;
  createdAt: Timestamp;
  status: ColumnType<string, string | undefined, string>;
  vaultId: number;
}

export interface ReferralCodesTable {
  id: string;
  code: string;
  creatorUserId: string;
  createdAt: Generated<Timestamp>;
  updatedAt: Generated<Timestamp>;
}
