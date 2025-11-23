-- Seed data for Lomi Social
-- This file contains initial data for gifts catalog and other static content

-- ============================================================================
-- GIFTS CATALOG
-- ============================================================================

INSERT INTO gifts (id, name_en, name_am, description_en, description_am, coin_price, birr_value, icon_url, animation_url, sound_url, has_special_effect, special_effect_duration_days, is_active, is_featured, display_order) VALUES

-- Coffee Ceremony
(
    uuid_generate_v4(),
    'Bunna Ceremony',
    'የቡና ሥነ ሥርዓት',
    'Traditional Ethiopian coffee ceremony - the ultimate sign of respect and love',
    'ባህላዊ የኢትዮጵያ የቡና ሥነ ሥርዓት - የአክብሮት እና የፍቅር ምልክት',
    100,
    10.00,
    '/gifts/icons/bunna_ceremony.png',
    '/gifts/animations/bunna_ceremony.webm',
    '/gifts/sounds/bunna_ceremony.mp3',
    FALSE,
    NULL,
    TRUE,
    TRUE,
    1
),

-- Doro Wot
(
    uuid_generate_v4(),
    'Doro Wot Plate',
    'የዶሮ ወጥ',
    'Delicious Ethiopian chicken stew - share a meal together',
    'ጣፋጭ የኢትዮጵያ የዶሮ ወጥ - አብረን እንብላ',
    150,
    15.00,
    '/gifts/icons/doro_wot.png',
    '/gifts/animations/doro_wot.webm',
    '/gifts/sounds/plate_serve.mp3',
    FALSE,
    NULL,
    TRUE,
    TRUE,
    2
),

-- Red Rose + Tej
(
    uuid_generate_v4(),
    'Red Rose & Tej',
    'ቀይ ሮዝ እና ጠጅ',
    'Classic romance - a beautiful rose with traditional honey wine',
    'ክላሲክ ፍቅር - ቆንጆ ሮዝ ከባህላዊ የማር ወይን ጋር',
    200,
    20.00,
    '/gifts/icons/rose_tej.png',
    '/gifts/animations/rose_tej.webm',
    '/gifts/sounds/romantic.mp3',
    FALSE,
    NULL,
    TRUE,
    TRUE,
    3
),

-- Habesha Dress
(
    uuid_generate_v4(),
    'Habesha Dress',
    'የሐበሻ ልብስ',
    'Beautiful traditional dress - their avatar wears it for 7 days!',
    'ቆንጆ ባህላዊ ልብስ - አቫታራቸው ለ7 ቀናት ይለብሳል!',
    500,
    50.00,
    '/gifts/icons/habesha_dress.png',
    '/gifts/animations/habesha_dress.webm',
    '/gifts/sounds/dress_sparkle.mp3',
    TRUE,
    7,
    TRUE,
    TRUE,
    4
),

-- Golden Key of Love
(
    uuid_generate_v4(),
    'Golden "Ye Fikir Key"',
    'የወርቅ "የፍቅር ቁልፍ"',
    'The ultimate gift - unlock their heart with the golden key of love',
    'ከፍተኛው ስጦታ - የፍቅር የወርቅ ቁልፍ በመጠቀም ልባቸውን ክፈት',
    1000,
    100.00,
    '/gifts/icons/golden_key.png',
    '/gifts/animations/golden_key.webm',
    '/gifts/sounds/golden_key.mp3',
    FALSE,
    NULL,
    TRUE,
    TRUE,
    5
),

-- Additional Gifts
(
    uuid_generate_v4(),
    'Injera Basket',
    'የእንጀራ ቅርጫት',
    'Share the staple of Ethiopian cuisine',
    'የኢትዮጵያ ምግብ መሰረት አካፍል',
    80,
    8.00,
    '/gifts/icons/injera_basket.png',
    '/gifts/animations/injera_basket.webm',
    NULL,
    FALSE,
    NULL,
    TRUE,
    FALSE,
    6
),

(
    uuid_generate_v4(),
    'Meskel Flower',
    'የመስቀል አበባ',
    'Beautiful yellow daisy - symbol of finding the true cross',
    'ቆንጆ ቢጫ አበባ - የእውነተኛ መስቀል ምልክት',
    120,
    12.00,
    '/gifts/icons/meskel_flower.png',
    '/gifts/animations/meskel_flower.webm',
    '/gifts/sounds/flower.mp3',
    FALSE,
    NULL,
    TRUE,
    FALSE,
    7
),

(
    uuid_generate_v4(),
    'Shemma Scarf',
    'ሸማ',
    'Traditional Ethiopian cotton scarf - wrap them in warmth',
    'ባህላዊ የኢትዮጵያ የጥጥ ሸማ - በሙቀት ይጠቅልሏቸው',
    250,
    25.00,
    '/gifts/icons/shemma.png',
    '/gifts/animations/shemma.webm',
    NULL,
    FALSE,
    NULL,
    TRUE,
    FALSE,
    8
),

(
    uuid_generate_v4(),
    'Teff Grain Blessing',
    'የጤፍ በረከት',
    'Ancient grain of Ethiopia - blessing of abundance',
    'የኢትዮጵያ ጥንታዊ እህል - የብዛት በረከት',
    90,
    9.00,
    '/gifts/icons/teff.png',
    '/gifts/animations/teff.webm',
    NULL,
    FALSE,
    NULL,
    TRUE,
    FALSE,
    9
),

