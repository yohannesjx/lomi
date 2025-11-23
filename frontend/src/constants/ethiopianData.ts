// Ethiopian Cities - Top 40 by population
// Used for city selection during onboarding

export interface EthiopianCity {
    name: string;
    region: string;
    lat: number;
    lng: number;
}

export const ETHIOPIAN_CITIES: EthiopianCity[] = [
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
    { name: 'Shashemene', region: 'Oromia', lat: 7.2000, lng: 38.6000 },
];

// Relationship goals for Ethiopian context
export interface RelationshipGoal {
    id: string;
    emoji: string;
    title: string;
    subtitle: string;
}

export const RELATIONSHIP_GOALS: RelationshipGoal[] = [
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
