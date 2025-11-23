# 🚀 Lomi Social - Premium Onboarding Implementation Plan

## Overview
Complete overhaul of onboarding flow to make it premium, modern, and perfectly adapted for Ethiopia in 2025.

## Implementation Checklist

### Phase 1: Database & Backend ✅
- [ ] Add location columns to users table (latitude, longitude)
- [ ] Add `has_seen_swipe_tutorial` flag
- [ ] Update user model
- [ ] Create migration file
- [ ] Add Ethiopian cities list endpoint

### Phase 2: Core Components 🔨
- [ ] Create LocationService (GPS + city detection)
- [ ] Create EthiopianCitiesPicker component
- [ ] Create RelationshipGoalPicker (6 options)
- [ ] Create MultiPhotoUploader (Instagram-style)
- [ ] Create ProgressBar component (animated)
- [ ] Create ConfettiAnimation component
- [ ] Create SwipeTutorialOverlay component

### Phase 3: Screen Updates 📱
- [ ] Fix all screens with KeyboardAvoidingView
- [ ] Add back buttons to all screens (except first)
- [ ] Update CityScreen with auto-detection + picker
- [ ] Update GenderPreferenceScreen with 6 new options
- [ ] Update PhotoUploadScreen with multi-select
- [ ] Update OnboardingCompleteScreen with confetti
- [ ] Add SwipeTutorial to SwipeScreen

### Phase 4: Polish & Animations ✨
- [ ] Add haptic feedback to all buttons
- [ ] Add scale animations to buttons
- [ ] Add progress bar to all screens
- [ ] Add smooth transitions between screens
- [ ] Test resume functionality

## Detailed Specifications

### 1. Database Schema Changes

```sql
-- Add to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS latitude DECIMAL(10,8);
ALTER TABLE users ADD COLUMN IF NOT EXISTS longitude DECIMAL(11,8);
ALTER TABLE users ADD COLUMN IF NOT EXISTS has_seen_swipe_tutorial BOOLEAN DEFAULT FALSE;

-- Create index for location-based queries
CREATE INDEX IF NOT EXISTS idx_users_location ON users(latitude, longitude);
```

### 2. Ethiopian Cities List (Top 40)

```typescript
export const ETHIOPIAN_CITIES = [
  { name: 'Addis Ababa', region: 'Addis Ababa', lat: 9.0320, lng: 38.7469 },
  { name: 'Dire Dawa', region: 'Dire Dawa', lat: 9.5930, lng: 41.8660 },
  { name: 'Mekelle', region: 'Tigray', lat: 13.4967, lng: 39.4753 },
  { name: 'Gondar', region: 'Amhara', lat: 12.6000, lng: 37.4667 },
  { name: 'Bahir Dar', region: 'Amhara', lat: 11.5933, lng: 37.3905 },
  { name: 'Hawassa', region: 'Sidama', lat: 7.0500, lng: 38.4833 },
  { name: 'Adama (Nazret)', region: 'Oromia', lat: 8.5400, lng: 39.2700 },
  { name: 'Jimma', region: 'Oromia', lat: 7.6700, lng: 36.8333 },
  { name: 'Jijiga', region: 'Somali', lat: 9.3500, lng: 42.8000 },
  { name: 'Dessie', region: 'Amhara', lat: 11.1333, lng: 39.6333 },
  { name: 'Bishoftu (Debre Zeit)', region: 'Oromia', lat: 8.7500, lng: 38.9833 },
  { name: 'Shashamane', region: 'Oromia', lat: 7.2000, lng: 38.6000 },
  { name: 'Harar', region: 'Harari', lat: 9.3100, lng: 42.1200 },
  { name: 'Dilla', region: 'SNNPR', lat: 6.4167, lng: 38.3167 },
  { name: 'Nekemte', region: 'Oromia', lat: 9.0833, lng: 36.5333 },
  { name: 'Debre Birhan', region: 'Amhara', lat: 9.6800, lng: 39.5300 },
  { name: 'Asella', region: 'Oromia', lat: 7.9500, lng: 39.1333 },
  { name: 'Debre Markos', region: 'Amhara', lat: 10.3500, lng: 37.7167 },
  { name: 'Kombolcha', region: 'Amhara', lat: 11.0833, lng: 39.7333 },
  { name: 'Arba Minch', region: 'SNNPR', lat: 6.0333, lng: 37.5500 },
  { name: 'Hosaena', region: 'SNNPR', lat: 7.5500, lng: 37.8500 },
  { name: 'Harar', region: 'Harari', lat: 9.3142, lng: 42.1185 },
  { name: 'Gambela', region: 'Gambela', lat: 8.2500, lng: 34.5833 },
  { name: 'Ambo', region: 'Oromia', lat: 8.9833, lng: 37.8500 },
  { name: 'Woldia', region: 'Amhara', lat: 11.8333, lng: 39.6000 },
  { name: 'Debre Tabor', region: 'Amhara', lat: 11.8500, lng: 38.0167 },
  { name: 'Adigrat', region: 'Tigray', lat: 14.2667, lng: 39.4500 },
  { name: 'Aksum', region: 'Tigray', lat: 14.1333, lng: 38.7167 },
  { name: 'Welkite', region: 'SNNPR', lat: 8.2833, lng: 37.7833 },
  { name: 'Burayu', region: 'Oromia', lat: 9.0667, lng: 38.6167 },
  { name: 'Sebeta', region: 'Oromia', lat: 8.9167, lng: 38.6167 },
  { name: 'Bale Robe', region: 'Oromia', lat: 7.1167, lng: 40.0000 },
  { name: 'Asosa', region: 'Benishangul-Gumuz', lat: 10.0667, lng: 34.5333 },
  { name: 'Semera', region: 'Afar', lat: 11.7833, lng: 41.0000 },
  { name: 'Metu', region: 'Oromia', lat: 8.3000, lng: 35.5833 },
  { name: 'Goba', region: 'Oromia', lat: 7.0000, lng: 39.9833 },
  { name: 'Bonga', region: 'SNNPR', lat: 7.2667, lng: 36.2333 },
  { name: 'Wolaita Sodo', region: 'SNNPR', lat: 6.8167, lng: 37.7500 },
  { name: 'Butajira', region: 'SNNPR', lat: 8.1167, lng: 38.3833 },
  { name: 'Durame', region: 'SNNPR', lat: 7.2333, lng: 37.8833 },
];
```

