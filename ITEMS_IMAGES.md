# Item Images Setup

This document explains how to set up and manage item images for the Freezino shop.

## Overview

The shop now supports 40+ meme and weird items across all categories (clothing, accessories, cars, houses). Each item has an associated image that is displayed in the shop.

## New Items Added

### Clothing (10 items)
- Tinfoil Hat - Protection from mind control!
- Banana Suit - Full body banana costume
- Dinosaur Pajamas - Cozy dino onesie
- Socks with Sandals - Dad-approved fashion
- Ugly Christmas Sweater - Party perfect
- Chicken Suit - Embrace your inner poultry
- Cat Meme T-Shirt - Classic internet culture
- Unicorn Onesie - Magical pajamas
- Inflatable T-Rex Costume - Battery-powered chaos!
- Pickle Costume - I'm Pickle Rick!

### Accessories (10 items)
- Rubber Duck - Debugging companion
- Fake Mustache - Instant disguise
- Monocle - Quite fancy!
- Power Ring - One ring to rule them all
- Magic Wand - May or may not work
- Glowing LED Glasses - Cyberpunk vibes
- Potato Clock - Powered by science!
- Finger Hands - Cursed tiny hands
- Googly Eyes Glasses - Silly specs
- Unicorn Horn Headband - Instant magic

### Cars (10 items)
- Oscar Mayer Wienermobile - Legendary hot dog car
- Flying Carpet - A whole new world!
- Racing Shopping Cart - Supermarket sweep champion
- Segway - Mall cop approved
- Power Wheels Barbie Jeep - 5mph dream machine
- Ice Cream Truck - Musical transportation
- Adult Tricycle - Three wheels of stability
- Rocket-Powered Skateboard - Back to the Future!
- Hoverboard - Doesn't hover, still cool
- Toy Rocket Ship - To the moon!

### Houses (10 items)
- Cardboard Box Mansion - Peak minimalism
- Epic Treehouse - No parents allowed!
- Cozy Cave - Return to monke
- Pillow Fort Palace - Maximum comfort
- Arctic Igloo - Natural AC
- Hobbit Hole - LOTR living
- UFO Landing Base - Aliens welcome
- Luxury Camping Tent - Glamping
- Van Down By The River - #VanLife
- Inflatable Bouncy Castle - Childhood dream home

## Setup Instructions

### 1. Download Images

Run the following command to download all item images from Unsplash:

```bash
make download-images
```

This will:
- Create the directory structure: `backend/static/images/{clothing,accessories,cars,houses}/`
- Download 40 images from Unsplash (free stock photos)
- Save them in the appropriate category folders

### 2. Image Storage

Images are stored in:
```
backend/static/images/
├── clothing/
├── accessories/
├── cars/
└── houses/
```

### 3. Backend Configuration

The backend (Go/Fiber) serves static images via:
- Route: `/images/*`
- Directory: `./static/images/`
- Configuration: `backend/cmd/server/main.go`

### 4. Frontend Configuration

The frontend (React/Vite) proxies image requests to the backend in development:
- Dev proxy: `/images` → `http://localhost:3000/images`
- Production: Images served directly from backend
- Configuration: `frontend/vite.config.ts`

### 5. Database Seeding

When you run the backend, it will automatically seed the database with all items including the new meme items. If you need to re-seed:

```bash
# Delete the database
rm backend/data/freezino.db

# Restart the backend - it will recreate and seed automatically
cd backend && make run
```

## Image URLs

All items use the path format: `/images/{category}/{item-name}.jpg`

Examples:
- `/images/clothing/banana-suit.jpg`
- `/images/accessories/rubber-duck.jpg`
- `/images/cars/hotdog-car.jpg`
- `/images/houses/cardboard-box.jpg`

## Troubleshooting

### Images not loading in development

1. Make sure the backend is running on port 3000
2. Check that images were downloaded: `ls backend/static/images/clothing/`
3. Verify the proxy is working: Open browser dev tools and check network tab

### Images not loading in production

1. Ensure images are deployed to the server in `backend/static/images/`
2. Check backend logs for static file serving errors
3. Verify the backend is serving files at `/images/*`

### Re-downloading images

If you need to re-download all images:

```bash
# Remove existing images
rm -rf backend/static/images/

# Download fresh copies
make download-images
```

## Image Sources

All images are from Unsplash (free stock photos) and are free to use under the Unsplash License:
- No attribution required (but appreciated!)
- Free for commercial and non-commercial use
- Modifications allowed

Image URLs are in: `scripts/download-images.sh`

## Adding New Items

To add new items with images:

1. Add the item to `backend/internal/database/seed.go`
2. Add the image URL to `scripts/download-images.sh`
3. Run `make download-images` to fetch the new image
4. Restart the backend to re-seed the database

## Production Deployment

When deploying to production:

1. Download images locally: `make download-images`
2. Copy images to the server: `scp -r backend/static/images/ user@server:/path/to/backend/static/`
3. Or run `make download-images` on the server directly

The deployment process should preserve the `static/` directory.
