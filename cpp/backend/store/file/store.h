#pragma once

#include <filesystem>

#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "backend/common/file.h"
#include "backend/common/page.h"
#include "backend/common/page_pool.h"
#include "backend/store/hash_tree.h"
#include "backend/structure.h"
#include "common/hash.h"
#include "common/status_util.h"
#include "common/type.h"

namespace carmen::backend::store {

// ----------------------------------------------------------------------------
//                              Declarations
// ----------------------------------------------------------------------------

namespace internal {

// The FileStoreBase is the common bases of file-backed implementations of a
// mutable key/value store. It provides mutation, lookup, and global state
// hashing support. Hashing can occure eager (before evicting pages) or lazy,
// when requesting hash computations.
template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size = 32, bool eager_hashing = true>
requires File<F<sizeof(ArrayPage<V, page_size / sizeof(V)>)>>
class FileStoreBase;

}  // namespace internal

// A FileStore implementation configured to perform eager hashing. Thus,
// before pages are evicted, hashes are computed. This slows down reads
// and updates, but improves hashing speed.
template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size = 32>
requires File<F<sizeof(ArrayPage<V, page_size / sizeof(V)>)>>
using EagerFileStore = internal::FileStoreBase<K, V, F, page_size, true>;

// A FileStore implementation configured to perform lazy hashing. Thus,
// pages are evicted without being hashes and need to be reloaded for computing
// hashes when needed. This speeds up read/write operations at the expense of
// hash performance.
template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size = 32>
requires File<F<sizeof(ArrayPage<V, page_size / sizeof(V)>)>>
using LazyFileStore = internal::FileStoreBase<K, V, F, page_size, false>;

// ----------------------------------------------------------------------------
//                              Definitions
// ----------------------------------------------------------------------------

namespace internal {

template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size, bool eager_hashing>
requires File<F<sizeof(ArrayPage<V, page_size / sizeof(V)>)>>
class FileStoreBase {
 public:
  // The value type used to index elements in this store.
  using key_type = K;

  // The type of value stored in this store.
  using value_type = V;

  // The page size in byte used by this store as configured. This may be less
  // then the actual page size, which may be larger due to alignment and padding
  // constraints.
  constexpr static std::size_t kPageSize = page_size;

  // A factory function creating an instance of this store type.
  static absl::StatusOr<FileStoreBase> Open(
      Context&, const std::filesystem::path& directory,
      std::size_t hash_branching_factor = 32) {
    // Make sure the directory exists.
    if (!CreateDirectory(directory)) {
      return absl::InternalError(
          absl::StrFormat("Unable to create parent directory %s", directory));
    }
    auto store = FileStoreBase(directory, hash_branching_factor);
    if (std::filesystem::exists(store.hash_file_)) {
      RETURN_IF_ERROR(store.hashes_->LoadFromFile(store.hash_file_));
    }
    return store;
  }

  // Supports instances to be moved.
  FileStoreBase(FileStoreBase&&) = default;

  // File stores are automatically closed on destruction.
  ~FileStoreBase() { Close().IgnoreError(); }

  // Updates the value associated to the given key.
  absl::Status Set(const K& key, V value);

  // Retrieves the value associated to the given key. If no values has
  // been previously set using the Set(..) function above, a zero-initialized
  // value is returned. The returned reference is only valid until the next
  // operation on the store.
  StatusOrRef<const V> Get(const K& key) const;

  // Computes a hash over the full content of this store.
  absl::StatusOr<Hash> GetHash() const;

  // Flushes internally buffered modified data to disk.
  absl::Status Flush();

  // Flushes the store and closes resource references.
  absl::Status Close();

  // Summarizes the memory usage of this instance.
  MemoryFootprint GetMemoryFootprint() const;

 private:
  using Page = ArrayPage<V, page_size / sizeof(V)>;
  using PagePool = PagePool<F<sizeof(Page)>>;

  // The actual size of a page, which may be larger than the specified page size
  // due to padding.
  constexpr static std::size_t kFilePageSize = sizeof(Page);

  // A listener to pool activities to react to loaded and evicted pages and
  // perform necessary hashing steps.
  class PoolListener : public PagePoolListener<kFilePageSize> {
   public:
    PoolListener(HashTree& hashes) : hashes_(hashes) {}

    void AfterLoad(PageId id, const RawPage<kFilePageSize>&) override {
      // When a page is loaded, make sure the HashTree is aware of it.
      hashes_.RegisterPage(id);
    }

    void BeforeEvict(PageId id, const RawPage<kFilePageSize>& page,
                     bool is_dirty) override {
      // Before we throw away a dirty page to make space for something else we
      // update the hash to avoid having to reload it again later.
      if (eager_hashing && is_dirty) {
        hashes_.UpdateHash(
            id, std::as_bytes(std::span(page.template As<Page>().AsArray())));
      }
    }