### 3. Relationship Goal Options (6 Cards)

```typescript
export const RELATIONSHIP_GOALS = [
  {
    id: 'friends',
    emoji: '☕',
    title: 'Just friends & coffee',
    subtitle: 'Casual hangouts',
  },
  {
    id: 'fun',
    emoji: '😏',
    title: 'Chat & fun',
    subtitle: 'Keep it light',
  },
  {
    id: 'dating',
    emoji: '💕',
    title: 'Dating & romance',
    subtitle: 'See where it goes',
  },
  {
    id: 'travel',
    emoji: '✈️',
    title: 'Travel partner',
    subtitle: 'Explore together',
  },
  {
    id: 'serious',
    emoji: '💍',
    title: 'Serious relationship',
    subtitle: 'Looking for marriage',
  },
  {
    id: 'open',
    emoji: '🌟',
    title: "Let's see where it goes",
    subtitle: 'No pressure',
  },
];
```

### 4. Component Architecture

```
frontend/src/
├── components/
│   ├── onboarding/
│   │   ├── ProgressBar.tsx (NEW)
│   │   ├── CityPicker.tsx (NEW)
│   │   ├── RelationshipGoalCard.tsx (NEW)
│   │   ├── MultiPhotoUploader.tsx (NEW)
│   │   └── ConfettiAnimation.tsx (NEW)
│   └── swipe/
│       └── SwipeTutorialOverlay.tsx (NEW)
├── services/
│   └── LocationService.ts (NEW)
├── constants/
│   └── ethiopianCities.ts (NEW)
└── screens/
    └── onboarding/
        ├── CityScreen.tsx (UPDATE)
        ├── GenderPreferenceScreen.tsx (UPDATE)
        ├── PhotoUploadScreen.tsx (UPDATE)
        └── OnboardingCompleteScreen.tsx (UPDATE)
```

### 5. KeyboardAvoidingView Pattern

```typescript
<SafeAreaView style={styles.container} edges={['top']}>
  <KeyboardAvoidingView
    behavior={Platform.OS === 'ios' ? 'padding' : undefined}
    style={styles.keyboardView}
    keyboardVerticalOffset={0}
  >
    <BackButton />
    <ProgressBar currentStep={X} totalSteps={8} />
    
    <ScrollView
      contentContainerStyle={styles.scrollContent}
      keyboardShouldPersistTaps="handled"
      showsVerticalScrollIndicator={false}
    >
      {/* Content */}
    </ScrollView>
    
    {/* Footer OUTSIDE ScrollView */}
    <View style={styles.footer}>
      <Button
        title="Continue"
        onPress={handleNext}
        hapticFeedback
        scaleAnimation
      />
    </View>
  </KeyboardAvoidingView>
</SafeAreaView>
```

