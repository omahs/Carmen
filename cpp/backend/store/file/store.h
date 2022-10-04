#pragma once

#include "backend/store/file/file.h"
#include "backend/store/file/hash_tree.h"
#include "backend/store/file/page_pool.h"
#include "common/hash.h"
#include "common/type.h"

namespace carmen::backend::store {

// The FileStore is a file-backed implementation of a mutable key/value store.
// It provides mutation, lookup, and global state hashing support.
template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size = 32>
requires File<F<page_size>, page_size>
class FileStore {
 public:
  // Creates a new, empty FileStore.
  FileStore();

  // Updates the value associated to the given key.
  void Set(const K& key, V value);

  // Retrieves the value associated to the given key. If no values has
  // been previously set using the Set(..) function above, a zero-initialized
  // value is returned. The returned reference is only valid until the next
  // operation on the store.
  const V& Get(const K& key) const;

  // Computes a hash over the full content of this store.
  Hash GetHash() const;

 private:
  using PagePool = PagePool<V, F, page_size>;

  // A listener to pool activities to react to loaded and evicted pages and
  // perform necessary hashing steps.
  class PoolListener : public PagePoolListener<V, page_size> {
   public:
    PoolListener(FileStore& store) : store_(store) {}

    void AfterLoad(PageId id, const Page<V, page_size>&) override {
      // When a page is loaded, make sure the HashTree is aware of it.
      store_.hashes_.RegisterPage(id);
    }

    void BeforeEvict(PageId id, const Page<V, page_size>& page,
                     bool is_dirty) override {
      // Before we throw away a dirty page to make space for something else we
      // update the hash to avoid having to reload it again later.
      if (is_dirty) {
        store_.hashes_.UpdateHash(id, page.AsRawData());
      }
    }

   private:
    FileStore& store_;
  };

  // An implementation of a PageSource passed to the HashTree to provide access
  // to pages through the page pool, and thus through its caching authority.
  class PageProvider : public PageSource {
   public:
    PageProvider(FileStore& store) : store_(store) {}

    std::span<const std::byte> GetPageData(PageId id) override {
      return store_.pool_.Get(id).AsRawData();
    }

   private:
    FileStore& store_;
  };

  // The number of elements per page, used for page and offset computaiton.
  constexpr static std::size_t kNumElementsPerPage =
      PagePool::Page::kNumElementsPerPage;

  // The page pool handling the in-memory buffer of pages fetched from disk.
  mutable PagePool pool_;

  // The data structure hanaging the hashing of states.
  mutable HashTree hashes_;
};

template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size>
requires File<F<page_size>, page_size>
FileStore<K, V, F, page_size>::FileStore()
    : hashes_(std::make_unique<PageProvider>(*this)) {
  pool_.AddListener(std::make_unique<PoolListener>(*this));
}

template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size>
requires File<F<page_size>, page_size>
void FileStore<K, V, F, page_size>::Set(const K& key, V value) {
  auto& trg = pool_.Get(key / kNumElementsPerPage)[key % kNumElementsPerPage];
  if (trg != value) {
    trg = value;
    pool_.MarkAsDirty(key / kNumElementsPerPage);
    hashes_.MarkDirty(key / kNumElementsPerPage);
  }
}

template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size>
requires File<F<page_size>, page_size>
const V& FileStore<K, V, F, page_size>::Get(const K& key) const {
  return pool_.Get(key / kNumElementsPerPage)[key % kNumElementsPerPage];
}

template <typename K, Trivial V, template <std::size_t> class F,
          std::size_t page_size>
requires File<F<page_size>, page_size> Hash
FileStore<K, V, F, page_size>::GetHash()
const { return hashes_.GetHash(); }

}  // namespace carmen::backend::store
