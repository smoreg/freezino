import { create } from 'zustand';
import api from '../services/api';
import type { Item, UserItem, ItemType, ItemRarity } from '../types';

interface ShopState {
  items: Item[];
  myItems: UserItem[];
  isLoading: boolean;
  error: string | null;

  // Filters
  filterType: ItemType | 'all';
  filterRarity: ItemRarity | 'all';
  minPrice: number;
  maxPrice: number;

  // Actions
  fetchItems: () => Promise<void>;
  fetchMyItems: () => Promise<void>;
  buyItem: (itemId: string) => Promise<void>;
  sellItem: (itemId: string) => Promise<void>;
  equipItem: (itemId: string) => Promise<void>;
  setFilterType: (type: ItemType | 'all') => void;
  setFilterRarity: (rarity: ItemRarity | 'all') => void;
  setPriceRange: (min: number, max: number) => void;
  resetFilters: () => void;
}

export const useShopStore = create<ShopState>((set, get) => ({
  items: [],
  myItems: [],
  isLoading: false,
  error: null,

  // Default filter values
  filterType: 'all',
  filterRarity: 'all',
  minPrice: 0,
  maxPrice: 1000000,

  fetchItems: async () => {
    set({ isLoading: true, error: null });

    try {
      const { filterType, filterRarity } = get();

      // Build query params
      const params = new URLSearchParams();
      if (filterType !== 'all') params.append('type', filterType);
      if (filterRarity !== 'all') params.append('rarity', filterRarity);

      const response = await api.get<Item[]>(`/shop/items?${params.toString()}`);

      set({
        items: response.data,
        isLoading: false
      });
    } catch (error: unknown) {
      set({
        error: (error as { response?: { data?: { message?: string } } }).response?.data?.message || 'Failed to fetch items',
        isLoading: false
      });
    }
  },

  fetchMyItems: async () => {
    try {
      const response = await api.get<UserItem[]>('/shop/my-items');
      set({ myItems: response.data });
    } catch (error: unknown) {
      console.error('Failed to fetch my items:', error);
    }
  },

  buyItem: async (itemId: string) => {
    try {
      await api.post(`/shop/buy/${itemId}`);

      // Refresh items and user items
      await get().fetchItems();
      await get().fetchMyItems();
    } catch (error: unknown) {
      throw new Error((error as { response?: { data?: { message?: string } } }).response?.data?.message || 'Failed to buy item');
    }
  },

  sellItem: async (itemId: string) => {
    try {
      await api.post(`/shop/sell/${itemId}`);

      // Refresh my items
      await get().fetchMyItems();
    } catch (error: unknown) {
      throw new Error((error as { response?: { data?: { message?: string } } }).response?.data?.message || 'Failed to sell item');
    }
  },

  equipItem: async (itemId: string) => {
    try {
      await api.post(`/shop/equip/${itemId}`);

      // Refresh my items
      await get().fetchMyItems();
    } catch (error: unknown) {
      throw new Error((error as { response?: { data?: { message?: string } } }).response?.data?.message || 'Failed to equip item');
    }
  },

  setFilterType: (type) => {
    set({ filterType: type });
    get().fetchItems();
  },

  setFilterRarity: (rarity) => {
    set({ filterRarity: rarity });
    get().fetchItems();
  },

  setPriceRange: (min, max) => {
    set({ minPrice: min, maxPrice: max });
  },

  resetFilters: () => {
    set({
      filterType: 'all',
      filterRarity: 'all',
      minPrice: 0,
      maxPrice: 1000000,
    });
    get().fetchItems();
  },
}));
