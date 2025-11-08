import { motion } from 'framer-motion';
import { useShopStore } from '../../store/shopStore';
import type { ItemType, ItemRarity } from '../../types';

const itemTypes: { value: ItemType | 'all'; label: string; emoji: string }[] = [
  { value: 'all', label: 'All', emoji: '🛍️' },
  { value: 'clothing', label: 'Clothing', emoji: '👔' },
  { value: 'car', label: 'Cars', emoji: '🚗' },
  { value: 'house', label: 'Houses', emoji: '🏠' },
  { value: 'accessory', label: 'Accessories', emoji: '💎' },
];

const rarities: { value: ItemRarity | 'all'; label: string; color: string }[] = [
  { value: 'all', label: 'All Rarities', color: 'text-gray-400' },
  { value: 'common', label: 'Common', color: 'text-gray-400' },
  { value: 'rare', label: 'Rare', color: 'text-blue-400' },
  { value: 'epic', label: 'Epic', color: 'text-purple-400' },
  { value: 'legendary', label: 'Legendary', color: 'text-yellow-400' },
];

export default function ShopFilters() {
  const { filterType, filterRarity, setFilterType, setFilterRarity, resetFilters } = useShopStore();

  return (
    <div className="bg-gray-800 rounded-lg p-6 shadow-xl mb-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
          <span>🔍</span>
          Filters
        </h2>

        <button
          onClick={resetFilters}
          className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-sm transition-colors"
        >
          Reset
        </button>
      </div>

      {/* Type Filter */}
      <div className="mb-6">
        <h3 className="text-sm font-semibold text-gray-400 uppercase mb-3">Category</h3>
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
          {itemTypes.map((type) => (
            <motion.button
              key={type.value}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={() => setFilterType(type.value)}
              className={`flex flex-col items-center gap-2 p-4 rounded-lg border-2 transition-all ${
                filterType === type.value
                  ? 'bg-yellow-600 border-yellow-500 text-white shadow-lg shadow-yellow-500/50'
                  : 'bg-gray-700 border-gray-600 text-gray-300 hover:bg-gray-600 hover:border-gray-500'
              }`}
            >
              <span className="text-3xl">{type.emoji}</span>
              <span className="text-sm font-semibold">{type.label}</span>
            </motion.button>
          ))}
        </div>
      </div>

      {/* Rarity Filter */}
      <div>
        <h3 className="text-sm font-semibold text-gray-400 uppercase mb-3">Rarity</h3>
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
          {rarities.map((rarity) => (
            <motion.button
              key={rarity.value}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={() => setFilterRarity(rarity.value)}
              className={`p-3 rounded-lg border-2 transition-all text-center ${
                filterRarity === rarity.value
                  ? 'bg-yellow-600 border-yellow-500 text-white shadow-lg shadow-yellow-500/50'
                  : 'bg-gray-700 border-gray-600 hover:bg-gray-600 hover:border-gray-500'
              }`}
            >
              <span className={`text-sm font-semibold ${filterRarity === rarity.value ? 'text-white' : rarity.color}`}>
                {rarity.label}
              </span>
            </motion.button>
          ))}
        </div>
      </div>
    </div>
  );
}
