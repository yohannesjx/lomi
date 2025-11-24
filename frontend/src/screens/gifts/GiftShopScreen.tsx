import React, { useState, useEffect } from 'react';
import {
    View,
    Text,
    StyleSheet,
    ScrollView,
    TouchableOpacity,
    Modal,
    Alert,
    ActivityIndicator,
    Dimensions,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { LinearGradient } from 'expo-linear-gradient';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { GiftService } from '../../api/services';

const { width } = Dimensions.get('window');
const GIFT_CARD_WIDTH = (width - SPACING.l * 3) / 2;

interface Gift {
    type: string;
    name: string;
    coin_price: number;
    etb_value: number;
    animation_url: string;
    sound_url: string;
}

interface CoinPack {
    id: string;
    name: string;
    etb_price: number;
    coins: number;
}

const COIN_PACKS: CoinPack[] = [
    { id: 'spark', name: 'Spark', etb_price: 55, coins: 600 },
    { id: 'flame', name: 'Flame', etb_price: 110, coins: 1300 },
    { id: 'blaze', name: 'Blaze', etb_price: 275, coins: 3500 },
    { id: 'inferno', name: 'Inferno', etb_price: 550, coins: 8000 },
    { id: 'galaxy', name: 'Galaxy', etb_price: 1100, coins: 18000 },
    { id: 'universe', name: 'Universe', etb_price: 5500, coins: 100000 },
];

export const GiftShopScreen = ({ navigation }: any) => {
    const [gifts, setGifts] = useState<Gift[]>([]);
    const [coinBalance, setCoinBalance] = useState(0);
    const [loading, setLoading] = useState(true);
    const [buyCoinsModalVisible, setBuyCoinsModalVisible] = useState(false);
    const [purchasing, setPurchasing] = useState<string | null>(null);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            const [giftsRes, balanceRes] = await Promise.all([
                GiftService.getShop(),
                GiftService.getWalletBalance(),
            ]);
            setGifts(giftsRes.gifts || []);
            setCoinBalance(balanceRes.coin_balance || 0);
        } catch (error: any) {
            console.error('Load data error:', error);
            Alert.alert('Error', 'Failed to load gift shop');
        } finally {
            setLoading(false);
        }
    };

    const handleBuyCoins = async (pack: CoinPack) => {
        setPurchasing(pack.id);
        try {
            const response = await GiftService.buyCoins(pack.id);
            // TODO: Open Telebirr payment URL
            Alert.alert(
                'Payment',
                `Redirecting to payment for ${pack.name} pack (${pack.etb_price} ETB)`,
                [
                    {
                        text: 'Cancel',
                        onPress: () => setPurchasing(null),
                    },
                    {
                        text: 'Pay',
                        onPress: () => {
                            // Open payment URL
                            if (response.payment_url) {
                                // In web: window.open(response.payment_url)
                                // In React Native: Linking.openURL(response.payment_url)
                                Alert.alert('Payment', 'Redirecting to Telebirr...');
                            }
                            setPurchasing(null);
                            setBuyCoinsModalVisible(false);
                        },
                    },
                ]
            );
        } catch (error: any) {
            Alert.alert('Error', error.message || 'Failed to initiate purchase');
            setPurchasing(null);
        }
    };

    const formatCoins = (coins: number) => {
        return new Intl.NumberFormat('en-US').format(coins);
    };

    if (loading) {
        return (
            <View style={styles.loadingContainer}>
                <ActivityIndicator size="large" color={COLORS.primary} />
            </View>
        );
    }

    return (
        <View style={styles.container}>
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                <ScrollView contentContainerStyle={styles.scrollContent}>
                    {/* Header with coin balance */}
                    <View style={styles.header}>
                        <Text style={styles.title}>Gift Shop</Text>
                        <View style={styles.balanceCard}>
                            <Text style={styles.balanceLabel}>Your Balance</Text>
                            <Text style={styles.balanceAmount}>
                                {formatCoins(coinBalance)} LC
                            </Text>
                            <Text style={styles.balanceETB}>
                                ≈ {formatCoins(coinBalance * 0.1)} ETB
                            </Text>
                        </View>
                        <TouchableOpacity
                            style={styles.buyCoinsButton}
                            onPress={() => setBuyCoinsModalVisible(true)}
                        >
                            <Text style={styles.buyCoinsButtonText}>Buy Coins</Text>
                        </TouchableOpacity>
                    </View>

                    {/* Gifts Grid */}
                    <View style={styles.giftsGrid}>
                        {gifts.map((gift) => (
                            <TouchableOpacity
                                key={gift.type}
                                style={styles.giftCard}
                                onPress={() => {
                                    // Navigate to gift detail or send gift
                                    navigation.navigate('GiftDetail', { gift });
                                }}
                            >
                                <View style={styles.giftIcon}>
                                    <Text style={styles.giftEmoji}>
                                        {getGiftEmoji(gift.type)}
                                    </Text>
                                </View>
                                <Text style={styles.giftName}>{gift.name}</Text>
                                <Text style={styles.giftPrice}>
                                    {formatCoins(gift.coin_price)} LC
                                </Text>
                                <Text style={styles.giftETB}>
                                    {gift.etb_value.toFixed(1)} ETB
                                </Text>
                            </TouchableOpacity>
                        ))}
                    </View>
                </ScrollView>
            </SafeAreaView>

            {/* Buy Coins Modal */}
            <Modal
                visible={buyCoinsModalVisible}
                transparent
                animationType="slide"
                onRequestClose={() => setBuyCoinsModalVisible(false)}
            >
                <View style={styles.modalOverlay}>
                    <View style={styles.modalContent}>
                        <Text style={styles.modalTitle}>Buy Coins</Text>
                        <ScrollView>
                            {COIN_PACKS.map((pack) => (
                                <TouchableOpacity
                                    key={pack.id}
                                    style={styles.packCard}
                                    onPress={() => handleBuyCoins(pack)}
                                    disabled={purchasing === pack.id}
                                >
                                    <View style={styles.packInfo}>
                                        <Text style={styles.packName}>{pack.name}</Text>
                                        <Text style={styles.packCoins}>
                                            {formatCoins(pack.coins)} LC
                                        </Text>
                                        <Text style={styles.packPrice}>
                                            {pack.etb_price} ETB
                                        </Text>
                                    </View>
                                    {purchasing === pack.id && (
                                        <ActivityIndicator
                                            size="small"
                                            color={COLORS.primary}
                                        />
                                    )}
                                </TouchableOpacity>
                            ))}
                        </ScrollView>
                        <TouchableOpacity
                            style={styles.closeButton}
                            onPress={() => setBuyCoinsModalVisible(false)}
                        >
                            <Text style={styles.closeButtonText}>Close</Text>
                        </TouchableOpacity>
                    </View>
                </View>
            </Modal>
        </View>
    );
};

