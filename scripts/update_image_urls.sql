-- Update Image URLs for All Shop Items
-- This script updates all placeholder image URLs to local image paths

-- CLOTHING - Common
UPDATE items SET image_url = '/images/clothing/plain-tshirt.jpg' WHERE name = 'Plain T-Shirt';
UPDATE items SET image_url = '/images/clothing/casual-jeans.jpg' WHERE name = 'Casual Jeans';
UPDATE items SET image_url = '/images/clothing/sneakers.jpg' WHERE name = 'Sneakers';
UPDATE items SET image_url = '/images/clothing/hoodie.jpg' WHERE name = 'Hoodie';
UPDATE items SET image_url = '/images/clothing/designer-shirt.jpg' WHERE name = 'Designer Shirt';

-- CLOTHING - Rare
UPDATE items SET image_url = '/images/clothing/leather-jacket.jpg' WHERE name = 'Leather Jacket';
UPDATE items SET image_url = '/images/clothing/designer-dress.jpg' WHERE name = 'Designer Dress';
UPDATE items SET image_url = '/images/clothing/business-suit.jpg' WHERE name = 'Business Suit';
UPDATE items SET image_url = '/images/clothing/evening-gown.jpg' WHERE name = 'Evening Gown';
UPDATE items SET image_url = '/images/clothing/tuxedo.jpg' WHERE name = 'Tuxedo';
UPDATE items SET image_url = '/images/clothing/designer-coat.jpg' WHERE name = 'Designer Coat';

-- CLOTHING - Epic
UPDATE items SET image_url = '/images/clothing/custom-suit.jpg' WHERE name = 'Custom Tailored Suit';
UPDATE items SET image_url = '/images/clothing/haute-couture.jpg' WHERE name = 'Haute Couture Dress';
UPDATE items SET image_url = '/images/clothing/fur-coat.jpg' WHERE name = 'Luxury Fur Coat';

-- CLOTHING - Legendary
UPDATE items SET image_url = '/images/clothing/limited-edition.jpg' WHERE name = 'Limited Edition Designer Collection';

-- ACCESSORIES - Common
UPDATE items SET image_url = '/images/accessories/sunglasses.jpg' WHERE name = 'Sunglasses';
UPDATE items SET image_url = '/images/accessories/wallet.jpg' WHERE name = 'Leather Wallet';
UPDATE items SET image_url = '/images/accessories/casual-watch.jpg' WHERE name = 'Casual Watch';
UPDATE items SET image_url = '/images/accessories/belt.jpg' WHERE name = 'Belt';
UPDATE items SET image_url = '/images/accessories/backpack.jpg' WHERE name = 'Backpack';

-- ACCESSORIES - Rare
UPDATE items SET image_url = '/images/accessories/designer-handbag.jpg' WHERE name = 'Designer Handbag';
UPDATE items SET image_url = '/images/accessories/gold-necklace.jpg' WHERE name = 'Gold Necklace';
UPDATE items SET image_url = '/images/accessories/briefcase.jpg' WHERE name = 'Designer Briefcase';
UPDATE items SET image_url = '/images/accessories/pearl-earrings.jpg' WHERE name = 'Pearl Earrings';
UPDATE items SET image_url = '/images/accessories/silver-bracelet.jpg' WHERE name = 'Silver Bracelet';
UPDATE items SET image_url = '/images/accessories/luxury-watch.jpg' WHERE name = 'Luxury Watch';

-- ACCESSORIES - Epic
UPDATE items SET image_url = '/images/accessories/diamond-ring.jpg' WHERE name = 'Diamond Ring';
UPDATE items SET image_url = '/images/accessories/cufflinks.jpg' WHERE name = 'Platinum Cufflinks';

-- ACCESSORIES - Legendary
UPDATE items SET image_url = '/images/accessories/collectible-watch.jpg' WHERE name = 'Rare Collectible Watch';
UPDATE items SET image_url = '/images/accessories/diamond-set.jpg' WHERE name = 'Diamond Necklace Set';

-- CARS - Common
UPDATE items SET image_url = '/images/cars/old-sedan.jpg' WHERE name = 'Old Sedan';
UPDATE items SET image_url = '/images/cars/compact-car.jpg' WHERE name = 'Compact Car';

-- CARS - Rare
UPDATE items SET image_url = '/images/cars/family-sedan.jpg' WHERE name = 'Family Sedan';
UPDATE items SET image_url = '/images/cars/used-suv.jpg' WHERE name = 'Used SUV';
UPDATE items SET image_url = '/images/cars/new-suv.jpg' WHERE name = 'New SUV';
UPDATE items SET image_url = '/images/cars/sports-coupe.jpg' WHERE name = 'Sports Coupe';

-- CARS - Epic
UPDATE items SET image_url = '/images/cars/electric-car.jpg' WHERE name = 'Electric Car';
UPDATE items SET image_url = '/images/cars/luxury-sedan.jpg' WHERE name = 'Luxury Sedan';
UPDATE items SET image_url = '/images/cars/tesla-model-s.jpg' WHERE name = 'Tesla Model S';

-- CARS - Legendary
UPDATE items SET image_url = '/images/cars/mercedes-s.jpg' WHERE name = 'Mercedes S-Class';
UPDATE items SET image_url = '/images/cars/porsche-911.jpg' WHERE name = 'Porsche 911';
UPDATE items SET image_url = '/images/cars/ferrari-f8.jpg' WHERE name = 'Ferrari F8';
UPDATE items SET image_url = '/images/cars/lamborghini.jpg' WHERE name = 'Lamborghini Aventador';

-- HOUSES - Common
UPDATE items SET image_url = '/images/houses/studio-apartment.jpg' WHERE name = 'Studio Apartment';
UPDATE items SET image_url = '/images/houses/small-apartment.jpg' WHERE name = 'Small Apartment';

-- HOUSES - Rare
UPDATE items SET image_url = '/images/houses/suburban-house.jpg' WHERE name = 'Suburban House';
UPDATE items SET image_url = '/images/houses/city-condo.jpg' WHERE name = 'City Condo';
UPDATE items SET image_url = '/images/houses/family-home.jpg' WHERE name = 'Family Home';
UPDATE items SET image_url = '/images/houses/lake-house.jpg' WHERE name = 'Lake House';

-- HOUSES - Epic
UPDATE items SET image_url = '/images/houses/beach-house.jpg' WHERE name = 'Beach House';
UPDATE items SET image_url = '/images/houses/penthouse.jpg' WHERE name = 'Luxury Penthouse';

-- HOUSES - Legendary
UPDATE items SET image_url = '/images/houses/mansion.jpg' WHERE name = 'Modern Mansion';
UPDATE items SET image_url = '/images/houses/estate.jpg' WHERE name = 'Private Estate';
UPDATE items SET image_url = '/images/houses/island-villa.jpg' WHERE name = 'Island Villa';

-- Verify the updates
SELECT COUNT(*) as updated_items FROM items WHERE image_url LIKE '/images/%';
