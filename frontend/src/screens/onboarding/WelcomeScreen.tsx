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
            // Wait a bit for Telegram WebApp to fully initialize
            const webApp = getTelegramWebApp();
            
            // Debug logging
            console.log('🔍 Telegram WebApp check:', {
                exists: !!webApp,
                initData: webApp?.initData ? 'present' : 'missing',
                initDataLength: webApp?.initData?.length || 0,
                platform: webApp?.platform,
                version: webApp?.version,
            });
            
            // Try to get initData with retry
            let initData = getTelegramInitData();
            
            // If not available, wait a bit and try again (Telegram might still be loading)
            if (!initData && webApp) {
                console.log('⏳ Waiting for Telegram initData...');
                await new Promise(resolve => setTimeout(resolve, 500));
                initData = getTelegramInitData();
            }
            
            if (!initData) {
                // Check if we're actually in Telegram
                const isInTelegram = webApp !== null || 
                    (typeof window !== 'undefined' && 
                     (window.location.search.includes('tgWebApp') || 
                      window.location.hash.includes('tgWebApp') ||
                      navigator.userAgent.includes('Telegram')));
                
                if (!isInTelegram) {
                    const errorMsg = 'Please open this app from Telegram. Go to your bot and click the menu button, then select the Mini App.';
                    console.error('❌', errorMsg);
                    alert(errorMsg);
                    return;
                }
                
                // We're in Telegram but initData is missing
                // This can happen if the app is opened in a browser or if Telegram hasn't provided initData yet
                const errorMsg = 'Telegram authentication data is missing. This might happen if:\n\n1. The app was opened in a browser instead of Telegram\n2. The Mini App URL is incorrect in BotFather\n3. Telegram needs to be updated\n\nPlease try:\n- Closing and reopening from Telegram\n- Checking BotFather Mini App settings\n- Updating Telegram app';
                console.error('❌', errorMsg);
                console.error('Debug info:', {
                    webAppExists: !!webApp,
                    userAgent: navigator.userAgent,
                    url: window.location.href,
                    search: window.location.search,
                    hash: window.location.hash,
                });
                alert(errorMsg);
                return;
            }
            
            console.log('✅ InitData found, attempting login...');
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
            
            // Show error message (use alert for maximum compatibility)
            const webApp = getTelegramWebApp();
            const errorMsg = error?.message || 'Login failed. Please try again.';
            
            // Use alert for now - showConfirm/showAlert have compatibility issues in some versions
            alert(errorMsg);
            
            // Optional: Try to use Telegram's native alert if available and working
            // if (webApp && typeof webApp.showAlert === 'function') {
            //     try {
            //         webApp.showAlert(errorMsg);
            //     } catch (e) {
            //         alert(errorMsg);
            //     }
            // } else {
            //     alert(errorMsg);
            // }
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
