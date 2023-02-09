#include "state/state.h"

#include "common/account_state.h"
#include "common/file_util.h"
#include "common/status_test_util.h"
#include "common/type.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "state/configurations.h"

namespace carmen {

using ::testing::ElementsAre;
using ::testing::ElementsAreArray;
using ::testing::IsOkAndHolds;

template <typename T>
class StateTest : public testing::Test {};

TYPED_TEST_SUITE_P(StateTest);

TYPED_TEST_P(StateTest, DefaultAccountStateIsUnknown) {
  Address a{0x01};
  Address b{0x02};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kUnknown));
  EXPECT_THAT(state.GetAccountState(b), IsOkAndHolds(AccountState::kUnknown));
}

TYPED_TEST_P(StateTest, AccountsCanBeCreatedAndAreDifferentiated) {
  Address a{0x01};
  Address b{0x02};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kUnknown));
  EXPECT_THAT(state.GetAccountState(b), IsOkAndHolds(AccountState::kUnknown));

  EXPECT_OK(state.CreateAccount(a));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kExists));
  EXPECT_THAT(state.GetAccountState(b), IsOkAndHolds(AccountState::kUnknown));

  EXPECT_OK(state.CreateAccount(b));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kExists));
  EXPECT_THAT(state.GetAccountState(b), IsOkAndHolds(AccountState::kExists));
}

TYPED_TEST_P(StateTest, CreatingAnAccountDeletesItsStorage) {
  Address a{0x01};
  Key k{0x01, 0x02};
  Value v{0x02, 0x03, 0x04};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));

  // Initially, the storage is empty, but can be written to.
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(Value{}));
  EXPECT_OK(state.SetStorageValue(a, k, v));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(v));

  // The account creation purges the storage.
  EXPECT_OK(state.CreateAccount(a));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(Value{}));
  EXPECT_OK(state.SetStorageValue(a, k, v));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(v));

  // At this point the account is re-created, storage should still be purged.
  EXPECT_OK(state.CreateAccount(a));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(Value{}));
}

TYPED_TEST_P(StateTest, AccountsCanBeDeleted) {
  Address a{0x01};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kUnknown));

  EXPECT_OK(state.CreateAccount(a));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kExists));

  EXPECT_OK(state.DeleteAccount(a));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kUnknown));
}

TYPED_TEST_P(StateTest, DeletingAnUnknownAccountDoesNotCreateIt) {
  Address a{0x01};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kUnknown));

  EXPECT_OK(state.DeleteAccount(a));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kUnknown));
}

TYPED_TEST_P(StateTest, DeletedAccountsCanBeRecreated) {
  Address a{0x01};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kUnknown));
  EXPECT_OK(state.CreateAccount(a));
  EXPECT_OK(state.DeleteAccount(a));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kUnknown));
  EXPECT_OK(state.CreateAccount(a));
  EXPECT_THAT(state.GetAccountState(a), IsOkAndHolds(AccountState::kExists));
}

TYPED_TEST_P(StateTest, DeletingAnAccountDeletesItsStorage) {
  Address a{0x01};
  Key k{0x01, 0x02};
  Value v{0x02, 0x03, 0x04};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));

  EXPECT_OK(state.CreateAccount(a));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(Value{}));
  EXPECT_OK(state.SetStorageValue(a, k, v));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(v));

  EXPECT_OK(state.DeleteAccount(a));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(Value{}));
}

TYPED_TEST_P(StateTest, DefaultBalanceIsZero) {
  Address a{0x01};
  Address b{0x02};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetBalance(a), IsOkAndHolds(Balance{}));
  EXPECT_THAT(state.GetBalance(b), IsOkAndHolds(Balance{}));
}

TYPED_TEST_P(StateTest, BalancesCanBeUpdated) {
  Address a{0x01};
  Address b{0x02};
  Balance zero{};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetBalance(a), IsOkAndHolds(zero));
  EXPECT_THAT(state.GetBalance(b), IsOkAndHolds(zero));

  EXPECT_OK(state.SetBalance(a, Balance{0x12}));
  EXPECT_THAT(state.GetBalance(a), IsOkAndHolds(Balance{0x12}));
  EXPECT_THAT(state.GetBalance(b), IsOkAndHolds(zero));

  EXPECT_OK(state.SetBalance(b, Balance{0x14}));
  EXPECT_THAT(state.GetBalance(a), IsOkAndHolds(Balance{0x12}));
  EXPECT_THAT(state.GetBalance(b), IsOkAndHolds(Balance{0x14}));
}

TYPED_TEST_P(StateTest, BalancesAreCoveredByGlobalStateHash) {
  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  ASSERT_OK_AND_ASSIGN(auto base_hash, state.GetHash());
  EXPECT_OK(state.SetBalance({}, Balance{0x12}));
  ASSERT_OK_AND_ASSIGN(auto value_12_hash, state.GetHash());
  EXPECT_NE(base_hash, value_12_hash);
  EXPECT_OK(state.SetBalance({}, Balance{0x14}));
  ASSERT_OK_AND_ASSIGN(auto value_14_hash, state.GetHash());
  EXPECT_NE(base_hash, value_14_hash);

  // Resetting value gets us original hash.
  EXPECT_OK(state.SetBalance({}, Balance{0x12}));
  EXPECT_THAT(state.GetHash(), IsOkAndHolds(value_12_hash));
}