   private:
    HashTree& hashes_;
  };

  // An implementation of a PageSource passed to the HashTree to provide access
  // to pages through the page pool, and thus through its caching authority.
  class PageProvider : public PageSource {
   public:
    PageProvider(PagePool& pool) : pool_(pool) {}

    std::span<const std::byte> GetPageData(PageId id) override {
      return std::as_bytes(std::span(pool_.template Get<Page>(id).AsArray()));
    }

   private:
    PagePool& pool_;
  };

  // The number of elements per page, used for page and offset computaiton.
  constexpr static std::size_t kNumElementsPerPage = Page::kNumElementsPerPage;

  // Creates a new file store maintaining its content in the given directory and
  // using the provided branching factor for its hash computation.
  FileStoreBase(std::filesystem::path directory,
                std::size_t hash_branching_factor);

  // The page pool handling the in-memory buffer of pages fetched from disk. The
  // pool is placed in a unique pointer to ensure pointer stability when the
  // store is moved.
  mutable std::unique_ptr<PagePool> pool_;

  // The data structure hanaging the hashing of states. The hashes are placed in
  // a unique pointer to ensure pointer stability when the store is moved.
  mutable std::unique_ptr<HashTree> hashes_;

  // The name of the file to safe hashes to.
  std::filesystem::path hash_file_;
};

// Since the template parameter declaration list and concept requirements to
// preceed the subsequent member function declarations is rather extensive, a
// macro covering it defined to aid readability.
#define FILE_STORE_TEMPLATE_PARAMETER                              \
  template <typename K, Trivial V, template <std::size_t> class F, \
            std::size_t page_size, bool eager_hashing>             \
  requires File<F<sizeof(ArrayPage<V, page_size / sizeof(V)>)>>

FILE_STORE_TEMPLATE_PARAMETER
FileStoreBase<K, V, F, page_size, eager_hashing>::FileStoreBase(
    std::filesystem::path directory, std::size_t hash_branching_factor)
    : pool_(std::make_unique<PagePool>(
          std::make_unique<F<kFilePageSize>>(directory / "data.dat"))),
      hashes_(std::make_unique<HashTree>(std::make_unique<PageProvider>(*pool_),
                                         hash_branching_factor)),
      hash_file_(directory / "hash.dat") {
  pool_->AddListener(std::make_unique<PoolListener>(*hashes_));
}

FILE_STORE_TEMPLATE_PARAMETER absl::Status
FileStoreBase<K, V, F, page_size, eager_hashing>::Set(const K& key, V value) {
  auto& trg = pool_->template Get<Page>(
      key / kNumElementsPerPage)[key % kNumElementsPerPage];
  if (trg != value) {
    trg = value;
    pool_->MarkAsDirty(key / kNumElementsPerPage);
    hashes_->MarkDirty(key / kNumElementsPerPage);
  }
  return absl::OkStatus();
}

FILE_STORE_TEMPLATE_PARAMETER
StatusOrRef<const V> FileStoreBase<K, V, F, page_size, eager_hashing>::Get(
    const K& key) const {
  return pool_->template Get<Page>(
      key / kNumElementsPerPage)[key % kNumElementsPerPage];
}

FILE_STORE_TEMPLATE_PARAMETER
absl::StatusOr<Hash> FileStoreBase<K, V, F, page_size, eager_hashing>::GetHash()
    const {
  return hashes_->GetHash();
}

FILE_STORE_TEMPLATE_PARAMETER absl::Status
FileStoreBase<K, V, F, page_size, eager_hashing>::Flush() {
  if (pool_) pool_->Flush();
  if (hashes_) {
    RETURN_IF_ERROR(hashes_->SaveToFile(hash_file_));
  }
  return absl::OkStatus();
}

FILE_STORE_TEMPLATE_PARAMETER absl::Status
FileStoreBase<K, V, F, page_size, eager_hashing>::Close() {
  RETURN_IF_ERROR(Flush());
  if (pool_) pool_->Close();
  return absl::OkStatus();
}

FILE_STORE_TEMPLATE_PARAMETER MemoryFootprint
FileStoreBase<K, V, F, page_size, eager_hashing>::GetMemoryFootprint() const {
  MemoryFootprint res(*this);
  res.Add("pool", pool_->GetMemoryFootprint());
  res.Add("hashes", hashes_->GetMemoryFootprint());
  return res;
}

#undef FILE_STORE_TEMPLATE_PARAMETER

}  // namespace internal
}  // namespace carmen::backend::store
