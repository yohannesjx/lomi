import React, { useEffect, useState } from 'react';
import { createStackNavigator } from '@react-navigation/stack';
import { useOnboardingStore } from '../store/onboardingStore';
import { useAuthStore } from '../store/authStore';
import { UserService } from '../api/services';
import { View, StyleSheet, Platform, ActivityIndicator } from 'react-native';
import { COLORS } from '../theme/colors';

// ToastAndroid is only available on Android
let ToastAndroid: any = null;
if (Platform.OS === 'android') {
    try {
        ToastAndroid = require('react-native').ToastAndroid;
    } catch (e) {
        // Not available
    }
}

// Onboarding Screens
import { ProfileSetupScreen } from '../screens/onboarding/ProfileSetupScreen';
import { GenderPreferenceScreen } from '../screens/onboarding/GenderPreferenceScreen';
import { PhotoUploadScreen } from '../screens/onboarding/PhotoUploadScreen';
import { InterestsScreen } from '../screens/onboarding/InterestsScreen';

// Additional screens we'll need to create
import { CityScreen } from '../screens/onboarding/CityScreen';
import { ReligionScreen } from '../screens/onboarding/ReligionScreen';
import { VideoScreen } from '../screens/onboarding/VideoScreen';
import { BioScreen } from '../screens/onboarding/BioScreen';
import { OnboardingCompleteScreen } from '../screens/onboarding/OnboardingCompleteScreen';
import { PhotoModerationStatusScreen } from '../screens/moderation/PhotoModerationStatusScreen';

const Stack = createStackNavigator();

// Total number of onboarding steps (including completion screen)
export const TOTAL_ONBOARDING_STEPS = 8;

// Step mapping: onboarding_step -> screen name
const STEP_TO_SCREEN: Record<number, string> = {
    0: 'ProfileSetup',      // Age & Gender - Step 1
    1: 'City',              // City - Step 2
    2: 'GenderPreference',  // Looking for + Goal - Step 3
    3: 'Religion',          // Religion - Step 4
    4: 'PhotoUpload',       // Photos (at least 3) - Step 5
    5: 'Video',             // Video (optional) - Step 6
    6: 'Bio',               // Bio & Interests - Step 7
    7: 'OnboardingComplete', // Completion screen - Step 8
};

// Screen order for navigation
const SCREEN_ORDER = [
    'ProfileSetup',
    'City',
    'GenderPreference',
    'Religion',
    'PhotoUpload',
    'Video',
    'Bio',
    'OnboardingComplete',
];