const getGiftEmoji = (type: string): string => {
    const emojiMap: Record<string, string> = {
        rose: '🌹',
        heart: '❤️',
        diamond_ring: '💍',
        fireworks: '🎆',
        yacht: '🛥️',
        sports_car: '🏎️',
        private_jet: '✈️',
        castle: '🏰',
        universe: '🌌',
        lomi_crown: '👑',
    };
    return emojiMap[type] || '🎁';
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    safeArea: {
        flex: 1,
    },
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: COLORS.background,
    },
    scrollContent: {
        padding: SPACING.l,
    },
    header: {
        marginBottom: SPACING.xl,
    },
    title: {
        fontSize: 32,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.l,
    },
    balanceCard: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.l,
        marginBottom: SPACING.m,
        borderWidth: 1,
        borderColor: COLORS.primary,
    },
    balanceLabel: {
        fontSize: 14,
        color: COLORS.textSecondary,
        marginBottom: SPACING.xs,
    },
    balanceAmount: {
        fontSize: 36,
        fontWeight: 'bold',
        color: COLORS.primary,
        marginBottom: SPACING.xs,
    },
    balanceETB: {
        fontSize: 16,
        color: COLORS.textSecondary,
    },
    buyCoinsButton: {
        backgroundColor: COLORS.primary,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        alignItems: 'center',
    },
    buyCoinsButtonText: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.background,
    },
    giftsGrid: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        justifyContent: 'space-between',
    },
    giftCard: {
        width: GIFT_CARD_WIDTH,
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        marginBottom: SPACING.m,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    giftIcon: {
        width: 80,
        height: 80,
        borderRadius: 40,
        backgroundColor: COLORS.surfaceHighlight,
        justifyContent: 'center',
        alignItems: 'center',
        marginBottom: SPACING.s,
    },
    giftEmoji: {
        fontSize: 40,
    },
    giftName: {
        fontSize: 16,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.xs,
    },
    giftPrice: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.primary,
        marginBottom: SPACING.xs,
    },
    giftETB: {
        fontSize: 12,
        color: COLORS.textSecondary,
    },
    modalOverlay: {
        flex: 1,
        backgroundColor: 'rgba(0,0,0,0.8)',
        justifyContent: 'flex-end',
    },
    modalContent: {
        backgroundColor: COLORS.surface,
        borderTopLeftRadius: SIZES.radiusL,
        borderTopRightRadius: SIZES.radiusL,
        padding: SPACING.l,
        maxHeight: '80%',
    },
    modalTitle: {
        fontSize: 24,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.l,
    },
    packCard: {
        backgroundColor: COLORS.surfaceHighlight,
        borderRadius: SIZES.radiusM,
        padding: SPACING.l,
        marginBottom: SPACING.m,
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    packInfo: {
        flex: 1,
    },
    packName: {
        fontSize: 20,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.xs,
    },
    packCoins: {
        fontSize: 18,
        color: COLORS.primary,
        marginBottom: SPACING.xs,
    },
    packPrice: {
        fontSize: 16,
        color: COLORS.textSecondary,
    },
    closeButton: {
        backgroundColor: COLORS.surfaceHighlight,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        alignItems: 'center',
        marginTop: SPACING.m,
    },
    closeButtonText: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
});

