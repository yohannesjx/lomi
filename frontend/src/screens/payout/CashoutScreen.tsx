import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, ScrollView, TouchableOpacity, TextInput, Alert, ActivityIndicator } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { PayoutService } from '../../api/services';
import { LinearGradient } from 'expo-linear-gradient';

const MIN_PAYOUT = 1000; // Minimum 1000 Birr
const PLATFORM_FEE = 25; // 25% platform fee

export const CashoutScreen = ({ navigation }: any) => {
    const [availableBalance, setAvailableBalance] = useState(0);
    const [pendingPayouts, setPendingPayouts] = useState(0);
    const [payoutAmount, setPayoutAmount] = useState('');
    const [selectedMethod, setSelectedMethod] = useState<'telebirr' | 'bank'>('telebirr');
    const [phoneNumber, setPhoneNumber] = useState('');
    const [accountName, setAccountName] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        loadBalance();
    }, []);

    const loadBalance = async () => {
        try {
            const response = await PayoutService.getBalance();
            setAvailableBalance(response.available_balance || 0);
            setPendingPayouts(response.pending_payouts || 0);
        } catch (error: any) {
            console.error('Load balance error:', error);
        }
    };

    const calculateNetAmount = (amount: number) => {
        const fee = (amount * PLATFORM_FEE) / 100;
        return amount - fee;
    };

    const handleQuickAmount = (amount: number) => {
        if (amount <= availableBalance) {
            setPayoutAmount(amount.toString());
        }
    };

    const handleSubmit = async () => {
        const amount = parseFloat(payoutAmount);

        if (!amount || amount < MIN_PAYOUT) {
            Alert.alert('Invalid Amount', `Minimum payout is ${MIN_PAYOUT} Birr`);
            return;
        }

        if (amount > availableBalance) {
            Alert.alert('Insufficient Balance', `You have ${availableBalance.toFixed(2)} Birr available`);
            return;
        }

        if (selectedMethod === 'telebirr' && !phoneNumber.trim()) {
            Alert.alert('Required', 'Please enter your Telebirr phone number');
            return;
        }

        if (selectedMethod === 'bank' && (!phoneNumber.trim() || !accountName.trim())) {
            Alert.alert('Required', 'Please enter bank account details');
            return;
        }

        setIsSubmitting(true);
        try {
            const response = await PayoutService.requestPayout({
                amount: amount,
                payment_method: selectedMethod === 'telebirr' ? 'telebirr' : 'cbe_birr',
                payment_account: phoneNumber,
                payment_account_name: accountName || undefined,
            });

            // Navigate to mock Telebirr flow
            navigation.navigate('TelebirrPayout', {
                payout: response.payout,
                netAmount: calculateNetAmount(amount),
            });
        } catch (error: any) {
            Alert.alert('Error', error.response?.data?.error || 'Failed to create payout request');
        } finally {
            setIsSubmitting(false);
        }
    };

    const netAmount = payoutAmount ? calculateNetAmount(parseFloat(payoutAmount)) : 0;
    const feeAmount = payoutAmount ? (parseFloat(payoutAmount) * PLATFORM_FEE) / 100 : 0;

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            <View style={styles.header}>
                <TouchableOpacity onPress={() => navigation.goBack()} style={styles.backButton}>
                    <Text style={styles.backIcon}>←</Text>
                </TouchableOpacity>
                <Text style={styles.headerTitle}>Cashout 💰</Text>
                <View style={styles.placeholder} />
            </View>

            <View style={styles.contentWrapper}>
                <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
                    {/* Balance Card */}
                    <LinearGradient
                        colors={[COLORS.primary, COLORS.primaryDark]}
                        style={styles.balanceCard}
                    >
                        <Text style={styles.balanceLabel}>Available Balance</Text>
                        <Text style={styles.balanceAmount}>🎁 {availableBalance.toFixed(2)} Birr</Text>
                        {pendingPayouts > 0 && (
                            <Text style={styles.pendingText}>
                                {pendingPayouts.toFixed(2)} Birr pending
                            </Text>
                        )}
                    </LinearGradient>

                    {/* Info Banner */}
                    <View style={styles.infoBanner}>
                        <Text style={styles.infoText}>
                            💡 Payouts are processed every Monday. Minimum withdrawal: {MIN_PAYOUT} Birr
                        </Text>
                    </View>

                    {/* Quick Amounts */}
                    <View style={styles.section}>
                        <Text style={styles.sectionTitle}>Quick Amount</Text>
                        <View style={styles.quickAmounts}>
                            {[1000, 2000, 5000, 10000].map((amount) => (
                                <TouchableOpacity
                                    key={amount}
                                    style={[
                                        styles.quickAmountBtn,
                                        payoutAmount === amount.toString() && styles.quickAmountBtnActive,
                                        amount > availableBalance && styles.quickAmountBtnDisabled,
                                    ]}
                                    onPress={() => handleQuickAmount(amount)}
                                    disabled={amount > availableBalance}
                                >
                                    <Text
                                        style={[
                                            styles.quickAmountText,
                                            payoutAmount === amount.toString() && styles.quickAmountTextActive,
                                            amount > availableBalance && styles.quickAmountTextDisabled,
                                        ]}
                                    >
                                        {amount}
                                    </Text>
                                </TouchableOpacity>
                            ))}
                        </View>
                    </View>

                    {/* Custom Amount */}
                    <View style={styles.section}>
                        <Text style={styles.sectionTitle}>Custom Amount</Text>
                        <View style={styles.inputContainer}>
                            <Text style={styles.inputLabel}>Amount (Birr)</Text>
                            <TextInput
                                style={styles.input}
                                placeholder={`Min ${MIN_PAYOUT} Birr`}
                                placeholderTextColor={COLORS.textSecondary}
                                value={payoutAmount}
                                onChangeText={setPayoutAmount}
                                keyboardType="numeric"
                            />
                        </View>
                    </View>

                    {/* Payment Method */}
                    <View style={styles.section}>
                        <Text style={styles.sectionTitle}>Payment Method</Text>
                        <TouchableOpacity
                            style={[
                                styles.methodCard,
                                selectedMethod === 'telebirr' && styles.methodCardSelected,
                            ]}
                            onPress={() => setSelectedMethod('telebirr')}
                        >
                            <View style={styles.methodLeft}>
                                <Text style={styles.methodIcon}>📱</Text>
                                <View>
                                    <Text style={styles.methodName}>Telebirr</Text>
                                    <Text style={styles.methodDesc}>Instant transfer</Text>
                                </View>
                            </View>
                            {selectedMethod === 'telebirr' && (
                                <View style={styles.checkmark}>
                                    <Text style={styles.checkmarkText}>✓</Text>
                                </View>
                            )}
                        </TouchableOpacity>

                        <TouchableOpacity
                            style={[
                                styles.methodCard,
                                selectedMethod === 'bank' && styles.methodCardSelected,
                            ]}
                            onPress={() => setSelectedMethod('bank')}
                        >
                            <View style={styles.methodLeft}>
                                <Text style={styles.methodIcon}>🏦</Text>
                                <View>
                                    <Text style={styles.methodName}>Bank Transfer</Text>
                                    <Text style={styles.methodDesc}>CBE, Awash, etc.</Text>
                                </View>
                            </View>
                            {selectedMethod === 'bank' && (
                                <View style={styles.checkmark}>
                                    <Text style={styles.checkmarkText}>✓</Text>
                                </View>
                            )}
                        </TouchableOpacity>
                    </View>

                    {/* Account Details */}
                    <View style={styles.section}>
                        <Text style={styles.sectionTitle}>Account Details</Text>
                        <View style={styles.inputContainer}>
                            <Text style={styles.inputLabel}>
                                {selectedMethod === 'telebirr' ? 'Telebirr Phone Number' : 'Account Number'}
                            </Text>
                            <TextInput
                                style={styles.input}
                                placeholder={selectedMethod === 'telebirr' ? '0912345678' : 'Account number'}
                                placeholderTextColor={COLORS.textSecondary}
                                value={phoneNumber}
                                onChangeText={setPhoneNumber}
                                keyboardType="phone-pad"
                            />
                        </View>
                        {selectedMethod === 'bank' && (
                            <View style={styles.inputContainer}>
                                <Text style={styles.inputLabel}>Account Name</Text>
                                <TextInput
                                    style={styles.input}
                                    placeholder="Full name on account"
                                    placeholderTextColor={COLORS.textSecondary}
                                    value={accountName}
                                    onChangeText={setAccountName}
                                />
                            </View>
                        )}
                    </View>

                    {/* Summary */}
                    {payoutAmount && parseFloat(payoutAmount) >= MIN_PAYOUT && (
                        <View style={styles.summaryCard}>
                            <Text style={styles.summaryTitle}>Payout Summary</Text>
                            <View style={styles.summaryRow}>
                                <Text style={styles.summaryLabel}>Requested Amount:</Text>
                                <Text style={styles.summaryValue}>{parseFloat(payoutAmount).toFixed(2)} Birr</Text>
                            </View>
                            <View style={styles.summaryRow}>
                                <Text style={styles.summaryLabel}>Platform Fee ({PLATFORM_FEE}%):</Text>
                                <Text style={styles.summaryValueFee}>-{feeAmount.toFixed(2)} Birr</Text>
                            </View>
                            <View style={styles.summaryDivider} />
                            <View style={styles.summaryRow}>
                                <Text style={styles.summaryLabelBold}>You'll Receive:</Text>
                                <Text style={styles.summaryValueBold}>{netAmount.toFixed(2)} Birr</Text>
                            </View>
                        </View>
                    )}
                </ScrollView>

                {/* Submit Button - Sticky */}
                <View style={styles.footer}>
                    <TouchableOpacity
                        style={[
                            styles.submitButton,
                            (!payoutAmount || parseFloat(payoutAmount) < MIN_PAYOUT || isSubmitting) &&
                                styles.submitButtonDisabled,
                        ]}
                        onPress={handleSubmit}
                        disabled={!payoutAmount || parseFloat(payoutAmount) < MIN_PAYOUT || isSubmitting}
                    >
                        {isSubmitting ? (
                            <ActivityIndicator color={COLORS.background} />
                        ) : (
                            <Text style={styles.submitButtonText}>Request Payout</Text>
                        )}
                    </TouchableOpacity>
                </View>
            </View>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    header: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: SPACING.m,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.surfaceHighlight,
        backgroundColor: COLORS.surface,
    },
    backButton: {
        padding: SPACING.s,
    },
    backIcon: {
        fontSize: 24,
        color: COLORS.textPrimary,
    },
    headerTitle: {
        fontSize: 20,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    placeholder: {
        width: 40,
    },
    contentWrapper: {
        flex: 1,
    },
    content: {
        flex: 1,
    },
    balanceCard: {
        margin: SPACING.l,
        padding: SPACING.xl,
        borderRadius: SIZES.radiusM,
        alignItems: 'center',
    },
    balanceLabel: {
        fontSize: 14,
        color: COLORS.background,
        opacity: 0.9,
        marginBottom: SPACING.xs,
    },
    balanceAmount: {
        fontSize: 36,
        fontWeight: 'bold',
        color: COLORS.background,
        marginBottom: SPACING.xs,
    },
    pendingText: {
        fontSize: 12,
        color: COLORS.background,
        opacity: 0.8,
    },
    infoBanner: {
        backgroundColor: 'rgba(255, 215, 0, 0.1)',
        marginHorizontal: SPACING.l,
        marginBottom: SPACING.l,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        borderWidth: 1,
        borderColor: COLORS.gold,
    },
    infoText: {
        fontSize: 12,
        color: COLORS.gold,
        textAlign: 'center',
    },
    section: {
        paddingHorizontal: SPACING.l,
        marginBottom: SPACING.xl,
    },
    sectionTitle: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.m,
    },
    quickAmounts: {
        flexDirection: 'row',
        gap: SPACING.m,
        flexWrap: 'wrap',
    },
    quickAmountBtn: {
        flex: 1,
        minWidth: '22%',
        backgroundColor: COLORS.surface,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        borderWidth: 2,
        borderColor: COLORS.surfaceHighlight,
        alignItems: 'center',
    },
    quickAmountBtnActive: {
        borderColor: COLORS.primary,
        backgroundColor: 'rgba(167, 255, 131, 0.1)',
    },
    quickAmountBtnDisabled: {
        opacity: 0.3,
    },
    quickAmountText: {
        fontSize: 16,
        fontWeight: '600',
        color: COLORS.textPrimary,
    },
    quickAmountTextActive: {
        color: COLORS.primary,
        fontWeight: 'bold',
    },
    quickAmountTextDisabled: {
        color: COLORS.textSecondary,
    },
    inputContainer: {
        marginBottom: SPACING.m,
    },
    inputLabel: {
        fontSize: 14,
        color: COLORS.textSecondary,
        marginBottom: SPACING.xs,
    },
    input: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        color: COLORS.textPrimary,
        fontSize: 16,
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    methodCard: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        marginBottom: SPACING.m,
        borderWidth: 2,
        borderColor: COLORS.surfaceHighlight,
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    methodCardSelected: {
        borderColor: COLORS.primary,
        backgroundColor: 'rgba(167, 255, 131, 0.1)',
    },
    methodLeft: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    methodIcon: {
        fontSize: 32,
        marginRight: SPACING.m,
    },
    methodName: {
        fontSize: 18,
        fontWeight: '600',
        color: COLORS.textPrimary,
    },
    methodDesc: {
        fontSize: 12,
        color: COLORS.textSecondary,
        marginTop: 2,
    },
    checkmark: {
        width: 28,
        height: 28,
        borderRadius: 14,
        backgroundColor: COLORS.primary,
        alignItems: 'center',
        justifyContent: 'center',
    },
    checkmarkText: {
        color: COLORS.background,
        fontSize: 18,
        fontWeight: 'bold',
    },
    summaryCard: {
        backgroundColor: COLORS.surface,
        marginHorizontal: SPACING.l,
        marginBottom: SPACING.xl,
        padding: SPACING.l,
        borderRadius: SIZES.radiusM,
    },
    summaryTitle: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.m,
    },
    summaryRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: SPACING.s,
    },
    summaryLabel: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    summaryValue: {
        fontSize: 14,
        fontWeight: '600',
        color: COLORS.textPrimary,
    },
    summaryValueFee: {
        fontSize: 14,
        fontWeight: '600',
        color: COLORS.error,
    },
    summaryLabelBold: {
        fontSize: 16,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    summaryValueBold: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.primary,
    },
    summaryDivider: {
        height: 1,
        backgroundColor: COLORS.surfaceHighlight,
        marginVertical: SPACING.m,
    },
    footer: {
        padding: SPACING.l,
        borderTopWidth: 1,
        borderTopColor: COLORS.surfaceHighlight,
        backgroundColor: COLORS.surface,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: -2 },
        shadowOpacity: 0.1,
        shadowRadius: 4,
        elevation: 5,
    },
    submitButton: {
        backgroundColor: COLORS.primary,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        alignItems: 'center',
        justifyContent: 'center',
    },
    submitButtonDisabled: {
        opacity: 0.5,
    },
    submitButtonText: {
        color: COLORS.background,
        fontSize: 18,
        fontWeight: 'bold',
    },
});

