alter table user_balance
  add constraint balance_non_negative check (balance >= 0);