import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { LinearGradient } from 'expo-linear-gradient';

export const PayoutThankYouScreen = ({ route, navigation }: any) => {
    const { payout, netAmount } = route.params;

    return (
        <SafeAreaView style={styles.container} edges={['top', 'bottom']}>
            <ScrollView contentContainerStyle={styles.scrollContent} showsVerticalScrollIndicator={false}>
                <View style={styles.content}>
                    {/* Success Icon */}
                    <View style={styles.iconContainer}>
                        <LinearGradient
                            colors={[COLORS.primary, COLORS.primaryDark]}
                            style={styles.iconGradient}
                        >
                            <Text style={styles.iconText}>🎉</Text>
                        </LinearGradient>
                    </View>

                    {/* Title */}
                    <Text style={styles.title}>Thank You! 🙏</Text>
                    <Text style={styles.subtitle}>Your payout request has been submitted</Text>

                    {/* Summary Card */}
                    <View style={styles.summaryCard}>
                        <Text style={styles.summaryTitle}>Payout Summary</Text>
                        
                        <View style={styles.summaryRow}>
                            <Text style={styles.summaryLabel}>Requested Amount</Text>
                            <Text style={styles.summaryValue}>{payout.gift_balance_amount.toFixed(2)} Birr</Text>
                        </View>
                        
                        <View style={styles.summaryRow}>
                            <Text style={styles.summaryLabel}>Platform Fee ({payout.platform_fee_percentage}%)</Text>
                            <Text style={styles.summaryValueFee}>-{payout.platform_fee_amount.toFixed(2)} Birr</Text>
                        </View>
                        
                        <View style={styles.divider} />
                        
                        <View style={styles.summaryRow}>
                            <Text style={styles.summaryLabelBold}>You'll Receive</Text>
                            <Text style={styles.summaryValueBold}>{netAmount.toFixed(2)} Birr</Text>
                        </View>
                    </View>

                    {/* Info Cards */}
                    <View style={styles.infoCard}>
                        <Text style={styles.infoIcon}>📅</Text>
                        <View style={styles.infoContent}>
                            <Text style={styles.infoTitle}>Processing Schedule</Text>
                            <Text style={styles.infoText}>
                                Payouts are processed every Monday. You'll receive your payment within 24-48 hours after processing.
                            </Text>
                        </View>
                    </View>

                    <View style={styles.infoCard}>
                        <Text style={styles.infoIcon}>📱</Text>
                        <View style={styles.infoContent}>
                            <Text style={styles.infoTitle}>Payment Method</Text>
                            <Text style={styles.infoText}>
                                {payout.payment_method === 'telebirr' ? 'Telebirr' : 'Bank Transfer'}
                            </Text>
                            <Text style={styles.infoTextSmall}>
                                To: {payout.payment_account}
                            </Text>
                        </View>
                    </View>

                    <View style={styles.infoCard}>
                        <Text style={styles.infoIcon}>📊</Text>
                        <View style={styles.infoContent}>
                            <Text style={styles.infoTitle}>Track Your Payout</Text>
                            <Text style={styles.infoText}>
                                Check your payout history in your profile to see the status of this request.
                            </Text>
                        </View>
                    </View>

                    {/* Leaderboard CTA */}
                    <TouchableOpacity
                        style={styles.leaderboardCard}
                        onPress={() => navigation.navigate('Leaderboard')}
                    >
                        <Text style={styles.leaderboardIcon}>🏆</Text>
                        <View style={styles.leaderboardContent}>
                            <Text style={styles.leaderboardTitle}>Top Gifted This Week</Text>
                            <Text style={styles.leaderboardText}>See who's earning the most! 👀</Text>
                        </View>
                        <Text style={styles.leaderboardArrow}>→</Text>
                    </TouchableOpacity>
                </View>
            </ScrollView>

            {/* Action Buttons */}
            <View style={styles.footer}>
                <TouchableOpacity
                    style={styles.primaryButton}
                    onPress={() => navigation.navigate('Main')}
                >
                    <Text style={styles.primaryButtonText}>Back to Home</Text>
                </TouchableOpacity>
                <TouchableOpacity
                    style={styles.secondaryButton}
                    onPress={() => navigation.navigate('Profile')}
                >
                    <Text style={styles.secondaryButtonText}>View History</Text>
                </TouchableOpacity>
            </View>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    scrollContent: {
        flexGrow: 1,
    },
    content: {
        padding: SPACING.xl,
        alignItems: 'center',
    },
    iconContainer: {
        marginBottom: SPACING.xl,
    },
    iconGradient: {
        width: 120,
        height: 120,
        borderRadius: 60,
        alignItems: 'center',
        justifyContent: 'center',
    },
    iconText: {
        fontSize: 64,
    },
    title: {
        fontSize: 32,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.s,
        textAlign: 'center',
    },
    subtitle: {
        fontSize: 16,
        color: COLORS.textSecondary,
        marginBottom: SPACING.xl,
        textAlign: 'center',
    },
    summaryCard: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.xl,
        width: '100%',
        marginBottom: SPACING.l,
    },
    summaryTitle: {
        fontSize: 20,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.m,
    },
    summaryRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: SPACING.m,
    },
    summaryLabel: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    summaryValue: {
        fontSize: 16,
        fontWeight: '600',
        color: COLORS.textPrimary,
    },
    summaryValueFee: {
        fontSize: 16,
        fontWeight: '600',
        color: COLORS.error,
    },
    summaryLabelBold: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    summaryValueBold: {
        fontSize: 24,
        fontWeight: 'bold',
        color: COLORS.primary,
    },
    divider: {
        height: 1,
        backgroundColor: COLORS.surfaceHighlight,
        marginVertical: SPACING.m,
    },
    infoCard: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.l,
        width: '100%',
        marginBottom: SPACING.m,
        flexDirection: 'row',
        alignItems: 'flex-start',
    },
    infoIcon: {
        fontSize: 32,
        marginRight: SPACING.m,
    },
    infoContent: {
        flex: 1,
    },
    infoTitle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.xs,
    },
    infoText: {
        fontSize: 14,
        color: COLORS.textSecondary,
        lineHeight: 20,
    },
    infoTextSmall: {
        fontSize: 12,
        color: COLORS.textSecondary,
        marginTop: SPACING.xs,
    },
    leaderboardCard: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.l,
        width: '100%',
        flexDirection: 'row',
        alignItems: 'center',
        borderWidth: 2,
        borderColor: COLORS.gold,
        marginTop: SPACING.m,
    },
    leaderboardIcon: {
        fontSize: 40,
        marginRight: SPACING.m,
    },
    leaderboardContent: {
        flex: 1,
    },
    leaderboardTitle: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.gold,
        marginBottom: SPACING.xs,
    },
    leaderboardText: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    leaderboardArrow: {
        fontSize: 24,
        color: COLORS.gold,
    },
    footer: {
        padding: SPACING.l,
        borderTopWidth: 1,
        borderTopColor: COLORS.surfaceHighlight,
        backgroundColor: COLORS.surface,
    },
    primaryButton: {
        backgroundColor: COLORS.primary,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        alignItems: 'center',
        marginBottom: SPACING.m,
    },
    primaryButtonText: {
        color: COLORS.background,
        fontSize: 18,
        fontWeight: 'bold',
    },
    secondaryButton: {
        backgroundColor: COLORS.surface,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    secondaryButtonText: {
        color: COLORS.textPrimary,
        fontSize: 16,
        fontWeight: '600',
    },
});

