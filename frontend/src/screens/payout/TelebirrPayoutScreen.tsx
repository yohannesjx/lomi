import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ActivityIndicator, Animated } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { LinearGradient } from 'expo-linear-gradient';

export const TelebirrPayoutScreen = ({ route, navigation }: any) => {
    const { payout, netAmount } = route.params;
    const [step, setStep] = useState<'processing' | 'success'>('processing');
    const [progress, setProgress] = useState(0);
    const fadeAnim = React.useRef(new Animated.Value(0)).current;

    useEffect(() => {
        // Simulate Telebirr processing
        const timer = setInterval(() => {
            setProgress((prev) => {
                if (prev >= 100) {
                    clearInterval(timer);
                    setStep('success');
                    Animated.timing(fadeAnim, {
                        toValue: 1,
                        duration: 500,
                        useNativeDriver: true,
                    }).start();
                    return 100;
                }
                return prev + 10;
            });
        }, 300);

        return () => clearInterval(timer);
    }, []);

    const handleContinue = () => {
        navigation.navigate('PayoutThankYou', {
            payout: payout,
            netAmount: netAmount,
        });
    };

    if (step === 'processing') {
        return (
            <SafeAreaView style={styles.container} edges={['top', 'bottom']}>
                <View style={styles.content}>
                    <View style={styles.logoContainer}>
                        <Text style={styles.logo}>📱</Text>
                        <Text style={styles.logoText}>Telebirr</Text>
                    </View>

                    <View style={styles.processingContainer}>
                        <ActivityIndicator size="large" color={COLORS.primary} />
                        <Text style={styles.processingText}>Processing payout...</Text>
                        <Text style={styles.progressText}>{progress}%</Text>

                        <View style={styles.progressBar}>
                            <View style={[styles.progressFill, { width: `${progress}%` }]} />
                        </View>
                    </View>

                    <View style={styles.infoBox}>
                        <Text style={styles.infoTitle}>Processing Details</Text>
                        <View style={styles.infoRow}>
                            <Text style={styles.infoLabel}>Amount:</Text>
                            <Text style={styles.infoValue}>{netAmount.toFixed(2)} Birr</Text>
                        </View>
                        <View style={styles.infoRow}>
                            <Text style={styles.infoLabel}>To:</Text>
                            <Text style={styles.infoValue}>{payout.payment_account}</Text>
                        </View>
                        <View style={styles.infoRow}>
                            <Text style={styles.infoLabel}>Status:</Text>
                            <Text style={styles.infoValue}>Processing...</Text>
                        </View>
                    </View>
                </View>
            </SafeAreaView>
        );
    }

    return (
        <SafeAreaView style={styles.container} edges={['top', 'bottom']}>
            <Animated.View style={[styles.content, { opacity: fadeAnim }]}>
                <View style={styles.successContainer}>
                    <View style={styles.successIcon}>
                        <Text style={styles.successCheckmark}>✓</Text>
                    </View>
                    <Text style={styles.successTitle}>Payout Request Submitted!</Text>
                    <Text style={styles.successMessage}>
                        Your payout request has been received and will be processed on Monday.
                    </Text>

                    <View style={styles.amountCard}>
                        <Text style={styles.amountLabel}>You'll Receive</Text>
                        <Text style={styles.amountValue}>{netAmount.toFixed(2)} Birr</Text>
                        <Text style={styles.amountNote}>
                            Platform fee ({payout.platform_fee_percentage}%) already deducted
                        </Text>
                    </View>

                    <TouchableOpacity style={styles.continueButton} onPress={handleContinue}>
                        <Text style={styles.continueButtonText}>Continue</Text>
                    </TouchableOpacity>
                </View>
            </Animated.View>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    content: {
        flex: 1,
        padding: SPACING.xl,
        justifyContent: 'center',
    },
    logoContainer: {
        alignItems: 'center',
        marginBottom: SPACING.xxl,
    },
    logo: {
        fontSize: 64,
        marginBottom: SPACING.m,
    },
    logoText: {
        fontSize: 32,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    processingContainer: {
        alignItems: 'center',
        marginBottom: SPACING.xl,
    },
    processingText: {
        fontSize: 18,
        color: COLORS.textPrimary,
        marginTop: SPACING.l,
        marginBottom: SPACING.m,
    },
    progressText: {
        fontSize: 24,
        fontWeight: 'bold',
        color: COLORS.primary,
        marginBottom: SPACING.m,
    },
    progressBar: {
        width: '100%',
        height: 8,
        backgroundColor: COLORS.surfaceHighlight,
        borderRadius: 4,
        overflow: 'hidden',
        marginTop: SPACING.m,
    },
    progressFill: {
        height: '100%',
        backgroundColor: COLORS.primary,
        borderRadius: 4,
    },
    infoBox: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.l,
        marginTop: SPACING.xl,
    },
    infoTitle: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.m,
    },
    infoRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        marginBottom: SPACING.s,
    },
    infoLabel: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    infoValue: {
        fontSize: 14,
        fontWeight: '600',
        color: COLORS.textPrimary,
    },
    successContainer: {
        alignItems: 'center',
    },
    successIcon: {
        width: 100,
        height: 100,
        borderRadius: 50,
        backgroundColor: COLORS.primary,
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: SPACING.xl,
    },
    successCheckmark: {
        fontSize: 60,
        color: COLORS.background,
        fontWeight: 'bold',
    },
    successTitle: {
        fontSize: 28,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.m,
        textAlign: 'center',
    },
    successMessage: {
        fontSize: 16,
        color: COLORS.textSecondary,
        textAlign: 'center',
        marginBottom: SPACING.xl,
        lineHeight: 24,
    },
    amountCard: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.xl,
        alignItems: 'center',
        width: '100%',
        marginBottom: SPACING.xl,
    },
    amountLabel: {
        fontSize: 14,
        color: COLORS.textSecondary,
        marginBottom: SPACING.xs,
    },
    amountValue: {
        fontSize: 36,
        fontWeight: 'bold',
        color: COLORS.primary,
        marginBottom: SPACING.xs,
    },
    amountNote: {
        fontSize: 12,
        color: COLORS.textSecondary,
    },
    continueButton: {
        backgroundColor: COLORS.primary,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        width: '100%',
        alignItems: 'center',
    },
    continueButtonText: {
        color: COLORS.background,
        fontSize: 18,
        fontWeight: 'bold',
    },
});

