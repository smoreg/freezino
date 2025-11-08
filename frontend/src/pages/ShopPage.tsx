import { useEffect, useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import { useShopStore } from '../store/shopStore';
import { useAuthStore } from '../store/authStore';
import ItemCard from '../components/shop/ItemCard';
import ShopFilters from '../components/shop/ShopFilters';
import BuyModal from '../components/shop/BuyModal';
import type { Item } from '../types';
import { PageTransition, rotateVariants } from '../components/animations';

export default function ShopPage() {
  const { items, myItems, isLoading, fetchItems, fetchMyItems, minPrice, maxPrice } = useShopStore();
  const user = useAuthStore((state) => state.user);

  const [selectedItem, setSelectedItem] = useState<Item | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Fetch items on mount
  useEffect(() => {
    fetchItems();
    fetchMyItems();
  }, [fetchItems, fetchMyItems]);

  // Filter items by price range
  const filteredItems = useMemo(() => {
    return items.filter((item) => item.price >= minPrice && item.price <= maxPrice);
  }, [items, minPrice, maxPrice]);

  // Check if item is owned
  const isItemOwned = (itemId: string) => {
    return myItems.some((userItem) => userItem.item_id === itemId);
  };

  const handleBuyClick = (item: Item) => {
    setSelectedItem(item);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedItem(null);
  };

  return (
    <PageTransition>
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 py-8 px-4">
        <div className="max-w-7xl mx-auto">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-2 flex items-center gap-3">
            <span>🛍️</span>
            Shop
          </h1>
          <p className="text-gray-400 text-lg">
            Spend your virtual dollars on exclusive items
          </p>

          {/* Balance Display */}
          <div className="mt-4 inline-block bg-gradient-to-r from-yellow-600 to-yellow-500 rounded-lg px-6 py-3 shadow-lg">
            <div className="flex items-center gap-2">
              <span className="text-white font-semibold">Your Balance:</span>
              <span className="text-white text-2xl font-bold">${user?.balance.toLocaleString()}</span>
            </div>
          </div>
        </motion.div>

        {/* Filters */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <ShopFilters />
        </motion.div>

        {/* Items Grid */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
        >
          {isLoading ? (
            <div className="flex items-center justify-center py-20">
              <div className="text-center">
                <motion.div
                  className="inline-block rounded-full h-16 w-16 border-t-4 border-b-4 border-yellow-500 mb-4"
                  variants={rotateVariants}
                  initial="initial"
                  animate="animate"
                />
                <p className="text-gray-400 text-lg">Loading items...</p>
              </div>
            </div>
          ) : filteredItems.length === 0 ? (
            <div className="text-center py-20">
              <div className="text-6xl mb-4">📦</div>
              <h3 className="text-2xl font-bold text-white mb-2">No items found</h3>
              <p className="text-gray-400">Try adjusting your filters</p>
            </div>
          ) : (
            <>
              {/* Results Count */}
              <div className="mb-4">
                <p className="text-gray-400">
                  Showing <span className="text-white font-semibold">{filteredItems.length}</span> items
                </p>
              </div>

              {/* Items Grid */}
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                {filteredItems.map((item, index) => (
                  <motion.div
                    key={item.id}
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: index * 0.05, duration: 0.3 }}
                  >
                    <ItemCard
                      item={item}
                      onBuy={handleBuyClick}
                      owned={isItemOwned(item.id)}
                    />
                  </motion.div>
                ))}
              </div>
            </>
          )}
        </motion.div>

        {/* Buy Modal */}
        <BuyModal item={selectedItem} isOpen={isModalOpen} onClose={handleCloseModal} />
      </div>
      </div>
    </PageTransition>
  );
}