TYPED_TEST_P(StateTest, DefaultNonceIsZero) {
  Address a{0x01};
  Address b{0x02};
  Nonce zero{};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetNonce(a), IsOkAndHolds(zero));
  EXPECT_THAT(state.GetNonce(b), IsOkAndHolds(zero));
}

TYPED_TEST_P(StateTest, NoncesCanBeUpdated) {
  Address a{0x01};
  Address b{0x02};
  Nonce zero{};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetNonce(a), IsOkAndHolds(zero));
  EXPECT_THAT(state.GetNonce(b), IsOkAndHolds(zero));

  EXPECT_OK(state.SetNonce(a, Nonce{0x12}));
  EXPECT_THAT(state.GetNonce(a), IsOkAndHolds(Nonce{0x12}));
  EXPECT_THAT(state.GetNonce(b), IsOkAndHolds(zero));

  EXPECT_OK(state.SetNonce(b, Nonce{0x14}));
  EXPECT_THAT(state.GetNonce(a), IsOkAndHolds(Nonce{0x12}));
  EXPECT_THAT(state.GetNonce(b), IsOkAndHolds(Nonce{0x14}));
}

TYPED_TEST_P(StateTest, NoncesAreCoveredByGlobalStateHash) {
  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  ASSERT_OK_AND_ASSIGN(auto base_hash, state.GetHash());
  EXPECT_OK(state.SetNonce({}, Nonce{0x12}));
  ASSERT_OK_AND_ASSIGN(auto value_12_hash, state.GetHash());
  EXPECT_NE(base_hash, value_12_hash);
  EXPECT_OK(state.SetNonce({}, Nonce{0x14}));
  ASSERT_OK_AND_ASSIGN(auto value_14_hash, state.GetHash());
  EXPECT_NE(base_hash, value_14_hash);

  // Resetting value gets us original hash.
  EXPECT_OK(state.SetNonce({}, Nonce{0x12}));
  EXPECT_THAT(state.GetHash(), IsOkAndHolds(value_12_hash));
}

TYPED_TEST_P(StateTest, DefaultCodeIsEmpty) {
  Address a{0x01};
  Address b{0x02};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetCode(a), IsOkAndHolds(ElementsAre()));
  EXPECT_THAT(state.GetCode(b), IsOkAndHolds(ElementsAre()));
}

TYPED_TEST_P(StateTest, CodesCanBeUpdated) {
  Address a{0x01};
  Address b{0x02};
  std::vector<std::byte> code1{std::byte{1}, std::byte{2}};
  std::vector<std::byte> code2{std::byte{3}, std::byte{4}};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetCode(a), IsOkAndHolds(ElementsAre()));
  EXPECT_THAT(state.GetCode(b), IsOkAndHolds(ElementsAre()));

  EXPECT_OK(state.SetCode(a, code1));
  EXPECT_THAT(state.GetCode(a), IsOkAndHolds(ElementsAreArray(code1)));
  EXPECT_THAT(state.GetCode(b), IsOkAndHolds(ElementsAre()));

  EXPECT_OK(state.SetCode(b, code2));
  EXPECT_THAT(state.GetCode(a), IsOkAndHolds(ElementsAreArray(code1)));
  EXPECT_THAT(state.GetCode(b), IsOkAndHolds(ElementsAreArray(code2)));

  EXPECT_OK(state.SetCode(a, code2));
  EXPECT_THAT(state.GetCode(a), IsOkAndHolds(ElementsAreArray(code2)));
  EXPECT_THAT(state.GetCode(b), IsOkAndHolds(ElementsAreArray(code2)));
}

TYPED_TEST_P(StateTest, UpdatingCodesUpdatesCodeHashes) {
  const Hash hash_of_empty_code = GetKeccak256Hash({});

  Address a{0x01};
  std::vector<std::byte> code{std::byte{1}, std::byte{2}};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_THAT(state.GetCodeHash(a), IsOkAndHolds(hash_of_empty_code));

  EXPECT_OK(state.SetCode(a, code));
  EXPECT_THAT(state.GetCodeHash(a),
              IsOkAndHolds(GetKeccak256Hash(std::span(code))));

  // Resetting code to zero updates the hash accordingly.
  EXPECT_OK(state.SetCode(a, {}));
  EXPECT_THAT(state.GetCodeHash(a), IsOkAndHolds(hash_of_empty_code));
}

TYPED_TEST_P(StateTest, CodesAreCoveredByGlobalStateHash) {
  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  ASSERT_OK_AND_ASSIGN(auto base_hash, state.GetHash());
  EXPECT_OK(state.SetCode({}, std::vector{std::byte{12}}));
  ASSERT_OK_AND_ASSIGN(auto value_12_hash, state.GetHash());
  EXPECT_NE(base_hash, value_12_hash);
  EXPECT_OK(state.SetCode({}, std::vector{std::byte{14}}));
  ASSERT_OK_AND_ASSIGN(auto value_14_hash, state.GetHash());
  EXPECT_NE(base_hash, value_14_hash);

  // Resetting value gets us original hash.
  EXPECT_OK(state.SetCode({}, std::vector{std::byte{12}}));
  EXPECT_THAT(state.GetHash(), IsOkAndHolds(value_12_hash));
}

