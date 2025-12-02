import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity, Platform, Alert, StatusBar } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Button } from '../../components/ui/Button';
import { BackButton } from '../../components/ui/BackButton';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { UserService } from '../../api/services';
import { useOnboardingStore } from '../../store/onboardingStore';
import { useAuthStore } from '../../store/authStore';
import { TOTAL_ONBOARDING_STEPS } from '../../navigation/OnboardingNavigator';

const ETHIOPIAN_CITIES = [
    'Addis Ababa',
    'Dire Dawa',
    'Mekelle',
    'Gondar',
    'Hawassa',
    'Bahir Dar',
    'Dessie',
    'Jimma',
    'Jijiga',
    'Shashamane',
    'Bishoftu',
    'Arba Minch',
    'Hosaena',
    'Harar',
    'Dilla',
    'Nekemte',
    'Debre Birhan',
    'Asella',
    'Debre Markos',
    'Kombolcha',
];

export const CityScreen = ({ navigation }: any) => {
    const [selectedCity, setSelectedCity] = useState('');
    const [showCityGrid, setShowCityGrid] = useState(false);
    const [isDetecting, setIsDetecting] = useState(false);
    const [isSaving, setIsSaving] = useState(false);
    const { updateStep } = useOnboardingStore();
    const { user } = useAuthStore();

    useEffect(() => {
        if (user?.city && user.city !== 'Not Set') {
            setSelectedCity(user.city);
            setShowCityGrid(true);
        }
    }, [user]);

    const handleDetectLocation = async () => {
        setIsDetecting(true);
        try {
            if (!navigator.geolocation) {
                setShowCityGrid(true);
                setIsDetecting(false);
                return;
            }

            navigator.geolocation.getCurrentPosition(
                async (position) => {
                    const { latitude, longitude } = position.coords;

                    try {
                        const response = await fetch(
                            `https://nominatim.openstreetmap.org/reverse?format=json&lat=${latitude}&lon=${longitude}`
                        );
                        const data = await response.json();

                        const detectedCity = data.address?.city ||
                            data.address?.town ||
                            data.address?.village ||
                            data.address?.state ||
                            '';

                        if (detectedCity) {
                            setSelectedCity(detectedCity);
                        }
                        setShowCityGrid(true);
                    } catch (error) {
                        setShowCityGrid(true);
                    }
                    setIsDetecting(false);
                },
                (error) => {
                    setShowCityGrid(true);
                    setIsDetecting(false);
                },
                { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 }
            );
        } catch (error) {
            setShowCityGrid(true);
            setIsDetecting(false);
        }
    };

    const handleNext = async () => {
        if (!selectedCity) return;

        setIsSaving(true);
        try {
            await UserService.updateProfile({ city: selectedCity });
            await updateStep(2);
            navigation.navigate('GenderPreference');
        } catch (error: any) {
            console.error('Save city error:', error);
            Alert.alert('Error', 'Failed to save city. Please try again.');
        } finally {
            setIsSaving(false);
        }
    };

    return (
        <View style={styles.container}>
            <StatusBar barStyle="dark-content" backgroundColor={COLORS.background} />
            <SafeAreaView style={styles.safeArea} edges={['bottom']}>
                <BackButton />
                <ScrollView
                    contentContainerStyle={styles.scrollContent}
                    showsVerticalScrollIndicator={false}
                >
                    <View style={styles.header}>
                        <Text style={styles.title}>Where are you located?</Text>
                        <Text style={styles.subtitle}>
                            {showCityGrid ? 'Select your city' : 'We can detect your location automatically'}
                        </Text>
                    </View>

                    {!showCityGrid ? (
                        <View style={styles.detectContainer}>
                            <Button
                                title={isDetecting ? 'Detecting...' : '📍 Detect My Location'}
                                onPress={handleDetectLocation}
                                isLoading={isDetecting}
                                disabled={isDetecting}
                                size="large"
                            />
                            <TouchableOpacity
                                style={styles.skipButton}
                                onPress={() => setShowCityGrid(true)}
                            >
                                <Text style={styles.skipText}>Or choose manually</Text>
                            </TouchableOpacity>
                        </View>
                    ) : (
                        <View style={styles.citiesGrid}>
                            {ETHIOPIAN_CITIES.map((city) => (
                                <TouchableOpacity
                                    key={city}
                                    style={[
                                        styles.cityButton,
                                        selectedCity === city && styles.cityButtonSelected,
                                    ]}
                                    onPress={() => setSelectedCity(city)}
                                >
                                    <Text
                                        style={[
                                            styles.cityText,
                                            selectedCity === city && styles.cityTextSelected,
                                        ]}
                                    >
                                        {city}
                                    </Text>
                                </TouchableOpacity>
                            ))}
                        </View>
                    )}
                </ScrollView>

                <View style={styles.footer}>
                    <Button
                        title="Continue"
                        onPress={handleNext}
                        disabled={!selectedCity || isSaving}
                        isLoading={isSaving}
                        size="large"
                    />
                </View>
            </SafeAreaView>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    safeArea: {
        flex: 1,
    },
    scrollContent: {
        flexGrow: 1,
        padding: SPACING.l,
        paddingBottom: SPACING.xl,
    },
    header: {
        marginBottom: SPACING.xl,
    },
    title: {
        fontSize: 28,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.xs,
    },
    subtitle: {
        fontSize: 16,
        color: COLORS.textSecondary,
    },
    detectContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        gap: SPACING.l,
    },
    skipButton: {
        padding: SPACING.m,
    },
    skipText: {
        color: COLORS.textSecondary,
        fontSize: 14,
        textDecorationLine: 'underline',
    },
    citiesGrid: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        gap: SPACING.m,
    },
    cityButton: {
        paddingHorizontal: SPACING.m,
        paddingVertical: SPACING.s,
        borderRadius: SIZES.radiusS,
        backgroundColor: COLORS.surface,
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    cityButtonSelected: {
        backgroundColor: COLORS.primary,
        borderColor: COLORS.primary,
    },
    cityText: {
        fontSize: 14,
        color: COLORS.textPrimary,
        fontWeight: '500',
    },
    cityTextSelected: {
        color: '#FFFFFF',
    },
    footer: {
        padding: SPACING.l,
        paddingBottom: Platform.OS === 'ios' ? SPACING.m : SPACING.l,
        backgroundColor: COLORS.background,
        borderTopWidth: 1,
        borderTopColor: COLORS.surfaceHighlight,
    },
});
