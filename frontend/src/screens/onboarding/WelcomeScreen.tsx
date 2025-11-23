import React, { useEffect } from 'react';
import { View, Text, StyleSheet, Image, Dimensions, Platform, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { LinearGradient } from 'expo-linear-gradient';
import { Button } from '../../components/ui/Button';
import { COLORS, SPACING } from '../../theme/colors';
import { getTelegramInitData, initializeTelegramWebApp, getTelegramWebApp } from '../../utils/telegram';
import { useAuthStore } from '../../store/authStore';

const { width, height } = Dimensions.get('window');

export const WelcomeScreen = ({ navigation }: any) => {
    const { login, isLoading } = useAuthStore();

    useEffect(() => {
        // Initialize Telegram WebApp
        if (Platform.OS === 'web') {
            initializeTelegramWebApp();
        }
        
        // Check if already authenticated
        const { isAuthenticated, user, loadTokens } = useAuthStore.getState();
        loadTokens().then(() => {
            const state = useAuthStore.getState();
            if (state.isAuthenticated && state.user) {
                if (state.user.has_profile) {
                    navigation.navigate('Main');
                } else {
                    navigation.navigate('ProfileSetup');
                }
            }
        });
    }, []);

    const handleLogin = async () => {
        try {
            const initData = getTelegramInitData();
            
            if (!initData) {
                // Telegram WebApp is only available when opened from Telegram
                // This requires a publicly accessible HTTPS URL
                if (__DEV__) {
                    console.warn('⚠️ Telegram WebApp not available in local development.');
                    console.log('ℹ️  This is expected. Telegram auth only works when:');
                    console.log('   1. App is deployed to a public HTTPS URL');
                    console.log('   2. App is opened from within Telegram');
                    console.log('   3. Mini App is configured in BotFather');
                    console.log('');
                    console.log('💡 For now, you can test the UI flow without authentication.');
                    
                    // Auto-navigate to profile setup in dev mode
                    // This allows testing the UI without Telegram
                    navigation.navigate('ProfileSetup');
                    return;
                } else {
                    // In production, require Telegram
                    throw new Error('Telegram WebApp is required. Please open this app from Telegram.');
                }
            }
            
            await login(initData);
            
            // Check if user has completed profile
            const user = useAuthStore.getState().user;
            if (user?.has_profile) {
                navigation.navigate('Main');
            } else {
                navigation.navigate('ProfileSetup');
            }
        } catch (error: any) {
            console.error('Login error:', error);
            
            // Show error message
            const webApp = getTelegramWebApp();
            if (webApp) {
                webApp.showAlert(error?.message || 'Login failed. Please try again.');
            } else {
                alert(error?.message || 'Login failed. Please try again.');
            }
        }
    };

    return (
        <View style={styles.container}>
            {/* Background Image / Gradient */}
            <Image
                source={{ uri: 'https://images.unsplash.com/photo-1566737236500-c8ac43014a67?q=80&w=1000&auto=format&fit=crop' }}
                style={styles.backgroundImage}
                resizeMode="cover"
            />

            <LinearGradient
                colors={['transparent', 'rgba(0,0,0,0.8)', '#000000']}
                style={styles.gradient}
            />

            <SafeAreaView style={styles.content}>
                <View style={styles.header}>
                    <View style={styles.logoContainer}>
                        <Text style={styles.logoEmoji}>🍋</Text>
                    </View>
                    <Text style={styles.appName}>Lomi Social</Text>
                </View>

                <View style={styles.footer}>
                    <Text style={styles.tagline}>
                        Find your <Text style={styles.highlight}>Lomi</Text> in Ethiopia
                    </Text>
                    <Text style={styles.description}>
                        The most beautiful way to meet Habesha singles.
                        Serious dating, culture, and fun.
                    </Text>

                    <Button
                        title="Continue with Telegram"
                        onPress={handleLogin}
                        style={styles.button}
                        size="large"
                    />

                    <Text style={styles.terms}>
                        By continuing, you agree to our Terms & Privacy Policy.
                    </Text>
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
    backgroundImage: {
        width: width,
        height: height * 0.7,
        position: 'absolute',
        top: 0,
    },
    gradient: {
        position: 'absolute',
        left: 0,
        right: 0,
        bottom: 0,
        height: height,
    },
    content: {
        flex: 1,
        justifyContent: 'space-between',
    },
    header: {
        alignItems: 'center',
        marginTop: SPACING.xl,
    },
    logoContainer: {
        width: 80,
        height: 80,
        borderRadius: 40,
        backgroundColor: 'rgba(167, 255, 131, 0.2)',
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: SPACING.m,
        borderWidth: 1,
        borderColor: COLORS.primary,
    },
    logoEmoji: {
        fontSize: 40,
    },
    appName: {
        fontSize: 24,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        letterSpacing: 1,
    },
    footer: {
        padding: SPACING.l,
        paddingBottom: SPACING.xl,
    },
    tagline: {
        fontSize: 36,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.m,
        lineHeight: 44,
    },
    highlight: {
        color: COLORS.primary,
    },
    description: {
        fontSize: 16,
        color: COLORS.textSecondary,
        marginBottom: SPACING.xl,
        lineHeight: 24,
    },
    button: {
        marginBottom: SPACING.m,
        shadowColor: COLORS.primary,
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.3,
        shadowRadius: 8,
        elevation: 5,
    },
    terms: {
        fontSize: 12,
        color: COLORS.textTertiary,
        textAlign: 'center',
    },
});