### 6. Multi-Photo Upload Flow

```typescript
// 1. User taps "Add Photos"
const pickMultiplePhotos = async () => {
  const result = await ImagePicker.launchImageLibraryAsync({
    mediaTypes: ImagePicker.MediaTypeOptions.Images,
    allowsMultipleSelection: true, // KEY!
    selectionLimit: 9,
    quality: 0.8,
  });
  
  if (!result.canceled) {
    // result.assets is array of selected photos
    const newPhotos = result.assets.map((asset, index) => ({
      uri: asset.uri,
      order: index,
      uploading: true,
    }));
    
    // Start uploading all in parallel
    uploadPhotosInBackground(newPhotos);
  }
};

// 2. Show grid with progress rings
{photos.map((photo, index) => (
  <PhotoCard
    key={index}
    photo={photo}
    onDelete={() => removePhoto(index)}
    onReorder={(newIndex) => reorderPhoto(index, newIndex)}
  />
))}
```

### 7. Confetti Animation

```typescript
import LottieView from 'lottie-react-native';

const ConfettiAnimation = ({ onComplete }) => {
  useEffect(() => {
    const timer = setTimeout(() => {
      onComplete();
    }, 2800); // Exactly 2.8 seconds
    
    return () => clearTimeout(timer);
  }, []);
  
  return (
    <View style={StyleSheet.absoluteFill}>
      <LottieView
        source={require('../assets/confetti.json')}
        autoPlay
        loop={false}
        style={{ width: '100%', height: '100%' }}
      />
    </View>
  );
};
```

### 8. Swipe Tutorial Overlay

```typescript
const SwipeTutorialOverlay = ({ onComplete }) => {
  const handAnim = useRef(new Animated.Value(0)).current;
  
  useEffect(() => {
    // Animate hand: right → left → disappear
    Animated.sequence([
      // Swipe right
      Animated.timing(handAnim, {
        toValue: 1,
        duration: 800,
        useNativeDriver: true,
      }),
      Animated.delay(500),
      // Reset
      Animated.timing(handAnim, {
        toValue: 0,
        duration: 0,
        useNativeDriver: true,
      }),
      // Swipe left
      Animated.timing(handAnim, {
        toValue: -1,
        duration: 800,
        useNativeDriver: true,
      }),
      Animated.delay(500),
    ]).start(() => {
      onComplete();
    });
  }, []);
  
  return (
    <View style={styles.overlay}>
      <Animated.View
        style={{
          transform: [{
            translateX: handAnim.interpolate({
              inputRange: [-1, 0, 1],
              outputRange: [-150, 0, 150],
            }),
          }],
        }}
      >
        <Text style={styles.hand}>👆</Text>
      </Animated.View>
      <Text style={styles.tutorialText}>
        Swipe right to Like 💚    Swipe left to Pass ❌
      </Text>
    </View>
  );
};
```

## Implementation Order

1. **Database Migration** (5 min)
2. **Ethiopian Cities Constant** (5 min)
3. **LocationService** (15 min)
4. **ProgressBar Component** (10 min)
5. **CityPicker Component** (20 min)
6. **RelationshipGoalPicker** (20 min)
7. **MultiPhotoUploader** (30 min)
8. **Update All Screens with KeyboardAvoidingView** (20 min)
9. **ConfettiAnimation** (15 min)
10. **SwipeTutorialOverlay** (15 min)
11. **Haptic Feedback + Animations** (10 min)
12. **Testing & Polish** (30 min)

**Total Estimated Time:** 3-4 hours

## Success Criteria

- ✅ 90% of users complete onboarding in < 60 seconds
- ✅ Buttons always visible above keyboard
- ✅ Back button works on every screen
- ✅ Location auto-detected or easy city selection
- ✅ Multi-photo upload feels like Instagram
- ✅ Confetti animation is smooth and beautiful
- ✅ Swipe tutorial shows only once
- ✅ Progress bar animates smoothly
- ✅ All buttons have haptic feedback

## Next Steps

1. Run database migration
2. Install required packages (expo-location, lottie-react-native, expo-haptics)
3. Create components in order
4. Update screens
5. Test thoroughly
6. Deploy

Let's build this! 🚀