export const OnboardingNavigator: React.FC<{ navigation: any }> = ({ navigation }) => {
    const { onboardingStep, onboardingCompleted, fetchStatus, isLoading } = useOnboardingStore();
    const { isAuthenticated, user } = useAuthStore();
    const [initialRoute, setInitialRoute] = useState<string | null>(null);
    const [hasShownWelcomeBack, setHasShownWelcomeBack] = useState(false);
    const [hasFetchedStatus, setHasFetchedStatus] = useState(false);
    const [hasCheckedModeration, setHasCheckedModeration] = useState(false);

    useEffect(() => {
        // Fetch onboarding status when component mounts
        if (isAuthenticated && !hasFetchedStatus) {
            fetchStatus().then(() => {
                setHasFetchedStatus(true);
            }).catch(() => {
                setHasFetchedStatus(true); // Still set to true even on error to prevent infinite loading
            });
        }
    }, [isAuthenticated, fetchStatus, hasFetchedStatus]);

    useEffect(() => {
        // Determine initial route based on onboarding step
        // Use user.onboarding_step as fallback if store hasn't loaded yet
        const currentStep = user?.onboarding_step ?? onboardingStep ?? 0;
        const currentCompleted = user?.onboarding_completed ?? onboardingCompleted ?? false;
        
        // Only set route if we've either fetched status or have user data
        if ((!isLoading && hasFetchedStatus) || (user && user.onboarding_step !== undefined)) {
            // If step is 4 (PhotoUpload), check if there are pending photos
            // If so, show PhotoStatus instead so user can see moderation status
            if (currentStep === 4 && isAuthenticated && !hasCheckedModeration) {
                setHasCheckedModeration(true);
                UserService.getModerationStatus()
                    .then((status) => {
                        const pending = status.summary?.pending ?? 0;
                        const totalPhotos = status.summary?.total_photos ?? 0;
                        const approved = status.summary?.approved ?? 0;
                        
                        // If user has uploaded photos (pending or approved), show PhotoStatus
                        // This handles the case where user uploaded photos and refreshed
                        if (totalPhotos > 0 || pending > 0 || approved > 0) {
                            console.log('📸 User has photos (pending or approved), showing PhotoStatus screen');
                            setInitialRoute('PhotoStatus');
                        } else {
                            // No photos uploaded yet, show PhotoUpload screen
                            console.log('📸 No photos uploaded yet, showing PhotoUpload screen');
                            setInitialRoute(STEP_TO_SCREEN[currentStep] || 'ProfileSetup');
                        }
                    })
                    .catch((error) => {
                        console.warn('⚠️ Failed to check moderation status, defaulting to PhotoUpload:', error);
                        // On error, default to PhotoUpload screen
                        setInitialRoute(STEP_TO_SCREEN[currentStep] || 'ProfileSetup');
                    });
            } else {
                const targetScreen = STEP_TO_SCREEN[currentStep] || 'ProfileSetup';
                setInitialRoute(targetScreen);
            }

            // Show welcome back toast if resuming (step > 0)
            if (currentStep > 0 && !currentCompleted && !hasShownWelcomeBack) {
                setHasShownWelcomeBack(true);
                const message = 'Welcome back! Continuing where you left off...';
                if (Platform.OS === 'android' && ToastAndroid) {
                    ToastAndroid.show(message, ToastAndroid.SHORT);
                } else {
                    // For iOS/web, use console or alert
                    console.log('👋', message);
                }
            }
        }
    }, [onboardingStep, onboardingCompleted, isLoading, hasShownWelcomeBack, user, hasFetchedStatus, isAuthenticated, hasCheckedModeration]);

    // If onboarding is completed, navigate to Main
    useEffect(() => {
        if (onboardingCompleted && navigation) {
            navigation.reset({
                index: 0,
                routes: [{ name: 'Main' }],
            });
        }
    }, [onboardingCompleted, navigation]);

    // Don't render until we know the initial route
    if (!initialRoute || isLoading) {
        return (
            <View style={styles.loadingContainer}>
                <ActivityIndicator size="large" color={COLORS.primary} />
            </View>
        );
    }

    const currentStep = user?.onboarding_step ?? onboardingStep ?? 0;

    return (
        <View style={styles.container}>
            <Stack.Navigator
                initialRouteName={initialRoute}
                screenOptions={{
                    headerShown: false,
                    cardStyleInterpolator: ({ current, next, layouts }) => {
                        return {
                            cardStyle: {
                                transform: [
                                    {
                                        translateX: current.progress.interpolate({
                                            inputRange: [0, 1],
                                            outputRange: [layouts.screen.width, 0],
                                        }),
                                    },
                                ],
                            },
                        };
                    },
                }}
            >
                <Stack.Screen name="ProfileSetup" component={ProfileSetupScreen} />
                <Stack.Screen name="City" component={CityScreen} />
                <Stack.Screen name="GenderPreference" component={GenderPreferenceScreen} />
                <Stack.Screen name="Religion" component={ReligionScreen} />
                <Stack.Screen name="PhotoUpload" component={PhotoUploadScreen} />
                <Stack.Screen name="PhotoStatus" component={PhotoModerationStatusScreen} />
                <Stack.Screen name="Video" component={VideoScreen} />
                <Stack.Screen name="Bio" component={BioScreen} />
                <Stack.Screen name="OnboardingComplete" component={OnboardingCompleteScreen} />
            </Stack.Navigator>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
    },
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: COLORS.background,
    },
});

