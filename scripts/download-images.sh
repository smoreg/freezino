#!/bin/bash

# Script to download item images from Unsplash and other free sources
# Run with: bash scripts/download-images.sh

set -e

# Create directories
mkdir -p backend/static/images/clothing
mkdir -p backend/static/images/accessories
mkdir -p backend/static/images/cars
mkdir -p backend/static/images/houses

echo "Downloading item images..."

# CLOTHING - Meme Items
curl -L "https://images.unsplash.com/photo-1578632292335-df3abbb0d586?w=400" -o "backend/static/images/clothing/tinfoil-hat.jpg"
curl -L "https://images.unsplash.com/photo-1626808642875-0aa545482dfb?w=400" -o "backend/static/images/clothing/banana-suit.jpg"
curl -L "https://images.unsplash.com/photo-1629367494173-c78a56567877?w=400" -o "backend/static/images/clothing/dino-pajamas.jpg"
curl -L "https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=400" -o "backend/static/images/clothing/socks-sandals.jpg"
curl -L "https://images.unsplash.com/photo-1576566588028-4147f3842f27?w=400" -o "backend/static/images/clothing/ugly-sweater.jpg"
curl -L "https://picsum.photos/400/400?random=1" -o "backend/static/images/clothing/chicken-suit.jpg"
curl -L "https://images.unsplash.com/photo-1618354691373-d851c5c3a990?w=400" -o "backend/static/images/clothing/cat-tshirt.jpg"
curl -L "https://images.unsplash.com/photo-1620799140408-edc6dcb6d633?w=400" -o "backend/static/images/clothing/onesie.jpg"
curl -L "https://images.unsplash.com/photo-1581235720704-06d3acfcb36f?w=400" -o "backend/static/images/clothing/inflatable-dinosaur.jpg"
curl -L "https://images.unsplash.com/photo-1612036782180-6f0b6cd846fe?w=400" -o "backend/static/images/clothing/pickle-costume.jpg"

# ACCESSORIES - Meme Items
curl -L "https://images.unsplash.com/photo-1563492065599-3520f775eeed?w=400" -o "backend/static/images/accessories/rubber-duck.jpg"
curl -L "https://images.unsplash.com/photo-1614732414444-096e5f1122d5?w=400" -o "backend/static/images/accessories/fake-mustache.jpg"
curl -L "https://images.unsplash.com/photo-1509695507497-903c140c43b0?w=400" -o "backend/static/images/accessories/monocle.jpg"
curl -L "https://images.unsplash.com/photo-1605100804763-247f67b3557e?w=400" -o "backend/static/images/accessories/power-ring.jpg"
curl -L "https://images.unsplash.com/photo-1580128660010-fd027e1e587a?w=400" -o "backend/static/images/accessories/magic-wand.jpg"
curl -L "https://images.unsplash.com/photo-1511499767150-a48a237f0083?w=400" -o "backend/static/images/accessories/glowing-glasses.jpg"
curl -L "https://images.unsplash.com/photo-1551033406-611cf9a28f67?w=400" -o "backend/static/images/accessories/potato-clock.jpg"
curl -L "https://images.unsplash.com/photo-1514064019862-23e2a332a6a6?w=400" -o "backend/static/images/accessories/finger-hands.jpg"
curl -L "https://images.unsplash.com/photo-1574169208507-84376144848b?w=400" -o "backend/static/images/accessories/googly-eyes.jpg"
curl -L "https://images.unsplash.com/photo-1587402092301-725e37c70fd8?w=400" -o "backend/static/images/accessories/unicorn-horn.jpg"

# CARS - Meme Items
curl -L "https://images.unsplash.com/photo-1568605117036-5fe5e7bab0b7?w=400" -o "backend/static/images/cars/hotdog-car.jpg"
curl -L "https://images.unsplash.com/photo-1578894381163-e72c17f2d45f?w=400" -o "backend/static/images/cars/magic-carpet.jpg"
curl -L "https://images.unsplash.com/photo-1534452203293-494d7ddbf7e0?w=400" -o "backend/static/images/cars/shopping-cart.jpg"
curl -L "https://images.unsplash.com/photo-1572120360610-d971b9d7767c?w=400" -o "backend/static/images/cars/segway.jpg"
curl -L "https://images.unsplash.com/photo-1581235720704-06d3acfcb36f?w=400" -o "backend/static/images/cars/power-wheels.jpg"
curl -L "https://images.unsplash.com/photo-1563729784474-d77dbb933a9e?w=400" -o "backend/static/images/cars/ice-cream-truck.jpg"
curl -L "https://images.unsplash.com/photo-1474487548417-781cb71495f3?w=400" -o "backend/static/images/cars/tricycle.jpg"
curl -L "https://images.unsplash.com/photo-1541123603104-512919d6a96c?w=400" -o "backend/static/images/cars/skateboard.jpg"
curl -L "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400" -o "backend/static/images/cars/hoverboard.jpg"
curl -L "https://images.unsplash.com/photo-1502877338535-766e1452684a?w=400" -o "backend/static/images/cars/rocket.jpg"

# HOUSES - Meme Items
curl -L "https://picsum.photos/400/400?random=2" -o "backend/static/images/houses/cardboard-box.jpg"
curl -L "https://images.unsplash.com/photo-1520034475321-cbe63696469a?w=400" -o "backend/static/images/houses/treehouse.jpg"
curl -L "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=400" -o "backend/static/images/houses/cave.jpg"
curl -L "https://picsum.photos/400/400?random=3" -o "backend/static/images/houses/pillow-fort.jpg"
curl -L "https://images.unsplash.com/photo-1548777123-e216912df7d8?w=400" -o "backend/static/images/houses/igloo.jpg"
curl -L "https://images.unsplash.com/photo-1530789253388-582c481c54b0?w=400" -o "backend/static/images/houses/hobbit-hole.jpg"
curl -L "https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=400" -o "backend/static/images/houses/ufo.jpg"
curl -L "https://images.unsplash.com/photo-1523755231516-e43fd2e8dca5?w=400" -o "backend/static/images/houses/tent.jpg"
curl -L "https://images.unsplash.com/photo-1478131143081-80f7f84ca84d?w=400" -o "backend/static/images/houses/van.jpg"
curl -L "https://images.unsplash.com/photo-1464146072230-91cabc968266?w=400" -o "backend/static/images/houses/castle.jpg"

echo "✓ All images downloaded successfully!"
echo "Images saved to backend/static/images/"