(
    uuid_generate_v4(),
    'Masinko Serenade',
    'የማሲንቆ ዜማ',
    'Traditional one-string violin music - a romantic serenade',
    'ባህላዊ አንድ ገመድ ቫዮሊን ሙዚቃ - የፍቅር ዜማ',
    180,
    18.00,
    '/gifts/icons/masinko.png',
    '/gifts/animations/masinko.webm',
    '/gifts/sounds/masinko.mp3',
    FALSE,
    NULL,
    TRUE,
    FALSE,
    10
),

(
    uuid_generate_v4(),
    'Lalibela Crown',
    'የላሊበላ አክሊል',
    'Royal crown from the rock-hewn churches - treat them like royalty',
    'ከድንጋይ የተቆረጹ ቤተክርስቲያኖች የንጉሣዊ አክሊል - እንደ ንጉሣዊ አድርጋቸው',
    800,
    80.00,
    '/gifts/icons/lalibela_crown.png',
    '/gifts/animations/lalibela_crown.webm',
    '/gifts/sounds/crown.mp3',
    TRUE,
    3,
    TRUE,
    TRUE,
    11
),

(
    uuid_generate_v4(),
    'Firfir Breakfast',
    'የፍርፍር ቁርስ',
    'Start their day right with delicious firfir',
    'ቀንን በጣፋጭ ፍርፍር ጀምር',
    70,
    7.00,
    '/gifts/icons/firfir.png',
    '/gifts/animations/firfir.webm',
    NULL,
    FALSE,
    NULL,
    TRUE,
    FALSE,
    12
),

(
    uuid_generate_v4(),
    'Timkat Blessing',
    'የጥምቀት በረከት',
    'Holy water blessing from Timkat celebration',
    'ከጥምቀት በዓል የተቀደሰ ውሃ በረከት',
    300,
    30.00,
    '/gifts/icons/timkat.png',
    '/gifts/animations/timkat.webm',
    '/gifts/sounds/blessing.mp3',
    FALSE,
    NULL,
    TRUE,
    FALSE,
    13
),

(
    uuid_generate_v4(),
    'Kolo Snack',
    'ቆሎ',
    'Traditional roasted barley snack - simple and sweet',
    'ባህላዊ የተጠበሰ ገብስ መክሰስ - ቀላል እና ጣፋጭ',
    50,
    5.00,
    '/gifts/icons/kolo.png',
    '/gifts/animations/kolo.webm',
    NULL,
    FALSE,
    NULL,
    TRUE,
    FALSE,
    14
),

(
    uuid_generate_v4(),
    'Buna Jebena',
    'የቡና ጀበና',
    'Traditional coffee pot - the heart of Ethiopian hospitality',
    'ባህላዊ የቡና ጀበና - የኢትዮጵያ እንግዳ ተቀባይነት ልብ',
    220,
    22.00,
    '/gifts/icons/jebena.png',
    '/gifts/animations/jebena.webm',
    '/gifts/sounds/coffee_pour.mp3',
    FALSE,
    NULL,
    TRUE,
    FALSE,
    15
);

-- ============================================================================
-- SAMPLE ADMIN USER (for testing)
-- Password: Admin@123 (hashed with bcrypt)
-- ============================================================================

-- Note: In production, create admin users through a secure CLI tool
-- This is just for initial development

INSERT INTO admin_users (email, password_hash, full_name, role, is_active) VALUES
(
    'admin@lomi.social',
    '$2a$10$rQZ9vXqZ9vXqZ9vXqZ9vXOqZ9vXqZ9vXqZ9vXqZ9vXqZ9vXqZ9vXq',  -- This is a placeholder, replace with actual bcrypt hash
    'System Administrator',
    'super_admin',
    TRUE
);

-- ============================================================================
-- REWARD CHANNELS (Earn Coins by Subscribing)
-- ============================================================================

INSERT INTO reward_channels (id, channel_username, channel_name, channel_link, coin_reward, icon_url, is_active, display_order) VALUES

(
    uuid_generate_v4(),
    'lomi_updates',
    'Lomi Social Updates',
    'https://t.me/lomi_updates',
    50,
    '/channels/icons/lomi_updates.png',
    TRUE,
    1
),

(
    uuid_generate_v4(),
    'lomi_dating_tips',
    'Lomi Dating Tips',
    'https://t.me/lomi_dating_tips',
    75,
    '/channels/icons/dating_tips.png',
    TRUE,
    2
),

(
    uuid_generate_v4(),
    'habesha_music',
    'Habesha Music',
    'https://t.me/habesha_music',
    100,
    '/channels/icons/music.png',
    TRUE,
    3
),

(
    uuid_generate_v4(),
    'ethiopian_culture',
    'Ethiopian Culture',
    'https://t.me/ethiopian_culture',
    60,
    '/channels/icons/culture.png',
    TRUE,
    4
),

(
    uuid_generate_v4(),
    'addis_life',
    'Addis Ababa Life',
    'https://t.me/addis_life',
    80,
    '/channels/icons/addis.png',
    TRUE,
    5
);

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE gifts IS 'Seed data includes 15 Ethiopian-themed gifts ranging from 50-1000 coins';
COMMENT ON TABLE reward_channels IS 'Seed data includes 5 reward channels for earning coins via Telegram subscriptions';