TYPED_TEST_P(StateTest, LookingUpMissingCodeDoesNotChangeGlobalHash) {
  Address a{0x01};
  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  ASSERT_OK_AND_ASSIGN(auto base_hash, state.GetHash());
  EXPECT_OK(state.GetCode(a));
  EXPECT_THAT(state.GetHash(), IsOkAndHolds(base_hash));
}

TYPED_TEST_P(StateTest, ValuesAddedCanBeRetrieved) {
  Address a;
  Key k;
  Value v{0x01, 0x02};

  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_OK(state.SetStorageValue(a, k, v));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(v));

  v = Value{0x03};
  EXPECT_OK(state.SetStorageValue(a, k, v));
  EXPECT_THAT(state.GetStorageValue(a, k), IsOkAndHolds(v));
}

TYPED_TEST_P(StateTest, UpdatesCanBeApplied) {
  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  EXPECT_OK(state.CreateAccount(Address{0x02}));

  Update update;
  update.Create(Address{0x01});
  update.Delete(Address{0x02});
  update.Set(Address{0x03}, Balance{0xB1});
  update.Set(Address{0x04}, Nonce{0xA1});
  update.Set(Address{0x05}, Key{0x06}, Value{0x07});
  update.Set(Address{0x06}, Code{0x01, 0x02});

  EXPECT_OK(state.Apply(12, update));

  EXPECT_THAT(state.GetAccountState(Address{0x00}), AccountState::kUnknown);
  EXPECT_THAT(state.GetAccountState(Address{0x01}), AccountState::kExists);
  EXPECT_THAT(state.GetAccountState(Address{0x02}), AccountState::kUnknown);

  EXPECT_THAT(state.GetBalance(Address{0x03}), IsOkAndHolds(Balance{0xB1}));
  EXPECT_THAT(state.GetNonce(Address{0x04}), IsOkAndHolds(Nonce{0xA1}));
  EXPECT_THAT(state.GetStorageValue(Address{0x05}, Key{0x06}),
              IsOkAndHolds(Value{0x07}));
  EXPECT_THAT(state.GetCode(Address{0x06}),
              IsOkAndHolds(ElementsAre(std::byte{0x01}, std::byte{0x02})));
}

TYPED_TEST_P(StateTest, UpdatesCanBeAppliedWithArchive) {
  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir, /*with_archive=*/true));
  EXPECT_OK(state.CreateAccount(Address{0x02}));

  Update update;
  update.Create(Address{0x01});
  update.Delete(Address{0x02});
  update.Set(Address{0x03}, Balance{0xB1});
  update.Set(Address{0x04}, Nonce{0xA1});
  update.Set(Address{0x05}, Key{0x06}, Value{0x07});
  update.Set(Address{0x06}, Code{0x01, 0x02});

  EXPECT_OK(state.Apply(12, update));

  // TODO: once history is accessible, check results.
}

TYPED_TEST_P(StateTest, CanProduceAMemoryFootprint) {
  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir));
  auto usage = state.GetMemoryFootprint();
  EXPECT_GT(usage.GetTotal(), Memory());
}

TYPED_TEST_P(StateTest, CanBeOpenedWithArchive) {
  TempDir dir;
  ASSERT_OK_AND_ASSIGN(auto state, TypeParam::Open(dir, /*with_archive=*/true));
}

REGISTER_TYPED_TEST_SUITE_P(
    StateTest, AccountsCanBeDeleted, AccountsCanBeCreatedAndAreDifferentiated,
    CreatingAnAccountDeletesItsStorage, BalancesAreCoveredByGlobalStateHash,
    BalancesCanBeUpdated, CodesAreCoveredByGlobalStateHash, CodesCanBeUpdated,
    DefaultAccountStateIsUnknown, DefaultBalanceIsZero, DefaultCodeIsEmpty,
    DefaultNonceIsZero, DeletedAccountsCanBeRecreated,
    DeletingAnAccountDeletesItsStorage, DeletingAnUnknownAccountDoesNotCreateIt,
    LookingUpMissingCodeDoesNotChangeGlobalHash,
    NoncesAreCoveredByGlobalStateHash, NoncesCanBeUpdated,
    UpdatingCodesUpdatesCodeHashes, ValuesAddedCanBeRetrieved,
    UpdatesCanBeApplied, UpdatesCanBeAppliedWithArchive,
    CanProduceAMemoryFootprint, CanBeOpenedWithArchive);

using StateConfigurations =
    ::testing::Types<InMemoryState, FileBasedState, LevelDbBasedState>;

INSTANTIATE_TYPED_TEST_SUITE_P(Config, StateTest, StateConfigurations);

}  // namespace carmen
